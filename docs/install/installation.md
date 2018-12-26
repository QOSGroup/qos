# 安装

## Install Go
参照[官方文档](https://golang.org/doc/install)安装最新Go(1.11+)，设置$GOPATH, $GOBIN, and $PATH
```
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
```

## Go modules
包依赖管理采用go modules

设置环境变量：
```
echo "export GO111MODULE=on" >> ~/.bash_profile
```
或在相应ide开启go modules支持

## Install QOS
```
mkdir -p $GOPATH/src/github.com/QOSGroup
cd $GOPATH/src/github.com/QOSGroup
git clone https://github.com/QOSGroup/qos
make install
```
会安装qosd,qoscli到GOBIN目录下，运行以下指令：
```
qosd version
qoscli version
```
确保一切正常。