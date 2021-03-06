# Test Cases

- [show-address](./TestCase01.md)

# Description
>     Shows this node's tendermint validator consensus address.

>     显示此节点的tendermint验证人共识地址。

# Usage
```
  qosd tendermint show-address [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag     | Required | Input Type | Default Input | Input Range | Description |
|:---------|:---------|:---------|:-----------|:--------------|:------------|:------------|
| `-h`     | `--help` | ✖        | -          | -             | -           | 帮助文档        |
| -        | `--json` | ✖        | -          | -             | -           | (主要参数)获取计算机可分析的输出 |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
