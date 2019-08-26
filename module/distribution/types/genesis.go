package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

type GenesisState struct {
	CommunityFeePool         btypes.BigInt                 `json:"community_fee_pool"`
	LastBlockProposer        btypes.ConsAddress                `json:"last_block_proposer"`
	PreDistributionQOSAmount btypes.BigInt                 `json:"pre_distribute_amount"`
	ValidatorHistoryPeriods  []ValidatorHistoryPeriodState `json:"validators_history_period"`
	ValidatorCurrentPeriods  []ValidatorCurrentPeriodState `json:"validators_current_period"`
	DelegatorEarningInfos    []DelegatorEarningStartState  `json:"delegator_earning_info"`
	DelegatorIncomeHeights   []DelegatorIncomeHeightState  `json:"delegator_income_height"`
	ValidatorEcoFeePools     []ValidatorEcoFeePoolState    `json:"validator_eco_fee_pools"`
	Params                   Params                        `json:"params"`
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
	ConsPubKey      string     `json:"consensus_pubkey"`
	Period          uint64          `json:"period"`
	Summary         qtypes.Fraction `json:"summary"`
}

type ValidatorCurrentPeriodState struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	ConsPubKey      string     `json:"consensus_pubkey"`
	CurrentPeriodSummary ValidatorCurrentPeriodSummary `json:"current_period_summary"`
}

type DelegatorEarningStartState struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	ConsPubKey      string     `json:"consensus_pubkey"`
	DeleAddress                btypes.AccAddress             `json:"delegator_address"`
	DelegatorEarningsStartInfo DelegatorEarningsStartInfo `json:"earning_start_info"`
}

type DelegatorIncomeHeightState struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	ConsPubKey      string     `json:"consensus_pubkey"`
	DeleAddress     btypes.AccAddress `json:"delegator_address"`
	Height          uint64         `json:"height"`
}

type ValidatorEcoFeePoolState struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	EcoFeePool       ValidatorEcoFeePool `json:"eco_fee_pool"`
}

func ValidateGenesis(_ GenesisState) error { return nil }
