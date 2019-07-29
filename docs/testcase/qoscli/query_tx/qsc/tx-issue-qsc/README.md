# Test Cases

- [缺失必须参数`--amount`，`--banker`，`--qsc-name`](./TestCase01.md)
- [参数`--amount`，`--banker`，`--qsc-name`不合法](./TestCase02.md)
- [正常发放QSC联盟币](./TestCase03.md)

# Description
>     issue qsc.

>     发放QSC联盟币。

针对使用包含`Banker`公钥创建的联盟币，可向`Banker`地址发放（增发）对应联盟币。

# Example

向联盟币AOE `Banker`中发放（增发）10000AOE：

```bash
$ qoscli tx issue-qsc --qsc-name AOE --banker ATM --amount 10000
Password to sign with 'ATM':<输入ATM本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"223"}
```

可通过[账户查询](#账户（account）)查看`ATM`账户所持有AOE数量。

# Usage
```
  qoscli tx issue-qsc [flags]
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
| -        | `--amount`          | ✔        | int        | -                         | -           | (主要参数)发送给银行`banker`的Coin数量                             |
| -        | `--banker`          | ✔        | string     | -                         | -           | (主要参数)银行`banker`账户本地密钥库名字或账户地址                         |
| -        | `--qsc-name`        | ✔        | string     | -                         | -           | (主要参数)QSC名称                                            |

# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |