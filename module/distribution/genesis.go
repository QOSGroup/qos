package distribution

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/mapper"
	"github.com/QOSGroup/qos/module/eco/types"
	qtypes "github.com/QOSGroup/qos/types"
)

type GenesisState struct {
	CommunityFeePool         btypes.BigInt                 `json:"community_fee_pool"`
	LastBlockProposer        btypes.Address                `json:"last_block_proposer"`
	PreDistributionQOSAmount btypes.BigInt                 `json:"pre_distribute_amount"`
	ValidatorHistoryPeriods  []ValidatorHistoryPeriodState `json:"validators_history_period"`
	ValidatorCurrentPeriods  []ValidatorCurrentPeriodState `json:"validators_current_period"`
	DelegatorEarningInfos    []DelegatorEarningStartState  `json:"delegators_earning_info"`
	DelegatorIncomeHeights   []DelegatorIncomeHeightState  `json:"delegators_income_height"`
	Params                   types.DistributionParams      `json:"params"`
}

func NewGenesisState(communityFeePool btypes.BigInt,
	lastBlockProposer btypes.Address,
	preDistributionQOSAmount btypes.BigInt,
	validatorHistoryPeriods []ValidatorHistoryPeriodState,
	validatorCurrentPeriods []ValidatorCurrentPeriodState,
	delegatorEarningInfos []DelegatorEarningStartState,
	delegatorIncomeHeights []DelegatorIncomeHeightState,
	params types.DistributionParams) GenesisState {
	return GenesisState{
		CommunityFeePool:         communityFeePool,
		LastBlockProposer:        lastBlockProposer,
		PreDistributionQOSAmount: preDistributionQOSAmount,
		ValidatorHistoryPeriods:  validatorHistoryPeriods,
		ValidatorCurrentPeriods:  validatorCurrentPeriods,
		DelegatorEarningInfos:    delegatorEarningInfos,
		DelegatorIncomeHeights:   delegatorIncomeHeights,
		Params:                   params,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		CommunityFeePool:         btypes.ZeroInt(),
		PreDistributionQOSAmount: btypes.ZeroInt(),
		Params:                   types.DefaultDistributionParams(),
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {

	distributionMapper := mapper.GetDistributionMapper(ctx)

	feePool := data.CommunityFeePool
	distributionMapper.SetCommunityFeePool(feePool.NilToZero())

	proposer := data.LastBlockProposer
	if !proposer.Empty() {
		distributionMapper.SetLastBlockProposer(proposer)
	}

	distributionMapper.SetPreDistributionQOS(data.PreDistributionQOSAmount.NilToZero())
	distributionMapper.SetParams(data.Params)

	for _, validatorHistoryPeriodState := range data.ValidatorHistoryPeriods {
		key := types.BuildValidatorHistoryPeriodSummaryKey(validatorHistoryPeriodState.ValAddress, validatorHistoryPeriodState.Period)
		distributionMapper.Set(key, validatorHistoryPeriodState.Summary)
	}

	for _, validatorCurrentPeriodState := range data.ValidatorCurrentPeriods {
		key := types.BuildValidatorCurrentPeriodSummaryKey(validatorCurrentPeriodState.ValAddress)
		distributionMapper.Set(key, validatorCurrentPeriodState.CurrentPeriodSummary)
	}

	for _, delegatorEarningInfoState := range data.DelegatorEarningInfos {
		key := types.BuildDelegatorEarningStartInfoKey(delegatorEarningInfoState.ValAddress, delegatorEarningInfoState.DeleAddress)
		distributionMapper.Set(key, delegatorEarningInfoState.DelegatorEarningsStartInfo)
	}

	for _, delegatorIncomeHeightState := range data.DelegatorIncomeHeights {
		key := types.BuildDelegatorPeriodIncomeKey(delegatorIncomeHeightState.ValAddress, delegatorIncomeHeightState.DeleAddress, delegatorIncomeHeightState.Height)
		distributionMapper.Set(key, true)
	}

}

func ExportGenesis(ctx context.Context) GenesisState {

	distributionMapper := mapper.GetDistributionMapper(ctx)

	feePool := distributionMapper.GetCommunityFeePool()
	lastBlockProposer := distributionMapper.GetLastBlockProposer()
	preDistributionQOS := distributionMapper.GetPreDistributionQOS()
	params := distributionMapper.GetParams()

	var validatorHistoryPeriods []ValidatorHistoryPeriodState
	distributionMapper.IteratorValidatorsHistoryPeriod(func(valAddr btypes.Address, period uint64, frac qtypes.Fraction) {
		vhps := ValidatorHistoryPeriodState{
			ValAddress: valAddr,
			Period:     period,
			Summary:    frac,
		}
		validatorHistoryPeriods = append(validatorHistoryPeriods, vhps)
	})

	var validatorCurrentPeriods []ValidatorCurrentPeriodState
	distributionMapper.IteratorValidatorsCurrentPeriod(func(valAddr btypes.Address, vcps types.ValidatorCurrentPeriodSummary) {
		vcpsState := ValidatorCurrentPeriodState{
			ValAddress:           valAddr,
			CurrentPeriodSummary: vcps,
		}
		validatorCurrentPeriods = append(validatorCurrentPeriods, vcpsState)
	})

	var delegatorEarningInfos []DelegatorEarningStartState
	distributionMapper.IteratorDelegatorsEarningStartInfo(func(valAddr btypes.Address, deleAddr btypes.Address, desi types.DelegatorEarningsStartInfo) {
		dess := DelegatorEarningStartState{
			ValAddress:                 valAddr,
			DeleAddress:                deleAddr,
			DelegatorEarningsStartInfo: desi,
		}
		delegatorEarningInfos = append(delegatorEarningInfos, dess)
	})

	var delegatorIncomeHeights []DelegatorIncomeHeightState
	distributionMapper.IteratorDelegatorsIncomeHeight(func(valAddr btypes.Address, deleAddr btypes.Address, height uint64) {
		dihs := DelegatorIncomeHeightState{
			ValAddress:  valAddr,
			DeleAddress: deleAddr,
			Height:      height,
		}
		delegatorIncomeHeights = append(delegatorIncomeHeights, dihs)
	})

	return NewGenesisState(feePool,
		lastBlockProposer,
		preDistributionQOS,
		validatorHistoryPeriods,
		validatorCurrentPeriods,
		delegatorEarningInfos,
		delegatorIncomeHeights,
		params,
	)
}

type ValidatorHistoryPeriodState struct {
	ValAddress btypes.Address  `json:"validator_address"`
	Period     uint64          `json:"period"`
	Summary    qtypes.Fraction `json:"summary"`
}

type ValidatorCurrentPeriodState struct {
	ValAddress           btypes.Address                      `json:"validator_address"`
	CurrentPeriodSummary types.ValidatorCurrentPeriodSummary `json:"current_period_summary"`
}

type DelegatorEarningStartState struct {
	ValAddress                 btypes.Address                   `json:"validator_address"`
	DeleAddress                btypes.Address                   `json:"delegator_address"`
	DelegatorEarningsStartInfo types.DelegatorEarningsStartInfo `json:"earning_start_info"`
}

type DelegatorIncomeHeightState struct {
	ValAddress  btypes.Address `json:"validator_address"`
	DeleAddress btypes.Address `json:"delegator_address"`
	Height      uint64         `json:"height"`
}
