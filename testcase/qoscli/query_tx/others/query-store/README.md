# Test Cases

- [缺失必要参数`--path`与`--data`](./TestCase01.md)
- [提供必要参数`--path`与`--data`](./TestCase02.md)

# Description
>     Query store data by low level.

>     按低级(low level)查询存储数据。

# Example

QOS网络的存储内容均可通过下面指令查找：

`qoscli query store --path /store/<store_key>/subspace --data <query_data>`

主要参数：

- `--path`  存储位置
- `--data`  查询内容，以<query_data>开头的数据会被查出来

查询QOS网络中存储的ROOT CA 信息：

```bash
$ qoscli query store --path /store/acc/subspace --data account --indent
```

执行结果：

```bash
[
  {
    "key": "account:\ufffdjw15+RMW:wS\ufffd\ufffd\ufffd\ufffd\u0003_\ufffd\ufffd",
    "value": {
      "type": "qos/types/QOSAccount",
      "value": {
        "base_account": {
          "account_address": "address1s348wvf49dfy64e6wafc90lcavp4lrd6xzhzhk",
          "public_key": null,
          "nonce": "0"
        },
        "qos": "10000000000",
        "qscs": null
      }
    }
  }
]
```

# Usage
```
  qoscli query store [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag       | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:-----------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`   | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--indent` | ✖        | -          | -                         | -           | 向JSON响应添加缩进                             |
| `-n`     | `--node`   | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>` |
| -        | `--data`   | ✔        | string     | -                         | -           | (主要参数)查询内容，以<query_data>开头的数据会被查出来        |
| -        | `--path`   | ✔        | string     | -                         | -           | (主要参数)存储位置                                |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |