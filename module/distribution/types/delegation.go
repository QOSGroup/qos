package types

import btypes "github.com/QOSGroup/qbase/types"

//DelegatorEarningsStartInfo delegator计算收益信息
type DelegatorEarningsStartInfo struct {
	PreviousPeriod        int64         `json:"previous_period"`       // 前收益计算点
	BondToken             btypes.BigInt `json:"bond_token"`            // 绑定tokens
	CurrentStartingHeight int64         `json:"earns_starting_height"` // 当前计算周期起始高度
	FirstDelegateHeight   int64         `json:"first_delegate_height"` // 首次委托高度
	HistoricalRewardFees  btypes.BigInt `json:"historical_rewards"`    // 累计未发放奖励
	LastIncomeCalHeight   int64         `json:"last_income_calHeight"` // 最后收益计算高度
	LastIncomeCalFees     btypes.BigInt `json:"last_income_calFees"`   // 最后一次发放收益
}
