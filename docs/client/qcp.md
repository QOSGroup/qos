# QCP命令行工具

qos基于[qbase](https://www.github.com/QOSGroup/qbase)，提供跨链交易(QCP)查询支持
```
qoscli query qcp --help
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
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors

Use "qoscli query qcp [command] --help" for more information about a command.
```
