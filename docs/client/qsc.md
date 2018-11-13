# QSC命令行工具

[QSC](../spec/txs/qsc.md)，创建联盟币，发放联盟币。

```
qoscli tx qsc
QSC subcommands

Usage:
  qoscli qsc [command]

Available Commands:
  query       query qsc info
  create      create qsc
  issue       issue qsc

Flags:
  -h, --help   help for qsc

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli qsc [command] --help" for more information about a command.
```

创建QSC需要申请[CA]()

## create

```
qoscli tx create-qsc --help
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
- accounts  初始发放地址币值集合，[addr1],[amount];[addr2],[amount2],...，eg：address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,100;address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,100。
该参数可为空，即只创建联盟币
```
qoscli tx create-qsc --creator=Arya --path-qsc="qsc.crt" --path-bank "banker.crt" 
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
qoscli query qsc QSC
``` 

## issue

```
qoscli tx issue-qsc --help
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

```
qoscli tx issue-qsc --qsc-name=QSC --banker=QSCBanker --amount=10000
```  
