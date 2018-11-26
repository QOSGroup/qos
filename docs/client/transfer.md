# 转账命令行工具

转账命令行工具:

* `qoscli tx transfer`

查询账户信息:

* `qoscli query account`

[转账设计](../spec/txs/transfer.md)

```
$ qoscli tx transfer --help
Transfer QOS and QSCs

Usage:
  qoscli tx transfer [flags]

Flags:
      --async              broadcast transactions asynchronously
      --chain-id string    Chain ID of tendermint node
  -h, --help               help for transfer
      --max-gas int        gas limit to set per tx
      --node string        <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --receivers string   Receivers, eg: address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,10qos,100qstar
      --senders string     Senders, eg: Arya,10qos,100qstar
      --trust-node         Trust connected full node (don't verify proofs for responses) (default true)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：
- senders   发送集合，账户参数传入存在keys中的name
- receivers 接收集合，账户参数传入地址

Arya向Sansa转账1个qos，1个qstar
```
$ qoscli tx transfer --senders=Arya,1qos,1qstar --receivers=address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh,1qos,1qstar
Password to sign with 'Arya':
{"check_tx":{},"deliver_tx":{},"hash":"21ECB72C8F51B3BD8E3CB9D59765003B9D78BE75","height":"40"}
```
可通过[账户命令行](account.md)查看账户状态，tx & txs命令行工具查看交易信息。
