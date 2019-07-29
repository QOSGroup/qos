# Test Cases

- 暂无

# Description
>     Add guardian.

>     添加特权用户(guardian)。

# Example

> 下面实例中假设:
> - `Arya` 地址为: `address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`
> - `Sansa` 地址为: `address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

在`genesis.json`中配置的特权账户可通过下面的添加指令添加新的特权账户：

`qoscli tx add-guardian --address <new_guardian_key_name_or_account_address> --creator <creator_key_name_or_account_address> --description <description>`

主要参数：

- `--address`         新特权账户，账户地址或密钥库中密钥名字
- `--creator`         创建账户，账户地址或密钥库中密钥名字
- `--description`     描述

`Arya`添加`Sansa`为特权账户：
```bash
$ qoscli tx add-guardian --address Sansa --creator Arya --description 'set Sansa to be a guardian'
```

# Usage
```
  qoscli tx add-guardian [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                | Required | Input Type | Default Input             | Input Range | Description                                            |
|:---------|:--------------------|:---------|:-----------|:--------------------------|:------------|:-------------------------------------------------------|
| `-h`     | `--help`            | ✖        | -          | -                         | -           | 帮助文档                                                   |
| -        | `--async`           | ✖        | -          | -                         | -           | 是否异步广播交易                                               |
| -        | `--chain-id`        | ✖        | string     | -                         | -           | Tendermint节点的链ID                                       |
| -        | `--indent`          | ✖        | -          | -                         | -           | 向JSON响应添加缩进                                            |
| -        | `--max-gas`         | ✖        | int        | `100000`                  | -           | 每个Tx设置的气体限制值                                           |
| `-n`     | `--node`            | ✖        | string     | `"tcp://localhost:26657"` | -           | 为此链提供的Tendermint RPC接口: `<host>:<port>`                |
| -        | `--nonce`           | ✖        | int        | -                         | -           | 要签署Tx的帐户nonce                                          |
| -        | `--nonce-node`      | ✖        | string     | -                         | -           | 用于其他链查询账户nonce的Tendermint RPC接口: `tcp://<host>:<port>` |
| -        | `--qcp`             | ✖        | -          | -                         | -           | 是否启用QCP模式(qcp mode), 发送QCP Tx                          |
| -        | `--qcp-blockheight` | ✖        | int        | -                         | -           | QCP模式Flag标志: 原始Tx块高度，块高度必须大于0                          |
| -        | `--qcp-extends`     | ✖        | string     | -                         | -           | QCP模式Flag标志: QCP Tx扩展信息                                |
| -        | `--qcp-from`        | ✖        | string     | -                         | -           | QCP模式Flag标志: QCP Tx源链ID                                |
| -        | `--qcp-seq`         | ✖        | int        | -                         | -           | QCP模式Flag标志: QCP顺序                                     |
| -        | `--qcp-signer`      | ✖        | string     | -                         | -           | QCP模式Flag标志: QCP Tx签名者key名称                            |
| -        | `--qcp-txindex`     | ✖        | int        | -                         | -           | QCP模式Flag标志: 原始Tx索引                                    |
| -        | `--trust-node`      | ✖        | -          | -                         | -           | 是否信任连接的完整节点（不验证其响应证据）                                  |
| -        | `--address`         | ✖        | string     | -                         | -           | (主要参数)特权用户(guardian)地址                                       |
| -        | `--creator`         | ✖        | string     | -                         | -           | (主要参数)特权用户创造者地址。                                             |
| -        | `--description`     | ✖        | string     | -                         | -           | (主要参数)提案描述                                                   |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |