package mapper

import (
	"fmt"

	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qos/module/eco/types"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
)

//删除和validator相关的eco数据
//CONTRACT:
//delegator当前收益和收益发放信息数据不删除, 只是将bondTokens重置为0
//发放收益时,若delegator非validator的委托人, 或validator 不存在 则可以将delegator的收益相关数据删除
//发放收益时,validator的汇总数据可能会不存在
func (ecoMapper EcoMapper) RemoveValidator(valAddr btypes.Address, height uint64) error {

	distributionMapper := ecoMapper.DistributionMapper
	delegationMapper := ecoMapper.DelegationMapper
	validatorMapper := ecoMapper.ValidatorMapper
	voteInfoMapper := ecoMapper.VoteInfoMapper

	// 删除validator相关数据
	validator, ok := validatorMapper.KickValidator(valAddr)
	if !ok {
		return fmt.Errorf("validator:%s not exsits", valAddr)
	}

	//1. validator的汇总收益增加
	endPeriod := distributionMapper.incrementValidatorPeriod(validator)

	//2. 计算所有delegator的收益信息,并将delegator绑定的token置为0
	prefixKey := append(types.GetDelegatorEarningsStartInfoPrefixKey(), valAddr...)
	iter := store.KVStorePrefixIterator(distributionMapper.GetStore(), prefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorEarningsStartInfo
		distributionMapper.BaseMapper.DecodeObject(iter.Value(), &info)

		unbondToken := info.BondToken
		_, deleAddr := types.GetDelegatorEarningStartInfoAddr(iter.Key())
		rewards := distributionMapper.calculateRewardsBetweenPeriod(valAddr, info.PreviousPeriod, endPeriod, unbondToken)

		info.BondToken = uint64(0)
		info.StartingHeight = height
		info.PreviousPeriod = endPeriod
		info.HistoricalRewardFees = info.HistoricalRewardFees.Add(rewards)

		distributionMapper.Set(types.BuildDelegatorEarningStartInfoKey(valAddr, deleAddr), info)

		// 删除delegate数据,增加unbond数据
		delegationMapper.DelDelegationInfo(deleAddr, valAddr)

		//TODO: unbond height
		unbondHeight := uint64(10)
		delegationMapper.AddDelegatorUnbondingQOSatHeight(height+unbondHeight, deleAddr, unbondToken)
	}

	//删除validator汇总收益数据
	distributionMapper.DeleteValidatorPeriodSummaryInfo(valAddr)

	//删除validator 投票数据
	voteInfoMapper.DelValidatorVoteInfo(valAddr)
	voteInfoMapper.ClearValidatorVoteInfoInWindow(valAddr)

	return nil
}

type EcoMapper struct {
	DistributionMapper *DistributionMapper
	DelegationMapper   *DelegationMapper
	ValidatorMapper    *ValidatorMapper
	VoteInfoMapper     *VoteInfoMapper
}

func GetEcoMapper(ctx context.Context) EcoMapper {

	distributionMapper := GetDistributionMapper(ctx)
	delegationMapper := GetDelegationMapper(ctx)
	validatorMapper := GetValidatorMapper(ctx)
	voteInfoMapper := GetVoteInfoMapper(ctx)

	return EcoMapper{
		DistributionMapper: distributionMapper,
		DelegationMapper:   delegationMapper,
		ValidatorMapper:    validatorMapper,
		VoteInfoMapper:     voteInfoMapper,
	}
}
