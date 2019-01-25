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

func NewDelegationMapper() *DelegationMapper {
	var delegationMapper = DelegationMapper{}
	delegationMapper.BaseMapper = mapper.NewBaseMapper(nil, stake_types.DelegationMapperName)
	return &delegationMapper
}

func (mapper *DelegationMapper) Copy() mapper.IMapper {
	delegationMapper := &DelegationMapper{}
	delegationMapper.BaseMapper = mapper.BaseMapper.Copy()
	return delegationMapper
}

func (mapper *DelegationMapper) SetDelegationInfo(info stake_types.DelegationInfo) {
	mapper.Set(stake_types.BuildDelegationByDelValKey(info.DelegatorAddr, info.ValidatorAddr), info)
	mapper.Set(stake_types.BuildDelegationByValDelKey(info.ValidatorAddr, info.DelegatorAddr), true)
}

func (mapper *DelegationMapper) GetDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) (info stake_types.DelegationInfo, exist bool) {
	exist = mapper.Get(stake_types.BuildDelegationByDelValKey(delAddr, valAddr), &info)
	return
}

func (mapper *DelegationMapper) DelDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) {
	mapper.Del(stake_types.BuildDelegationByDelValKey(delAddr, valAddr))
	mapper.Del(stake_types.BuildDelegationByValDelKey(valAddr, delAddr))
}

func (mapper *DelegationMapper) setDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address, amount uint64) {
	mapper.Set(stake_types.BuildUnbondingDelegationByHeightDelKey(height, delAddr), amount)
}

func (mapper *DelegationMapper) GetDelegatorUnbondingQOSatHeight(height uint64, delAdd btypes.Address) (amount uint64, exist bool) {
	exist = mapper.Get(stake_types.BuildUnbondingDelegationByHeightDelKey(height, delAdd), &amount)
	return
}

func (mapper *DelegationMapper) AddDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address, add_amount uint64) {
	amount, exist := mapper.GetDelegatorUnbondingQOSatHeight(height, delAddr)
	if exist {
		add_amount += amount
	}
	mapper.setDelegatorUnbondingQOSatHeight(height, delAddr, add_amount)
}

func (mapper *DelegationMapper) RemoveDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address) {
	mapper.Del(stake_types.BuildUnbondingDelegationByHeightDelKey(height, delAddr))
}
