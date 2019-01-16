package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

const (
	AddrLen = 20
)

type DelegatorEarningsStartInfo struct {
	PreviousPeriod       uint64        `json:"previous_period"`
	BondToken            uint64        `json:"bond_token"`
	StartingHeight       uint64        `json:"starting_height"`
	HistoricalRewardFees btypes.BigInt `json:"historical_rewards"`
}

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
