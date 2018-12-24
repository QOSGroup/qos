# QOS公链经济模型

QOS公链是基于[授权股权证明Delegated Proof-of-Stake]https://multicoin.capital/wp-content/uploads/2018/03/DPoS_-Features-and-Tradeoffs.pdf 和[拜占庭容错共识算法]https://en.wikipedia.org/wiki/Byzantine_fault_tolerance 的双层链机制的区块链基础设施。

![QOS经济模型概览](https://github.com/QOSGroup/qos/tree/master/docs/client/validators/eco_overview.png)

## QOS公链节点的构成

### 轻节点客户端（light-client）

*未来的版本支持

QOS轻节点可以执行QOScli支持的交易，不需要全部账本，仅验证少量头部信息及交易，需要较小资源

### 全节点（full-node）

和其他区块链网络相同，QOS公链全节点指包含全部账本的QOS节点。

QOS公链验证人，必须是QOS公链的全节点，但全节点需要发出[创建验证人交易]，并符合[一定条件]，才能成为验证人。

### 验证人（Validator）

QOS公链中有一个验证人节点的集合，验证人节点担当了BFT共识算法的具体实现——网络中的每一块都需要收集至少2/3的验证人节点签名。QOS公链中的每一块包含零到多条交易，验证人节点对块中的交易进行校验，对校验通过的块用自己的私钥签名，并广播到网络中去。

QOS公链验证人节点通过绑定一定的QOS，同时承担了DPOS算法的实现——依照其绑定的QOS数量，获得QOS网络挖矿的收益。详见[QOS公链挖矿机制]https://github.com/QOSGroup/qos/tree/master/docs/client/validators/eco_module.md#QOS公链挖矿机制

希望了解更多验证人节点的信息或希望成为QOS验证人，请查阅[验证人节点详解]https://github.com/QOSGroup/qos/tree/master/docs/client/validators/all_about_validators.md

## QOS公链挖矿机制
根据[白皮书]https://github.com/QOSGroup/whitepaper ，QOS公链的挖矿数额是按年度固定的，在主网上线的第一年内，每产生一个区块产生的QOS数量大体相同。
QOS网络中的全部[活跃的验证人]都可以依据其绑定的QOS数量占网络中总的绑定QOS的比例获得挖矿收益。

## QOS公链代理机制
在后续的QOS测试网版本中，对于没有能力成为验证人的节点，也将可以通过将其账户中的QOS委托给验证人的方式获得挖矿收益。
每个QOS验证人将可以发布一系列的委托合约，合约规定委托人通过将QOS在一定时间段内交予验证人作为绑定的QOS参与挖矿，获得的收益如何分配。