# Test Cases

- 暂无

# Description
>     Query all delegations made to one validator.

>     查询一个验证人所接收的所有委托。

# Example

> 下面实例中假设:
> - `Arya` 地址为: `address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`
> - `Sansa` 地址为: `address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

查询`Arya`验证节点上的所有代理信息：
```bash
$ qoscli query delegations-to Arya
```

查询结果示例：
```bash
[
  {
    "delegator_address": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
    "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
    },
    "delegate_amount": "100",
    "is_compound": false
  }
  ...
]
```

# Usage
```
  qoscli query delegations-to [validator-owner] [flags]
```

`[validator-owner]`为接受委托的验证人账户地址或密钥库中密钥名字

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag           | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:---------------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`       | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--chain-id`   | ✖        | string     | -                         | -           | Tendermint节点的链ID                        |
| -        | `--height`     | ✖        | int        | -                         | -           | (可选)要查询的块高度，省略以获取最新的可证明块                |
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
