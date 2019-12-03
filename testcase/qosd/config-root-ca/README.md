# Test Cases

- [不指定任何参数](./TestCase01.md)
- [指定QCP参数](./TestCase02.md)
- [指定QSC参数](./TestCase03.md)

# Description
>     Config pubKey of root CA for QCP and QSC

>     为QCP和QSC配置根CA(root CA)的公钥(pubKey)

设置Root CA公钥信息，用于[联盟币](qoscli.md#联盟币（qsc）)和[联盟链](qoscli.md#联盟链（qcp）)涉及到证书操作的校验。

# Usage
```
  qosd config-root-ca [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag     | Required | Input Type | Default Input | Input Range | Description                  |
|:---------|:---------|:---------|:-----------|:--------------|:------------|:-----------------------------|
| `-h`     | `--help` | ✖        | -          | -             | -           | 帮助文档                         |
| -        | `--qcp`  | ✔        | string     | -             | -           | (主要参数)QCP根证书公钥文件root.pub文件路径 |
| -        | `--qsc`  | ✔        | string     | -             | -           | (主要参数)QCP根证书公钥文件root.pub文件路径 |

# Global Flags

| ShortCut | Flag          | Required | Input Type | Default Input                    | Input Range | Description  |
|:---------|:--------------|:---------|:-----------|:---------------------------------|:------------|:-------------|
| -        | `--home`      | ✖        | string     | `/.qosd`                         | -           | 配置和数据的目录     |
| -        | `--log_level` | ✖        | string     | `"main:info,state:info,*:error"` | -           | 日志级别         |
| -        | `--trace`     | ✖        | -          | -                                | -           | 打印出错时的完整堆栈跟踪 |
