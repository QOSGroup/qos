# 存储

`MapperName`为`governance`

## 提议

全网递增提议ID:

- newProposalID: `newProposalID -> amino(newProposalID)`

每提交一个提议，`newProposalID`加1。


存储结构：
```go
type Proposal struct {
	ProposalContent `json:"proposal_content"` // 提议内容，不同提议类型具体结构不同

	ProposalID int64 `json:"proposal_id"` //  提议ID

	Status           ProposalStatus `json:"proposal_status"`    //  提议状态
	FinalTallyResult TallyResult    `json:"final_tally_result"` //  投票结果统计

	SubmitTime     time.Time     `json:"submit_time"`      // 提议时间
	DepositEndTime time.Time     `json:"deposit_end_time"` // 质押最迟结束时间
	TotalDeposit   btypes.BigInt `json:"total_deposit"`    // 累计质押

	VotingStartTime   time.Time `json:"voting_start_time"` // 投票开始时间
	VotingStartHeight int64     `json:"voting_start_height"` // 投票开始高度
	VotingEndTime     time.Time `json:"voting_end_time"` // 投票结束时间
}
```

- 提议 `proposals:{proposal_id} -> amino(Proposal)`
- 质押阶段提议 `inactiveProposalQueue:{deposit_end_time}:{proposal_id} -> amino(proposal_id)`
- 投票阶段提议 `activeProposalQueue:{voting_end_time}:{proposal_id} -> amino(proposal_id)`

## 质押

存储结构:
```go
type Deposit struct {
	Depositor  types.AccAddress `json:"depositor"`   //  质押账户地址
	ProposalID int64            `json:"proposal_id"` //  提议ID
	Amount     types.BigInt     `json:"amount"`      //  质押QOS
}
```

- 质押 `deposits:{proposal_id}:{depositor} -> amino(Deposit)`

## 投票

存储结构：
```go
type Vote struct {
	Voter      types.AccAddress `json:"voter"`       //  投票账户地址
	ProposalID int64            `json:"proposal_id"` //  提议ID
	Option     VoteOption       `json:"option"`      //  投票选项
}
```

- 投票 `votes:{proposal_id}:{voter} -> amino(Vote)`