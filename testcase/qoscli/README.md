# Description
>     QOS light-client

>     QOS轻客户端, 是用于与QOS网络交互的命令行工具

# Usage
```
  qoscli [command]
```

# Available Commands

| Command          | Alias      | Has-Subcommand | Description       |
|:-----------------|:-----------|:---------------|:------------------|
| `qoscli config`  | -          | ✖              | 创建或查询CLI配置文件      |
| `qoscli keys`    | -          | ✔              | 密钥管理工具: 添加或查看本地私钥 |
| `qoscli query`   | `qoscli q` | ✔              | 查询子命令             |
| `qoscli tx`      | -          | ✔              | 交易子命令             |
| `qoscli version` | -          | ✖              | 打印应用程序版本          |
| `qoscli help`    | -          | ✖              | 关于任何命令的帮助         |


# Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| `-h`     | `--help`     | ✖        | -          | -             | -                 | 帮助文档     |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |
