package types

import btypes "github.com/QOSGroup/qbase/types"

//DelegatorEarningsStartInfo delegator计算收益信息
type DelegatorEarningsStartInfo struct {
	PreviousPeriod        int64         `json:"previous_period"`
	BondToken             btypes.BigInt `json:"bond_token"`
	CurrentStartingHeight int64         `json:"earns_starting_height"`
	FirstDelegateHeight   int64         `json:"first_delegate_height"`
	HistoricalRewardFees  btypes.BigInt `json:"historical_rewards"`
	LastIncomeCalHeight   int64         `json:"last_income_calHeight"`
	LastIncomeCalFees     btypes.BigInt `json:"last_income_calFees"`
}
