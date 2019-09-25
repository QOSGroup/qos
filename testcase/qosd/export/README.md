# Test Cases

- [不指定参数](./TestCase01.md)
- [指定参数height](./TestCase02.md)
- [指定参数height和for-zero-height](./TestCase03.md)
- [指定参数o](./TestCase04.md)

# Description
>     Export state to JSON

>     将状态(state)导出到JSON 

导出区块高度为4的状态数据：
```
qosd export --height 4
```
导出完成后, 默认会在`$HOME/.qosd`下生成以`genesis-<height>-<timestamp>.json`命名的json文件。

# Usage
```
  qosd export [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                | Required | Input Type | Default Input | Input Range | Description                         |
|:---------|:--------------------|:---------|:-----------|:--------------|:------------|:------------------------------------|
| `-h`     | `--help`            | ✖        | -          | -             | -           | 帮助文档                                |
| -        | `--for-zero-height` | ✖        | -          | -             | -           | (主要参数)是否导出状态从0高度重新启动网络（执行预处理）       |
| -        | `--height`          | ✖        | int        | `-1`          | -           | (主要参数)从特定高度导出状态, 指定导出区块高度（-1表示最新高度） |
| -        | `--o`               | ✖        | string     | `"/.qosd"`    | -           | (主要参数)导出JSON文件的目录                   |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
