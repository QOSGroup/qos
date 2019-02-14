package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

const (
	AddrLen = 20
)

type DelegationInfo struct {
	DelegatorAddr btypes.Address `json:"delegator_addr"`
	ValidatorAddr btypes.Address `json:"validator_addr"`
	Amount        uint64         `json:"delegate_amount"` // 委托数量。TODO 注意溢出处理
	IsCompound    bool           `json:"is_compound"`     // 是否复投
}

func NewDelegationInfo(delAddr btypes.Address, valAddr btypes.Address, amount uint64, isCompound bool) DelegationInfo {
	return DelegationInfo{delAddr, valAddr, amount, isCompound}
}

//DelegatorEarningsStartInfo delegator计算收益信息
type DelegatorEarningsStartInfo struct {
	PreviousPeriod        uint64        `json:"previous_period"`
	BondToken             uint64        `json:"bond_token"`
	CurrentStartingHeight uint64        `json:"earns_starting_height"`
	FirstDelegateHeight   uint64        `json:"first_delegate_height"`
	HistoricalRewardFees  btypes.BigInt `json:"historical_rewards"`
	LastIncomeCalHeight   uint64        `json:"last_income_calHeight"`
	LastIncomeCalFees     btypes.BigInt `json:"last_income_calFees"`
}

//ValidatorCurrentPeriodSummary validator当前周期收益信息
type ValidatorCurrentPeriodSummary struct {
	Fees   btypes.BigInt `json:"fees"`
	Period uint64        `json:"period"`
}

type ValidatorVoteInfo struct {
	StartHeight         uint64 `json:"startHeight"`
	IndexOffset         uint64 `json:"indexOffset"` //统计截止高度=StartHeight+IndexOffset-1
	MissedBlocksCounter uint64 `json:"missedBlocksCounter"`
}

func NewValidatorVoteInfo(startHeight, indexOffset, missedBlocksCounter uint64) ValidatorVoteInfo {
	return ValidatorVoteInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		MissedBlocksCounter: missedBlocksCounter,
	}
}
