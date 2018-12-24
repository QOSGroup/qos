# QOS验证人节点详解

:::

*随着QOS版本迭代，本文档亦在更新中*

文中涉及参数（以$开头）的具体设置，可能与本文举例中的具体数字不同，测试网执行的参数可详见[测试网的genesis.json文件配置](https://github.com/QOSGroup/testnets)

:::

## QOS验证人的权利

* 对交易进行验证
* 获得挖矿收益
* 通过制定代理合约受益*（待实现功能）*
* 获得交易费用*（待实现功能）*

## QOS验证人的义务

* 保证稳定在线
* 对交易进行验证
* 保证自己的私钥安全
* 参与社区治理*（待实现功能）*

## 如何成为QOS验证人

### 硬件要求

要成为QOS验证人，必须首先成为一个全节点，保证持续稳定的在线运行来进行区块内交易的校验签名及广播，并为此采取必要的安全策略。通过测试，我们推荐验证人节点的硬件达到以下要求：

* 可以使用云服务器或独立机房，可持续不间断运行
* 带宽4M及以上，低延时公共网络
* 2核CPU，4G内存，100G硬盘存储空间

### 验证人的数量限制

QOS网络中将以验证人绑定QOS总数即权重从大到小排序，总数不超过$max_validator_cnt

*在测试网中，$max_validator_cnt=10000，相当于无限制*

### 验证人节点的几种状态

![验证人状态转换](https://github.com/QOSGroup/static/blob/master/validator_status.png?raw=true)

* **活跃状态**

保持不间断地验证区块交易，以私钥签名并广播的状态。
普通全节点，通过发出[create-validator](https://github.com/QOSGroup/qos/tree/master/docs/client/validators/all_about_validators.md#create-validator)交易，或者一个非活跃状态的验证人，通过[active-validator](https://github.com/QOSGroup/qos/tree/master/docs/client/validators/all_about_validators.md#active-validator)交易，可能转为活跃状态。

但并非任意全节点都可以通过以上方式成为活跃验证人，由于网络限制了总验证人数量，在一个特定时间，QOS网络以过去的$voting_status_len个块中，验证过并有签名的块数至少要达到$voting_status_least，来明确一个验证人节点是否活跃。我们称$voting_status_len为验证人保活窗口。

例如，测试网中的保活窗口宽度$voting_status_len=10000，最小保活块数$voting_status_least=5000

如果验证人未能达到这个要求，将被强制切换到[非活跃状态]

一个新创建或者重新激活的验证人，如果经历的总块数尚不足窗口宽度，但漏签块数已达$voting_status_least，也将被切换到非活跃状态

活跃状态的验证人，可以进行区块验证，可以提交区块，获得挖矿收益，可以通过达成代理合约获得收益，也可以获得交易费用。

* **非活跃状态**

由于未达到活跃窗口要求，或者通过发出[revoke-validator](https://github.com/QOSGroup/qos/tree/master/docs/client/validators/all_about_validators.md#revoke-validator)交易主动要求，验证人将转为非活跃状态。非活跃状态是验证人从活跃状态到退出状态之间所必须经历的中间态。

非活跃状态最久能够维持观察期即$survival_secs秒，非活跃的验证人如果什么都不做，经过$survival_secs后将自动退出，失去其验证人身份。

非活跃状态的验证人，不能进行区块验证，不能提交区块，不能获得挖矿收益和交易费用，不能达成代理合约，需要渡过观察期退出后，通过代理合约绑定的QOS才能回到投资者账户上。

* **退出状态**

退出状态的验证人将其上绑定的QOS自动返还给各投资者，自绑定的部分也会回到验证节点的所有者（owner）账户上。

退出后的验证人的权益与普通节点无异。

### 验证人节点的权重（voting power）

作为一个DPOS区块链网络，QOS网络中的验证人节点需要绑定一定量的QOS来构成其权益。

QOS目前规定验证人必须有一定的自绑定QOS来初始化运行验证人节点。创建后，其绑定的QOS可以来自于验证人所有者（owner）自己的账户，在createValidatorTX初始化时绑定，或者后期再绑定给自己*（后续版本）*；也可以通过发布和签订代理合约（delegation contract），来吸纳不具备代理人资格的节点的投资*（后续版本）*。

* 参与挖矿收益的分配

每出一个新块时，验证人的权重决定了其分配挖矿收益的比例，如下：


![挖矿分配](https://github.com/QOSGroup/static/blob/master/voting_power.png?raw=true)

* 社区自治的话语权

进行社区自治投票时，验证人的权重决定其决定的话语权比例。但普通节点也有社区自治的投票权，当验证人绑定的QOS来自普通节点的委托协议时，投资者的意志将覆盖验证人这部分权重*(待实现功能)*。

## 验证人交易类型

### create-validator

全节点通过发出create-validator交易来成为验证人，该交易需要提供以下参数：

- name 验证人的名字，必须提供

- owner 验证人节点所有者，对应keybase中的用户名或者地址（以"address"开头）

- pubkey 验证人节点公钥(ed25519)

- tokens 初始化自绑定的Token数量

- description 描述信息，可选

命令格式：

```
qoscli create-validator --name validatorName --owner ownerName --pubkey "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA=" --tokens 100
```

### revoke-validator

活跃的验证人放弃验证人身份，转为非活跃状态，该交易需要提供参数：

- owner 验证人节点所有者，对应keybase中的用户名或者地址（以"address"开头）

命令格式:

```
qoscli revoke-validator --owner ownerName
```

### active-validator

非活跃状态的验证人恢复活跃状态，该交易需要提供参数：

- owner 验证人节点所有者，对应keybase中的用户名或者地址（以"address"开头）

```
qoscli revoke-validator --owner ownerName
```
