# Governance


QOS的链上治理提案类型包括以下几种：

* 普通文本：提议者可以自由定义提案内容
* 参数修改：提议者针对某个配置参数提出修改意见
* 社区基金提取：提议者对如何适用社区基金提出意见，提案通过则可以将一定比例的社区基金提取到某一账户

## 提案过程

### 提案抵押

为了避免攻击，每个提案都需要抵押一定的QOS作为`Deposit`。
提交提议者至少抵押30%的`$MinDeposit`，当抵押金超过`$MinDeposit`，才能进入投票阶段。该提议超过`$MaxDepositPeriod`，还未进超过`$MinDeposit`，则提议会被删除，且押金不会返还。
提议进入投票阶段后，依然可以进行抵押。

* `MinDeposit` QOS最小抵押
* `MaxDepositPeriod` 抵押最大时长
* `VotingPeriod` 投票时长
* `Quorum` voting power最小比例
* `Threshold` 判定通过需要的`Yes`比例
* `Veto` 判定强烈不同意需要的`Veto`比例
* `Penalty` 不投票验证节点的惩罚比例

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

* 无效：参与投票的voting power/总voting power < `$participation`

* 通过：投`Yes`的voting power/总voting power > `$Threshold`

* 强烈反对：投`NoWithVeto`的voting power/总voting power > `$Veto`

* 未通过：除以上以外其他结果

### 销毁机制

提案`通过`或`未通过`，都要销毁抵押`Deposit`的20%，作为治理的费用，把剩余的`Deposit`原路返回。

但如果投票结果为`强烈反对``，则抵押`Deposit`会被全部收入到社区基金账户。

### 惩罚机制

如果一个账户提议进入投票阶段，他是验证人，然后该提议进入统计阶段，他还是验证人，但是他并没有投票，则会按`Penalty`的比例被惩罚。一个验证节点上的任何一个验证人或委托人投过票就算该验证节点参与过投票。

