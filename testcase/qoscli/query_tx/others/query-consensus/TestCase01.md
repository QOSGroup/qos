# Description
```
直接调用
```
# Input
- 原始输出：
```
$ qoscli query consensus
```
- 格式化输出：
```
$ qoscli query consensus --indent
```
# Output
- 原始输出：
```
$ qoscli query consensus
{"type":"abci/consensus/ConsensusParams","value":{"block":{"max_bytes":"1048576","max_gas":"-1"},"evidence":{"max_age":"100000"},"validator":{"pub_key_types":["ed25519"]}}}
```
- 格式化输出：
```
$ qoscli query consensus --indent
{
  "type": "abci/consensus/ConsensusParams",
  "value": {
    "block": {
      "max_bytes": "1048576",
      "max_gas": "-1"
    },
    "evidence": {
      "max_age": "100000"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    }
  }
}
```