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

func (mapper *Mapper) GetDelegationsByValidator(valAddr btypes.Address) (infos []types.DelegationInfo) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BuildDelegationByValidatorPrefix(valAddr))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		valAddr, delAddr := types.GetDelegationValDelKeyAddress(iter.Key())
		if info, exists := mapper.GetDelegationInfo(delAddr, valAddr); exists {
			infos = append(infos, info)
		}
	}

	return
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

func (mapper *Mapper) Delegate(ctx context.Context, info types.DelegationInfo, reDelegate bool) {
	if !reDelegate {
		am := baseabci.GetAccountMapper(ctx)
		delegator := am.GetAccount(info.DelegatorAddr).(*qtypes.QOSAccount)
		delegator.MustMinusQOS(btypes.NewInt(int64(info.Amount)))
		am.SetAccount(delegator)
	}

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
	mapper.AddUnbondingDelegation(types.NewUnbondingDelegationInfo(info.DelegatorAddr, info.ValidatorAddr, uint64(ctx.BlockHeight()), unbondHeight, tokens))
	mapper.SetDelegationInfo(info)
}

func (mapper *Mapper) ReDelegate(ctx context.Context, delegation types.DelegationInfo, info types.RedelegationInfo) {
	// update origin delegation
	delegation.Amount -= info.Amount
	mapper.BeforeDelegationModified(ctx, delegation.ValidatorAddr, delegation.DelegatorAddr, delegation.Amount)
	mapper.SetDelegationInfo(delegation)

	// save redelegation
	mapper.AddRedelegation(info)
}

func (mapper *Mapper) IterateUnbondingDelegations(fn func([]types.UnbondingDelegationInfo)) {
	unbondings := []types.UnbondingDelegationInfo{}
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.UnbondingHeightDelegatorValidatorKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var unbonding types.UnbondingDelegationInfo
		mapper.DecodeObject(iter.Value(), &unbonding)
		unbondings = append(unbondings, unbonding)
	}
	fn(unbondings)
}

func (mapper *Mapper) GetUnbondingDelegationsByDelegator(delegator btypes.Address) (unbondings []types.UnbondingDelegationInfo) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BuildUnbondingByDelegatorPrefix(delegator))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		delAddr, height, valAddr := types.GetUnbondingDelegationDelegatorHeightValidator(key)
		ubonding, exists := mapper.GetUnbondingDelegation(height, delAddr, valAddr)
		if exists {
			unbondings = append(unbondings, ubonding)
		}
	}

	return
}

func (mapper *Mapper) GetUnbondingDelegationsByValidator(validator btypes.Address) (unbondings []types.UnbondingDelegationInfo) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BuildUnbondingByValidatorPrefix(validator))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		valAddr, height, delAddr := types.GetUnbondingDelegationValidatorHeightDelegator(key)
		ubonding, exists := mapper.GetUnbondingDelegation(height, delAddr, valAddr)
		if exists {
			unbondings = append(unbondings, ubonding)
		}
	}

	return
}

func (mapper *Mapper) SetUnbondingDelegation(unbonding types.UnbondingDelegationInfo) {
	mapper.Set(types.BuildUnbondingHeightDelegatorValidatorKey(unbonding.CompleteHeight, unbonding.DelegatorAddr, unbonding.ValidatorAddr), unbonding)
	mapper.Set(types.BuildUnbondingDelegatorHeightValidatorKey(unbonding.DelegatorAddr, unbonding.CompleteHeight, unbonding.ValidatorAddr), true)
	mapper.Set(types.BuildUnbondingValidatorHeightDelegatorKey(unbonding.ValidatorAddr, unbonding.CompleteHeight, unbonding.DelegatorAddr), true)
}

func (mapper *Mapper) GetUnbondingDelegation(height uint64, delAddr btypes.Address, valAddr btypes.Address) (unbonding types.UnbondingDelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildUnbondingHeightDelegatorValidatorKey(height, delAddr, valAddr), &unbonding)
	return
}

func (mapper *Mapper) AddUnbondingDelegation(unbonding types.UnbondingDelegationInfo) {
	origin, exist := mapper.GetUnbondingDelegation(unbonding.CompleteHeight, unbonding.DelegatorAddr, unbonding.ValidatorAddr)
	if exist {
		origin.Amount += unbonding.Amount
		unbonding = origin
	}
	mapper.SetUnbondingDelegation(unbonding)
}

func (mapper *Mapper) AddUnbondingDelegations(unbondingsAdd []types.UnbondingDelegationInfo) {
	for _, unbonding := range unbondingsAdd {
		mapper.AddUnbondingDelegation(unbonding)
	}
}

func (mapper *Mapper) RemoveUnbondingDelegation(height uint64, delAddr btypes.Address, valAddr btypes.Address) {
	mapper.Del(types.BuildUnbondingHeightDelegatorValidatorKey(height, delAddr, valAddr))
	mapper.Del(types.BuildUnbondingDelegatorHeightValidatorKey(delAddr, height, valAddr))
	mapper.Del(types.BuildUnbondingValidatorHeightDelegatorKey(valAddr, height, delAddr))
}

