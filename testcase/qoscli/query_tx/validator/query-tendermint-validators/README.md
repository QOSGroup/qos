# Test Cases

- 暂无

# Description
>     Get tendermint validator set at given height.

>     获取给定高度的tendermint验证人集合。

# Example

`qoscli query tendermint-validators <height>`

查询最新高度所有验证节点：
```bash
$ qoscli query tendermint-validators --indent
```

执行结果：
```bash
current query height: 260
[
  {
    "Address": "address1axqkgynrrdp2uwfpw60lm80pyx48g4pz5xj3er",
    "VotingPower": "1000",
    "PubKey": {
      "type": "tendermint/PubKeyEd25519",
      "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
    }
  }
]
```

# Usage
```
  qoscli query tendermint-validators [height] [flags]
```

# Available Commands

>无可用命令

# Flags


| ShortCut | Flag       | Required | Input Type | Default Input             | Input Range | Description                             |
|:---------|:-----------|:---------|:-----------|:--------------------------|:------------|:----------------------------------------|
| `-h`     | `--help`   | ✖        | -          | -                         | -           | 帮助文档                                    |
| -        | `--indent` | ✖        | -          | -                         | -           | 向JSON响应添加缩进                             |
| `-n`     | `--node`   | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>` |



# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |