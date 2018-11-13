# 账户命令行工具

[账户设计](../spec/account.md)，可通过用户名或账户地址查询。

```
qoscli query account --help
query account by address or name

Usage:
  qoscli query account [flags]

Flags:
      --addr string       address of account
      --chain-id string   Chain ID of tendermint node
      --height int        block height to query, omit to get most recent provable block
  -h, --help              help for account
      --name string       name of account
      --node string       <host>:<port> to tendermint rpc interface for this chain (default "tcp://localhost:26657")
      --trust-node        Trust connected full node (don't verify proofs for responses)

Global Flags:
  -e, --encoding string   Binary encoding (hex|b64|btc) (default "hex")
      --home string       directory for config and data (default "/home/imuge/.qoscli")
  -o, --output string     Output format (text|json) (default "text")
      --trace             print out full stack trace on errors
```

* by name
```
qoscli query account --name=Arya
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
```
qoscli query account --addr=address1evmncf3z99a4uhq5n5yjwputfqmtjsuknv43fn
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