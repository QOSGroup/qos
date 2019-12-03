# Test Cases

- [缺失参数[name or address]](./TestCase01.md)
- [指定的[name or address]在密钥库中不存在](./TestCase02.md)
- [查询已存在的[name or address]](./TestCase03.md)

# Description
>     Query account info by address or name.

>     按地址(address)或名称(name)查询帐户信息。

# Example

查询账户
`qoscli query account <key_name_or_account_address>`

<key_name_or_account_address>为本地密钥库存储的密钥名字或对应账户的地址。

假设本地密钥库中`Arya`地址为`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`，且QOS网络中已经创建了`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`对应账号，可执行：
```bash
qoscli query account Arya --indent
```
或
```bash
qoscli query account address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --indent
```
输出类似如下信息：
```bash
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "dfYz3Zg+g1VFU52frAiKyXRU4wVulJMYgIuboPuBtZ4="
      },
      "nonce": "0"
    },
    "qos": "10000",
    "qscs": [
        {
            "coin_name": "AOE",
            "amount": "10000"
        }
    ]
  }
}
```
可以看到`Arya`持有10000个QOS、10000个AOE，更多账户说明请阅读[QOS账户设计](../spec/account.md)文档。

# Usage
```
  qoscli query account [name or address] [flags]
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

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |