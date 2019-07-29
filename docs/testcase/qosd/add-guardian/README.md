# Test Cases

- [缺失必选参数address](./TestCase01.md)
- [填写必选参数address](./TestCase02.md)
- [填写必选参数address和可选参数description](./TestCase03.md)

# Description
>     Add guardian to genesis

>     添加特权帐户(Guardian)到genesis.json

添加特权账户至`genesis.json`文件, 例如：
```bash
$ qosd add-guardian --address address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --description "this is the description"
```
会在`genesis.json`文件`app-state`中`guardian`部分添加地址为`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`的特权账户。

# Usage
```
  qosd add-guardian [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag            | Required | Input Type | Default Input | Input Range | Description                                        |
|:---------|:----------------|:---------|:-----------|:--------------|:------------|:---------------------------------------------------|
| `-h`     | `--help`        | ✖        | -          | -             | -           | 帮助文档                                               |
| -        | `--home-client` | ✖        | string     | `/.qoscli`    | -           | (主要参数)keybase所在目录                                  |
| -        | `--address`     | ✔        | string     | -             | -           | (主要参数)特权帐户地址, 可接收`TaxUsageProposal`提议从社区费池提取的QOS代币 |
| -        | `--description` | ✖        | string     | -             | -           | (主要参数)描述                                           |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
