 # Test Cases

- [缺失参数name](./TestCase01.md)
- [指定的name在密钥库中不存在](./TestCase02.md)
- [导出已存在的密钥](./TestCase03.md)
- [导出已存在的密钥（只导出公钥）](./TestCase04.md)

# Description

>     Export key for the given name.

>     导出给定名称的密钥。

# Usage
```
  qoscli keys export [name] [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag       | Required | Input Type | Default Input | Input Range | Description |
|:---------|:-----------|:---------|:-----------|:--------------|:------------|:------------|
| `-h`     | `--help`   | ✖        | -          | -             | -           | 帮助文档        |
| -        | `--pubkey` | ✖        | -          | -             | -           | (主要参数)只导出公钥 |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |