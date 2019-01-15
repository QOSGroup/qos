package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

type DelegationInfo struct{
	DelegatorAddr btypes.Address	`json:"delegator_addr"`
	ValidatorAddr btypes.Address	`json:"validator_addr"`
	Amount uint64					`json:"delegate_amount"`	// 委托数量。TODO 注意溢出处理
	IsCompound bool					`json:"is_compound"`		// 是否复投
}

func NewDelegationInfo(delAddr btypes.Address, valAddr btypes.Address, amount uint64, isCompound bool) DelegationInfo{
	return DelegationInfo{ delAddr,valAddr, amount, isCompound}
}
