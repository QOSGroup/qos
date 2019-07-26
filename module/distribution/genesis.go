package distribution

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/mapper"
	"github.com/QOSGroup/qos/module/distribution/types"
	"github.com/QOSGroup/qos/module/stake"
	qtypes "github.com/QOSGroup/qos/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {

	distributionMapper := mapper.GetMapper(ctx)

	feePool := data.CommunityFeePool
	distributionMapper.SetCommunityFeePool(feePool.NilToZero())

	proposer := data.LastBlockProposer
	if !proposer.Empty() {
		distributionMapper.SetLastBlockProposer(proposer)
	}

	distributionMapper.SetPreDistributionQOS(data.PreDistributionQOSAmount.NilToZero())
	distributionMapper.SetParams(ctx, data.Params)

	for _, validatorHistoryPeriodState := range data.ValidatorHistoryPeriods {
		key := types.BuildValidatorHistoryPeriodSummaryKey(btypes.Address(validatorHistoryPeriodState.ValidatorPubKey.Address()), validatorHistoryPeriodState.Period)
		distributionMapper.Set(key, validatorHistoryPeriodState.Summary)
	}

	for _, validatorCurrentPeriodState := range data.ValidatorCurrentPeriods {
		key := types.BuildValidatorCurrentPeriodSummaryKey(btypes.Address(validatorCurrentPeriodState.ValidatorPubKey.Address()))
		distributionMapper.Set(key, validatorCurrentPeriodState.CurrentPeriodSummary)
	}

	for _, delegatorEarningInfoState := range data.DelegatorEarningInfos {
		key := types.BuildDelegatorEarningStartInfoKey(btypes.Address(delegatorEarningInfoState.ValidatorPubKey.Address()), delegatorEarningInfoState.DeleAddress)
		distributionMapper.Set(key, delegatorEarningInfoState.DelegatorEarningsStartInfo)
	}

	for _, delegatorIncomeHeightState := range data.DelegatorIncomeHeights {
		key := types.BuildDelegatorPeriodIncomeKey(btypes.Address(delegatorIncomeHeightState.ValidatorPubKey.Address()), delegatorIncomeHeightState.DeleAddress, delegatorIncomeHeightState.Height)
		distributionMapper.Set(key, true)
	}

	for _, validatorFeePoolState := range data.ValidatorEcoFeePools {
		key := types.BuildValidatorEcoFeePoolKey(validatorFeePoolState.ValidatorAddress)
		distributionMapper.Set(key, validatorFeePoolState.EcoFeePool)
	}
}

func ExportGenesis(ctx context.Context) types.GenesisState {

	distributionMapper := mapper.GetMapper(ctx)
	validatorMapper := stake.GetMapper(ctx)

	feePool := distributionMapper.GetCommunityFeePool()
	lastBlockProposer := distributionMapper.GetLastBlockProposer()
	preDistributionQOS := distributionMapper.GetPreDistributionQOS()
	params := distributionMapper.GetParams(ctx)

	var validatorHistoryPeriods []types.ValidatorHistoryPeriodState
	distributionMapper.IteratorValidatorsHistoryPeriod(func(valAddr btypes.Address, period uint64, frac qtypes.Fraction) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			vhps := types.ValidatorHistoryPeriodState{
				ValidatorPubKey: validator.ValidatorPubKey,
				Period:          period,
				Summary:         frac,
			}
			validatorHistoryPeriods = append(validatorHistoryPeriods, vhps)
		}
	})

	var validatorCurrentPeriods []types.ValidatorCurrentPeriodState
	distributionMapper.IteratorValidatorsCurrentPeriod(func(valAddr btypes.Address, vcps types.ValidatorCurrentPeriodSummary) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			vcpsState := types.ValidatorCurrentPeriodState{
				ValidatorPubKey:      validator.ValidatorPubKey,
				CurrentPeriodSummary: vcps,
			}
			validatorCurrentPeriods = append(validatorCurrentPeriods, vcpsState)
		}
	})

	var delegatorEarningInfos []types.DelegatorEarningStartState
	distributionMapper.IteratorDelegatorEarningStartInfo(func(valAddr btypes.Address, deleAddr btypes.Address, desi types.DelegatorEarningsStartInfo) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			dess := types.DelegatorEarningStartState{
				ValidatorPubKey:            validator.ValidatorPubKey,
				DeleAddress:                deleAddr,
				DelegatorEarningsStartInfo: desi,
			}
			delegatorEarningInfos = append(delegatorEarningInfos, dess)
		}
	})

	var delegatorIncomeHeights []types.DelegatorIncomeHeightState
	distributionMapper.IteratorDelegatorsIncomeHeight(func(valAddr btypes.Address, deleAddr btypes.Address, height uint64) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			dihs := types.DelegatorIncomeHeightState{
				ValidatorPubKey: validator.ValidatorPubKey,
				DeleAddress:     deleAddr,
				Height:          height,
			}
			delegatorIncomeHeights = append(delegatorIncomeHeights, dihs)
		}
	})

	var validatorEcoFeePools []types.ValidatorEcoFeePoolState
	distributionMapper.IteratorValidatorEcoFeePools(func(validatorAddr btypes.Address, pool types.ValidatorEcoFeePool) {
		validatorEcoFeePools = append(validatorEcoFeePools, types.ValidatorEcoFeePoolState{
			ValidatorAddress: validatorAddr,
			EcoFeePool:       pool,
		})
	})

	return types.NewGenesisState(feePool,
		lastBlockProposer,
		preDistributionQOS,
		validatorHistoryPeriods,
		validatorCurrentPeriods,
		delegatorEarningInfos,
		delegatorIncomeHeights,
		validatorEcoFeePools,
		params,
	)
	return types.GenesisState{}
}
