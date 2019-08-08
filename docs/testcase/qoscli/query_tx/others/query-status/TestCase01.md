# Description
```
直接调用
```
# Input
- 原始输出：
```
$ qoscli query status
```
- 格式化输出：
```
$ qoscli query status --indent
```
# Output
- 原始输出：
```
$ qoscli query status
{"node_info":{"protocol_version":{"p2p":"7","block":"10","app":"0"},"id":"84efd4503c755e9bb9887917b4862f0798828acf","listen_addr":"tcp://0.0.0.0:26656","network":"test-chain","version":"0.32.0","channels":"4020212223303800","moniker":"testnet","other":{"tx_index":"on","rpc_address":"tcp://127.0.0.1:26657"}},"sync_info":{"latest_block_hash":"3552756AB9F15F613372EF6B360D79374BDFB664F2371D66608454A96522AD15","latest_app_hash":"D76A8BA5ACDE861FA2E67BD60B9EB171038D3DA72CBDAEBA9D5FF057AFD50562","latest_block_height":"2644","latest_block_time":"2019-08-06T03:18:46.1793191Z","catching_up":false},"validator_info":{"address":"02608373282DF4C924009356D94DF68A1D89F35A","pub_key":{"type":"tendermint/PubKeyEd25519","value":"stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="},"voting_power":"100000"}}
```
- 格式化输出：
```
$ qoscli query status --indent
{
  "node_info": {
    "protocol_version": {
      "p2p": "7",
      "block": "10",
      "app": "0"
    },
    "id": "84efd4503c755e9bb9887917b4862f0798828acf",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "test-chain",
    "version": "0.32.0",
    "channels": "4020212223303800",
    "moniker": "testnet",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://127.0.0.1:26657"
    }
  },
  "sync_info": {
    "latest_block_hash": "EC2E51FE4F699F72D76870C17A0FFD4AE971D4B45452AF71EB8CC9B2ABB33750",
    "latest_app_hash": "02F1C0394C257BB958D969088E26C2BFCCC58A596FC70A2F4D91AAB0C1F80951",
    "latest_block_height": "2653",
    "latest_block_time": "2019-08-06T03:19:31.2816995Z",
    "catching_up": false
  },
  "validator_info": {
    "address": "02608373282DF4C924009356D94DF68A1D89F35A",
    "pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "stxAw3cY2oTc5abe/8190af7FXlmxWUz/vIhkn/cgKw="
    },
    "voting_power": "100000"
  }
}
```