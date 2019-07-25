package mapper

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func (mapper *Mapper) SetDelegationInfo(info types.DelegationInfo) {
	mapper.Set(types.BuildDelegationByDelValKey(info.DelegatorAddr, info.ValidatorAddr), info)
	mapper.Set(types.BuildDelegationByValDelKey(info.ValidatorAddr, info.DelegatorAddr), true)
}

func (mapper *Mapper) GetDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) (info types.DelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildDelegationByDelValKey(delAddr, valAddr), &info)
	return
}

func (mapper *Mapper) DelDelegationInfo(delAddr btypes.Address, valAddr btypes.Address) {
	mapper.Del(types.BuildDelegationByDelValKey(delAddr, valAddr))
	mapper.Del(types.BuildDelegationByValDelKey(valAddr, delAddr))
}

func (mapper *Mapper) IterateDelegationsValDeleAddr(valAddr btypes.Address, fn func(btypes.Address, btypes.Address)) {

	var prefixKey []byte

	if valAddr.Empty() {
		prefixKey = types.DelegationByValDelKey
	} else {
		prefixKey = append(types.DelegationByValDelKey, valAddr...)
	}

	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), prefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		_, deleAddr := types.GetDelegationValDelKeyAddress(k)
		fn(valAddr, deleAddr)
	}
}

//------------------------------genesisi export

func (mapper *Mapper) IterateDelegationsInfo(deleAddr btypes.Address, fn func(types.DelegationInfo)) {

	var prefixKey []byte

	if deleAddr.Empty() {
		prefixKey = types.DelegationByDelValKey
	} else {
		prefixKey = append(types.DelegationByDelValKey, deleAddr...)
	}

	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), prefixKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info types.DelegationInfo
		mapper.DecodeObject(iter.Value(), &info)
		fn(info)
	}
}

func (mapper *Mapper) Delegate(ctx context.Context, info types.DelegationInfo) {
	am := baseabci.GetAccountMapper(ctx)
	delegator := am.GetAccount(info.DelegatorAddr).(*qtypes.QOSAccount)
	delegator.MustMinusQOS(btypes.NewInt(int64(info.Amount)))
	am.SetAccount(delegator)

	delegation, exists := mapper.GetDelegationInfo(info.DelegatorAddr, info.ValidatorAddr)
	if !exists {
		mapper.SetDelegationInfo(info)
		mapper.AfterDelegationCreated(ctx, info.ValidatorAddr, info.DelegatorAddr)
	} else {
		delegation.Amount += info.Amount
		delegation.IsCompound = info.IsCompound
		mapper.BeforeDelegationModified(ctx, info.ValidatorAddr, info.DelegatorAddr, delegation.Amount, false)
		mapper.SetDelegationInfo(delegation)
	}

}

func (mapper *Mapper) UnbondTokens(ctx context.Context, info types.DelegationInfo, tokens uint64) {
	info.Amount = info.Amount - tokens
	mapper.BeforeDelegationModified(ctx, info.ValidatorAddr, info.DelegatorAddr, info.Amount, false)
	mapper.SetDelegationInfo(info)
}

func (mapper *Mapper) ReDelegate(ctx context.Context, delegation types.DelegationInfo, info types.ReDelegateInfo) {
	// update origin delegation
	delegation.Amount -= info.Amount
	mapper.BeforeDelegationModified(ctx, delegation.ValidatorAddr, delegation.DelegatorAddr, delegation.Amount, true)
	mapper.SetDelegationInfo(delegation)

	// save new delegation
	reDelegation, exists := mapper.GetDelegationInfo(info.DelegatorAddr, info.ToValidator)
	if !exists {
		mapper.SetDelegationInfo(types.NewDelegationInfo(info.DelegatorAddr, info.ToValidator, info.Amount, info.IsCompound))
		mapper.AfterDelegationCreated(ctx, info.ToValidator, info.DelegatorAddr)
	} else {
		reDelegation.Amount += info.Amount
		reDelegation.IsCompound = info.IsCompound
		mapper.BeforeDelegationModified(ctx, reDelegation.ValidatorAddr, reDelegation.DelegatorAddr, reDelegation.Amount, false)
		mapper.SetDelegationInfo(reDelegation)
	}
}
