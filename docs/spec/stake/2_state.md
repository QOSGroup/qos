# 存储

## 验证节点

验证节点存储结构如下：
```go
type Validator struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	Owner           btypes.AccAddress `json:"owner"`
	ConsPubKey      crypto.PubKey     `json:"consensus_pubkey"`
	BondTokens      btypes.BigInt     `json:"bond_tokens"`
	Description     Description       `json:"description"`
	Commission      Commission        `json:"commission"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactive_code"`
	InactiveTime   time.Time    `json:"inactive_time"`
	InactiveHeight int64        `json:"inactive_height"`

	MinPeriod  int64 `json:"min_period"`
	BondHeight int64 `json:"bond_height"`
}
```

验证节点投票窗口汇总信息：
```go
type ValidatorVoteInfo struct {
	StartHeight         int64 `json:"start_height"`
	IndexOffset         int64 `json:"index_offset"` //统计截止高度=StartHeight+IndexOffset-1
	MissedBlocksCounter int64 `json:"missed_blocks_counter"`
}
```

- 验证节点信息： `0x01 + valAddress -> amino(Validator)`
- 验证节点地址与共识地址对应关系： `0x02 + consensusAddress -> amino(ValAddress)`
- 处于`inactive`状态的验证节点： `0x04 + inactiveTime + validatorAddress -> amino(inactiveTime)`
- 处于`active`状态的验证节点： `0x05 + votePower + validatorAddress -> amino(true)`
- 验证节点投票统计信息： `0x11 + validatorAddress -> amino(ValidatorVoteInfo)`
- 验证节点投票信息： `0x12 + validatorAddress + indenx -> amino(true/false)`

## 委托

委托信息的存储结构如下：
```go
type DelegationInfo struct {
	DelegatorAddr btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Amount        btypes.BigInt     `json:"delegate_amount"` // 委托数量
	IsCompound    bool              `json:"is_compound"`     // 是否复投
}
```

解除委托信息的存储结构如下:
```go
type UnbondingDelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`
	Height         int64             `json:"height"`
	CompleteHeight int64             `json:"complete_height"`
	Amount         btypes.BigInt     `json:"delegate_amount"`
}
```

转委托信息的存储结构如下：
```go
type RedelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	FromValidator  btypes.ValAddress `json:"from_validator"`
	ToValidator    btypes.ValAddress `json:"to_validator"`
	Amount         btypes.BigInt     `json:"delegate_amount"`
	Height         int64             `json:"height"`
	CompleteHeight int64             `json:"complete_height"`
	IsCompound     bool              `json:"is_compound"` // 是否复投
}
```

- 委托信息： `0x31 + delegatorAddress + validatorAddress -> amino(DelegationInfo)`
- 验证节点与委托账户记录： `0x32 + validatorAddress + delegatorAddress -> amino(true)`
- 解绑信息： `0x41 + height + delegatorAddress + validatorAddress -> amino(UnbondingDelegationInfo)`
- 解绑信息： `0x42 + delegatorAddress + height + validatorAddress -> amino(true)`
- 解绑信息： `0x43 + validatorAddress + height + delegatorAddress -> amino(true)`
- 转委托信息： `0x51 + height + delegatorAddress + validatorAddress -> amino(RedelegationInfo)`
- 转委托信息： `0x52 + delegatorAddress + height + validatorAddress -> amino(true)`
- 转委托信息： `0x53 + validatorAddress + height + delegatorAddress -> amino(true)`


