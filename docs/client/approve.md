# 预授权命令行工具

预授权命令行工具包含以下命令:

* `qoscli tx create-approve`
* `qoscli tx increase-approve`
* `qoscli tx decrease-approve`
* `qoscli tx use-approve`
* `qoscli tx cancel-approve`
* `qoscli query approve`


[预授权设计](../spec/txs/approve.md)


## create

```bash
$ qoscli tx create-approve --help
Create approve

Usage:
  qoscli tx create-approve [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --coins string          Coins to this approve. ex: 10qos,100qstars,50qsc
      --from string           Name or Address of approve creator
  -h, --help                  help for create-approve
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --to string             Name or Address of approve receiver
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- from  授权账户，keys中保存的name
- to    被授权账户地址
- coins 授权币种、币值列表，[amount1][coin1],[amount2][coin2],...，以半角逗号相隔，eg: 10qos,100qsc1,100qsc2

Arya向Sansa授权100个qos，100个qstar
```
$ qoscli tx create-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --coins 100qos,100qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"9917953D8CDE80F457CD072DBCE73A36449B7A7C","height":"333"}
```

## query

查询预授权
```bash
$ qoscli query approve --help
Query approve by from and to

Usage:
  qoscli query approve [flags]

Flags:
      --chain-id string   Chain ID of tendermint node
      --from string       Name or Address of approve creator
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for approve
      --indent            add indent to json response
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --to string         Name or Address of approve receiver
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- from  授权账户地址
- to    被授权账户地址
```bash
$ qoscli query approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --indent
{
  "from": "address1evmncf3z99a4uhq5n5yjwputfqmtjsuknv43fn",
  "to": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "qstar",
      "amount": "100"
    }
  ]
}
```

## increase

```bash
$ qoscli tx increase-approve --help
Increase approve

Usage:
  qoscli tx increase-approve [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --coins string          Coins to this approve. ex: 10qos,100qstars,50qsc
      --from string           Name or Address of approve creator
  -h, --help                  help for increase-approve
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --to string             Name or Address of approve receiver
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- from  授权账户，keys中保存的name
- to    被授权账户地址
- coins 授权币种、币值列表，[amount1][coin1],[amount2][coin2],...，以半角逗号相隔，eg: 10qos,100qsc1,100qsc2

Arya向Sansa增加授权100个qos，100个qstar
```bash
$ qoscli tx increase-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --coins 100qos,100qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"406"}
```

## decrease

```bash
$ qoscli tx decrease-approve --help
Decrease approve

Usage:
  qoscli tx decrease-approve [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --coins string          Coins to this approve. ex: 10qos,100qstars,50qsc
      --from string           Name or Address of approve creator
  -h, --help                  help for decrease-approve
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --to string             Name or Address of approve receiver
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- from  授权账户，keys中保存的name
- to    被授权账户地址
- coins 授权币种、币值列表，[amount1][coin1],[amount2][coin2],...，以半角逗号相隔，eg: 10qos,100qsc1,100qsc2

Arya向Sansa减少授权100个qos，100个qstar
```bash
$ qoscli tx decrease-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --coins 100qos,100qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"9DC18AD3CB0B59FCD354C267D8C22A1CC75E5624","height":"414"}
```

## use

```bash
$ qoscli tx use-approve --help
Use approve

Usage:
  qoscli tx use-approve [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --coins string          Coins to this approve. ex: 10qos,100qstars,50qsc
      --from string           Name or Address of approve creator
  -h, --help                  help for use-approve
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --to string             Name or Address of approve receiver
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- from  授权账户地址
- to    被授权账户，keys中保存的name
- coins 授权币种、币值列表，[amount1][coin1],[amount2][coin2],...，以半角逗号相隔，eg: 10qos,100qsc1,100qsc2

Sansa使用Arya向自己授权的10个qos，10个qstar
```bash
$ qoscli tx use-approve --from address1evmncf3z99a4uhq5n5yjwputfqmtjsuknv43fn --to Sansa --coins 10qos,10qstar
Password to sign with 'Sansa':
{"check_tx":{},"deliver_tx":{},"hash":"0573760D6B316E6695FBB63A56F2A20C0635FCAE","height":"437"}
```

## cancel

```bash
$ qoscli tx cancel-approve --help
Cancel approve

Usage:
  qoscli tx cancel-approve [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --from string           Name or Address of approve creator
  -h, --help                  help for cancel-approve
      --indent                add indent to json response
      --max-gas int           gas limit to set per tx
      --node string           <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --nonce int             account nonce to sign the tx
      --nonce-node string     tcp://<host>:<port> to tendermint rpc interface for some chain to query account nonce
      --qcp                   enable qcp mode. send qcp tx
      --qcp-blockheight int   qcp mode flag. original tx blockheight, blockheight must greater than 0
      --qcp-extends string    qcp mode flag. qcp tx extends info
      --qcp-from string       qcp mode flag. qcp tx source chainID
      --qcp-seq int           qcp mode flag.  qcp in sequence
      --qcp-signer string     qcp mode flag. qcp tx signer key name
      --qcp-txindex int       qcp mode flag. original tx index
      --to string             Name or Address of approve receiver
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- from  授权账户，keys中保存的name
- to    被授权账户地址

Arya取消向Sansa授权任何资产
```bash
$ qoscli tx cancel-approve --from Arya --to Sansa
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"484"}
```
