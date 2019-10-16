# Parameters

```go
type Params struct {
	// 普通提议参数
	// 抵押参数
	NormalMinDeposit             btypes.BigInt `json:"normal_min_deposit"`               //  最小抵押
	NormalMinProposerDepositRate qtypes.Dec    `json:"normal_min_proposer_deposit_rate"` //  提议者最小初始抵押
	NormalMaxDepositPeriod       time.Duration `json:"normal_max_deposit_period"`        //  最长抵押时间
	// 投票参数
	NormalVotingPeriod time.Duration `json:"normal_voting_period"` //  投票期时长
	// 记票参数
	NormalQuorum    qtypes.Dec `json:"normal_quorum"`    //  最小当前网络voting power投票比例
	NormalThreshold qtypes.Dec `json:"normal_threshold"` //  提议通过，最小投Yes比例
	NormalVeto      qtypes.Dec `json:"normal_veto"`      //  强烈反对，最小投Veto比例
	NormalPenalty   qtypes.Dec `json:"normal_penalty"`   //  微投票验证节点惩罚比例
	// 扣留参数
	NormalBurnRate qtypes.Dec `json:"normal_burn_rate"` // 抵押扣留比例

	// 重要提议参数
	// 抵押参数
	ImportantMinDeposit             btypes.BigInt `json:"important_min_deposit"`               //  最小抵押
	ImportantMinProposerDepositRate qtypes.Dec    `json:"important_min_proposer_deposit_rate"` //  提议者最小初始抵押
	ImportantMaxDepositPeriod       time.Duration `json:"important_max_deposit_period"`        //  最长抵押时间
	// 投票参数
	ImportantVotingPeriod time.Duration `json:"important_voting_period"` //  投票期时长
	// 记票参数
	ImportantQuorum    qtypes.Dec `json:"important_quorum"`    //  最小当前网络voting power投票比例
	ImportantThreshold qtypes.Dec `json:"important_threshold"` //  提议通过，最小投Yes比例
	ImportantVeto      qtypes.Dec `json:"important_veto"`      //  强烈反对，最小投Veto比例
	ImportantPenalty   qtypes.Dec `json:"important_penalty"`   //  微投票验证节点惩罚比例
	// 扣留参数
	ImportantBurnRate qtypes.Dec `json:"important_burn_rate"` // 抵押扣留比例

	// 危险提议参数
	// 抵押参数
	CriticalMinDeposit             btypes.BigInt `json:"critical_min_deposit"`               //  最小抵押
	CriticalMinProposerDepositRate qtypes.Dec    `json:"critical_min_proposer_deposit_rate"` //  提议者最小初始抵押
	CriticalMaxDepositPeriod       time.Duration `json:"critical_max_deposit_period"`        //  最长抵押时间
	// 投票参数
	CriticalVotingPeriod time.Duration `json:"critical_voting_period"` //  投票期时长
	// 记票参数
	CriticalQuorum    qtypes.Dec `json:"critical_quorum"`    //  最小当前网络voting power投票比例
	CriticalThreshold qtypes.Dec `json:"critical_threshold"` //  提议通过，最小投Yes比例
	CriticalVeto      qtypes.Dec `json:"critical_veto"`      //  强烈反对，最小投Veto比例
	CriticalPenalty   qtypes.Dec `json:"critical_penalty"`   //  微投票验证节点惩罚比例
	// 扣留参数
	CriticalBurnRate qtypes.Dec `json:"critical_burn_rate"` // 抵押扣留比例
}
```

所有参数均可通过[治理](../gov)模块提交参数修改提议进行修改。