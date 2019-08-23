# Test Cases

- 暂无

# Description
>     testnet will create "v" number of directories and populate each with necessary files (private validator, genesis, config, etc.).
>     Note, strict routability for addresses is turned off in the config file.

>     testnet将创建“v”数量的目录，并用必要的文件（私有验证器、Genesis、config等）填充每个目录。
>     注意，在配置文件中关闭了地址的严格可路由性。

批量生成集群配置文件.

# Example
```
  qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos
```

# Usage
```
  qosd testnet [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                    | Required | Input Type | Default Input   | Input Range | Description                                                                         |
|:---------|:------------------------|:---------|:-----------|:----------------|:------------|:------------------------------------------------------------------------------------|
| `-h`     | `--help`                | ✖        | -          | -               | -           | 帮助文档                                                                                |
| -        | `--chain-id`            | ✖        | string     | -               | -           | 区块链ID                                                                               |
| -        | `--compound`            | ✖        | -          | `true`          | -           | 验证人的收入是否计算为复利                                                                       |
| -        | `--genesis-accounts`    | ✖        | string     | -               | -           | 将Genesis帐户添加到`genesis.json`, 多个帐户以';'分隔                                             |
| -        | `--guardians`           | ✖        | string     | -               | -           | 将Guardian帐户添加到`genesis.json`, 多个帐户以','分隔                                            |
| -        | `--home-client`         | ✖        | string     | `"/.qoscli"`    | -           | keybase所在目录                                                                         |
| -        | `--hostname-prefix`     | ✖        | string     | `"node"`        | -           | 主机名前缀("node"导致持久化peer列表: id0@node0:26656, id1@node1:26656, ...)                     |
| -        | `--node-dir-prefix`     | ✖        | string     | `"node"`        | -           | 节点目录名前缀("node"导致: node0, node1, ...)                                                |
| -        | `--o`                   | ✖        | string     | `"./mytestnet"` | -           | 存储testnet初始化数据的目标路径                                                                 |
| -        | `--qcp-root-ca`         | ✖        | string     | -               | -           | 为QCP配置根CA的pubkey                                                                    |
| -        | `--qsc-root-ca`         | ✖        | string     | -               | -           | 为QSC配置根CA的pubkey                                                                    |
| -        | `--starting-ip-address` | ✖        | string     | -               | -           | 起始IP地址("192.168.0.1"导致持久化peer列表: ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...) |
| -        | `--v`                   | ✖        | int        | `4`             | -           | 用于初始化testnet的验证人数量                                                                  |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
