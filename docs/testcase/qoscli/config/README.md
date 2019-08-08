# Description
>     Create or query a CLI configuration file

>     创建或查询CLI配置文件

# Usage
```
  qoscli config <key> [value] [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag      | Required | Input Type | Default Input | Input Range | Description                  |
|:---------|:----------|:---------|:-----------|:--------------|:------------|:-----------------------------|
| `-h`     | `--help`  | ✖        | -          | -             | -           | 帮助文档                         |
| `-p`     | `--print` | ✖        | -          | -             | -           | (主要参数)打印配置项的值, 如果未设置，则打印其默认值 |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |
