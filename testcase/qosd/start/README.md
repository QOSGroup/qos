# Test Cases

- 暂无

# Description
>     Run the full node.

>     启动全节点 

启动QOS网络.

# Usage
```
  qosd start [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                              | Required | Input Type | Default Input             | Input Range                               | Description                                                                                    |
|:---------|:----------------------------------|:---------|:-----------|:--------------------------|:------------------------------------------|:-----------------------------------------------------------------------------------------------|
| `-h`     | `--help`                          | ✖        | -          | -                         | -                                         | 帮助文档                                                                                       |
| -        | `--abci`                          | ✖        | string     | `"socket"`                | `"socket"`, `"grpc"`                      | 指定ABCI端口                                                                                       |
| -        | `--address`                       | ✖        | string     | `"tcp://0.0.0.0:26658"`   | -                                         | 监听地址                                                                                           |
| -        | `--consensus.create_empty_blocks` | ✖        | -          | `true`                    | -                                         | 将此值设置为false，以仅在存在tx或apphash更改时生成块。                                                             |
| -        | `--fast_sync`                     | ✖        | -          | `true`                    | -                                         | 是否启用快速区块链同步                                                                                    |
| -        | `--moniker`                       | ✖        | string     | `"GX8-CR6S1"`             | <your_computer_name>                      | 节点名称                                                                                           |
| -        | `--p2p.laddr`                     | ✖        | string     | `"tcp://0.0.0.0:26656"`   | -                                         | 节点监听地址（0.0.0.0:0表示任何接口、任何端口）                                                                   |
| -        | `--p2p.persistent_peers`          | ✖        | string     | -                         | -                                         | 逗号分隔的“ID@host:port”持久化peers                                                                    |
| -        | `--p2p.pex`                       | ✖        | -          | `true`                    | -                                         | 是否启用peer交换(PEX, Peer-Exchange)                                                                 |
| -        | `--p2p.private_peer_ids`          | ✖        | string     | -                         | -                                         | 逗号分隔的私有peerID                                                                                  |
| -        | `--p2p.seed_mode`                 | ✖        | -          | -                         | -                                         | 是否启用种子模式(seed mode)                                                                            |
| -        | `--p2p.seeds`                     | ✖        | string     | -                         | -                                         | 逗号分隔的"ID@host:port"种子节点                                                                        |
| -        | `--p2p.upnp`                      | ✖        | -          | -                         | -                                         | 是否启用UPNP端口转发                                                                                   |
| -        | `--priv_validator_laddr`          | ✖        | string     | -                         | -                                         | 监听外部私有验证人进程连接的套接字地址                                                                            |
| -        | `--proxy_app`                     | ✖        | string     | `"tcp://127.0.0.1:26658"` | -                                         | 代理应用程序地址，或者是'kvstore', 'persistent_kvstore', 'counter', 'counter_serial'其中之一， 或者是'noop'以用于本地测试 |
| -        | `--pruning`                       | ✖        | string     | `"syncable"`              | `"syncable"`, `"nothing"`, `"everything"` | 修剪策略                                                                                           |
| -        | `--rpc.grpc_laddr`                | ✖        | string     | -                         | -                                         | GRPC侦听地址（仅限BroadcastTX）。需要端口                                                                   |
| -        | `--rpc.laddr`                     | ✖        | string     | `"tcp://127.0.0.1:26657"` | -                                         | RPC侦听地址。需要端口                                                                                   |
| -        | `--rpc.unsafe `                   | ✖        | -          | -                         | -                                         | 是否启用不安全的RPC方法                                                                                  |
| -        | `--trace-store`                   | ✖        | string     | -                         | -                                         | 是否启用对输出文件的kvstore跟踪                                                                            |
| -        | `--with-tendermint`               | ✖        | -          | `true`                    | -                                         | 使用tendermint运行嵌入进程中的ABCI应用程序                                                                   |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
