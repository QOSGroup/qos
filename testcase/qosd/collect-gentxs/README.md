# Test Cases

- [gentx目录非空](./TestCase01.md)

# Description
>     Collect genesis txs and output a genesis.json file

>     收集创世交易(Genesis Txs)到genesis.json 

收集`gentx`目录下交易数据，填充到`genesis.json`中`app_state`下`gen_txs`中。

# Usage
```
  qosd collect-gentxs [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag          | Required | Input Type | Default Input                 | Input Range | Description     |
|:---------|:--------------|:---------|:-----------|:------------------------------|:------------|:----------------|
| `-h`     | `--help`      | ✖        | -          | -                             | -           | 帮助文档            |
| -        | `--gentx-dir` | ✖        | string     | `"$HOME/.qosd/config/gentx/"` | -           | (主要参数)gentx文件目录 |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
