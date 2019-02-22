package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
)

type DelegationMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*DelegationMapper)(nil)

func GetDelegationMapper(ctx context.Context) *DelegationMapper {
	return ctx.Mapper(ecotypes.DelegationMapperName).(*DelegationMapper)
}

func NewDelegationMapper() *DelegationMapper {
	var delegationMapper = DelegationMapper{}
	delegationMapper.BaseMapper = mapper.NewBaseMapper(nil, ecotypes.DelegationMapperName)
	return &delegationMapper
}

func (mapper *DelegationMapper) Copy() mapper.IMapper {
	delegationMapper := &DelegationMapper{}
	delegationMapper.BaseMapper = mapper.BaseMapper.Copy()
	return delegationMapper
}

func (mapper *DelegationMapper) SetDelegationInfo(info ecotypes.DelegationInfo) {
	mapper.Set(ecotypes.BuildDelegationByDelValKey(info.DelegatorAddr, info.ValidatorAddr), info)
	mapper.Set(ecotypes.BuildDelegationByValDelKey(info.ValidatorAddr, info.DelegatorAddr), true)
}

func (mapper *DelegationMapper) GetDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) (info ecotypes.DelegationInfo, exist bool) {
	exist = mapper.Get(ecotypes.BuildDelegationByDelValKey(delAddr, valAddr), &info)
	return
}

func (mapper *DelegationMapper) DelDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) {
	mapper.Del(ecotypes.BuildDelegationByDelValKey(delAddr, valAddr))
	mapper.Del(ecotypes.BuildDelegationByValDelKey(valAddr, delAddr))
}

func (mapper *DelegationMapper) SetDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address, amount uint64) {
	mapper.Set(ecotypes.BuildUnbondingDelegationByHeightDelKey(height, delAddr), amount)
}

func (mapper *DelegationMapper) GetDelegatorUnbondingQOSatHeight(height uint64, delAdd btypes.Address) (amount uint64, exist bool) {
	exist = mapper.Get(ecotypes.BuildUnbondingDelegationByHeightDelKey(height, delAdd), &amount)
	return
}

func (mapper *DelegationMapper) AddDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address, add_amount uint64) {
	amount, exist := mapper.GetDelegatorUnbondingQOSatHeight(height, delAddr)
	if exist {
		add_amount += amount
	}
	mapper.SetDelegatorUnbondingQOSatHeight(height, delAddr, add_amount)
}

func (mapper *DelegationMapper) RemoveDelegatorUnbondingQOSatHeight(height uint64, delAddr btypes.Address) {
	mapper.Del(ecotypes.BuildUnbondingDelegationByHeightDelKey(height, delAddr))
}

func (mapper *DelegationMapper) IterateDelegationsValDeleAddr(valAddr btypes.Address, fn func(btypes.Address, btypes.Address)) {

	var prefixKey []byte

	if valAddr.Empty() {
		prefixKey = ecotypes.DelegationByValDelKey
	} else {
		prefixKey = append(ecotypes.DelegationByValDelKey, valAddr...)
	}

	iter := store.KVStorePrefixIterator(mapper.GetStore(), prefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		_, deleAddr := ecotypes.GetDelegationValDelKeyAddress(k)
		fn(valAddr, deleAddr)
	}
}

//------------------------------genesisi export

func (mapper *DelegationMapper) IterateDelegationsInfo(deleAddr btypes.Address, fn func(ecotypes.DelegationInfo)) {

	var prefixKey []byte

	if deleAddr.Empty() {
		prefixKey = ecotypes.DelegationByDelValKey
	} else {
		prefixKey = append(ecotypes.DelegationByDelValKey, deleAddr...)
	}

	iter := store.KVStorePrefixIterator(mapper.GetStore(), prefixKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info ecotypes.DelegationInfo
		mapper.DecodeObject(iter.Value(), &info)
		fn(info)
	}
}

func (mapper *DelegationMapper) IterateDelegationsUnbondInfo(fn func(btypes.Address, uint64, uint64)) {
	iter := store.KVStorePrefixIterator(mapper.GetStore(), ecotypes.DelegatorUnbondingQOSatHeightKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		height, deleAddr := ecotypes.GetUnbondingDelegationHeightAddress(key)
		var amount uint64
		mapper.DecodeObject(iter.Value(), &amount)
		fn(deleAddr, height, amount)
	}
}
