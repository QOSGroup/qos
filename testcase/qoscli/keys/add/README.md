# Test Cases

- [缺失参数name](./TestCase01.md)
- [添加新密钥](./TestCase02.md)
- [使用种子短语恢复原有密钥](./TestCase03.md)

# Description

>     Add a public/private key pair to the key store.
>     If you select `--recover`, you can recover a key from the seed phrase, otherwise, a new key will be generated. 

>     将公钥/私钥对添加到密钥存储。
>     如果选择“--recover”，则可以从种子短语中恢复密钥，否则将生成一个新密钥。

# Usage
```
  qoscli keys add <name> [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag        | Required | Input Type | Default Input | Input Range | Description               |
|:---------|:------------|:---------|:-----------|:--------------|:------------|:--------------------------|
| `-h`     | `--help`    | ✖        | -          | -             | -           | 帮助文档                      |
| -        | `--recover` | ✖        | -          | -             | -           | (主要参数)提供种子短语以恢复现有密钥，而不是创建 |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |