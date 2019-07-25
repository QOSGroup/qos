# Install

This guide will explain how to install the qosd and qoscli entry points onto your system. 
With these installed on a server, you can participate in the mainnet as either a Full Node or a Validator.

Recommended server specs:
- 1+ CPU
- Memory: 2+GB
- Disk: 50+GB SSD
- Bandwidth: 4+Mbps
- Allow all incoming connections on TCP port 26656 and 26657

We provide three ways to install QOS:

## Download runnable files

Visit [DOWNLOAD](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md) page to download the specific version.

## Docker

Build from source code:
```bash
$ mkdir -p $GOPATH/src/github.com/QOSGroup
$ cd $GOPATH/src/github.com/QOSGroup
$ git clone https://github.com/QOSGroup/qos
$ cd qos
$ docker build -t qos .
```
or pull the official images：
```bash
$ docker pull qoschain/qos:latest
$ docker tag qoschain/qos:latest qos:latest
```

List local images：
```bash
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
qos                 latest              741b61ad6fdd        9 seconds ago       65.5MB
qoschain/qos        latest              1fe9ca3e4cac        3 minutes ago       65.5MB
```

Set alias:
```bash
$ alias qosd='docker run --rm -v $HOME/.qosd:/root/.qosd -v $HOME/.qoscli:/root/.qoscli -p 26657:26657 -p 26656:26656 --name qosd -d qos qosd'
$ alias qoscli='docker run --rm -v $HOME/.qosd:/root/.qosd -v $HOME/.qoscli:/root/.qoscli --link qosd:qosd qos qoscli --node qosd:26657'
```

## Compile source code

**Install Go**

Install go(1.11.5+) by following the [official docs](https://golang.org/doc/install). Do not forget to set your $GOPATH, $GOBIN, and $PATH environment variables.

***Go modules***

We use go modules to manage package dependency. Set environment variable GO111MODULE=on on your server.

***Install QOS***

Download source code:
```bash
$ mkdir -p $GOPATH/src/github.com/QOSGroup
$ cd $GOPATH/src/github.com/QOSGroup
$ git clone https://github.com/QOSGroup/qos
$ cd qos
```

Different QOS networks running different code, install the right version from [QOS testnet](https://github.com/QOSGroup/qos-testnets)，for example:
```bash
# source code for `capricorn-1000`
$ git checkout v0.0.3
$ make install
```
`qosd`,`qoscli` will be installed.


Now check your QOS version:
```bash
$ qosd version
$ qoscli version
```