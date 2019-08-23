# Test Cases

- 暂无

# Description
>     Query distribution validator period info.

>     查询分发(distribution)验证人周期信息。

# Example

`qoscli query validator-period --owner  <key_name_or_account_address>`

`key_name_or_account_address`为操作者账户地址或密钥库中密钥名字

查询`Arya`的节点漏块信息：
```bash
$ qoscli query validator-period --owner Arya
```

执行结果：
```bash
{
  "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "validator_pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
  },
  "fees": "0",
  "current_tokens": "4782741",
  "current_period": "15",
  "last_period": "14",
  "last_period_fraction": {
    "value": "1177.934327765593760252"
  }
}
```

# Usage
```
  qoscli query validator-period [flags]
```

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
| -        | `--owner`      | ✖        | string     | -                         | -           | (主要参数)`Owner`账户本地密钥库名字或账户地址             |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |