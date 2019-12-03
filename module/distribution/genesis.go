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

	// 初始化社区费池
	feePool := data.CommunityFeePool
	distributionMapper.SetCommunityFeePool(feePool.NilToZero())

	// 最新块提议验证节点共识地址
	proposer := data.LastBlockProposer
	if !proposer.Empty() {
		distributionMapper.SetLastBlockProposer(proposer)
	}

	// 待分配奖励
	distributionMapper.SetPreDistributionQOS(data.PreDistributionQOSAmount.NilToZero())

	// 初始化参数
	if err := data.Params.Validate(); err == nil {
		distributionMapper.SetParams(ctx, data.Params)
	} else {
		panic(err)
	}

	// 验证节点历史节点数据
	for _, validatorHistoryPeriodState := range data.ValidatorHistoryPeriods {
		key := types.BuildValidatorHistoryPeriodSummaryKey(validatorHistoryPeriodState.OperatorAddress, validatorHistoryPeriodState.Period)
		distributionMapper.Set(key, validatorHistoryPeriodState.Summary)
	}

	// 验证节点当前收益信息
	for _, validatorCurrentPeriodState := range data.ValidatorCurrentPeriods {
		key := types.BuildValidatorCurrentPeriodSummaryKey(validatorCurrentPeriodState.OperatorAddress)
		distributionMapper.Set(key, validatorCurrentPeriodState.CurrentPeriodSummary)
	}

	// 委托收益信息
	for _, delegatorEarningInfoState := range data.DelegatorEarningInfos {
		key := types.BuildDelegatorEarningStartInfoKey(delegatorEarningInfoState.OperatorAddress, delegatorEarningInfoState.DeleAddress)
		distributionMapper.Set(key, delegatorEarningInfoState.DelegatorEarningsStartInfo)
	}

	// 委托收益发放高度信息
	for _, delegatorIncomeHeightState := range data.DelegatorIncomeHeights {
		key := types.BuildDelegatorPeriodIncomeKey(delegatorIncomeHeightState.OperatorAddress, delegatorIncomeHeightState.DeleAddress, delegatorIncomeHeightState.Height)
		distributionMapper.Set(key, true)
	}

	// 验证节点委托共享费池
	for _, validatorFeePoolState := range data.ValidatorEcoFeePools {
		key := types.BuildValidatorEcoFeePoolKey(validatorFeePoolState.OperatorAddress)
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
	distributionMapper.IteratorValidatorsHistoryPeriod(func(valAddr btypes.ValAddress, period int64, frac qtypes.Fraction) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			vhps := types.ValidatorHistoryPeriodState{
				OperatorAddress: validator.OperatorAddress,
				ConsPubKey:      btypes.MustConsensusPubKeyString(validator.ConsPubKey),
				Period:          period,
				Summary:         frac,
			}
			validatorHistoryPeriods = append(validatorHistoryPeriods, vhps)
		}
	})

	var validatorCurrentPeriods []types.ValidatorCurrentPeriodState
	distributionMapper.IteratorValidatorsCurrentPeriod(func(valAddr btypes.ValAddress, vcps types.ValidatorCurrentPeriodSummary) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			vcpsState := types.ValidatorCurrentPeriodState{
				OperatorAddress:      validator.OperatorAddress,
				ConsPubKey:           btypes.MustConsensusPubKeyString(validator.ConsPubKey),
				CurrentPeriodSummary: vcps,
			}
			validatorCurrentPeriods = append(validatorCurrentPeriods, vcpsState)
		}
	})

	var delegatorEarningInfos []types.DelegatorEarningStartState
	distributionMapper.IteratorDelegatorEarningStartInfo(func(valAddr btypes.ValAddress, deleAddr btypes.AccAddress, desi types.DelegatorEarningsStartInfo) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			dess := types.DelegatorEarningStartState{
				OperatorAddress:            validator.OperatorAddress,
				ConsPubKey:                 btypes.MustConsensusPubKeyString(validator.ConsPubKey),
				DeleAddress:                deleAddr,
				DelegatorEarningsStartInfo: desi,
			}
			delegatorEarningInfos = append(delegatorEarningInfos, dess)
		}
	})

	var delegatorIncomeHeights []types.DelegatorIncomeHeightState
	distributionMapper.IteratorDelegatorsIncomeHeight(func(valAddr btypes.ValAddress, deleAddr btypes.AccAddress, height int64) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			dihs := types.DelegatorIncomeHeightState{
				OperatorAddress: validator.OperatorAddress,
				ConsPubKey:      btypes.MustConsensusPubKeyString(validator.ConsPubKey),
				DeleAddress:     deleAddr,
				Height:          height,
			}
			delegatorIncomeHeights = append(delegatorIncomeHeights, dihs)
		}
	})

	var validatorEcoFeePools []types.ValidatorEcoFeePoolState
	distributionMapper.IteratorValidatorEcoFeePools(func(validatorAddr btypes.ValAddress, pool types.ValidatorEcoFeePool) {
		validatorEcoFeePools = append(validatorEcoFeePools, types.ValidatorEcoFeePoolState{
			OperatorAddress: validatorAddr,
			EcoFeePool:      pool,
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
