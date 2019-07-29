package types

import btypes "github.com/QOSGroup/qbase/types"

type DelegationInfo struct {
	DelegatorAddr btypes.Address `json:"delegator_addr"`
	ValidatorAddr btypes.Address `json:"validator_addr"`
	Amount        uint64         `json:"delegate_amount"` // 委托数量。TODO 注意溢出处理
	IsCompound    bool           `json:"is_compound"`     // 是否复投
}

func NewDelegationInfo(delAddr btypes.Address, valAddr btypes.Address, amount uint64, isCompound bool) DelegationInfo {
	return DelegationInfo{delAddr, valAddr, amount, isCompound}
}

// unbond
type UnbondingDelegationInfo struct {
	DelegatorAddr  btypes.Address `json:"delegator_addr"`
	ValidatorAddr  btypes.Address `json:"validator_addr"`
	Height         uint64         `json:"height"`
	CompleteHeight uint64         `json:"complete_height"`
	Amount         uint64         `json:"delegate_amount"`
}

func NewUnbondingDelegationInfo(delAddr btypes.Address, valAddr btypes.Address, height uint64, completeHeight uint64, amount uint64) UnbondingDelegationInfo {
	return UnbondingDelegationInfo{delAddr, valAddr, height, completeHeight, amount}
}

// re delegate
type RedelegationInfo struct {
	DelegatorAddr  btypes.Address `json:"delegator_addr"`
	FromValidator  btypes.Address `json:"from_validator"`
	ToValidator    btypes.Address `json:"to_validator"`
	Amount         uint64         `json:"delegate_amount"`
	Height         uint64         `json:"height"`
	CompleteHeight uint64         `json:"complete_height"`
	IsCompound     bool           `json:"is_compound"` // 是否复投
}

func NewRedelegateInfo(delAddr btypes.Address, fromValAddr btypes.Address, toValAddr btypes.Address, amount uint64, height uint64, completeHeight uint64, isCompound bool) RedelegationInfo {
	return RedelegationInfo{delAddr, fromValAddr, toValAddr, amount, height, completeHeight, isCompound}
}