func (mapper *Mapper) IterateRedelegationsInfo(fn func([]types.RedelegationInfo)) {
	redelegations := []types.RedelegationInfo{}
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.RedelegationHeightDelegatorFromValidatorKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var redelegation types.RedelegationInfo
		mapper.DecodeObject(iter.Value(), &redelegation)
		redelegations = append(redelegations, redelegation)
	}
	fn(redelegations)
}

func (mapper *Mapper) GetRedelegationsByDelegator(delegator btypes.Address) (redelegations []types.RedelegationInfo) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BuildRedelegationByDelegatorPrefix(delegator))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		delAddr, height, valAddr := types.GetRedelegationDelegatorHeightFromValidator(key)
		redelegation, exists := mapper.GetRedelegation(height, delAddr, valAddr)
		if exists {
			redelegations = append(redelegations, redelegation)
		}
	}

	return
}

func (mapper *Mapper) GetRedelegationsByFromValidator(validator btypes.Address) (redelegations []types.RedelegationInfo) {
	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BuildRedelegationByFromValidatorPrefix(validator))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		valAddr, height, delAddr := types.GetRedelegationFromValidatorHeightDelegator(key)
		redelegation, exists := mapper.GetRedelegation(height, delAddr, valAddr)
		if exists {
			redelegations = append(redelegations, redelegation)
		}
	}

	return
}

func (mapper *Mapper) SetRedelegation(redelegation types.RedelegationInfo) {
	mapper.Set(types.BuildRedelegationHeightDelegatorFromValidatorKey(redelegation.CompleteHeight, redelegation.DelegatorAddr, redelegation.FromValidator), redelegation)
	mapper.Set(types.BuildRedelegationDelegatorHeightFromValidatorKey(redelegation.DelegatorAddr, redelegation.CompleteHeight, redelegation.FromValidator), true)
	mapper.Set(types.BuildRedelegationFromValidatorHeightDelegatorKey(redelegation.FromValidator, redelegation.CompleteHeight, redelegation.DelegatorAddr), true)
}

func (mapper *Mapper) GetRedelegation(height uint64, delAdd btypes.Address, valAddr btypes.Address) (reDelegation types.RedelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildRedelegationHeightDelegatorFromValidatorKey(height, delAdd, valAddr), &reDelegation)
	return
}

func (mapper *Mapper) AddRedelegation(redelegation types.RedelegationInfo) {
	origin, exist := mapper.GetRedelegation(redelegation.CompleteHeight, redelegation.DelegatorAddr, redelegation.FromValidator)
	if exist {
		redelegation.Amount += origin.Amount
	}
	mapper.SetRedelegation(redelegation)
}

func (mapper *Mapper) AddRedelegations(reDelegations []types.RedelegationInfo) {
	for _, reDelegation := range reDelegations {
		mapper.AddRedelegation(reDelegation)
	}
}

func (mapper *Mapper) RemoveRedelegation(height uint64, delAddr btypes.Address, valAddr btypes.Address) {
	mapper.Del(types.BuildRedelegationHeightDelegatorFromValidatorKey(height, delAddr, valAddr))
	mapper.Del(types.BuildRedelegationDelegatorHeightFromValidatorKey(delAddr, height, valAddr))
	mapper.Del(types.BuildRedelegationFromValidatorHeightDelegatorKey(valAddr, height, delAddr))
}

func (mapper *Mapper) SlashUnbondings(valAddr btypes.Address, infractionHeight int64, fraction qtypes.Dec, maxSlash int64) int64 {
	unbondings := mapper.GetUnbondingDelegationsByValidator(valAddr)
	for _, unbonding := range unbondings {
		if unbonding.Height >= uint64(infractionHeight) {
			if maxSlash <= 0 {
				break
			}
			amountSlash := fraction.MulInt(btypes.NewInt(int64(unbonding.Amount))).TruncateInt64()
			if maxSlash-amountSlash <= 0 {
				amountSlash = maxSlash
			}
			if amountSlash == int64(unbonding.Amount) {
				mapper.RemoveUnbondingDelegation(unbonding.CompleteHeight, unbonding.DelegatorAddr, unbonding.ValidatorAddr)
			} else {
				unbonding.Amount -= uint64(amountSlash)
				mapper.SetUnbondingDelegation(unbonding)
			}
			maxSlash -= amountSlash
		}
	}

	return maxSlash
}

func (mapper *Mapper) SlashRedelegations(valAddr btypes.Address, infractionHeight int64, fraction qtypes.Dec, maxSlash int64) int64 {
	redelegations := mapper.GetRedelegationsByFromValidator(valAddr)
	for _, redelegation := range redelegations {
		if redelegation.Height >= uint64(infractionHeight) {
			if maxSlash == 0 {
				break
			}
			amountSlash := fraction.MulInt(btypes.NewInt(int64(redelegation.Amount))).TruncateInt64()
			if maxSlash-amountSlash <= 0 {
				amountSlash = maxSlash
			}
			if amountSlash == int64(redelegation.Amount) {
				mapper.RemoveRedelegation(redelegation.CompleteHeight, redelegation.DelegatorAddr, redelegation.FromValidator)
			} else {
				redelegation.Amount -= uint64(amountSlash)
				mapper.SetRedelegation(redelegation)
			}
			maxSlash -= amountSlash
		}
	}

	return maxSlash
}
