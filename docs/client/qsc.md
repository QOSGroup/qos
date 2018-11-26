# QSC命令行工具

[QSC](../spec/txs/qsc.md)工具包含以下命令:

* `qoscli tx create-qsc`: 创建联盟币，发放联盟币。
* `qoscli tx issue-qsc`: 发行联盟币
* `qoscli query qsc`: 查询qsc信息


## create

> 创建QSC需要申请[CA]()

1. 创建QSC

```
$ qoscli tx create-qsc --help
create qsc

Usage:
  qoscli tx create-qsc [flags]

Flags:
      --accounts string    init accounts: Sansa,100;Lisa,100
      --async              broadcast transactions asynchronously
      --chain-id string    Chain ID of tendermint node
      --creator string     name of banker
      --desc string        description
      --extrate string     extrate: qos:qscxxx (default "1:280.0000")
  -h, --help               help for create-qsc
      --max-gas int        gas limit to set per tx
      --node string        <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --path-bank string   path of CA(banker)
      --path-qsc string    path of CA(qsc)
      --trust-node         Trust connected full node (don't verify proofs for responses) (default true)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- creator   创建账号名
- path-bank bank 证书位置
- path-qsc  qsc 证书位置
- qsc-chain qsc chainID
- accounts  初始发放地址币值集合，[addr1],[amount];[addr2],[amount2],...，eg：address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,100;address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,100。
该参数可为空，即只创建联盟币

> 可以通过`qoscli keys import`导入*creator*账户


```
$ qoscli tx create-qsc --creator=Arya --path-qsc="qsc.crt" --path-bank "banker.crt" --qsc-chain qunion-chain
```

2. 查询QOS绑定的chains

```
$ qoscli query store --path /store/qcp/subspace --data pubkey
```

3. 查询QOS绑定的QSCs

```
$ qoscli query store --path /store/qsc/subspace --data qsc
```



## query
```
qoscli query qsc --help
query qsc info by name

Usage:
  qoscli query qsc [qsc-name] [flags]

Flags:
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for qsc
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- qsc-name

```
$ qoscli query qsc QSC
```

## issue

```
$ qoscli tx issue-qsc --help
issue qsc

Usage:
  qoscli tx issue-qsc [flags]

Flags:
      --amount int        coin amount send to banker (default 100000)
      --async             broadcast transactions asynchronously
      --banker string     name of banker
      --chain-id string   Chain ID of tendermint node
  -h, --help              help for issue-qsc
      --max-gas int       gas limit to set per tx
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --qsc-name string   qsc name
      --trust-node        Trust connected full node (don't verify proofs for responses) (default true)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- qsc-name  qsc名
- banker    banker账户名
- amount    qsc币值

> 可以通过`qoscli keys import QSCBanker --file ~/banker.pri` 使用banker的私钥文件导入*QSCBanker*账户

```
$ qoscli tx issue-qsc --qsc-name=QSC --banker=QSCBanker --amount=10000
```
