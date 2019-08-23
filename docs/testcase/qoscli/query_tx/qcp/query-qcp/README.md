# Description

>     qcp subcommands. 

>     QCP子命令。

跨链协议是[qbase](https://www.github.com/QOSGroup/qbase)提供支持，主要有以下四个查询指令：
- `qoscli query qcp list`
- `qoscli query qcp out` 
- `qoscli query qcp in`
- `qoscli query qcp tx`

指令说明请参照[qbase-Qcp](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#Qcp)。

# Usage
```
  qoscli query qcp [command]
```

# Alias
```
  qoscli q qcp [command]
```

# Available Commands

| Command                 | Alias               | Has-Subcommand | Description        |
|:------------------------|:--------------------|:---------------|:-------------------|
| `qoscli query qcp list` | `qoscli q qcp list` | ✖              | 列出所有CrossQCP链的序列信息 |
| `qoscli query qcp out`  | `qoscli q qcp out`  | ✖              | 获取到OutChain的最大序列   |
| `qoscli query qcp in`   | `qoscli q qcp in`   | ✖              | 获取从InChain接收的最大序列  |
| `qoscli query qcp tx`   | `qoscli q qcp tx`   | ✖              | 查询QCP Out Tx       |

# Flags

| ShortCut | Flag      | Input Type | Default Input | Input Range | Description            |
|:---------|:----------|:-----------|:--------------|:------------|:-----------------------|
| `-h`     | `--help`  | -          | -             | -           | (可选)帮助文档                   |

# Global Flags

| ShortCut | Flag         | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |
