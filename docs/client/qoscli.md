# QOS Client

提供与QOS交互的命令行工具。

```
$ qoscli --help
QOS light-client

Usage:
  qoscli [command]

Available Commands:

  keys        keys management tools. Add or view local private keys
  qcp         qcp subcommands
  query       query subcommands
  tx          tx subcommands
  tendermint  tendermint subcommands
  version     Print the app version
  help        Help about any command

Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
  -h, --help              help for qoscli
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli [command] --help" for more information about a command.
```
所有操作实例基于[Single-node](../install/networks.md#single-node)网络

## keys

[keystore](https://www.github.com/QOSGroup/qbase)

添加地址信息Sansa：
```
$ qoscli keys add Sansa
Enter a passphrase for your key:12345678
Repeat the passphrase:12345678
NAME:	TYPE:	ADDRESS:						PUBKEY:
Sansa	local	address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh	PubKeyEd25519{143CEBE483744337D6A1C785FDAF552E0FDCFB06008D87A57E925B92CA3F3E66}
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

sentence swap network level reason jewel radio apple soap vessel symptom improve mimic early wise real float clarify forward turkey lake actress typical twin
```

## qcp

[跨链交易数据查询](qcp.md)

## query

查询相关命令
```
$ qoscli query --help
query subcommands

Usage:
  qoscli query [command]

Available Commands:
  account     query account by address or name
  approve     Query approve by from and to
  qsc         query qsc info by name

Flags:
  -h, --help   help for query

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli query [command] --help" for more information about a command.
```
### account

[链上账户状态查询](account.md)

### approve

[预授权查询](approve.md)

### qsc

[联盟链信息查询](qsc.md)

## tx
```
$ qoscli tx --help
tx subcommands

Usage:
  qoscli tx [command]

Available Commands:
  create-qsc       create qsc
  issue-qsc        issue qsc

  transfer         Transfer QOS and QSCs

  create-approve   Create approve
  increase-approve Increase approve
  decrease-approve Decrease approve
  use-approve      Use approve
  cancel-approve   Cancel approve

  create-validator Create validator.

Flags:
  -h, --help   help for tx

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli tx [command] --help" for more information about a command.
```

### qsc

[创建、发放联盟币](qsc.md)

### transfer

[转账](transfer.md)

### approve

[预授权](approve.md)

### validator

验证节点

## tendermint
```
$ qoscli tendermint --help
tendermint subcommands

Usage:
  qoscli tendermint [command]

Available Commands:
  status      Query remote node for status

  validators  Get validator set at given height
  block       Get block info at given height

  txs         Search for all transactions that match the given tags.
  tx          query match hash tx in all commit block

Flags:
  -h, --help   help for tendermint

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli tendermint [command] --help" for more information about a command.
```

### status

运行状态

### validators

验证节点信息

### block

区块信息

### txs & tx

tx数据查询
