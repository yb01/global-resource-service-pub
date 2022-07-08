#!/usr/bin/env bash

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
