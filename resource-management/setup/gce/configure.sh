#!/usr/bin/env bash
#
# Copyright 2022 Authors of Global Resource Service.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


set -o errexit
set -o nounset
set -o pipefail


# Use --retry-connrefused opt only if it's supported by curl.
CURL_RETRY_CONNREFUSED=""
if curl --help | grep -q -- '--retry-connrefused'; then
  CURL_RETRY_CONNREFUSED='--retry-connrefused'
fi



function validate-python {
  local ver=$(python3 -c"import sys; print(sys.version_info.major)")
  echo "python3 version: $ver"
  if [[ $ver -ne 3 ]]; then
    apt-get -y update
    apt-get install -y python3
    apt-get install -y python3-pip
    pip install pyyaml
  else
    echo "python3: $ver is running.."
  fi
}

function download-server-env {
  # Fetch server-env from GCE metadata server.
  (
    umask 077
    local -r tmp_server_env="/tmp/server-env.yaml"
    curl --fail --retry 5 --retry-delay 3 ${CURL_RETRY_CONNREFUSED} --silent --show-error \
      -H "X-Google-Metadata-Request: True" \
      -o "${tmp_server_env}" \
      http://metadata.google.internal/computeMetadata/v1/instance/attributes/server-env
    # Convert the yaml format file into a shell-style file.
    eval $(python3 -c '''
import pipes,sys,yaml
items = yaml.load(sys.stdin, Loader=yaml.BaseLoader).items()
for k, v in items:
  print("readonly {var}={value}".format(var = k, value = pipes.quote(str(v))))
''' < "${tmp_server_env}" > "${SERVER_HOME}/server-env")
    rm -f "${tmp_server_env}"
  )
}

# Get default service account credentials of the VM.
GCE_METADATA_INTERNAL="http://metadata.google.internal/computeMetadata/v1/instance"
function get-credentials {
  curl  --fail --retry 5 --retry-delay 3 ${CURL_RETRY_CONNREFUSED} --silent --show-error "${GCE_METADATA_INTERNAL}/service-accounts/default/token" -H "Metadata-Flavor: Google" -s | python3 -c \
    'import sys; import json; print(json.loads(sys.stdin.read())["access_token"])'
}

# intall-redis
function install-redis {
  local -r version="$1"
  if [ `uname -s` == "Linux" ]; then
    LINUX_OS=`uname -v |awk -F'-' '{print $2}' |awk '{print $1}'`
    if [ "$LINUX_OS" == "Ubuntu" ]; then
      UBUNTU_VERSION_ID=`grep VERSION_ID /etc/os-release |awk -F'"' '{print $2}'`

      echo "1. Install Redis on Ubuntu ......"
      REDIS_GPG_FILE=/usr/share/keyrings/redis-archive-keyring.gpg
      if [ -f $REDIS_GPG_FILE ]; then
          rm -f $REDIS_GPG_FILE
      fi 
      curl -fsSL https://packages.redis.io/gpg | gpg --dearmor -o $REDIS_GPG_FILE

      echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/redis.list

      apt-get update
      if [ "$UBUNTU_VERSION_ID" != "20.04" ]; then
          echo "The Ubuntu $UBUNTU_VERSION_ID is not currently supported and exit"
          return
      fi

      echo "Purge existing version of Redis ......"
      apt-get purge redis -y
      apt-get purge redis-server -y
      apt-get purge redis-tools -y

      echo "Install Redis 7.0.0 ......"
      apt-get install redis-tools=$version
      apt-get install redis-server=$version
      apt-get install redis=$version
      echo "End to install on Ubuntu ......"

      echo ""
      echo "2. Enable and Run Redis ......"
      echo "==============================="
      REDIS_CONF_Ubuntu=/etc/redis/redis.conf
      ls -alg $REDIS_CONF_Ubuntu

      sed -i -e "s/^supervised auto$/supervised systemd/g" $REDIS_CONF_Ubuntu
      egrep -v "(^#|^$)" $REDIS_CONF_Ubuntu |grep "supervised "

      sed -i -e "s/^appendonly no$/appendonly yes/g" $REDIS_CONF_Ubuntu
      egrep -v "(^#|^$)" $REDIS_CONF_Ubuntu |egrep "(appendonly |appendfsync )"

      ls -al /lib/systemd/system/ |grep redis

      systemctl restart redis-server.service
      systemctl status redis-server.service
    else
      echo ""
      echo "This Linux OS ($LinuxOS) is currently not supported and exit"
      return
    fi
  else
    echo ""
    echo "only ubuntu is currently supported"
    return
  fi

  echo ""
  echo "Sleeping for 5 seconds after Redis installation ......"
  sleep 5

  echo ""
  echo "3. Simply Test Redis ......"
  echo "==============================="
  which redis-cli
  echo "3.1) Test ping ......"
  redis-cli ping 

  echo ""
  echo "3.2) Test write key and value ......"
  redis-cli << EOF
SET server:name "fido"
GET server:name
EOF

  echo ""
  echo "3.3) Test write queue ......"
  redis-cli << EOF
lpush demos redis-macOS-demo
rpop demos
EOF

  echo ""
  echo "Sleep 5 seconds after Redis tests ..."
  sleep 5

  # Redis Persistence Options:
  #
  # 1.Redis Database File (RDB) persistence takes snapshots of the database at intervals corresponding to the save directives in the redis.conf file. The redis.conf file contains three default intervals. RDB persistence generates a compact file for data recovery. However, any writes since the last snapshot is lost.

  # 2. Append Only File (AOF) persistence appends every write operation to a log. Redis replays these transactions at startup to restore the database state. You can configure AOF persistence in the redis.conf file with the appendonly and appendfsync directives. This method is more durable and results in less data loss. Redis frequently rewrites the file so it is more concise, but AOF persistence results in larger files, and it is typically slower than the RDB approach

  echo ""
  echo "************************************************************"
  echo "*                                                          *"
  echo "* You are successful to install and configure Redis Server *"
  echo "*                                                          *"
  echo "************************************************************"
}

function setup-server-env {
  golang_version=${GOLANG_VERSION:-"1.17.11"}
  redis_version=${REDIS_VERSION:-"6:7.0.0-1rl1~focal1"}
  echo "Update apt."
  apt-get -y update

  echo "Install jq."
  apt-get -y install jq

  install-golang
  

  install-redis ${redis_version}
}

function install-golang {
  echo "Installinng golang."
  GOROOT="/usr/local/go"
  GOPATH="${SERVER_HOME}/go"
  
  wget https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz -P /tmp
  tar -C /usr/local -xzf /tmp/go${GOLANG_VERSION}.linux-amd64.tar.gz

  #export GOROOT=${GOROOT}
  #export GOPATH=${GOPATH}
  #export PATH=/usr/local/go/bin:$PATH
}


function gitclone-project {
  git --version &>/dev/null
  GIT_IS_AVAILABLE=$?
  if [ $GIT_IS_AVAILABLE -ne 0 ]; then
      echo "git doesn't exist, installing"
      apt-get -y update
      apt-get -y install git
  fi

  echo "git clone global resource service repo"
  if [ -d "${GIT_REPO}/global-resource-service" ]; then
      rm -r ${GIT_REPO}/global-resource-service
  fi
  mkdir -p ${GIT_REPO}
  cd ${GIT_REPO}
  git clone https://github.com/CentaurusInfra/global-resource-service.git
  cd ${GIT_REPO}/global-resource-service 
}

function set-broken-motd {
  cat > /etc/motd <<EOF
Broken (or in progress) resource management service setup! Check the initialization status
using the following commands.

Server instance:
  - sudo systemctl status server-installation
  - sudo systemctl status server-configuration

Simulator instance:
  - sudo systemctl status sim-installation
  - sudo systemctl status sim-configuration
EOF
}


######### Main Function ##########
# redirect stdout/stderr to a file
exec >> /var/log/server-init.log 2>&1
echo "Start to setup resource management service"
# if install fails, message-of-the-day (motd) will warn at login shell
set-broken-motd

SERVER_HOME="/home/grs"
SERVER_BIN="${SERVER_HOME}/bin"
GIT_REPO="${SERVER_HOME}/go/src"



#ensure-container-runtime
# validate or install python
validate-python
# download and source server-env
download-server-env
source "${SERVER_HOME}/server-env"

# setup server enviroment
setup-server-env

#gitclone-project

##TODO: add build to build cmd bin to avoid go run and git clone.
##TODO: add "too many open files" configuration

echo "Done for installing resource management server files, Please run and add 'export PATH=\$PATH:/usr/local/go/bin' into your shell profile."
