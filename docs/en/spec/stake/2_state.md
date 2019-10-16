# State

## Validator

validator:
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

validator vote:
```go
type ValidatorVoteInfo struct {
	StartHeight         int64 `json:"start_height"`
	IndexOffset         int64 `json:"index_offset"` 
	MissedBlocksCounter int64 `json:"missed_blocks_counter"`
}
```

- validator: `0x01 + valAddress -> amino(Validator)`
- validator: `0x02 + consensusAddress -> amino(ValAddress)`
- validator in `inactive` status: `0x04 + inactiveTime + validatorAddress -> amino(inactiveTime)`
- validator in `active` status: `0x05 + votePower + validatorAddress -> amino(true)`
- validator voteInfo: `0x11 + validatorAddress -> amino(ValidatorVoteInfo)`
- validator voteInfo: `0x12 + validatorAddress + indenx -> amino(true/false)`

## Delegation

delegation:
```go
type DelegationInfo struct {
	DelegatorAddr btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Amount        btypes.BigInt     `json:"delegate_amount"` 
	IsCompound    bool              `json:"is_compound"`    
}
```

unbonding delegation:
```go
type UnbondingDelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`
	Height         int64             `json:"height"`
	CompleteHeight int64             `json:"complete_height"`
	Amount         btypes.BigInt     `json:"delegate_amount"`
}
```

redelegation:
```go
type RedelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	FromValidator  btypes.ValAddress `json:"from_validator"`
	ToValidator    btypes.ValAddress `json:"to_validator"`
	Amount         btypes.BigInt     `json:"delegate_amount"`
	Height         int64             `json:"height"`
	CompleteHeight int64             `json:"complete_height"`
	IsCompound     bool              `json:"is_compound"` 
}
```

- delegation: `0x31 + delegatorAddress + validatorAddress -> amino(DelegationInfo)`
- delegation: `0x32 + validatorAddress + delegatorAddress -> amino(true)`
- unbonding: `0x41 + height + delegatorAddress + validatorAddress -> amino(UnbondingDelegationInfo)`
- unbonding: `0x42 + delegatorAddress + height + validatorAddress -> amino(true)`
- unbonding: `0x43 + validatorAddress + height + delegatorAddress -> amino(true)`
- redelegation: `0x51 + height + delegatorAddress + validatorAddress -> amino(RedelegationInfo)`
- redelegation: `0x52 + delegatorAddress + height + validatorAddress -> amino(true)`
- redelegation: `0x53 + validatorAddress + height + delegatorAddress -> amino(true)`


