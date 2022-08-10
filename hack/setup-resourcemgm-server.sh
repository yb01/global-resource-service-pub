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


# Convenience script to setup a fresh Linux installation for resource management service.

set -o errexit
set -o nounset
set -o pipefail

echo "The script is to help install prerequisites of resource management service"
echo "on a fresh Linux installation."

GOLANG_VERSION=${GOLANG_VERSION:-"1.17.11"}

echo "Update apt."
sudo apt-get -y update

echo "Install jq."
sudo apt-get -y install jq

echo "Install golang."
wget https://dl.google.com/go/go${GOLANG_VERSION}.linux-amd64.tar.gz -P /tmp
sudo tar -C /usr/local -xzf /tmp/go${GOLANG_VERSION}.linux-amd64.tar.gz
echo 'export PATH=/usr/local/go/bin:$PATH' >>~/.bash_profile
source ~/.bash_profile

echo "Done."
echo "Please run and add 'export PATH=\$PATH:/usr/local/go/bin' into your shell profile."
echo "You can proceed to run ./setup/grs-up.sh if you want to start resource management service."
