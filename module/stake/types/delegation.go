package types

import btypes "github.com/QOSGroup/qbase/types"

type DelegationInfo struct {
	DelegatorAddr btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr btypes.ValAddress `json:"validator_addr"`
	Amount        uint64            `json:"delegate_amount"` // 委托数量。TODO 注意溢出处理
	IsCompound    bool              `json:"is_compound"`     // 是否复投
}

func NewDelegationInfo(deleAddr btypes.AccAddress, valAddr btypes.ValAddress, amount uint64, isCompound bool) DelegationInfo {
	return DelegationInfo{deleAddr, valAddr, amount, isCompound}
}

// unbond
type UnbondingDelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`
	Height         uint64            `json:"height"`
	CompleteHeight uint64            `json:"complete_height"`
	Amount         uint64            `json:"delegate_amount"`
}

func NewUnbondingDelegationInfo(deleAddr btypes.AccAddress, valAddr btypes.ValAddress, height uint64, completeHeight uint64, amount uint64) UnbondingDelegationInfo {
	return UnbondingDelegationInfo{deleAddr, valAddr, height, completeHeight, amount}
}

// re delegate
type RedelegationInfo struct {
	DelegatorAddr  btypes.AccAddress `json:"delegator_addr"`
	FromValidator  btypes.ValAddress `json:"from_validator"`
	ToValidator    btypes.ValAddress `json:"to_validator"`
	Amount         uint64            `json:"delegate_amount"`
	Height         uint64            `json:"height"`
	CompleteHeight uint64            `json:"complete_height"`
	IsCompound     bool              `json:"is_compound"` // 是否复投
}

func NewRedelegateInfo(deleAddr btypes.AccAddress, fromValAddr, toValAddr btypes.ValAddress, amount uint64, height uint64, completeHeight uint64, isCompound bool) RedelegationInfo {
	return RedelegationInfo{deleAddr, fromValAddr, toValAddr, amount, height, completeHeight, isCompound}
}
