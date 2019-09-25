# Test Cases

- [缺失必选参数`--tags`](./TestCase01.md)
- [指定错误的可选参数`--page`和`--limit`](./TestCase02.md)
- [查询区块链中已存在的交易](./TestCase03.md)

# Description
>     Search for transactions that match the exact given tags where results are paginated..

>     分页查询与一组tag匹配的交易。

# Example

```
$ <appcli> query txs --tags '<tag1>:<value1>&<tag2>:<value2>' --page 1 --limit 30
```

执行交易后会同时会返回QOS为交易所打tag，通过交易tag可查询交易信息。

根据`approve-from`=`address1s348wvf49dfy64e6wafc90lcavp4lrd6xzhzhk`查询预授权交易信息：

```bash
$ qoscli query txs --tag "approve-from='address1s348wvf49dfy64e6wafc90lcavp4lrd6xzhzhk'" --indent
```
输出示例：

```bash
[
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
      "gasWanted": "100000",
      "gasUsed": "15220",
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
]

```

更多交易Tag请查阅[index](../spec/indexing.md)

# Usage
```
  qoscli query txs [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag           | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:---------------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`       | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--indent`     | ✖        | -          | -                         | -           | 向JSON响应添加缩进                             |
| `-n`     | `--node`       | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>` |
| -        | `--trust-node` | ✖        | -          | -                         | -           | 是否信任连接的完整节点（不验证其响应证据）                   |
| -        | `--limit`      | ✖        | uint32     | `100`                     | -           | (主要参数)每页返回的交易结果查询数                            |
| -        | `--page`       | ✖        | uint32     | `1`                       | -           | (主要参数)查询分页结果的特定页面                             |
| -        | `--tags`       | ✔        | string     | -                         | -           | (主要参数)必须匹配的`tag:value`标记列表                    |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |