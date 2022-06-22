#!/usr/bin/bash
#
# This script is used to quickly install configure Redis Server 
# on Ubuntun20.04/18.04/16.04 and MacOS(Darwin 20.6.0)
#
#    Running on Ubuntu 20.04: ./hack/redis_install.sh
#                              Redis 7.0.0 is installed as default
#    Running on Ubuntu 18.04:  /bin/bash ./hack/redis_install.sh
#                              Redis 7.0.0 is installed as default
#    Running on Ubuntu 16.04:  /bin/bash ./hack/redis_install.sh
#                              Redis 7.0.0 is installed as default
#
#    Running on MacOS:  /bin/bash ./hack/redis_install.sh
#                       Redis 7.0.0 is installed as default
#
# Reference: 
#    For Ubuntu: https://redis.io/docs/getting-started/installation/install-redis-on-linux/
#                https://www.linode.com/docs/guides/install-redis-ubuntu/
#
#    For MacOS:  https://redis.io/docs/getting-started/installation/install-redis-on-mac-os/
#    
#
export PATH=$PATH

if [ `uname -s` == "Linux" ]; then
  LINUX_OS=`uname -v |awk -F'-' '{print $2}' |awk '{print $1}'`

  if [ "$LINUX_OS" == "Ubuntu" ]; then
    UBUNTU_VERSION_ID=`grep VERSION_ID /etc/os-release |awk -F'"' '{print $2}'`

    echo "1. Install Redis on Ubuntu ......"
    REDIS_GPG_FILE=/usr/share/keyrings/redis-archive-keyring.gpg
    if [ -f $REDIS_GPG_FILE ]; then
      sudo rm -f $REDIS_GPG_FILE
    fi 
    curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o $REDIS_GPG_FILE

    echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list

    sudo apt-get update

    if [ "$UBUNTU_VERSION_ID" == "20.04" ]; then
      REDIS_VERSION="6:7.0.0-1rl1~focal1"
    elif [ "$UBUNTU_VERSION_ID" == "18.04" ]; then
      REDIS_VERSION="6:7.0.0-1rl1~bionic1"
    elif [ "$UBUNTU_VERSION_ID" == "16.04" ]; then
      REDIS_VERSION="6:7.0.0-1rl1~xenial1"
    else
      echo "The Ubuntu $UBUNTU_VERSION_ID is not currently supported and exit"
      exit 1
    fi

    sudo apt-get install redis=$REDIS_VERSION
    echo "End to install on Ubuntu ......"

    echo ""
    echo "2. Enable and Run Redis ......"
    echo "==============================="
    REDIS_CONF_Ubuntu=/etc/redis/redis.conf
    sudo ls -alg $REDIS_CONF_Ubuntu

    sudo sed -i -e "s/^supervised auto$/supervised systemd/g" $REDIS_CONF_Ubuntu
    sudo egrep -v "(^#|^$)" $REDIS_CONF_Ubuntu |grep "supervised "

    sudo sed -i -e "s/^appendonly no$/appendonly yes/g" $REDIS_CONF_Ubuntu
    sudo egrep -v "(^#|^$)" $REDIS_CONF_Ubuntu |egrep "(appendonly |appendfsync )"

    sudo ls -al /lib/systemd/system/ |grep redis

    sudo systemctl restart redis-server.service
    sudo systemctl status redis-server.service
  else
    echo ""
    echo "This Linux OS ($LinuxOS) is currently not supported and exit"
    exit 1
  fi
elif [ `uname -s` == "Darwin" ]; then
  echo "1. Install and configure Redis on MacOS ......"
  brew --version

  echo ""
  brew install redis=7.0
  brew services start redis
  brew services info redis --json

  echo "End to install Redis on MacOS ......"

  echo ""
  echo "2. Enable and Run Redis ......"
  echo "==============================="
  REDIS_CONF_MacOS=/usr/local/etc/redis.conf
  sed -i -e "s/^# supervised auto$/supervised systemd/g" $REDIS_CONF_MacOS
  egrep -v "(^#|^$)" $REDIS_CONF_MacOS |grep "supervised "

  #
  #Configure Redis Persistence using Append Only File (AOF)
  #
  sed -i -e "s/^appendonly no$/appendonly yes/g" $REDIS_CONF_MacOS
  egrep -v "(^#|^$)" $REDIS_CONF_MacOS |egrep "(appendonly |appendfsync )"

  brew services stop redis
  sleep 2
  brew services start redis
  brew services info redis --json
else
  echo ""
  echo "Unknown OS and exit"
  exit 1
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

exit 0
