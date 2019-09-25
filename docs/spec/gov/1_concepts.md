# 概念

## 提案过程

### 提案抵押

为了避免攻击，每个提案都需要抵押一定的QOS作为`Deposit`。
提交提议者至少抵押`$min_deposit * $min_proposer_deposit_rate`，当抵押金超过`$min_deposit`，才能进入投票阶段。该提议超过`$max_deposit_period`，总抵押还未超过`$min_deposit`，则提议会被删除，且押金不会返还。
提议进入投票阶段后，依然可以进行抵押。

* `min_deposit` QOS最小抵押
* `min_proposer_deposit_rate` 提议者最小抵押比例
* `max_deposit_period` 抵押最大时长
* `voting_period` 投票时长
* `quorum` voting power最小比例
* `threshold` 判定通过需要的`Yes`比例
* `veto` 判定强烈不同意需要的`Veto`比例
* `penalty` 不投票验证节点的惩罚比例
* `burn_rate` 提议通过(`PASS`)或不通过(`REJECT`)抵押销毁比例

### 投票

投票选项包括：
* `Yes`同意
* `Abstain`弃权
* `No`不同意
* `NoWithVeto`强烈不同意。

只有验证人和委托人投票有效，重复投票以最后一次投票所提交的选项计票。

### 计票

验证人/委托人的投票以其绑定的QOS量作为权重，绑定的QOS越多，拥有越多话语权。通过验证委托将QOS绑定才能在链上治理投票中拥有话语权，因此绑定的QOS也称为voting power。

统计结果以`选项投票者的voting power/全网所有voting power`比值进行计算。

计票会产生以下几种结果：

* 无效：参与投票的voting power/总voting power < `$quorum`

* 通过：投`Yes`的voting power/总voting power > `$threshold`

* 强烈反对：投`NoWithVeto`的voting power/总voting power > `$veto`

* 未通过：除以上以外其他结果

### 销毁机制

提案`通过`或`未通过`，都要销毁抵押`Deposit * $burn_rate`，作为治理的费用，把剩余的`Deposit`原路返回。

但如果投票结果为`强烈反对`，则抵押`Deposit`会被全部收入到社区基金账户。

### 惩罚机制

如果一个节点在同一议案的进入投票阶段和计票统计阶段的块高度上，都是验证人，但并没有参与过该议案的投票，则其上绑定的token会按`Penalty`的比例被惩罚。

需要注意，判断验证人是否投票的标准，含绑定该验证人的委托人，一个验证节点上的验证人或任何一个委托人投过票就算该验证节点参与过投票。反之，如果该验证人因为未投票而受到惩罚，绑定该验证人的委托人的token会受到同比例惩罚。


