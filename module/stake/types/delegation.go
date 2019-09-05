package types

import btypes "github.com/QOSGroup/qbase/types"

type DelegationInfo struct {
	DelegatorAddr btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Amount        btypes.BigInt     `json:"delegate_amount"` // 委托数量
	IsCompound    bool              `json:"is_compound"`     // 是否复投
}

func NewDelegationInfo(deleAddr btypes.AccAddress, valAddr btypes.ValAddress, amount btypes.BigInt, isCompound bool) DelegationInfo {
	return DelegationInfo{deleAddr, valAddr, amount, isCompound}
}

// unbond
type UnbondingDelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`
	Height         int64             `json:"height"`
	CompleteHeight int64             `json:"complete_height"`
	Amount         btypes.BigInt     `json:"delegate_amount"`
}

func NewUnbondingDelegationInfo(deleAddr btypes.AccAddress, valAddr btypes.ValAddress, height, completeHeight int64, amount btypes.BigInt) UnbondingDelegationInfo {
	return UnbondingDelegationInfo{deleAddr, valAddr, height, completeHeight, amount}
}

// re delegate
type RedelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	FromValidator  btypes.ValAddress `json:"from_validator"`
	ToValidator    btypes.ValAddress `json:"to_validator"`
	Amount         btypes.BigInt     `json:"delegate_amount"`
	Height         int64             `json:"height"`
	CompleteHeight int64             `json:"complete_height"`
	IsCompound     bool              `json:"is_compound"` // 是否复投
}

func NewRedelegateInfo(deleAddr btypes.AccAddress, fromValAddr, toValAddr btypes.ValAddress, amount btypes.BigInt, height, completeHeight int64, isCompound bool) RedelegationInfo {
	return RedelegationInfo{deleAddr, fromValAddr, toValAddr, amount, height, completeHeight, isCompound}
}
