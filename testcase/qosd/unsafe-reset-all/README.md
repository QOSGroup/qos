# Test Cases

- 暂无

# Description
>     Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state.

>     重置区块链数据库，删除通讯簿文件，并将priv_validator.json重置为genesis状态。

重置区块链数据库，删除地址簿文件，重置状态至初始状态。

# Usage
```
  qosd unsafe-reset-all [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag     | Required | Input Type | Default Input | Input Range | Description |
|:---------|:---------|:---------|:-----------|:--------------|:------------|:------------|
| `-h`     | `--help` | ✖        | -          | -             | -           | 帮助文档        |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
