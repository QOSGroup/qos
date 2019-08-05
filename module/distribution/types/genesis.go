package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	CommunityFeePool         btypes.BigInt                 `json:"community_fee_pool"`
	LastBlockProposer        btypes.Address                `json:"last_block_proposer"`
	PreDistributionQOSAmount btypes.BigInt                 `json:"pre_distribute_amount"`
	ValidatorHistoryPeriods  []ValidatorHistoryPeriodState `json:"validators_history_period"`
	ValidatorCurrentPeriods  []ValidatorCurrentPeriodState `json:"validators_current_period"`
	DelegatorEarningInfos    []DelegatorEarningStartState  `json:"delegator_earning_info"`
	DelegatorIncomeHeights   []DelegatorIncomeHeightState  `json:"delegator_income_height"`
	ValidatorEcoFeePools     []ValidatorEcoFeePoolState    `json:"validator_eco_fee_pools"`
	Params                   Params                        `json:"params"`
}

func NewGenesisState(communityFeePool btypes.BigInt,
	lastBlockProposer btypes.Address,
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
	ValidatorPubKey crypto.PubKey   `json:"validator_pubkey"`
	Period          uint64          `json:"period"`
	Summary         qtypes.Fraction `json:"summary"`
}

type ValidatorCurrentPeriodState struct {
	ValidatorPubKey      crypto.PubKey                 `json:"validator_pub_key"`
	CurrentPeriodSummary ValidatorCurrentPeriodSummary `json:"current_period_summary"`
}

type DelegatorEarningStartState struct {
	ValidatorPubKey            crypto.PubKey              `json:"validator_pub_key"`
	DeleAddress                btypes.Address             `json:"delegator_address"`
	DelegatorEarningsStartInfo DelegatorEarningsStartInfo `json:"earning_start_info"`
}

type DelegatorIncomeHeightState struct {
	ValidatorPubKey crypto.PubKey  `json:"validator_pub_key"`
	DeleAddress     btypes.Address `json:"delegator_address"`
	Height          uint64         `json:"height"`
}

type ValidatorEcoFeePoolState struct {
	ValidatorAddress btypes.Address      `json:"validator_address"`
	EcoFeePool       ValidatorEcoFeePool `json:"eco_fee_pool"`
}

func ValidateGenesis(_ GenesisState) error { return nil }