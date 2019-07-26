package mapper

import (
	"bytes"
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
		mapper.BeforeDelegationModified(ctx, info.ValidatorAddr, info.DelegatorAddr, delegation.Amount)
		mapper.SetDelegationInfo(delegation)
	}

}

func (mapper *Mapper) UnbondTokens(ctx context.Context, info types.DelegationInfo, tokens uint64) {
	info.Amount = info.Amount - tokens
	mapper.BeforeDelegationModified(ctx, info.ValidatorAddr, info.DelegatorAddr, info.Amount)
	unbondHeight := uint64(mapper.GetParams(ctx).DelegatorUnbondReturnHeight) + uint64(ctx.BlockHeight())
	mapper.AddUnbondingDelegation(unbondHeight, types.NewUnbondingDelegationInfo(info.DelegatorAddr, info.ValidatorAddr, uint64(ctx.BlockHeight()), tokens))
	mapper.SetDelegationInfo(info)
}

func (mapper *Mapper) ReDelegate(ctx context.Context, delegation types.DelegationInfo, info types.ReDelegationInfo) {
	// update origin delegation
	delegation.Amount -= info.Amount
	mapper.BeforeDelegationModified(ctx, delegation.ValidatorAddr, delegation.DelegatorAddr, delegation.Amount)
	mapper.SetDelegationInfo(delegation)

	// save new delegation
	reDelegation, exists := mapper.GetDelegationInfo(info.DelegatorAddr, info.ToValidator)
	if !exists {
		mapper.SetDelegationInfo(types.NewDelegationInfo(info.DelegatorAddr, info.ToValidator, info.Amount, info.IsCompound))
		mapper.AfterDelegationCreated(ctx, info.ToValidator, info.DelegatorAddr)
	} else {
		reDelegation.Amount += info.Amount
		reDelegation.IsCompound = info.IsCompound
		mapper.BeforeDelegationModified(ctx, reDelegation.ValidatorAddr, reDelegation.DelegatorAddr, reDelegation.Amount)
		mapper.SetDelegationInfo(reDelegation)
	}
}

func (mapper *Mapper) IterateUnbondingDelegations(fn func(btypes.Address, uint64, []types.UnbondingDelegationInfo)) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.UnbondingHeightDelegatorKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		height, delAddr := types.GetUnbondingDelegationHeightAddress(key)
		var unbondings []types.UnbondingDelegationInfo
		mapper.DecodeObject(iter.Value(), &unbondings)
		fn(delAddr, height, unbondings)
	}
}

func (mapper *Mapper) setUnbondingDelegations(height uint64, delAddr btypes.Address, unbondings []types.UnbondingDelegationInfo) {
	mapper.Set(types.BuildUnbondingHeightDelegatorKey(height, delAddr), unbondings)
	mapper.Set(types.BuildUnbondingDelegatorHeightKey(delAddr, height), true)
}

func (mapper *Mapper) getUnbondingDelegations(height uint64, delAdd btypes.Address) (unbondings []types.UnbondingDelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildUnbondingHeightDelegatorKey(height, delAdd), &unbondings)
	return
}

func (mapper *Mapper) AddUnbondingDelegation(height uint64, unbonding types.UnbondingDelegationInfo) {
	unbondings := []types.UnbondingDelegationInfo{}
	origins, exist := mapper.getUnbondingDelegations(height, unbonding.DelegatorAddr)
	if exist {
		added := false
		for _, ub := range origins {
			if bytes.Equal(ub.ValidatorAddr, unbonding.ValidatorAddr) {
				ub.Amount += unbonding.Amount
				added = true
			}
			unbondings = append(unbondings, ub)
		}
		if !added {
			unbondings = append(unbondings, unbonding)
		}
	} else {
		unbondings = append(unbondings, unbonding)
	}
	mapper.setUnbondingDelegations(height, unbonding.DelegatorAddr, unbondings)
}

func (mapper *Mapper) AddUnbondingDelegations(height uint64, unbondingsAdd []types.UnbondingDelegationInfo) {
	for _, unbonding := range unbondingsAdd {
		mapper.AddUnbondingDelegation(height, unbonding)
	}
}

func (mapper *Mapper) RemoveUnbondingDelegations(height uint64, delAddr btypes.Address) {
	mapper.Del(types.BuildUnbondingHeightDelegatorKey(height, delAddr))
	mapper.Del(types.BuildUnbondingDelegatorHeightKey(delAddr, height))
}

func (mapper *Mapper) IterateRedelegationsInfo(fn func(btypes.Address, uint64, uint64)) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.RedelegationHeightDelegatorKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		height, delAddr := types.GetRedelegationHeightAddress(key)
		var amount uint64
		mapper.DecodeObject(iter.Value(), &amount)
		fn(delAddr, height, amount)
	}
}

func (mapper *Mapper) setRedelegationsAtHeight(height uint64, delAddr btypes.Address, amount uint64) {
	mapper.Set(types.BuildRedelegationHeightDelegatorKey(height, delAddr), amount)
}

func (mapper *Mapper) getRedelegationsAtHeight(height uint64, delAdd btypes.Address) (amount uint64, exist bool) {
	exist = mapper.Get(types.BuildRedelegationHeightDelegatorKey(height, delAdd), &amount)
	return
}

func (mapper *Mapper) AddRedelegationAtHeight(height uint64, delAddr btypes.Address, add_amount uint64) {
	amount, exist := mapper.getRedelegationsAtHeight(height, delAddr)
	if exist {
		add_amount += amount
	}
	mapper.setRedelegationsAtHeight(height, delAddr, add_amount)
}

func (mapper *Mapper) RemoveRedelegationsAtHeight(height uint64, delAddr btypes.Address) {
	mapper.Del(types.BuildUnbondingHeightDelegatorKey(height, delAddr))
}
