# Quick Start

此文档介绍如何搭建自己的QOS网络，单节点或集群方式。

## Install

目前给常见操作系统提供官方预编译安装包，点击这里[下载](/DOWNLOAD.md)。

## Single-node
* init
```bash
$ qosd init --chain-id qos-test
{
  "chain_id": "qos-test",
  "node_id": "1c3100c28a44f1facf45aa83e9aa3d8ff8ac6b1f",
  "app_message": "null"
}
```
注意init 可添加--home flag指定配置文件地址，默认在$HOME/.qosd
`init`操作后,通过执行`qosd add-genesis-validator`添加validator

* add-genesis-accounts

使用`qosd add-genesis-accounts`初始化account账户到配置文件中.

> 使用`qoscli keys add `创建account账户

```bash

$ qoscli keys add qosInitAcc
Enter a passphrase for your key:
Repeat the passphrase:

$ qoscli keys list

NAME:   TYPE:   ADDRESS:                                                PUBKEY:
qosInitAcc      local   address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8  4MFA7MtUl1+Ak3WBtyKxGKvpcu4e5ky5TfAC26cN+mQ=

```

初始化账户
```bash
$ qosd add-genesis-accounts address1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8,1000000qos
```

* config-root-ca

root CA用于校验[QSC](../spec/txs/qsc.md)和[QCP](../spec/txs/qcp.md)，不存在相关业务时可不配置。CA的获取和使用请查阅[CA 文档](../spec/ca.md)

使用`qosd config-root-ca`初始化root CA公钥到配置文件.
```bash
$ qosd add-genesis-validator --help

Config root CA

Usage:
  qosd config-root-ca [root.pub] [flags]

Flags:
  -h, --help   help for config-root-ca

Global Flags:
      --home string        directory for config and data (default "$HOME/.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
      
```
设置roort CA
```bash
$ qosd config-root-ca root.pub
```

查看genesis.json内容，确认配置成功。

* add-genesis-validator

使用`qosd add-genesis-validator`初始化validator到配置文件中，只有配置了validator才能正常打块。

```bash

$ qosd add-genesis-validator --help

pubkey is a tendermint validator pubkey. the public key of the validator used in
Tendermint consensus.

home node's home directory.

owner is account address.

ex: pubkey: {"type":"tendermint/PubKeyEd25519","value":"VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="}

example:

         qoscli add-genesis-validator --home "/.qosd/" --name validatorName --owner address1vdp54s5za8tl4dmf9dcldfzn62y66m40ursfsa --pubkey "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA=" --tokens 100

Usage:
  qosd add-genesis-validator [flags]

Flags:
      --description string   description
  -h, --help                 help for add-genesis-validator
      --name string          name for validator
      --owner string         account address
      --pubkey string        tendermint consensus validator public key
      --tokens int           bond tokens amount

Global Flags:
      --home string        directory for config and data (default "$HOME//.qosd")
      --log_level string   Log level (default "main:info,state:info,*:error")
      --trace              print out full stack trace on errors
```    

查看priv_validator.json
```bash
$ cat  $HOME/.qosd/config/priv_validator.json
{                                                                             
  "address": "CBB6D9DF3C19A897AEED6E387992106C0B16DF51",                      
  "pub_key": {                                                                
    "type": "tendermint/PubKeyEd25519",                                       
    "value": "PJ58L4OuZp20opx2YhnMhkcTzdEWI+UayicuckdKaTo="                   
  },                                                                          
  "last_height": "0",                                                         
  "last_round": "0",                                                          
  "last_step": 0,                                                             
  "priv_key": {                                                               
    "type": "tendermint/PrivKeyEd25519",                                      
    "value": "jISQomswckTLAS2QzN0HNMrIhsrfibgIlFDIWrVLZs48nnwvg65mnbSinHZiGcyG
RxPN0RYj5RrKJy5yR0ppOg=="                                                     
  }                                                                           
}                                                                             
```

使用上面的初始化账户地址作为owner
```bash
$ qosd add-genesis-validator --name validatorName --owner qosInitAcc --pubkey "PJ58L4OuZp20opx2YhnMhkcTzdEWI+UayicuckdKaTo=" --tokens 10 --description "I am the first validator." --home "$HOME/.qosd/"

```

主要参数说明:
- owner is account keyname or address store in your local keystore, run `qoscli keys list` can find it.
- name is your validator's name, you can name it as you like.
- pubkey is the `value` part of validator's pubkey.
- tokens means the voting power, LTE the QOS amount in your account. 

* start
```bash
$ qosd start --with-tendermint
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
      --genesis-accounts string      Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'
  -h, --help                         help for testnet
      --hostname-prefix string       Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...) (default "node")
      --name string               Moniker
      --n int                        Number of non-validators to initialize the testnet with
      --node-dir-prefix string       Prefix the directory name for each node with (node results in node0, node1, ...) (default "node")
      --o string                     Directory to store initialization data for the testnet (default "./mytestnet")
      --p2p-port int                 P2P Port (default 26656)
      --populate-persistent-peers    Update config of each node with the list of persistent peers build using either hostname-prefix or starting-ip-address (default true)
      --root-ca string               Config pubKey of root CA
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
- name                miniker
- root-ca             CA公钥
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
$ qosd start --home <directory_for_config_and_data> --with-tendermint

```
