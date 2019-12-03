# ABCI

## BeginBlocker

### 双签惩罚

当证据池中存在`ABCIEvidenceTypeDuplicateVote`类型证据时，根据参数`max_evidence_age`验证证据时效，验证作案高度。
对于需要惩罚的双签节点，根据作案高度节点绑定`tokens`和参数`slash_fraction_double_sign`计算惩罚数量，优先从作案高度之后验证节点相关的`unbonding`和`redelegation`中扣除，
未扣满部分从验证节点当前委托中扣除。

### 漏签惩罚

记录验证节点投票信息。

QOS网络期望验证节点能够时时在线，参与打快和投票。根据参数`voting_status_len`和`voting_status_least`，要求验证节点在每`voting_status_len`的区块高度至少参与`voting_status_least`块投票。
否则将参照`slash_fraction_downtime`扣除验证节点当前绑定`tokens`。

## EndBlocker

### 到期解绑委托

对于到期的解除绑定信息，返回绑定`tokens`到委托账户，删除解绑信息。

### 到期转委托

对于到期转委托信息，创建委托账户到`to-validator`的委托关系，删除转委托信息。