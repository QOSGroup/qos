# QOS公链经济模型 v1
[toc]

QOS公链是基于[授权股权证明Delegated Proof-of-Stake](https://multicoin.capital/wp-content/uploads/2018/03/DPoS_-Features-and-Tradeoffs.pdf)和[拜占庭容错共识算法](https://en.wikipedia.org/wiki/Byzantine_fault_tolerance)的双层链机制的区块链基础设施。

![QOS经济模型概览](https://github.com/QOSGroup/static/blob/master/eco_overview.png?raw=true)

## 设计原则原则

* 不管是验证人还是委托人，币有相同的收益权
* 不管是验证人还是委托人，币有相同的治理权
* 惩罚验证人时，除了验证人自己的币，委托人的币也会受惩罚，因此委托人需要对验证人进行考察后进行委托；
* 一个验证人能无限接收委托人委托的币，这会增加验证人的权益，QOS公链对此不做机制上的干涉，但验证人的当前数据将会公开，由委托人自己决定是否要委托给该验证人；
* 验证人排名的依据：自抵押qos+委托人委托的qos。排名前`$max_validator_cnt`的验证人可以进行验证挖矿；
* 验证人挖矿收益的到帐方式：系统定期打款到帐，并可选择自动复投。

## 角色

### [轻钱包（litewallet）](https://github.com/QOSGroup/litewallet)

即QOS轻节点客户端，可以执行QOScli支持的交易，不需要全部账本，仅验证少量头部信息及交易，需要较小资源，用于普通用户/手机客户端。
目前除qos外，还支持Ethereum、cosmos等账户管理、委托功能。详见[litewallet repo](https://github.com/QOSGroup/litewallet)

### 全节点（full-node）

和其他区块链网络相同，QOS公链全节点指包含全部账本的QOS节点。

### 验证人（Validator）

QOS公链中有一个验证人节点的集合，验证人节点担当了BFT共识算法的具体实现——网络中的每一块都需要收集至少2/3的验证人节点签名。QOS公链中的每一块包含零到多条交易，验证人节点对块中的交易进行校验，对校验通过的块用自己的私钥签名，并广播到网络中去。

QOS公链验证人，必须是QOS公链的全节点，但全节点需要发出[创建验证人交易](all_about_validators.md#create-validator)，并符合[一定条件](all_about_validators.md#如何成为QOS验证人)，才能成为验证人。

QOS公链验证人节点通过绑定一定的QOS，同时承担了DPOS算法的实现——依照其绑定的QOS数量，获得QOS网络挖矿的收益。详见[QOS公链挖矿机制](#QOS公链挖矿机制)

希望了解更多验证人节点的信息或希望成为QOS验证人，请查阅[验证人节点详解](all_about_validators.md)

### 委托人（Delegator）

对于自己没有能力或者意愿来自己运行一个验证节点，但希望得到挖矿收益的QOS持有者，可以选择一个验证人，通过委托（delegation）将QOS投入到该验证人的总绑定数中，增加验证人的投票权重，收到相应的挖矿收益作为回报。
关于委托收益的计算，详见[QOS公链代理机制](#QOS公链代理机制)

委托人可以不运行QOS全节点，通过轻钱包就可以进行委托操作。

委托人分享验证人出块的收益，意味着他们也分担验证人的责任和义务。当验证人因宕机/作恶而受到惩罚，其委托人也会受到相应的惩罚。

在社区自治（待实现功能）中，委托人和验证人拥有同等的投票权。

因此即使没有运行全节点，占网络最大数量的委托人依然担任着主动且重要的角色，即他们要选择可信、稳定的验证人，来增加这些验证人的投票权重，并关注验证人的动向，以维护网络的安全和稳定。

## 模块

### QOS通胀机制

根据[白皮书](https://github.com/QOSGroup/whitepaper)，QOS公链的挖矿数额是按年度固定的，在主网上线的第一年内，每产生一个区块产生的QOS数量大体相同。

主网通胀计划：

时间|第一个四年|第二个四年|第三个四年|第四个四年|第五个四年|第六个四年|第七个四年
:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:
新铸币数量（亿）|25.5|12.75|6.375|3.1875|1.59375|0.796875|0.796875

我们将其中每个四年定义为一个inflation_phrase通胀阶段，由endtime和total_amount组成，applied_amount标识本阶段已经分发的QOS，一个阶段结束，即进入下一阶段。
测试网通胀依照测试目的另外制定，详情可见[测试网的genesis.json文件配置](https://github.com/QOSGroup/qos-testnets)中的"mint"-"params"-"inflation_phrases"，例如：

```
        "inflation_phrases": [
          {
            "endtime": "2023-01-01T00:00:00Z",
            "total_amount": "2500000000000",
            "applied_amount": "0"
          },
          {
            "endtime": "2027-01-01T00:00:00Z",
            "total_amount": "12750000000000",
            "applied_amount": "0"
          },
          {
            "endtime": "2031-01-01T00:00:00Z",
            "total_amount": "6375000000000",
            "applied_amount": "0"
          },
          {
            "endtime": "2035-01-01T00:00:00Z",
            "total_amount": "3185000000000",
            "applied_amount": "0"
          }
        ]
```

可以通过社区投票修改参数来制定通胀策略，修改通胀计划。

每一块通胀的QOS数：

![每一块的通胀总量](https://github.com/QOSGroup/static/blob/master/rewardPerBlock.png?raw=true)

### 社区基金

在QOS每一块通胀的qos中，将有`$community_reward_rate`的QOS归属于社区基金，社区基金用于社区运营建设、奖励开发者、奖励有价值的生态推广（如社区认可的QSC联盟链）等活动。

社区基金的账户公开透明，用户可以发起`TaxUsage`类型的自治提议，申请将部分社区基金打入某一QOS账户对其进行使用，每个参与验证委托的QOS持有者都有权对该提议进行投票表态。

### QOS验证/委托挖矿机制

#### 验证人

QOS公链中有一个验证人节点的集合，验证人节点担当了BFT共识算法的具体实现——网络中的每一块都需要收集至少2/3的验证人节点签名。QOS公链中的每一块包含零到多条交易，验证人节点对块中的交易进行校验，对校验通过的块用自己的私钥签名，并广播到网络中去。

每一块都有一个验证人来进行打块（proposer），该验证人会有4%的额外收益：

![出块验证人收益](https://github.com/QOSGroup/static/blob/master/proposerReward.png?raw=true)

验证人打块的机会是与其绑定QOS数成正比的，因此打块的额外收益不会改变每个验证人在网络中的投票权重。

QOS公链验证人节点通过绑定一定的QOS，同时承担了DPOS算法的实现——依照其绑定的QOS数量，获得QOS网络挖矿的收益。

![验证人（及其委托人）单块总收益](https://github.com/QOSGroup/static/blob/master/validatorReward.png?raw=true)

希望了解更多验证人节点的信息或希望成为QOS验证人，请查阅[验证人节点详解](all_about_validators.md)

#### 委托人

验证人所绑定的QOS由两部分组成：验证人自己绑定的（self-bond），委托人委托给验证人的(delegation-bond)

**验证人总绑定(投票权重) = 验证人自绑定QOS数量 + ∑委托人委托给该验证人QOS数量**

对于委托人，其委托的QOS可以从验证人的总收入中获得相应比例的收益。由于验证人付出了人力和物力，验证人可以从总收益中抽取一定比例的佣金，QOS网络中的验证人佣金是统一的，以参数`$validator_commission_rate`定义。

![验证人自身每块收益](https://github.com/QOSGroup/static/blob/master/validatorSelfReward.png?raw=true)

![委托人每块收益](https://github.com/QOSGroup/static/blob/master/delegatorReward.png?raw=true)

* 分配周期

创建delegate后，由`$delegator_income_period_height`参数定义之后的每多少块为一个*分配周期*（在capricorn-2000测试网中为30块），在每个周期交替时为委托人分配收益/处理请求。

委托人后期追加委托QOS、unbond等的操作不会影响分配周期。

当前周期对绑定QOS的增减，对配置参数的修改，到下一周期开始时生效。

委托人在一个周期内多次修改同一配置项（例如是否复投），以该周期内最后一次修改为准，应用到下一周期。

* 解除绑定

通过`TxUnbondDelegate`解绑QOS时的收益计算

```

     |                x                             y          |
     |  --------------------------------|----------------------|
     |                                 unbond -----------------|-------------------------|

上次收益发放                                             下次收益发放                unbond解绑

```

unbond后,对应的validator将会增加一个计费点,unbond金额将在`unbond周期`之后返还至delegator账户.
unbond操作立即生效, 先统计出当前收益,并追加到下次收益发放总额中.

下次收益发放时, 发放金额为 x + y：
当解绑QOS为当前所有绑定QOS时，y = 0;
解绑QOS为部分绑定QOS时， x > y > 0

* 复投

委托人可以指定、并后期修改是否复投（`is_compound`）。复投表示前一周期产生的收益自动绑定并参与到下一周期的委托中，若不复投则收益自动打入委托人账户。

复投可以持续、自动地扩大委托挖矿的投资规模，是一种好的选择，但需要注意，绑定的token赎回需要经过一个由参数`$unbond_return_height`定义块数的*冻结期*后才能回到委托人账户，盲目扩大委托绑定规模不利于流动性。

### 惩罚机制

验证人受到惩罚的原因是由于其有意作恶/无意犯错，或者没有履行到验证人的义务。
另一方面，QOS网络维护不仅需要验证人，也需要委托人。委托人并非完全被动、只追求收益的角色，而是能够祈祷对验证人进行主动筛选、监督的作用，并在社区自治中发出自己的声音。
在这个基本思想的指导下，在QOS中，验证人受到惩罚的形式是以消减其绑定QOS的方式，其委托人绑定的QOS也会受到等比例的惩罚。

#### 验证人作恶

在QOS网络中，验证人作恶是指：

* 双签（double-signing）

在QOS网络的同一高度上，同一验证人签名一次以上，并广播不同的信息到网络中。
在BFT网络中，双签被视为拜占庭节点，拜占庭节点超过2/3时，网络一定会分叉，因此我们将双签视为严重的错误，并施以较高的罚金——销毁验证人及其委托人`$`比例的绑定QOS。

在实际操作中，验证人双签往往由于无意的失误，包括：
私钥被盗
配置错误导致同一验证节点启动两次或以上

#### 验证人不参与社区自治的情形

在社区自治投票中，验证人及其委托人均未参与投票时，验证人（及其委托人）将会受到到`$penalty`比例的惩罚。

### gas

