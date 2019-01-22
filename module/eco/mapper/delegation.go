package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	stake_types "github.com/QOSGroup/qos/module/eco/types"
)

type DelegationMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*DelegationMapper)(nil)

func GetDelegationMapper(ctx context.Context) *DelegationMapper {
	return ctx.Mapper(stake_types.DelegationMapperName).(*DelegationMapper)
}

func (mapper *DelegationMapper) Copy() mapper.IMapper {
	delegationMapper := &DelegationMapper{}
	delegationMapper.BaseMapper = mapper.BaseMapper.Copy()
	return delegationMapper
}

func (mapper *DelegationMapper) SetDelegationInfo(info stake_types.DelegationInfo){
	mapper.Set(stake_types.BuildDelegationByDelValKey(info.DelegatorAddr, info.ValidatorAddr), info)
	mapper.Set(stake_types.BuildDelegationByValDelKey(info.ValidatorAddr, info.DelegatorAddr), nil)
}


func (mapper *DelegationMapper) GetDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) (info stake_types.DelegationInfo, exist bool) {
	exist = mapper.Get(stake_types.BuildDelegationByDelValKey(delAddr, valAddr), &info)
	return
}

func (mapper *DelegationMapper) DelDelegationInfo(delAddr btypes.Address, valAddr btypes.Address){
	mapper.Del(stake_types.BuildDelegationByDelValKey(delAddr, valAddr))
	mapper.Del(stake_types.BuildDelegationByValDelKey(valAddr, delAddr))
}

func (mapper *DelegationMapper) AddQOStoDelegationInfo(delAddr btypes.Address, valAddr btypes.Address, add_amount uint64) (amount uint64, exist bool) {
	delegationInfo, exist := mapper.GetDelegationInfo(delAddr, valAddr)
	if exist {
		delegationInfo.Amount += add_amount
		mapper.SetDelegationInfo(delegationInfo)
		amount = delegationInfo.Amount
	}
	// if the delegation doesn't exist, return (nil, false) and do nothing
	return
}

func (mapper *DelegationMapper) ReduceQOSfromDelegationInfo(delAddr btypes.Address, valAddr btypes.Address, reduce_amount uint64) (amount uint64, exist bool){
	delegationInfo, exist := mapper.GetDelegationInfo(delAddr, valAddr)
	if exist {
		if (delegationInfo.Amount > reduce_amount){
			delegationInfo.Amount -= reduce_amount
			mapper.SetDelegationInfo(delegationInfo)
			amount = delegationInfo.Amount
		}
		if (delegationInfo.Amount == reduce_amount){
			mapper.DelDelegationInfo(delAddr, valAddr)
			amount = 0
		}
	}
	// if the delegation doesn't exist, return (nil, false) and do nothing
	return
}

func (mapper *DelegationMapper) ChangeDelegationInfoCompound(delAddr btypes.Address, valAddr btypes.Address, isCompound bool) (exist bool){
	delegationInfo, exist := mapper.GetDelegationInfo(delAddr, valAddr)
	if exist {
		delegationInfo.IsCompound = isCompound
		mapper.SetDelegationInfo(delegationInfo)
	}
	return
}

func (mapper *DelegationMapper) SetDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address, amount uint64){
	mapper.Set(stake_types.BuildUnbondingDelegationByHeightDelKey(height, delAddr), amount)
}

func (mapper *DelegationMapper) GetDelegatorUnbondingQOSatHeight(height uint64, delAdd btypes.Address) (amount uint64, exist bool){
	exist = mapper.Get(stake_types.BuildUnbondingDelegationByHeightDelKey(height, delAdd), &amount)
	return
}

func (mapper *DelegationMapper) AddDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address, add_amount uint64){
	amount, exist := mapper.GetDelegatorUnbondingQOSatHeight(height, delAddr)
	if exist {
		add_amount += amount
	}
	mapper.SetDelegatorUnbondingQOSatHeight(height, delAddr, add_amount)
}

func (mapper *DelegationMapper) RemoveDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address) {
	mapper.Del(stake_types.BuildUnbondingDelegationByHeightDelKey(height, delAddr))
}

func (mapper *DelegationMapper) CreateDelegation(delAddr btypes.Address, valAddr btypes.Address, amount uint64, isCompound bool) (info stake_types.DelegationInfo){
	info = stake_types.NewDelegationInfo(delAddr, valAddr, amount, isCompound)
	//TODO：append to existing delegationInfo
	mapper.SetDelegationInfo(info)
	return
}