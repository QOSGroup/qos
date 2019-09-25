# Test Cases

- 暂无

# Description
>     Submit proposal.

>     发起提案。

# Example

> 下面实例中假设:
> - `Arya` 地址为: `address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`
> - `Sansa` 地址为: `address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

`qoscli tx submit-proposal 
    --title <proposal_title> 
    --proposal-type <proposal_type> 
    --proposer <proposer_key_name_or_account_address> 
    --deposit <deposit_amount_of_qos> 
    --description <description>`

- `--title`             标题
- `--proposal-type`     提议类型：`Text`、`ParameterChange`、`TaxUsage`
- `--proposer`          提议账户，账户地址或密钥库中密钥名字
- `--deposit`           提议押金，不能小于`MinDeposit`的三分之一
- `--description`       描述信息

`TaxUsage`类型提议特有参数：

- `--dest-address`      目标地址，用于接收QOS
- `--percent`           社区费池提取比例，小数0~1

`ParameterChange`类型提议特有参数：

- `--params`            参数列表，格式：'module:key_name:value,module:key_name:value，如：gov:min_deposit:10000


`Arya`提交一个文本提议：
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type Text --proposer Arya --deposit 10000000 --description 'this is the description'
```

`Arya`提交一个参数修改提议：
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type ParameterChange --proposer Arya --deposit 10000000 --description 'this is the description' --params gov:min_deposit:1000
```

假设`Arya`在QOS初始化时已经通过[添加特权账户](qosd.md#添加特权账户) 添加到了`genesis.json`，`Arya`提交一个提取费池提议：
```bash
$ qoscli tx submit-proposal --title 'update qos' --proposal-type TaxUsage --proposer Arya --deposit 10000000 --description 'this is the description' --dest-address Sansa --percent 0.5
```

# Usage
```
  qoscli tx submit-proposal [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                | Required | Input Type | Default Input             | Input Range                                 | Description                                                  |
|:---------|:--------------------|:---------|:-----------|:--------------------------|:--------------------------------------------|:-------------------------------------------------------------|
| `-h`     | `--help`            | ✖        | -          | -                         | -                                           | 帮助文档                                                         |
| -        | `--async`           | ✖        | -          | -                         | -                                           | 是否异步广播交易                                                     |
| -        | `--chain-id`        | ✖        | string     | -                         | -                                           | Tendermint节点的链ID                                             |
| -        | `--indent`          | ✖        | -          | -                         | -                                           | 向JSON响应添加缩进                                                  |
| -        | `--max-gas`         | ✖        | int        | `100000`                  | -                                           | 每个Tx设置的气体限制值                                                 |
| `-n`     | `--node`            | ✖        | string     | `"tcp://localhost:26657"` | -                                           | 为此链提供的Tendermint RPC接口: `<host>:<port>`                      |
| -        | `--nonce`           | ✖        | int        | -                         | -                                           | 要签署Tx的帐户nonce                                                |
| -        | `--nonce-node`      | ✖        | string     | -                         | -                                           | 用于其他链查询账户nonce的Tendermint RPC接口: `tcp://<host>:<port>`       |
| -        | `--qcp`             | ✖        | -          | -                         | -                                           | 是否启用QCP模式(qcp mode), 发送QCP Tx                                |
| -        | `--qcp-blockheight` | ✖        | int        | -                         | -                                           | QCP模式Flag标志: 原始Tx块高度，块高度必须大于0                                |
| -        | `--qcp-extends`     | ✖        | string     | -                         | -                                           | QCP模式Flag标志: QCP Tx扩展信息                                      |
| -        | `--qcp-from`        | ✖        | string     | -                         | -                                           | QCP模式Flag标志: QCP Tx源链ID                                      |
| -        | `--qcp-seq`         | ✖        | int        | -                         | -                                           | QCP模式Flag标志: QCP顺序                                           |
| -        | `--qcp-signer`      | ✖        | string     | -                         | -                                           | QCP模式Flag标志: QCP Tx签名者key名称                                  |
| -        | `--qcp-txindex`     | ✖        | int        | -                         | -                                           | QCP模式Flag标志: 原始Tx索引                                          |
| -        | `--trust-node`      | ✖        | -          | -                         | -                                           | 是否信任连接的完整节点（不验证其响应证据）                                        |
| -        | `--deposit`         | ✖        | uint       | -                         | -                                           | (主要参数)提案发起人支付的初始保证金。必须是严格的正数。                                |
| -        | `--description`     | ✖        | string     | -                         | -                                           | (主要参数)提案描述                                                   |
| -        | `--dest-address`    | ✖        | string     | -                         | -                                           | (主要参数)接收QOS的地址                                               |
| -        | `--params`          | ✖        | string     | -                         | -                                           | (主要参数)参数，格式为：`<module>/<key>:<value>,<module>/<key>:<value>` |
| -        | `--percent`         | ✖        | float      | -                         | -                                           | (主要参数)费用池(fee pool)中发送到目标地址(dest-address)的QOS百分比             |
| -        | `--proposal-type`   | ✖        | string     | `"Text"`                  | `"Text"`, `"ParameterChange"`, `"TaxUsage"` | (主要参数)验证人`Owner`账户本地密钥库名字或账户地址                               |
| -        | `--proposer`        | ✖        | string     | -                         | -                                           | (主要参数)提交提案的提案人                                               |
| -        | `--title`           | ✖        | string     | -                         | -                                           | (主要参数)提案标题                                                   |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |