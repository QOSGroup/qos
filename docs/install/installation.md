# 安装

可通过**下载可执行文件**和**编译源码**两种方式安装QOS.

* `qosd 创建、初始化、启动QOS网络命令工具`
* `qoscli 客户端命令行集合，执行交易和查询信息`

## Download runnable files

[Download](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md)

## Build from source code

### Install Go
参照[官方文档](https://golang.org/doc/install)安装最新Go(1.11+)，并正确设置GOPATH, GOROOT等相关环境变量。

### Go modules
包依赖管理采用go modules

设置GO111MODULE环境变量，或在相应ide开启go modules支持

### Install QOS

不同的QOS网络运行的qos代码可能不一致，编译前请切换到正确的qos代码版本。
```
mkdir -p $GOPATH/src/github.com/QOSGroup
cd $GOPATH/src/github.com/QOSGroup
git clone https://github.com/QOSGroup/qos
make install
```
执行以上命令会安装qosd, qoscli到GOPATH/bin目录下

运行以下指令：
```
qosd version
qoscli version
```
确保正确安装。