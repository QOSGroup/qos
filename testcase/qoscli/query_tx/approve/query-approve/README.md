# Test Cases

- [缺失必须参数`--from`，`--to`](./TestCase01.md)
- [参数`--from`，`--to`不合法](./TestCase02.md)
- [正常查询预授权](./TestCase03.md)

# Description
>     Query approve by from and to.

>     按来源(from)和目标(to)查询预授权(approve)。

# Example

> 下面实例中假设:
> - `Arya` 地址为: `address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`
> - `Sansa` 地址为: `address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

查询`Arya`对`Sansa`的预授权: 
```
$ qoscli query approve --from Arya --to Sansa
```
执行结果: 
```
$ qoscli query approve --from Arya --to Sansa
{
  "from": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "to": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "AOE",
      "amount": "100"
    }
  ]
}
```

# Usage
```
  qoscli query approve [flags]
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
| -        | `--from`       | ✔        | string     | -                         | -           | (主要参数)授权账户本地密钥库名字或账户地址                  |
| -        | `--to`         | ✔        | string     | -                         | -           | (主要参数)被授权账户本地密钥库名字或账户地址                 |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |