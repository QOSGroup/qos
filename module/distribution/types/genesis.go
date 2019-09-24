package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

type GenesisState struct {
	CommunityFeePool         btypes.BigInt                 `json:"community_fee_pool"`        // 社区费池
	LastBlockProposer        btypes.ConsAddress            `json:"last_block_proposer"`       // 最新区块提议验证节点共识地址
	PreDistributionQOSAmount btypes.BigInt                 `json:"pre_distribute_amount"`     // 代分发奖励
	ValidatorHistoryPeriods  []ValidatorHistoryPeriodState `json:"validators_history_period"` // 验证节点收益历史节点信息
	ValidatorCurrentPeriods  []ValidatorCurrentPeriodState `json:"validators_current_period"` // 验证节点当前收益信息
	DelegatorEarningInfos    []DelegatorEarningStartState  `json:"delegator_earning_info"`    // 委托收益信息
	DelegatorIncomeHeights   []DelegatorIncomeHeightState  `json:"delegator_income_height"`   // 委托收益发放高度
	ValidatorEcoFeePools     []ValidatorEcoFeePoolState    `json:"validator_eco_fee_pools"`   // 验证节点委托共享费池
	Params                   Params                        `json:"params"`                    // 参数
}

func NewGenesisState(communityFeePool btypes.BigInt,
	lastBlockProposer btypes.ConsAddress,
	preDistributionQOSAmount btypes.BigInt,
	validatorHistoryPeriods []ValidatorHistoryPeriodState,
	validatorCurrentPeriods []ValidatorCurrentPeriodState,
	delegatorEarningInfos []DelegatorEarningStartState,
	delegatorIncomeHeights []DelegatorIncomeHeightState,
	validatorEcoFeePools []ValidatorEcoFeePoolState,
	params Params) GenesisState {
	return GenesisState{
		CommunityFeePool:         communityFeePool,
		LastBlockProposer:        lastBlockProposer,
		PreDistributionQOSAmount: preDistributionQOSAmount,
		ValidatorHistoryPeriods:  validatorHistoryPeriods,
		ValidatorCurrentPeriods:  validatorCurrentPeriods,
		DelegatorEarningInfos:    delegatorEarningInfos,
		DelegatorIncomeHeights:   delegatorIncomeHeights,
		ValidatorEcoFeePools:     validatorEcoFeePools,
		Params:                   params,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		CommunityFeePool:         btypes.ZeroInt(),
		PreDistributionQOSAmount: btypes.ZeroInt(),
		Params:                   DefaultParams(),
	}
}

type ValidatorHistoryPeriodState struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	ConsPubKey      string            `json:"consensus_pubkey"`
	Period          int64             `json:"period"`
	Summary         qtypes.Fraction   `json:"summary"`
}

type ValidatorCurrentPeriodState struct {
	OperatorAddress      btypes.ValAddress             `json:"validator_address"`
	ConsPubKey           string                        `json:"consensus_pubkey"`
	CurrentPeriodSummary ValidatorCurrentPeriodSummary `json:"current_period_summary"`
}

type DelegatorEarningStartState struct {
	OperatorAddress            btypes.ValAddress          `json:"validator_address"`
	ConsPubKey                 string                     `json:"consensus_pubkey"`
	DeleAddress                btypes.AccAddress          `json:"delegator_address"`
	DelegatorEarningsStartInfo DelegatorEarningsStartInfo `json:"earning_start_info"`
}

type DelegatorIncomeHeightState struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	ConsPubKey      string            `json:"consensus_pubkey"`
	DeleAddress     btypes.AccAddress `json:"delegator_address"`
	Height          int64             `json:"height"`
}

type ValidatorEcoFeePoolState struct {
	OperatorAddress btypes.ValAddress   `json:"validator_address"`
	EcoFeePool      ValidatorEcoFeePool `json:"eco_fee_pool"`
}

func ValidateGenesis(data GenesisState) error {

	// validate params
	err := data.Params.Validate()
	if err != nil {
		return err
	}

	return nil
}
