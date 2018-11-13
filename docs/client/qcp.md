# QCP命令行工具

qos基于[qbase](https://www.github.com/QOSGroup/qbase)，提供跨链交易(QCP)支持
```
qoscli qcp --help
qcp subcommands

Usage:
  qoscli qcp [command]

Available Commands:
  outseq      Get max sequence to outChain
  outtx       Query qcp out tx 
  inseq       Get max sequence received from inChain

Flags:
  -h, --help   help for qcp

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli qcp [command] --help" for more information about a command.
```

* inseq

Get max sequence received from inChain
```
qoscli qcp inseq --chain-id=xxx
```
* outseq

Get max sequence  to outChain
```
qoscli qcp outseq --chain-id=xxx
```
* outtx

Query qcp out tx
```
qoscli qcp outtx --chain-id=xxx --seq=x
```