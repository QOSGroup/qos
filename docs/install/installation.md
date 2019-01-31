# 安装

推荐配置：
* 可以使用云服务器或独立机房，可持续不间断运行
* 带宽4M及以上，低延时公共网络
* 1核CPU，2G内存，50G硬盘存储空间

可通过**下载可执行文件**或**编译源码**两种方式安装QOS.

* `qosd 创建、初始化、启动QOS网络命令工具`
* `qoscli 客户端命令行集合，执行交易和查询信息`

## 下载可执行文件

[文件列表](https://github.com/QOSGroup/qos/blob/master/DOWNLOAD.md)页下载对应版本可执行文件

> 下载完成可跳过下面的编译源码步骤

## 编译源码

**安装 Go**

参照[官方文档](https://golang.org/doc/install)安装最新Go(1.11+)，并正确设置GOPATH, GOROOT等相关环境变量。

***Go modules***

包依赖管理采用go modules

设置GO111MODULE=on环境变量，或在相应ide开启go modules支持

***安装 QOS***

下载源码：
```bash
$ mkdir -p $GOPATH/src/github.com/QOSGroup
$ cd $GOPATH/src/github.com/QOSGroup
$ git clone https://github.com/QOSGroup/qos
$ cd qos
```

不同的QOS测试网络运行的qos代码可能不一样，编译前请切换到正确的qos代码版本。
按照[QOS测试网络](https://github.com/QOSGroup/qos-testnets)说明切换到正确代码版本，以测试网`capricorn-1000`版本要求为例：
```bash
$ git checkout v0.0.3
$ make install
```
执行以上命令会安装`qosd`,`qoscli`到GOPATH/bin目录下，中国大陆用户可能需要**科学上网**才能编译成功。


运行以下指令：
```bash
$ qosd version
$ qoscli version
```

确保正确安装。