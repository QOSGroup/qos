# Description
>     QOS Daemon (server)

>     QOS守护进程（服务端）

# Usage
```
  qosd [command]
```

# Available Commands

| Command                     | Alias | Has-Subcommand | Description                                        |
|:----------------------------|:------|:---------------|:---------------------------------------------------|
| `qosd add-genesis-accounts` | -     | ✖              | 添加创世帐户(Genesis Account)到genesis.json               |
| `qosd add-guardian`         | -     | ✖              | 添加特权帐户(Guardian)到genesis.json                      |
| `qosd collect-gentxs`       | -     | ✖              | 收集创世交易(Genesis Txs)到genesis.json                   |
| `qosd config-root-ca`       | -     | ✖              | 为QCP和QSC配置根CA(root CA)的公钥(pubKey)                  |
| `qosd export`               | -     | ✖              | 将状态(state)导出到JSON                                  |
| `qosd gentx`                | -     | ✖              | 生成一个带有自委托的创世交易(Genesis Tx)                         |
| `qosd help`                 | -     | ✖              | 关于任何命令的帮助                                          |
| `qosd init`                 | -     | ✖              | 初始化私有验证人、P2P、Genesis和应用程序配置文件                      |
| `qosd start`                | -     | ✖              | 运行完整节点                                             |
| `qosd tendermint`           | -     | ✔              | Tendermint子命令                                      |
| `qosd testnet`              | -     | ✖              | 为QOS测试网初始化文件                                       |
| `qosd unsafe-reset-all`     | -     | ✖              | 重置区块链数据库，删除通讯簿文件，并将priv_validator.json重置为genesis状态 |
| `qosd version`              | -     | ✖              | 打印应用程序版本                                           |

# Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| `-h`     | `--help`      | ✖        | -          | -                                | -           | (可选)帮助文档     |
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
