# Test Cases

- [缺失参数[hash]](./TestCase01.md)
- [指定的[hash]在区块链中不存在](./TestCase02.md)
- [查询区块链中已存在的交易[hash]](./TestCase03.md)

# Description
>     Query for a transaction by hash in a committed block.

>     在提交的块中按哈希查询交易。

# Example

执行交易后会返回交易hash，通过交易hash可查询交易详细信息。

根据hash `f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385` 查询交易信息：

```bash
$ qoscli query tx f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385 --indent
```
输出示例：

```bash
{
  "hash": "f5fc2c228cba754d5b95e49b02e81ff818f7b9140f1859d3797b09fb4aa12385",
  "height": "246",
  "tx": {
    "type": "qbase/txs/stdtx",
    "value": {
      "itx": [
        {
          "type": "approve/txs/TxCreateApprove",
          "value": {
            "Approve": {
              "from": "address1s348wvf49dfy64e6wafc90lcavp4lrd6xzhzhk",
              "to": "address1yqekgyy66v2cxzww6lqg6sdrsugjguxqws6mkf",
              "qos": "100",
              "qscs": null
            }
          }
        }
      ],
      "sigature": [
        {
          "pubkey": {
            "type": "tendermint/PubKeyEd25519",
            "value": "B/iatjhcJ4yFyHfGYKw2IneYGu2zG+ZOR8XmRUaji0A="
          },
          "signature": "VrsOsULJx86y8ch529zvl3Sh19TwGm/AldPlQhVWqhtg+calZmBrk25sD9HxCYijAt+ZUWMiLtPg3QZzCCqHAg==",
          "nonce": "1"
        }
      ],
      "chainid": "QOS",
      "maxgas": "100000"
    }
  },
  "result": {
    "gas_wanted": "100000",
    "gas_used": "15220",
    "tags": [
      {
        "key": "YWN0aW9u",
        "value": "Y3JlYXRlLWFwcHJvdmU="
      },
      {
        "key": "YXBwcm92ZS1mcm9t",
        "value": "YWRkcmVzczFzMzQ4d3ZmNDlkZnk2NGU2d2FmYzkwbGNhdnA0bHJkNnh6aHpoaw=="
      },
      {
        "key": "YXBwcm92ZS10bw==",
        "value": "YWRkcmVzczF5cWVrZ3l5NjZ2MmN4end3NmxxZzZzZHJzdWdqZ3V4cXdzNm1rZg=="
      }
    ]
  }
}
```

# Usage
```
  qoscli query tx [hash] [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag           | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:---------------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`       | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--chain-id`   | ✖        | string     | -                         | -           | Tendermint节点的链ID                        |
| -        | `--indent`     | ✖        | -          | -                         | -           | 向JSON响应添加缩进                             |
| `-n`     | `--node`       | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>` |
| -        | `--trust-node` | ✖        | -          | -                         | -           | 是否信任连接的完整节点（不验证其响应证据）                   |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |