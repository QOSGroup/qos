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

### Build from source code

>由于众所周知的网络原因，我们更推荐您[下载可执行版本](https://github.com/QOSGroup/qos/blob/master/docs/install/installation.md#download-runnable-files)

```
mkdir -p $GOPATH/src/github.com/QOSGroup
cd $GOPATH/src/github.com/QOSGroup
git clone https://github.com/QOSGroup/qos
make install
```
会安装qosd,qoscli到GOBIN目录下

### Download runnable files

根据您的系统，选择对应的[可执行文件](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md)

运行以下指令：
```
qosd version

qoscli version

```
显示版本号与网页一致，表示已成功安装