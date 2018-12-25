# 转账命令行工具

转账命令行工具:

* `qoscli tx transfer`

查询账户信息:

* `qoscli query account`

[转账设计](../spec/txs/transfer.md)

```bash
$ qoscli tx transfer --help
Transfer QOS and QSCs

Usage:
  qoscli tx transfer [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
  -h, --help                  help for transfer
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
      --receivers string      Receivers, eg: address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,10qos,100qstar. multiple users separated by ';'
      --senders string        Senders, eg: Arya,10qos,100qstar. multiple users separated by ';' 
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：
- senders   发送集合，账户传keystore name 或 address，多个账户半角分号分隔
- receivers 接收集合，账户传keystore name 或 address，多个账户半角分号分隔

Arya向Sansa转账1个qos，1个qstar
```bash
$ qoscli tx transfer --senders Arya,1qos,1qstar --receivers address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh,1qos,1qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"21ECB72C8F51B3BD8E3CB9D59765003B9D78BE75","height":"40"}
```
可通过[账户命令行](account.md)查看账户状态，tx & txs命令行工具查看交易信息。
