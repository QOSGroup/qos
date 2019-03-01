# QOS Client

提供与QOS网络交互的命令行工具`qoscli`，主要提供以下命令行功能：
* `keys`        [本地密钥库](#密钥（keys）)
* `query`       [信息查询](#查询（query）)
* `tx`          [执行交易](#交易（tx）)
* `tendermint`  [tendermint自带指令](#tendermint)
* `version`     版本信息

所有命令均可通过添加`--help`获取命令说明

命令中涉及的通用参数：

| 参数 | 默认值 | 说明 |
| :--- | :---: | :--- |
|--nonce | 0 | account nonce to sign the tx |
|--max-gas| 0 | gas limit to set per tx |
|--chain-id| "" | Chain ID of tendermint node |
|--node| tcp://localhost:26657 | tcp://\<host\>:\<port\> to tendermint rpc interface for this chain |
|--height| 0 | block height to query, omit to get most recent provable block |
|--async| false | broadcast transactions asynchronously |
|--trust-node| false | Trust connected full node |
|--qcp| false | enable qcp mode. send qcp tx |
|--qcp-signer| "" | qcp mode flag. qcp tx signer key name |
|--qcp-seq| 0 | qcp mode flag.  qcp in sequence |
|--qcp-from| "" | qcp mode flag. qcp tx source |
|--qcp-blockheight| 0 | qcp mode flag. original tx blockheight |
|--qcp-txindex| 0 | qcp mode flag. original tx index |
|--qcp-extends| "" | qcp mode flag. qcp tx extends info |
|--indent| false | add indent to json response |
|--nonce-node| "" | tcp://\<host\>:\<port\> to tendermint rpc interface for some chain to query account nonce |

更多说明参照[qbase-通用参数](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#客户端命令)

## 密钥（keys）

本地密钥库主要包含以下指令：
* `qoscli keys add`     [新增密钥](#新增（add）)
* `qoscli keys list`    [显示密钥列表](#列表（list）)
* `qoscli keys update`  [更新密钥保存密码](#更新（update）)
* `qoscli keys delete`  [从密钥库删除密钥](#删除（delete）)
* `qoscli keys import`  [导入密钥](#导入（import）)
* `qoscli keys export`  [导出密钥](#导出（export）)

> 密钥库为本地存储，默认存储位置为：$HOME/.qoscli/keys/，删除存储文件会清空本地存储所有私钥。通过`keys`相关指令操作密钥不影响QOS网络中账户状态，请妥善保管账户私钥信息。

### 新增（add）

`qoscli keys add <key_name>`

<key_name>可随意填写，仅作为本地密钥库密钥区分。

如下指令将生成一个名字为`Arya`的密钥到本地密钥库：
```bash
$ qoscli keys add Arya
Enter a passphrase for your key:<输入不少于8位的密码>
Repeat the passphrase:<重复上面输入的密码>
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy	dfYz3Zg+g1VFU52frAiKyXRU4wVulJMYgIuboPuBtZ4=
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

thought frame must space few omit muffin fix merge mail ivory clump unveil dirt gadget load glove hub inner final crime churn crop stone
```
其中`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`为适用于QOS网络的账户地址，`dfYz3Zg+g1VFU52frAiKyXRU4wVulJMYgIuboPuBtZ4=`为账户公钥信息，`thought frame must space few omit muffin fix merge mail ivory clump unveil dirt gadget load glove hub inner final crime churn crop stone`为助记词，可用于账户私钥找回，请妥善保管助记词。

### 列表（list）

`qoscli keys list`
```bash
$ qoscli keys list
NAME:	TYPE:	ADDRESS:						PUBKEY:
Arya	local	address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy	dfYz3Zg+g1VFU52frAiKyXRU4wVulJMYgIuboPuBtZ4=
```

### 更新（update）

`qoscli keys update <key_name>`

更新`Arya`存储密码：
```bash
$ qoscli keys update Arya
Enter the current passphrase:<输入当前密码>
Enter the new passphrase:<输入新密码>
Repeat the new passphrase:<重复新密码>
Password successfully updated!
```

### 导出（export）

`qoscli keys export <key_name>`

导出`Arya`密钥信息：
```bash
qoscli keys export Arya
Password to sign with 'Arya':<输入>
**Important** Don't leak your private key information to others.
Please keep your private key safely, otherwise your account will be attacked.

{"Name":"Arya","address":"address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy","pubkey":{"type":"tendermint/PubKeyEd25519","value":"dfYz3Zg+g1VFU52frAiKyXRU4wVulJMYgIuboPuBtZ4="},"privkey":{"type":"tendermint/PrivKeyEd25519","value":"bXeccNwvLk8w2cloSXtO6FKcHXTQ7sfEecpyPzcUzg119jPdmD6DVUVTnZ+sCIrJdFTjBW6UkxiAi5ug+4G1ng=="}}
```
导出的密钥是通过JSON序列化后的密钥信息，可以将JSON字符串中的`privkey`部分内容保存为文件并妥善保存，可用于密钥导入。

### 删除（delete）

`qoscli keys delete <key_name>`

删除`Arya`密钥信息：
```bash
$ qoscli keys delete Arya
DANGER - enter password to permanently delete key:<输入密码>
Password deleted forever (uh oh!)
```

### 导入（import）

`qoscli keys import Arya --file <私钥文件路径>`

导入上面通过`export`导出的私钥文件：
```bash
qoscli keys import Arya --file Arya.pri 
> Enter a passphrase for your key:<输入不少于8位的密码>
> Repeat the passphrase:<重复上面输入的密码>
```

## 查询（query）

* `qoscli query account`                [账户查询](#账户（account）)
* `qoscli query store`                  [存储查询](#存储（store）)
* `qoscli query consensus`              共识参数查询
* `qoscli query approve`                [预授权](#查询预授权)
* `qoscli query qcp`                    [跨链相关信息查询](#查询联盟链)
* `qoscli query qsc`                    [联盟币信息查询](#查询联盟币)
* `qoscli query validators`             [验证节点列表](#验证节点列表)
* `qoscli query validator`              [验证节点查询](#查询验证节点)
* `qoscli query validator-miss-vote`    [验证节点漏块信息](#查询验证节点漏块信息)
* `qoscli query validator-period`       [验证节点窗口信息](#验证节点窗口信息)
* `qoscli query community-fee-pool`     [社区收益池](#社区收益池)
* `qoscli query delegation`             [委托查询](#委托查询)
* `qoscli query delegations-to`         [验证节点委托列表](#验证节点委托列表)
* `qoscli query delegations`            [代理用户委托列表](#代理用户委托列表)
* `qoscli query delegator-income`       [委托收益查询](#委托收益查询)

查询的具体指令将在各自模块进行介绍。

### 账户（account）

查询账户
`qoscli query account <key_name_or_account_address>`

<key_name_or_account_address>为本地密钥库存储的密钥名字或对应账户的地址。

假设本地密钥库中`Arya`地址为`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`，且QOS网络中已经创建了`address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy`对应账号，可执行：
```bash
qoscli query account Arya --indent
```
或
```bash
qoscli query account address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --indent
```
输出类似如下信息：
```bash
{
  "type": "qbase/account/QOSAccount",
  "value": {
    "base_account": {
      "account_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
      "public_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "dfYz3Zg+g1VFU52frAiKyXRU4wVulJMYgIuboPuBtZ4="
      },
      "nonce": "0"
    },
    "qos": "10000",
    "qscs": [
        {
            "coin_name": "AOE",
            "amount": "10000"
        }
    ]
  }
}
```
可以看到`Arya`持有10000个QOS、10000个AOE，更多账户说明请阅读[QOS账户设计](../spec/account.md)文档。

### 存储（store）

QOS网络的存储内容均可通过下面指令查找：

`qoscli query store --path /store/<store_key>/subspace --data <query_data>`

主要参数：

- `--path`  存储位置
- `--data`  查询内容，以<query_data>开头的数据会被查出来

查询QOS网络中存储的ROOT CA 信息：

```bash
$ qoscli query store --path /store/base/subspace --data rootca --indent
```

执行结果：

```bash
[
  {
    "key": "rootca",
    "value": {
      "type": "tendermint/PubKeyEd25519",
      "value": "L+P3Vm8NQRDwVt4rHzlqtBJLGSsLZGLmmd4wLYrUe6U="
    }
  }
]
```

## 交易（tx）

QOS支持以下几种交易类型：

* `qoscli tx transfer`         [转账](#转账（transfer）)
* `qoscli tx create-approve`   [创建预授权](#创建预授权)
* `qoscli tx increase-approve` [增加预授权](#增加预授权)
* `qoscli tx decrease-approve` [减少预授权](#减少预授权)
* `qoscli tx use-approve`      [使用预授权](#使用预授权)
* `qoscli tx cancel-approve`   [取消预授权](#取消预授权)
* `qoscli tx create-qsc`       [创建联盟币](#创建联盟币)
* `qoscli tx issue-qsc`        [发放联盟币](#发放联盟币)
* `qoscli tx init-qcp`         [初始化联盟链](#初始化联盟链)
* `qoscli tx create-validator` [成为验证节点](#成为验证节点)
* `qoscli tx revoke-validator` [撤销验证节点](#撤销验证节点)
* `qoscli tx active-validator` [激活验证节点](#激活验证节点)
* `qoscli tx delegate`         [委托](#委托)
* `qoscli tx modify-compound`  [修改收益复投方式](#修改收益复投方式)
* `qoscli tx unbond`           [解除委托](#解除委托)
* `qoscli tx redelegate`       [变更委托验证节点](#变更委托验证节点)

分为**转账**、**预授权**、**联盟币**、**联盟链**、**验证节点**五大类。

### 转账（transfer）

查阅[转账设计](../spec/txs/transfer.md)了解QOS转账交易设计。

`qoscli tx transfer --senders <senders_and_coins> --receivers <receivers_and_coins>`

支持一次转账中包含多币种，多账户

主要参数：
- `--senders`   发送集合，账户传keystore name 或 address，多个账户半角分号分隔
- `--receivers` 接收集合，账户传keystore name 或 address，多个账户半角分号分隔

`Arya`向地址`address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`转账1个QOS，1个AOE
```bash
$ qoscli tx transfer --senders Arya,1QOS,1AOE --receivers address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh,1QOS,1AOE
Password to sign with 'Arya':<输入密码>
{"check_tx":{},"deliver_tx":{},"hash":"21ECB72C8F51B3BD8E3CB9D59765003B9D78BE75","height":"300"}
```

转账成功可通过[账户查询](#账户（account）)查看最新账户状态，交易执行可能会有一定时间的延迟。

### 预授权（approve）

[QOS预授权设计](../spec/txs/approve.md)包含以下操作指令：

* `qoscli tx create-approve`    [创建预授权](#创建预授权)
* `qoscli query approve`        [查询预授权](#查询预授权)
* `qoscli tx increase-approve`  [增加预授权](#增加预授权)
* `qoscli tx decrease-approve`  [减少预授权](#减少预授权)
* `qoscli tx use-approve`       [使用预授权](#使用预授权)
* `qoscli tx cancel-approve`    [取消预授权](#取消预授权)

> 下面实例中假设`Sansa`地址为`address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh`

#### 创建预授权

`qoscli tx create-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

主要参数：

- `--from`  授权账户本地密钥库名字或账户地址
- `--to`    被授权账户地址
- `--coins` 授权币种、币值列表，[amount1][coin1],[amount2][coin2],...，以半角逗号相隔

`Arya`向`Sansa`授权100个QOS，100个AOE：
```
$ qoscli tx create-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --coins 100QOS,100AOE
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"9917953D8CDE80F457CD072DBCE73A36449B7A7C","height":"333"}
```

#### 查询预授权

`qoscli query approve --from <key_name_or_account_address> --to <account_address>`

查询`Arya`对`Sansa`的预授权：
```bash
qoscli query approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh
```
执行结果：
```bash
{
  "from": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "to": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
  "qos": "100",
  "qscs": [
    {
      "coin_name": "AOE",
      "amount": "100"
    }
  ]
}
```

#### 增加预授权

`qoscli tx increase-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

`Arya`向`Sansa`增加授权100个QOS，100个AOE：
```bash
$ qoscli tx increase-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --coins 100QOS,100AOE
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"406"}
```

#### 减少预授权

`qoscli tx decrease-approve --from <key_name_or_account_address> --to <account_address> --coins <qos_and_qscs>`

`Arya`向`Sansa`减少授权10个QOS，10个AOE：
```bash
$ qoscli tx decrease-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh --coins 10QOS,10AOE
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"410"}
```

#### 使用预授权

`qoscli tx use-approve --from <account_address> --to <key_name_or_account_address> --coins <qos_and_qscs>`

`Sansa`使用`Arya`向自己预授权中的10个QOS，10个AOE：
```bash
$ qoscli tx use-approve --from address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --to Sansa --coins 10QOS,10AOE
Password to sign with 'Sansa':<输入Sansa本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"3C06676C53A5439D39CB4D0FBA3213C44DC1BA8E","height":"430"}
```

可通过[账户查询](#账户（account）)查看`Arya`和`Sansa`最新账户状态

#### 取消预授权

`qoscli tx cancel-approve --from <account_address> --to <key_name_or_account_address>'

`Arya`取消对`Sansa`的授权：
```bash
$ qoscli tx cancel-approve --from Arya --to address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"484"}
```

### 联盟币（qsc）

> 创建联盟币前需要申请[CA](../spec/ca.md)，点击[联盟币设计文档](../spec/txs/qsc.md)了解更多。

联盟币相关指令：
* `qoscli tx create-qsc`    [创建联盟币](#创建联盟币)
* `qoscli query qsc`        [查询联盟币](#查询联盟币)
* `qoscli tx issue-qsc`     [发放联盟币](#发放联盟币)

#### 创建联盟币

`qoscli tx create-qsc --creator <key_name_or_account_address> --qsc.crt <qsc.crt_file_path> --accounts <account_qsc_s>`

主要参数：

- `--creator`       创建账号
- `--qsc.crt`       证书位置
- `--accounts`      初始发放地址币值集合，[addr1],[amount];[addr2],[amount2],...，该参数可为空，即只创建联盟币

`Arya`在QOS网络中创建`QOE`，不含初始发放地址币值信息：
```bash
$ qoscli tx create-qsc --creator Arya --qsc.crt aoe.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```
> 假设`Arya`已在CA中心申请`aoe.crt`证书，`aoe.crt`中包含`banker`公钥，对应地址`address1rpmtqcexr8m20zpl92llnquhpzdua9stszmhyq`，已经导入到本地私钥库中，名字为`ATM`，。

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"200"}
```

#### 查询联盟币

`qoscli query qsc <qsc_name>`

`qsc_name`为联盟币名称

查询`AOE`信息：
```bash
$ qoscli query qsc QOE --indent
```
执行结果：
```bash
{
  "name": "AOE",
  "chain_id": "capricorn-1000",
  "extrate": "1:280.0000",
  "description": "",
  "banker": "address1rpmtqcexr8m20zpl92llnquhpzdua9stszmhyq"
}
```

#### 发放联盟币

针对使用包含`Banker`公钥创建的联盟币，可向`Banker`地址发放（增发）对应联盟币：

`qoscli tx issue-qsc --qsc-name <qsc_name> --banker <key_name_or_account_address> --amount <qsc_amount>`

主要参数：
- `--qsc-name`  联盟币名字
- `--banker`    Banker地址或私钥库中私钥名
- `--amount`    联盟币发放（增发）量

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

### 联盟链（qcp）

QOS跨链协议QCP，支持跨链交易

> 创建联盟链前需要申请[CA](../spec/ca.md)，点击[联盟链设计文档](../spec/txs/qcp.md)了解更多。

联盟链相关指令：
* `qoscli tx init-qcp`: [初始化联盟链](#初始化联盟链)
* `qoscli query qcp`:   [查询qcp信息](#查询联盟链)

#### 初始化联盟链

`qoscli tx init-qcp --creator <key_name_or_account_address> --qcp.crt <qcp.crt_file_path>`

主要参数：

- `--creator`       创建账号
- `--qcp.crt`       证书位置

> 假设`Arya`已在CA中心申请`qcp.crt`证书，`qcp.crt`中联盟链ID为`aoe-1000`

`Arya`在QOS网络中初始化联盟链信息：
```bash
$ qoscli tx init-qcp --creator Arya --qcp.crt qcp.crt
Password to sign with 'Arya':<输入Arya本地密钥库密码>
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"243"}
```

#### 查询联盟链

跨链协议是[qbase](https://www.github.com/QOSGroup/qbase)提供支持，主要有以下四个查询指令：
- `qoscli query qcp list`
- `qoscli query qcp out` 
- `qoscli query qcp in`
- `qoscli query qcp tx`

指令说明请参照[qbase-Qcp](https://github.com/QOSGroup/qbase/blob/master/docs/client/command.md#Qcp)。

### 验证节点（validator）

验证节点相关概念和机制请参阅[验证人详解](../spec/validators/all_about_validators.md)和[QOS经济模型](../spec/validators/eco_module.md)。验证节点包含以下子命令：

* `qoscli tx create-validator`          [成为验证节点](#成为验证节点)
* `qoscli query validator`              [查询验证节点](#查询验证节点)
* `qoscli query validators`             [验证节点列表](#验证节点列表)
* `qoscli query validator-miss-vote`    [验证节点漏块信息](#查询验证节点漏块信息)
* `qoscli query community-fee-pool`     [社区收益池](#社区收益池)
* `qoscli tx revoke-validator`          [撤消验证节点](#撤销验证节点)
* `qoscli tx active-validator`          [激活验证节点](#激活验证节点)

#### 成为验证节点

`qoscli tx create-validator --name <validator_name> --owner <key_name_or_account_address> --tokens <tokens>`

主要参数：

- `--owner`         操作者账户地址或密钥库中密钥名字
- `--name`          验证节点名字
- `--nodeHome`      节点配置文件和数据所在目录，默认：`$HOME/.qosd`
- `--tokens`        绑定tokens，不能大于操作者持有QOS数量
- `--compound`      是否收益复投
- `--description`   备注

创建的validator基于本地的配置文件取`$HOME/.qosd/config/priv_validator.json`内信息，如果更改过默认位置，请使用`--home`指定`config`所在目录。

`Arya`初始化了一个[全节点](../install/testnet/fullnode.md)，可通过下面指令成为验证节点：
```bash
$ qoscli tx create-validator --name "Arya's node" --owner Arya --tokens 1000
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"258"}
```

执行成为验证节点命令后将从`Arya`账户扣除1000QOS，绑定到验证节点中，验证节点参与投票、打块所获得的挖矿收益将直接增加到`Arya`账户。

#### 查询验证节点

`qoscli query validator [validator-owner]`

`validator-owner`为账户地址或本地秘钥库名字

可根据操作者查找与其绑定的验证节点信息。

```bash
$ qoscli query validator address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy --indent
```

执行结果：
```bash
{
  "name": "Arya's node",
  "owner": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "validatorPubkey": {
    "type": "tendermint/PubKeyEd25519",
    "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
  },
  "bondTokens": "1000",
  "description": "",
  "status": 0,
  "inactiveCode": 0,
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "bondHeight": "258"
}
```

#### 验证节点列表

`qoscli query validators`

查询所有验证节点：
```bash
$ qoscli query validators --indent
```

执行结果：
```bash
validators: 
{
  "name": "Arya's node",
  "owner": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "validatorPubkey": {
    "type": "tendermint/PubKeyEd25519",
    "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
  },
  "bondTokens": "1000",
  "description": "",
  "status": 0,
  "inactiveCode": 0,
  "inactiveTime": "0001-01-01T00:00:00Z",
  "inactiveHeight": "0",
  "bondHeight": "258"
}
```

#### 查询验证节点漏块信息

`qoscli query validator-miss-vote [validator-owner]`

`validator-owner`为操作者账户地址或密钥库中密钥名字

查询`Arya`的节点漏块信息：
```bash
$ qoscli query validator-miss-vote Arya
```

执行结果：
```bash
{"startHeight":"258","endHeight":"387","missCount":0,"voteDetail":[]}
```

#### 验证节点窗口信息
`qoscli query validator-period --owner  <key_name_or_account_address>`

`key_name_or_account_address`为操作者账户地址或密钥库中密钥名字

查询`Arya`的节点漏块信息：
```bash
$ qoscli query validator-period --owner Arya
```

执行结果：
```bash
{
  "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "validator_pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
  },
  "fees": "0",
  "current_tokens": "4782741",
  "current_period": "15",
  "last_period": "14",
  "last_period_fraction": {
    "value": "1177.934327765593760252"
  }
}
```

#### 社区收益池
`qoscli query community-fee-pool`

查询社区收益：
```bash
$ qoscli query community-fee-pool
```

执行结果：
```bash
123456
```

#### 撤销验证节点

`qoscli tx revoke-validator --owner <key_name_or_account_address>`

`key_name_or_account_address`为操作者账户地址或密钥库中密钥名字

`Arya`将自己的节点撤销为为验证节点：
```bash
$ qoscli tx revoke-validator --owner Arya
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"268"}
```

执行撤销命令后`Arya`的节点将处于pending状态，不再参与投票和打块。

#### 激活验证节点

`qoscli tx active-validator --owner <key_name_or_account_address>`

`key_name_or_account_address`为操作者账户地址或密钥库中密钥名字

`Arya`将自己处于pending状态的节点重新激活为验证节点：
```bash
$ qoscli tx active-validator --owner Arya
```

执行结果：
```bash
{"check_tx":{},"deliver_tx":{},"hash":"BA45F8416780C76468C925E34372B05F5A7FEAAC","height":"275"}
```

执行成功，`Arya`的节点将继续参与投票、打块等共识职能，并获得挖矿奖励。



### 委托（delegate）

* `qoscli tx delegate`              [委托](#委托)
* `qoscli query delegation`         [委托查询](#委托查询)
* `qoscli query delegations-to`     [验证节点委托列表](#验证节点委托列表)
* `qoscli query delegations`        [代理用户委托列表](#代理用户委托列表)
* `qoscli query delegator-income`   [委托收益查询](#委托收益查询)
* `qoscli tx modify-compound`       [修改收益复投方式](#修改收益复投方式)
* `qoscli tx unbond`                [解除委托](#解除委托)
* `qoscli tx redelegate`            [变更委托验证节点](#变更委托验证节点)

#### 委托

`qoscli tx delegate --owner <validator_key_name_or_account_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --compound <compound_or_not>`

主要参数：

- `--owner`         代理验证节点操作账户地址或密钥库中密钥名字
- `--delegator`     被代理账户地址或秘钥库中秘钥名字
- `--tokens`        绑定tokens，不能大于`delegator`持有QOS数量
- `--compound`      收益是否复投，默认`false`

`Sansa`将自己的100个QOS代理给`Arya`创建的验证节点：
```bash
$ qoscli tx delegate --owner Arya --delegator Sansa --tokens 100
```

#### 委托查询

`qoscli query delegation --owner <validator_key_name_or_account_address> --delegator <delegator_key_name_or_account_address>`

主要参数：

- `--owner`         代理验证节点操作账户地址或密钥库中密钥名字
- `--delegator`     被代理账户地址或秘钥库中秘钥名字

`Sansa`在`Arya`上的代理信息：
```bash
$ qoscli query delegation --owner Arya --delegator Sansa
```

查询结果：
```bash
{
  "delegator_address": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
  "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "validator_pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
  },
  "delegate_amount": "100",
  "is_compound": false
}
```

#### 验证节点委托列表

`qoscli query delegations-to [validator-owner]`

主要参数：

- `validator-owner`     代理验证节点操作账户地址或密钥库中密钥名字

`Arya`验证节点上的所有代理信息：
```bash
$ qoscli query delegations-to Arya
```

查询结果示例：
```bash
[
  {
    "delegator_address": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
    "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
    },
    "delegate_amount": "100",
    "is_compound": false
  }
  ...
]
```

#### 代理用户委托列表

`qoscli query delegations [delegator]`

主要参数：

- `delegator`     被代理账户地址或秘钥库中秘钥名字

`Sansa`的所有代理信息：
```bash
$ qoscli query delegations Sansa
```

查询结果：
```bash
[
  {
    "delegator_address": "address1t7eadnyl8g6ct9xyrasvz4rdztvkeqpc0hzujh",
    "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
    "validator_pub_key": {
      "type": "tendermint/PubKeyEd25519",
      "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
    },
    "delegate_amount": "100",
    "is_compound": false
  }
]
```

#### 委托收益查询

`qoscli query delegator-income --owner <validator_key_name_or_account_address> --delegator <delegator_key_name_or_account_address`

主要参数：

- `--owner`         代理验证节点操作账户地址或密钥库中密钥名字
- `--delegator`     被代理账户地址或秘钥库中秘钥名字

`Sansa`查询代理给`Arya`的收益信息：
```bash
$ qoscli query delegator-income --owner Arya --delegator Sansa
```

查询结果：
```bash
{
  "owner_address": "address1ctmavdk57x0q7c9t98v7u79607222ars4qczcy",
  "validator_pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="
  },
  "previous_validaotr_period": "1",
  "bond_token": "100",
  "earns_starting_height": "101",
  "first_delegate_height": "1",
  "historical_rewards": "0",
  "last_income_calHeight": "101",
  "last_income_calFees": "0"
}
```

#### 修改收益复投方式

`qoscli tx modify-compound --owner <validator_key_name_or_account_address> --delegator <delegator_key_name_or_account_address> --compound <compound_or_not>`

主要参数：

- `--owner`         代理验证节点操作账户地址或密钥库中密钥名字
- `--delegator`     被代理账户地址或秘钥库中秘钥名字
- `--compound`      收益是否复投，默认`false`

`Sansa`将收益设置为复投方式：
```bash
$ qoscli tx modify-compound --owner Arya --delegator Sansa --compound
```

#### 解除委托

`qoscli tx unbond --owner <validator_key_name_or_account_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --all <unbond_all>`

主要参数：

- `--owner`         代理验证节点操作账户地址或密钥库中密钥名字
- `--delegator`     被代理账户地址或秘钥库中秘钥名字
- `--tokens`        解绑tokens，不能大于目前代理的QOS数量
- `--all`           是否取消全部QOS代理，默认false

`Sansa`解除代理给`Arya`的50个QOS：
```bash
$ qoscli tx unbond --owner Arya --delegator Sansa --tokens 50
```

#### 变更委托验证节点

`qoscli tx redelegate --from-owner <validator_key_name_or_account_address> --to-owner <validator_key_name_or_account_address> --delegator <delegator_key_name_or_account_address> --tokens <tokens> --all <unbond_all>`

主要参数：

- `--from-owner`    代理验证节点操作账户地址或密钥库中密钥名字
- `--to-owner`      新的代理验证节点操作账户地址或密钥库中密钥名字
- `--delegator`     被代理账户地址或秘钥库中秘钥名字
- `--tokens`        解绑并代理给新代理的tokens，不能大于目前代理的QOS数量
- `--compound`      新代理收益是否复投，默认`false`
- `--all`           是否从`from-owner`完全解绑，全部代理给`to-owner`，默认false

`Sansa`将代理给`Arya`的10个QOS转移到`John`操作的验证节点上：
```bash
$ qoscli tx redelegate --from-owner Arya --to-owner John --delegator Sansa --tokens 10
```

## tendermint

QOS中包含的tendermint提供的基础指令：

* `qoscli tendermint status`      查询节点状态
* `qoscli tendermint validators`  获取指定高度验证节点集合
* `qoscli tendermint block`       获取指定高度区块信息
* `qoscli tendermint txs`         根据标签查找交易
* `qoscli tendermint tx`          根据交易hash查询交易信息

更多tendermint使用说明参照[tendermint 官方文档](https://tendermint.com/docs/)