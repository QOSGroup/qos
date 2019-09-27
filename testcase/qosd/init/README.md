# Test Cases

- [缺失必选参数moniker](./TestCase01.md)
- [填写必选参数moniker](./TestCase02.md)
- [填写必选参数moniker和可选参数chain-id](./TestCase03.md)
- [填写必选参数moniker和可选参数overwrite](./TestCase04.md)
- [填写全部参数](./TestCase05.md)

# Description
>     Initialize validators's and node's configuration files.

>     初始化验证人和节点的配置文件。

初始化`genesis`、`priv-validator`、`p2p-node`文件，会在`$HOME/.qosd/`下创建`data`和`config`两个目录。
- `data`为空目录，用于存储网络启动后保存的数据，
- `config`中会生成`config.toml`，`genesis.json`，`node_key.json`，`priv_validator.json`四个文件。

# Usage
```
  qosd init [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag          | Required | Input Type | Default Input | Input Range | Description                                                        |
|:---------|:--------------|:---------|:-----------|:--------------|:------------|:-------------------------------------------------------------------|
| `-h`     | `--help`      | ✖        | -          | -             | -           | 帮助文档                                                               |
| -        | `--chain-id`  | ✖        | string     | -             | -           | (主要参数)`genesis.json`文件中的`chain-id`, 如果留白则随机生成，链ID一致的节点才能组成同一个P2P网络 |
| -        | `--moniker`   | ✔        | string     | -             | -           | (主要参数)设置验证人在P2P网络中的名称，与`config.toml`中`moniker`配置项对应，可后期修改          |
| `-o`     | `--overwrite` | ✖        | -          | -             | -           | (主要参数)是否覆盖已存在的`genesis.json`文件                                     |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
