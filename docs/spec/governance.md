# Governance

QOS包括以下治理策略：

1. 普通文本提议的链上治理
2. 参数修改提议的链上治理
3. 社区费池提取提议的链上治理

## 交互流程

### 提议参数

* `MinDeposit` QOS最小抵押
* `MaxDepositPeriod` 抵押最大时长 
* `VotingPeriod` 投票时长
* `Quorum` voting power最小比例
* `Threshold` 判定通过需要的`Yes`比例
* `Veto` 判定强烈不同意需要的`Veto`比例
* `Penalty` 不投票验证节点的惩罚比例

### 抵押阶段
提交提议者至少抵押30%的 `MinDeposit` ，当抵押金超过 `MinDeposit` ,才能进入投票阶段。该提议超过 `MaxDepositPeriod` ，还未进超过 `MinDeposit`，则提议会被删除，并不会返还抵押金。 
可以对进入投票阶段的提议再进行抵押。

### 投票阶段
只有验证人和委托人投票有效，重复投票以最后一次投票统计。投票选项有：`Yes`同意, `Abstain`弃权,`No`不同意,`NoWithVeto`强烈不同意。

### 统计阶段

统计结果有三类：同意，不同意，强烈不同意。

在所有有效投票者的`voting_power`占系统总的`voting_power`的比例超过`participation`的前提下,如果强烈反对的`voting_power`占所有投票者的`voting_power` 超过 `veto`, 结果是强烈不同意。如果没有超过且赞同的`voting_power`占所有投票者的`voting_power` 超过 `threshold`，提议结果是同意。其他情况皆为不同意。

### 销毁机制

提议通过或未通过，都要销毁`Deposit`的20%，作为治理的费用，把剩余的`Deposit`按比例原路退回。但如果是强烈不同意，则把`Deposit`全部累加到社区费池中。

### 惩罚机制

如果一个账户提议进入投票阶段，他是验证人，然后该提议进入统计阶段，他还是验证人，但是他并没有投票，则会按`Penalty`的比例被惩罚。一个验证节点上的任何一个验证人或委托人投过票就算该验证节点参与过投票。

## 提议类型

以下三种提议[操作说明](../client/qoscli.md#提交提议)

### TxProposal

好的想法和建议可以通过提交到线上治理，投票表决后QOS社区会执行相应内容实现。

```go
type ProposalType byte

const (
	ProposalTypeNil             ProposalType = 0x00
	ProposalTypeText            ProposalType = 0x01
	ProposalTypeParameterChange ProposalType = 0x02
	ProposalTypeTaxUsage        ProposalType = 0x03
)

type TxProposal struct {
	Title          string              `json:"title"`          
	Description    string              `json:"description"`    
	ProposalType   gtypes.ProposalType `json:"proposal_type"`  
	Proposer       btypes.Address      `json:"proposer"`       
	InitialDeposit uint64              `json:"initial_deposit"`
}
```
字段说明：
- Title 标题
- Description 描述
- ProposalType 类型
- Proposer 提议账户地址
- InitialDeposit 初始抵押


### TxParameterChange

提议修改QOS运行网络参数配置，提议通过后新参数会实时生效。

可修改参数通过[参数查询](../client/qoscli.md#参数查询)获取，包含`gov`、`distribution`、`stake`三个模块参数配置。

```go
type TxParameterChange struct {
	TxProposal
	Params []gtypes.Param `json:"params"`
}

type Param struct {
	Module string `json:"module"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
```

字段说明：
- TxProposal 提议基础信息
- Params 参数列表

### TxTaxUsage

社区费池QOS提取，接收地址仅能为Genesis [Guardian](guardian.md)

```go
type TxTaxUsage struct {
	TxProposal
	DestAddress btypes.Address `json:"dest_address"`
	Percent     types.Dec      `json:"percent"`
}
```

字段说明：
- TxProposal 提议基础信息
- DestAddress 接收地址
- Percent 提取比例

