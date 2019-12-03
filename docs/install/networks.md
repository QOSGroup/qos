# Networks

此文档介绍如何搭建自己的QOS网络，单节点或集群方式。

## Single-node
* init

参照[初始化](../command/qosd.md#初始化) 执行：
```bash
$ qosd init --moniker moniker --chain-id qos-test
{
 "moniker": "moniker",
 "chain_id": "qos-test",
 "node_id": "66853240dc1b26e6f6b35afcf008658823542076",
 "gentxs_dir": "",
 "app_message": {
  "accounts": null,
  "mint": {
   "params": {
    "total_amount": "10000000000",
    "total_block": "6307200"
   }
  },
  "stake": {
   "params": {
    "max_validator_cnt": 10,
    "voting_status_len": 100,
    "voting_status_least": 50,
    "survival_secs": 600
   },
   "validators": null
  },
  "qcp": {
   "ca_root_pub_key": null
  },
  "qsc": {
   "ca_root_pub_key": null
  }
 }
}
```
默认在$HOME/.qosd目录下生成配置文件。

* add-genesis-accounts

使用`qosd add-genesis-accounts`初始化account账户到配置文件中.

> 使用`qoscli keys add `创建account公私钥和地址信息

```bash

$ qoscli keys add qosInitAcc
Enter a passphrase for your key:
Repeat the passphrase:

$ qoscli keys list

NAME:   TYPE:   ADDRESS:                                                PUBKEY:
qosInitAcc      local   qosacc1hqcz9hhxa7qqxghc276vxxgcd3qkr279nz5gfq  qosaccpub1zcjduepqfzd5r2hzdnz58pjc9xuw5r2ez8f4khhwtekfxdjyvkvhrly6rxzqll3fgz

```
更多本地秘钥库相关指令参照[qoscli keys](../command/qoscli.md#密钥)

参照[设置账户](../command/qosd.md#设置账户) 初始化账户信息：
```bash
$ qosd add-genesis-accounts qosacc1hqcz9hhxa7qqxghc276vxxgcd3qkr279nz5gfq,49000000000000qos
```

* config-root-ca

root CA用于校验[QSC](../spec/qsc)和[QCP](../spec/qcp)，不存在相关业务时**可不配置**。CA的获取和使用请查阅[CA 文档](../spec/ca.md)

使用`qosd config-root-ca`初始化root CA公钥到配置文件。

```bash
$ qosd config-root-ca --qcp <qcp-root.pub> --qsc <qsc-root.pub>
```

更多操作说明查看[设置CA](../command/qosd.md#设置ca)

查看genesis.json内容，确认配置成功。

* create-validator

使用`qosd gentx`和`qosd collect-gentxs`初始化validator到配置文件中，只有配置了validator才能正常运行QOS网络。

使用上面的初始化账户地址作为owner
```bash
$ qosd gentx --moniker validatorName --owner qosacc1hqcz9hhxa7qqxghc276vxxgcd3qkr279nz5gfq --tokens 10
```

主要参数说明:
- `--owner`         操作者账户地址
- `--moniker`       验证节点名字
- `--logo`          logo
- `--website`       网址
- `--details`       详细描述信息
- `--tokens`        绑定tokens，不能大于操作者持有QOS数量
- `--compound`      收益复投方式，默认false，即收益不复投

更多操作说明参照[生成创世交易](../command/qosd.md#生成创世交易)

运行：
```bash
$ qosd collect-gentxs
```
将创建验证节点交易写入`genesis.json`文件中。

* start
```bash
$ qosd start --log_level debug
```
如果一切正常，会看到控制台输出打块信息

## Cluster

* qosd testnet
[qosd-testnet](../command/qosd.md#初始化测试网络)命令可以批量生成一个测试网络多个验证节点配置信息

假设第一台机器IP: 192.168.1.100
```bash
$ qosd testnet --v 4 --name capricorn --starting-ip-address 192.168.1.100
Successfully initialized 4 node directories

```
会在当前目录下生成mytestnet文件夹，分别放置node0-3配置文件。
其中priv_validator_owner.json为对应validator owner私钥，可通过`qoscli keys import`导入。

* start
启动前请确保按照[安装说明](installation.md)在四台机器上正确安装QOS。
拷贝node0-3至不同机器，分别执行：
```bash
$ qosd start --home <directory_for_config_and_data>
```
