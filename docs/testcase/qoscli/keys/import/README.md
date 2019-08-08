# Test Cases

- [缺失参数name](./TestCase01.md)
- [指定的name在密钥库中已存在](./TestCase02.md)
- [导入密钥库中不存在的密钥](./TestCase03.md)
- [导入密钥库中不存在的密钥（从CA PRI文件导入）](./TestCase04.md)

# Description

>     Interactive command to import a new private key, encrypt it, and save to disk.

>     交互式命令，用于导入新的私钥、对其加密并保存到磁盘。

# Usage
```
  qoscli keys import [name] [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag     | Required | Input Type | Default Input | Input Range | Description         |
|:---------|:---------|:---------|:-----------|:--------------|:------------|:--------------------|
| `-h`     | `--help` | ✖        | -          | -             | -           | 帮助文档                |
| -        | `--file` | ✖        | string     | -             | -           | (主要参数)从CA PRI文件导入私钥 |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |