# Test Cases

- [缺失必须参数`--creator`，`--qsc.crt`](./TestCase01.md)
- [参数`--creator`，`--qsc.crt`不合法](./TestCase02.md)
- [正常创建QSC联盟币](./TestCase03.md)

# Description
>     create qsc.

>     创建QSC联盟币。

# Example

`Arya`在QOS网络中创建`AOE`，不含初始发放地址币值信息：
```bash
$ qoscli tx create-qsc --creator Arya --qsc.crt aoe.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
> 假设`Arya`已在CA中心申请`aoe.crt`证书，`aoe.crt`中包含`banker`公钥，对应地址`address1rpmtqcexr8m20zpl92llnquhpzdua9stszmhyq`，已经导入到本地私钥库中，名字为`ATM`，。

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"200"}
```

# Usage
```
  qoscli tx create-qsc [flags]
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
| -        | `--accounts`        | ✖        | string     | -                         | -           | (主要参数)初始帐户，例如: address1,100;address2,100               |
| -        | `--creator`         | ✔        | string     | -                         | -           | (主要参数)创建者账户本地密钥库名字或账户地址                                |
| -        | `--desc`            | ✖        | string     | -                         | -           | (主要参数)描述                                               |
| -        | `--extrate`         | ✖        | string     | `"1"`                     | -           | (主要参数)extrate: qos:qscxxx                              |
| -        | `--qsc.crt`         | ✔        | string     | -                         | -           | (主要参数)CA路径（qsc）                                        |


# Global Flags

| ShortCut | Flag         | Required | Input Type | Default Input | Input Range       | Description  |
|:---------|:-------------|:---------|:-----------|:--------------|:------------------|:-------------|
| `-e`     | `--encoding` | ✖        | string     | `hex`         | `hex`/`b64`/`btc` | 二进制编码        |
| -        | `--home`     | ✖        | string     | `/.qoscli`    | -                 | 配置和数据的目录     |
| `-o`     | `--output`   | ✖        | string     | `text`        | `text`/`json`     | 输出格式         |
| -        | `--trace`    | ✖        | -          | -             | -                 | 打印出错时的完整堆栈跟踪 |