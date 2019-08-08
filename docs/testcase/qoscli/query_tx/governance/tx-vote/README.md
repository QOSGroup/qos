# Test Cases

- 暂无

# Description
>     vote.

>     向提案投票。

# Example

> 下面实例中假设:
> - `Arya` 地址为: `address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`
> - `Sansa` 地址为: `address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

进入投票阶段的提议可通过下面指令进行投票操作：

`qoscli tx vote --proposal-id <proposal_id> --voter <voter_key_name_or_account_address> --option <vote_option>`

主要参数：

- `--proposal-id`       提议ID
- `--voter`             投票账户，地址或密钥库名字
- `--option`            投票选项，可选值：`Yes`,`Abstain`,`No`,`NoWithVeto`

`Arya`给1号提议投票`Yes`：
```bash
$ qoscli tx vote --proposal-id 1 --voter Arya --option Yes
```

# Usage
```
  qoscli tx vote [flags]
```

# Available Commands

>无可用命令

# Flags

| ShortCut | Flag                | Required | Input Type | Default Input             | Input Range                                  | Description                                            |
|:---------|:--------------------|:---------|:-----------|:--------------------------|:---------------------------------------------|:-------------------------------------------------------|
| `-h`     | `--help`            | ✖        | -          | -                         | -                                            | 帮助文档                                                   |
| -        | `--async`           | ✖        | -          | -                         | -                                            | 是否异步广播交易                                               |
| -        | `--chain-id`        | ✖        | string     | -                         | -                                            | Tendermint节点的链ID                                       |
| -        | `--indent`          | ✖        | -          | -                         | -                                            | 向JSON响应添加缩进                                            |
| -        | `--max-gas`         | ✖        | int        | `100000`                  | -                                            | 每个Tx设置的气体限制值                                           |
| `-n`     | `--node`            | ✖        | string     | `"tcp://localhost:26657"` | -                                            | 为此链提供的Tendermint RPC接口: `<host>:<port>`                |
| -        | `--nonce`           | ✖        | int        | -                         | -                                            | 要签署Tx的帐户nonce                                          |
| -        | `--nonce-node`      | ✖        | string     | -                         | -                                            | 用于其他链查询账户nonce的Tendermint RPC接口: `tcp://<host>:<port>` |
| -        | `--qcp`             | ✖        | -          | -                         | -                                            | 是否启用QCP模式(qcp mode), 发送QCP Tx                          |
| -        | `--qcp-blockheight` | ✖        | int        | -                         | -                                            | QCP模式Flag标志: 原始Tx块高度，块高度必须大于0                          |
| -        | `--qcp-extends`     | ✖        | string     | -                         | -                                            | QCP模式Flag标志: QCP Tx扩展信息                                |
| -        | `--qcp-from`        | ✖        | string     | -                         | -                                            | QCP模式Flag标志: QCP Tx源链ID                                |
| -        | `--qcp-seq`         | ✖        | int        | -                         | -                                            | QCP模式Flag标志: QCP顺序                                     |
| -        | `--qcp-signer`      | ✖        | string     | -                         | -                                            | QCP模式Flag标志: QCP Tx签名者key名称                            |
| -        | `--qcp-txindex`     | ✖        | int        | -                         | -                                            | QCP模式Flag标志: 原始Tx索引                                    |
| -        | `--trust-node`      | ✖        | -          | -                         | -                                            | 是否信任连接的完整节点（不验证其响应证据）                                  |
| -        | `--option`          | ✖        | string     | -                         | `"Yes"`, `"Abstain"`, `"No"`, `"NoWithVeto"` | (主要参数)提案投票选项                                             |
| -        | `--proposal-id`     | ✖        | uint       | -                         | -                                            | (主要参数)提案ID                                               |
| -        | `--voter`           | ✖        | string     | -                         | -                                            | (主要参数)投票人                                                |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |