# Networks

此文档介绍如何搭建自己的QOS网络，单节点或集群方式。

## Single-node
* init

参照[初始化](../client/qosd.md#初始化) 执行：
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
qosInitAcc      local   address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8  4MFA7MtUl1+Ak3WBtyKxGKvpcu4e5ky5TfAC26cN+mQ=

```
更多本地秘钥库相关指令参照[qoscli keys](../client/qoscli.md#密钥（keys）)

参照[设置账户](../client/qosd.md#设置账户) 初始化账户信息：
```bash
$ qosd add-genesis-accounts address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8,1000000qos
```

* config-root-ca

root CA用于校验[QSC](../spec/txs/qsc.md)和[QCP](../spec/txs/qcp.md)，不存在相关业务时**可不配置**。CA的获取和使用请查阅[CA 文档](../spec/ca.md)

使用`qosd config-root-ca`初始化root CA公钥到配置文件。

```bash
$ qosd config-root-ca --qcp <qcp-root.pub> --qsc <qsc-root.pub>
```

更多操作说明查看[设置CA](../client/qosd.md#设置ca) 

查看genesis.json内容，确认配置成功。

* add-genesis-validator

使用`qosd add-genesis-validator`初始化validator到配置文件中，只有配置了validator才能正常运行QOS网络。

使用上面的初始化账户地址作为owner
```bash
$ qosd add-genesis-validator --name validatorName --owner address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8 --tokens 10 --description "I am the first validator."
```

主要参数说明:
- `--owner`         操作者账户地址
- `--name`          验证节点名字
- `--tokens`        绑定tokens，不能大于操作者持有QOS数量
- `--description`   备注
- `--compound`      收益复投方式，默认false，即收益不复投

更多操作说明参照[设置验证节点](../client/qosd.md#设置验证节点)

* start
```bash
$ qosd start --log_level debug
```
如果一切正常，会看到控制台输出打块信息

## Cluster

### qosd testnet
`qosd testnet`命令行工具，可批量生成集群配置文件，相关命令参考：
```bash
$ qosd testnet --help
testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars

Usage:
  qosd testnet [flags]

Flags:
      --chain-id string              Chain ID
      --compound                     whether the validator's income is calculated as compound interest, default: true (default true)
      --genesis-accounts string      Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'
  -h, --help                         help for testnet
      --hostname-prefix string       Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")
      --moniker string               Moniker
      --n int                        Number of non-validators to initialize the testnet with
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
      --o string                     Directory to store initialization data for the testnet (default "./mytestnet")
      --p2p-port int                 P2P Port (default 26656)
      --populate-persistent-peers    Update config of each node with the list of persistent peers build using either hostname-prefix or starting-ip-address (default true)
      --qcp-root-ca string           Config pubKey of root CA for QSC
      --qsc-root-ca string           Config pubKey of root CA for QCP
      --starting-ip-address string   Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)
      --v int                        Number of validators to initialize the testnet with (default 4)

Global Flags:
      --home string        directory for config and data (default "$HOME/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors

```

主要参数说明：
- chain-id            链ID
- genesis-accounts    初始账户
- hostname-prefix     hostName前缀
- miniker             miniker
- qcp-root-ca         pubKey of root CA for QCP
- qsc-root-ca         pubKey of root CA for QSC
- compound            收益复投方式，默认true，即收益参与复投
- starting-ip-address 起始IP地址

假设第一台机器IP: 192.168.1.100
```bash
$ qosd testnet --v 4 --name capricorn --starting-ip-address 192.168.1.100
Successfully initialized 4 node directories

```
会在当前目录下生成mytestnet文件夹，分别放置node0-3配置文件。
其中priv_validator_owner.json为对应validator owner私钥，可通过`qoscli keys import`导入。

`qosd testnet` 默认初始化owner 1000000QOS，validator bond 1000tokens。

### start
启动前请确保按照[安装说明](installation.md)在四台机器上正确安装QOS。
拷贝node0-3至不同机器，分别执行：
```bash
$ qosd start --home <directory_for_config_and_data>

```