package types

import btypes "github.com/QOSGroup/qbase/types"

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
