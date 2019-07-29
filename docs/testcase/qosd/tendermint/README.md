# Description
>     Tendermint subcommands.

>     Tendermint子命令. 

查询tendermint node相关信息.

# Usage
```
  qosd tendermint [command]
```

# Available Commands

| Command               | Alias | Has-Subcommand | Description             |
|:----------------------|:------|:---------------|:------------------------|
| `qosd tendermint show-address`   | -     | ✖              | 显示此节点的tendermint验证人共识地址 |
| `qosd tendermint show-node-id`   | -     | ✖              | 显示此节点的ID                |
| `qosd tendermint show-validator` | -     | ✖              | 显示此节点的tendermint验证人信息   |

# Flags

| ShortCut | Flag     | Input Type | Default Input | Input Range | Description |
|:---------|:---------|:-----------|:--------------|:------------|:------------|
| `-h`     | `--help` | -          | -             | -           | (可选)帮助文档        |

# Global Flags

| ShortCut | Flag          | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |

