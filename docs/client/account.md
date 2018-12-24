# 账户命令行工具

[账户设计](../spec/account.md)，可通过用户名或账户地址查询。

```bash
$ qoscli query account --help
qoscli query account [name or address] [flags]

Usage:
  qoscli query account [flags]

Flags:
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for account
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

* by name
```bash
$ qoscli query account Arya
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1evmncf3z99a4uhq5n5yjwputfqmtjsuknv43fn",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "vPVE7hwjYcpnyAKaMKI7Y2FVTfFwImzL+4vUJHc1YEE="
      },
      "nonce": "1"
    },
    "qos": "100000000",
    "qscs": [
      {
        "coin_name": "qstar",
        "amount": "100000000"
      }
    ]
  }
}
```

* by address
```bash
$ qoscli query account address1evmncf3z99a4uhq5n5yjwputfqmtjsuknv43fn
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1evmncf3z99a4uhq5n5yjwputfqmtjsuknv43fn",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "vPVE7hwjYcpnyAKaMKI7Y2FVTfFwImzL+4vUJHc1YEE="
      },
      "nonce": "1"
    },
    "qos": "100000000",
    "qscs": [
      {
        "coin_name": "qstar",
        "amount": "100000000"
      }
    ]
  }
}
```
