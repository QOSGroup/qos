# Test Cases

- [缺失必选参数](./TestCase01.md)
- [仅指定必选参数](./TestCase02.md)
- [指定可选参数](./TestCase03.md)
- [必选参数owner指定的账户不存在](./TestCase04.md)
- [可选参数ip格式不合法](./TestCase05.md)

# Description
>     This command is an alias of the 'gaiad tx create-validator' command'.

>     生成一个带有自委托的创世交易(Genesis Tx) 

生成创建验证节点交易[TxCreateValidator](../../../spec/staking.md#TxCreateValidator)

默认会在`$HOME/.qosd/config/gentx`目录下生成以`nodeID@IP`为文件名的已签名的交易数据文件。

# Example
```
  qosd gentx --moniker validatorName --owner ownerName --tokens 100
```

# Usage
```
  qosd gentx [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag            | Required | Input Type | Default Input | Input Range | Description                             |
|:---------|:----------------|:---------|:-----------|:--------------|:------------|:----------------------------------------|
| `-h`     | `--help`        | ✖        | -          | -             | -           | 帮助文档                                    |
| -        | `--compound`    | ✖        | -          | -             | -           | (主要参数)作为一个自委托者，收入是否计算为复利                |
| -        | `--details`     | ✖        | string     | -             | -           | (主要参数)验证人详细描述信息, `len(details) <= 1000` |
| -        | `--home-client` | ✖        | string     | `"/.qoscli"`  | -           | (主要参数)节点keybase所在目录                     |
| -        | `--home-node`   | ✖        | string     | `"/.qosd"`    | -           | (主要参数)节点配置文件和数据所在目录                     |
| -        | `--ip`          | ✖        | string     | `"127.0.0.1"` | -           | (主要参数)验证节点IP                            |
| -        | `--logo`        | ✖        | string     | -             | -           | (主要参数)logo链接， `len(logo) <= 255`        |
| -        | `--moniker`     | ✔        | string     | -             | -           | (主要参数)验证节点名字，`len(moniker) <= 300`      |
| -        | `--owner`       | ✔        | string     | -             | -           | (主要参数)操作者账户地址或密钥库中密钥名字                  |
| -        | `--tokens`      | ✔        | int        | -             | -           | (主要参数)绑定的代币(tokens)数量，不能大于操作者持有QOS数量    |
| -        | `--website`     | ✖        | string     | -             | -           | (主要参数)验证人网址， `len(website) <= 255`      |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
