# Test Cases

- [获取列表](./TestCase01.md)
- [获取列表(JSON格式)](./TestCase02.md)

# Description

>     Return a list of all public keys stored by this key manager along with their associated name and address.

>     返回此密钥管理器存储的所有公钥的列表, 及其关联的名称和地址。

# Usage
```
  qoscli keys list [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag     | Required | Input Type | Default Input | Input Range | Description |
|:---------|:---------|:---------|:-----------|:--------------|:------------|:------------|
| `-h`     | `--help` | ✖        | -          | -             | -           | 帮助文档        |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |