# QCP命令行工具

[QCP](../spec/txs/qcp.md)工具包含以下命令:

* `qoscli tx init-qcp`: 创建联盟链
* `qoscli query qcp`: 查询qcp信息

## init
> 初始化QCP需要申请[CA](../spec/ca.md)

```bash
$ qoscli tx init-qcp --help
init qcp

Usage:
  qoscli tx init-qcp [flags]

Flags:
      --async                 broadcast transactions asynchronously
      --chain-id string       Chain ID of tendermint node
      --creator string        address or name of creator
  -h, --help                  help for init-qcp
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
      --qcp.crt string        path of CA(QCP)
      --trust-node            Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```
主要参数：

- creator       创建账号
- qcp.crt       证书位置

> 可以通过`qoscli keys import`导入*creator*账户

```bash
$ qoscli tx init-qcp --creator qosInitAcc --qcp.crt qcp.crt
```

## query

```bash
$ qoscli query qcp --help
qcp subcommands

Usage:
  qoscli query qcp [command]

Available Commands:
  list        List all crossQcp chain's sequence info
  out         Get max sequence to outChain
  in          Get max sequence received from inChain
  tx          Query qcp out tx

Flags:
  -h, --help   help for qcp

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "$HOME/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli query qcp [command] --help" for more information about a command.
```
