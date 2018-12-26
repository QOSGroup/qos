# QOS Client

提供与QOS交互的命令行工具。

```bash
$ qoscli --help
QOS light-client

Usage:
  qoscli [command]

Available Commands:

  keys        keys management tools. Add or view local private keys
  query       query(alias `q`) subcommands.
  tx          tx subcommands
  tendermint  tendermint(alias `t`)  subcommands
  version     Print the app version
  help        Help about any command

Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
  -h, --help              help for qoscli
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli [command] --help" for more information about a command.
```

## keys

[keystore](https://www.github.com/QOSGroup/qbase)

添加地址信息Sansa：
```bash
$ qoscli keys add Sansa
Enter a passphrase for your key:12345678
Repeat the passphrase:12345678
NAME:	TYPE:	ADDRESS:						PUBKEY:
Sansa	local	address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh	PubKeyEd25519{143CEBE483744337D6A1C785FDAF552E0FDCFB06008D87A57E925B92CA3F3E66}
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

sentence swap network level reason jewel radio apple soap vessel symptom improve mimic early wise real float clarify forward turkey lake actress typical twin
```


## query

查询相关命令
```bash
$ qoscli query --help
query(alias `q`) subcommands.

Usage:
  qoscli query [command]

Aliases:
  query, q

Available Commands:
  account     Query account info by address or name
  store       Query store data by low level
  qcp         qcp subcommands
  approve     Query approve by from and to
  qsc         query qsc info by name

Flags:
  -h, --help   help for query

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli query [command] --help" for more information about a command.
```
### account

[链上账户状态查询](account.md)

### store
通过abci query直接查询store的底层方法，path参数指定查询路径，data指定查询key
```bash
$ qoscli query store --help
Query store data by low level

Usage:
  qoscli query store [flags]

Flags:
      --data string   store query data
  -h, --help          help for store
      --indent        print indent result json
  -n, --node string   Node to connect to (default "tcp://localhost:26657")
      --path string   store query path

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

查询指定qsc
```bash
qoscli query store --path /store/qsc/key --data qsc/<qsc_name>
```
查询所有账户
```bash
qoscli query store --path /store/acc/subspace --data account:
```

### qcp

[跨链交易数据查询](qcp.md)

### approve

[预授权查询](approve.md)

### qsc

[联盟链信息查询](qsc.md)

## tx
```bash
$ qoscli tx --help
tx subcommands

Usage:
  qoscli tx [command]

Available Commands:
  create-qsc       create qsc
  issue-qsc        issue qsc

  init-qcp         init qcp

  transfer         Transfer QOS and QSCs

  create-approve   Create approve
  increase-approve Increase approve
  decrease-approve Decrease approve
  use-approve      Use approve
  cancel-approve   Cancel approve

  create-validator Create validator
  revoke-validator Revoke validator
  active-validator Active validator

Flags:
  -h, --help   help for tx

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
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

[验证节点](validator.md)

## tendermint
```bash
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
      --home string       directory for config and data (default "$HOME/.qoscli")
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
