## Set up developer environment
Recommended dev machine configuration: ubuntu 20.04, 8 cpu or above, disk size 100GB or up.
 
### Clone repo
```
$ mkdir -p go/src/
$ cd go/src/
$ git clone https://github.com/CentaurusInfra/global-resource-service.git
```

### Install Golang
```
$ sudo apt-get update
# $ sudo apt-get -y upgrade // optional
$ cd /tmp
$ wget https://dl.google.com/go/go1.17.11.linux-amd64.tar.gz
$ sudo tar -C /usr/local -xvf go1.17.11.linux-amd64.tar.gz
$ rm go1.17.11.linux-amd64.tar.gz
```
Add the following lines to ~/.profile
```
GOROOT=/usr/local/go
GOPATH=$HOME/go
PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```
Update the current shell session
```
$ source ~/.profile
```