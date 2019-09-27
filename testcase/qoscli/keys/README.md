# Description

>     Keys allows you to manage your local keystore for tendermint. 
>     These keys may be in any format supported by go-crypto and can be used by light-clients, full nodes, or any other application that needs to sign with a private key.

>     密钥允许您为TenderMint管理本地密钥库。
>     这些密钥可以是Go Crypto支持的任何格式，并且可以由轻型客户端、完整节点或任何其他需要用私钥签名的应用程序使用。

# Usage
```
  qoscli keys [command]
```

# Available Commands

| Command              | Alias | Has-Subcommand | Description               |
|:---------------------|:------|:---------------|:--------------------------|
| `qoscli keys add`    | -     | ✖              | 创建新密钥，或从种子导入              |
| `qoscli keys list`   | -     | ✖              | 列出所有密钥                    |
| `qoscli keys delete` | -     | ✖              | 删除给定的密钥                   |
| `qoscli keys update` | -     | ✖              | 更改用于保护私钥的密码               |
| `qoscli keys export` | -     | ✖              | 导出给定名称的密钥                 |
| `qoscli keys import` | -     | ✖              | 交互式命令，用于导入新的私钥、对其加密并保存到磁盘 |

# Flags

| ShortCut | Flag      | Input Type | Default Input | Input Range | Description            |
|:---------|:----------|:-----------|:--------------|:------------|:-----------------------|
| `-h`     | `--help`  | -          | -             | -           | (可选)帮助文档                   |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |