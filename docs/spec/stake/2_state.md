# 存储

## 验证人

验证人存储结构如下：
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

## 委托信息

委托信息的存储结构如下：
```go
type DelegationInfo struct {
	DelegatorAddr btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Amount        btypes.BigInt     `json:"delegate_amount"` // 委托数量
	IsCompound    bool              `json:"is_compound"`     // 是否复投
}
```

## 转委托信息

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

## 解除委托信息

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



