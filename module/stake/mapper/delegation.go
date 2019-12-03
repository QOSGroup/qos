package mapper

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
)

// 保存委托信息
func (mapper *Mapper) SetDelegationInfo(info types.DelegationInfo) {
	mapper.Set(types.BuildDelegationByDelValKey(info.DelegatorAddr, info.ValidatorAddr), info)
	mapper.Set(types.BuildDelegationByValDelKey(info.ValidatorAddr, info.DelegatorAddr), true)
}

// 获取委托信息
func (mapper *Mapper) GetDelegationInfo(delAddr btypes.AccAddress, valAddr btypes.ValAddress) (info types.DelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildDelegationByDelValKey(delAddr, valAddr), &info)
	return
}

// 删除委托信息
func (mapper *Mapper) DelDelegationInfo(delAddr btypes.AccAddress, valAddr btypes.ValAddress) {
	mapper.Del(types.BuildDelegationByDelValKey(delAddr, valAddr))
	mapper.Del(types.BuildDelegationByValDelKey(valAddr, delAddr))
}

// 根据验证节点地址获取委托列表
func (mapper *Mapper) GetDelegationsByValidator(valAddr btypes.ValAddress) (infos []types.DelegationInfo) {
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

// 遍历委托，当valAddr为空时，遍历所有委托，否则遍历此valAddr下委托
func (mapper *Mapper) IterateDelegationsValDeleAddr(valAddr btypes.ValAddress, fn func(btypes.ValAddress, btypes.AccAddress)) {

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

// 遍历委托，当deleAddr为空时，遍历所有委托，否则遍历此deleAddr的所有委托
func (mapper *Mapper) IterateDelegationsInfo(deleAddr btypes.AccAddress, fn func(types.DelegationInfo)) {

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

// 委托
func (mapper *Mapper) Delegate(ctx context.Context, info types.DelegationInfo, reDelegate bool) {
	if !reDelegate {
		am := baseabci.GetAccountMapper(ctx)
		delegator := am.GetAccount(info.DelegatorAddr).(*qtypes.QOSAccount)
		delegator.MustMinusQOS(info.Amount)
		am.SetAccount(delegator)
	}

	delegation, exists := mapper.GetDelegationInfo(info.DelegatorAddr, info.ValidatorAddr)
	if !exists {
		// 委托存在，保存最新委托信息，然后执行AfterDelegationCreated hooks方法
		mapper.SetDelegationInfo(info)
		mapper.AfterDelegationCreated(ctx, info.ValidatorAddr, info.DelegatorAddr)
	} else {
		// 委托不存在，执行BeforeDelegationModified hooks方法，然后保存委托信息
		delegation.Amount = delegation.Amount.Add(info.Amount)
		delegation.IsCompound = info.IsCompound
		mapper.BeforeDelegationModified(ctx, info.ValidatorAddr, info.DelegatorAddr, delegation.Amount)
		mapper.SetDelegationInfo(delegation)
	}

}

// 解除委托， 执行BeforeDelegationModified hooks方法，保存最新delegation
func (mapper *Mapper) UnbondTokens(ctx context.Context, info types.DelegationInfo, tokens btypes.BigInt) {
	info.Amount = info.Amount.Sub(tokens)
	mapper.BeforeDelegationModified(ctx, info.ValidatorAddr, info.DelegatorAddr, info.Amount)
	unbondHeight := mapper.GetParams(ctx).DelegatorUnbondFrozenHeight + ctx.BlockHeight()
	mapper.AddUnbondingDelegation(types.NewUnbondingDelegationInfo(info.DelegatorAddr, info.ValidatorAddr, ctx.BlockHeight(), unbondHeight, tokens))
	mapper.SetDelegationInfo(info)
}

// 转委托，执行BeforeDelegationModified hooks方法，保存转委托信息
func (mapper *Mapper) ReDelegate(ctx context.Context, delegation types.DelegationInfo, info types.RedelegationInfo) {
	// update origin delegation
	delegation.Amount = delegation.Amount.Sub(info.Amount)
	mapper.BeforeDelegationModified(ctx, delegation.ValidatorAddr, delegation.DelegatorAddr, delegation.Amount)
	mapper.SetDelegationInfo(delegation)

	// save redelegation
	mapper.AddRedelegation(info)
}

// 遍历解除委托
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

// 根据委托账户获取解除委托信息列表
func (mapper *Mapper) GetUnbondingDelegationsByDelegator(delegator btypes.AccAddress) (unbondings []types.UnbondingDelegationInfo) {
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

// 根据验证节点地址获取解除委托信息
func (mapper *Mapper) GetUnbondingDelegationsByValidator(validator btypes.ValAddress) (unbondings []types.UnbondingDelegationInfo) {
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

// 保存解除委托信息
func (mapper *Mapper) SetUnbondingDelegation(unbonding types.UnbondingDelegationInfo) {
	mapper.Set(types.BuildUnbondingHeightDelegatorValidatorKey(unbonding.CompleteHeight, unbonding.DelegatorAddr, unbonding.ValidatorAddr), unbonding)
	mapper.Set(types.BuildUnbondingDelegatorHeightValidatorKey(unbonding.DelegatorAddr, unbonding.CompleteHeight, unbonding.ValidatorAddr), true)
	mapper.Set(types.BuildUnbondingValidatorHeightDelegatorKey(unbonding.ValidatorAddr, unbonding.CompleteHeight, unbonding.DelegatorAddr), true)
}

// 获取解除委托信息，完成时间+委托地址+验证节点地址
func (mapper *Mapper) GetUnbondingDelegation(height int64, delAddr btypes.AccAddress, valAddr btypes.ValAddress) (unbonding types.UnbondingDelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildUnbondingHeightDelegatorValidatorKey(height, delAddr, valAddr), &unbonding)
	return
}

// 新增解除委托信息，若存在累加tokens数量
func (mapper *Mapper) AddUnbondingDelegation(unbonding types.UnbondingDelegationInfo) {
	origin, exist := mapper.GetUnbondingDelegation(unbonding.CompleteHeight, unbonding.DelegatorAddr, unbonding.ValidatorAddr)
	if exist {
		origin.Amount = origin.Amount.Add(unbonding.Amount)
		unbonding = origin
	}
	mapper.SetUnbondingDelegation(unbonding)
}

// 添加解除委托信息
func (mapper *Mapper) AddUnbondingDelegations(unbondingsAdd []types.UnbondingDelegationInfo) {
	for _, unbonding := range unbondingsAdd {
		mapper.AddUnbondingDelegation(unbonding)
	}
}

// 删除解除委托信息
func (mapper *Mapper) RemoveUnbondingDelegation(height int64, delAddr btypes.AccAddress, valAddr btypes.ValAddress) {
	mapper.Del(types.BuildUnbondingHeightDelegatorValidatorKey(height, delAddr, valAddr))
	mapper.Del(types.BuildUnbondingDelegatorHeightValidatorKey(delAddr, height, valAddr))
	mapper.Del(types.BuildUnbondingValidatorHeightDelegatorKey(valAddr, height, delAddr))
}

// 遍历转委托信息
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

// 根据委托地址获取转委托信息
func (mapper *Mapper) GetRedelegationsByDelegator(delegator btypes.AccAddress) (redelegations []types.RedelegationInfo) {
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

// 根据验证节点地址遍历转委托信息
func (mapper *Mapper) GetRedelegationsByFromValidator(validator btypes.ValAddress) (redelegations []types.RedelegationInfo) {
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

// 保存转委托信息
func (mapper *Mapper) SetRedelegation(redelegation types.RedelegationInfo) {
	mapper.Set(types.BuildRedelegationHeightDelegatorFromValidatorKey(redelegation.CompleteHeight, redelegation.DelegatorAddr, redelegation.FromValidator), redelegation)
	mapper.Set(types.BuildRedelegationDelegatorHeightFromValidatorKey(redelegation.DelegatorAddr, redelegation.CompleteHeight, redelegation.FromValidator), true)
	mapper.Set(types.BuildRedelegationFromValidatorHeightDelegatorKey(redelegation.FromValidator, redelegation.CompleteHeight, redelegation.DelegatorAddr), true)
}

// 获取转委托信息，完成时间+委托地址+验证节点地址
func (mapper *Mapper) GetRedelegation(height int64, delAdd btypes.AccAddress, valAddr btypes.ValAddress) (reDelegation types.RedelegationInfo, exist bool) {
	exist = mapper.Get(types.BuildRedelegationHeightDelegatorFromValidatorKey(height, delAdd, valAddr), &reDelegation)
	return
}

// 新增转委托，若存在则累加tokens数量
func (mapper *Mapper) AddRedelegation(redelegation types.RedelegationInfo) {
	origin, exist := mapper.GetRedelegation(redelegation.CompleteHeight, redelegation.DelegatorAddr, redelegation.FromValidator)
	if exist {
		redelegation.Amount = redelegation.Amount.Add(origin.Amount)
	}
	mapper.SetRedelegation(redelegation)
}

// 添加转委托信息
func (mapper *Mapper) AddRedelegations(reDelegations []types.RedelegationInfo) {
	for _, reDelegation := range reDelegations {
		mapper.AddRedelegation(reDelegation)
	}
}

// 删除转委托信息
func (mapper *Mapper) RemoveRedelegation(height int64, delAddr btypes.AccAddress, valAddr btypes.ValAddress) {
	mapper.Del(types.BuildRedelegationHeightDelegatorFromValidatorKey(height, delAddr, valAddr))
	mapper.Del(types.BuildRedelegationDelegatorHeightFromValidatorKey(delAddr, height, valAddr))
	mapper.Del(types.BuildRedelegationFromValidatorHeightDelegatorKey(valAddr, height, delAddr))
}

// 从解除委托信息中扣除惩罚
func (mapper *Mapper) SlashUnbondings(valAddr btypes.ValAddress, infractionHeight int64, fraction qtypes.Dec, maxSlash btypes.BigInt) btypes.BigInt {
	if !maxSlash.GT(btypes.ZeroInt()) {
		return btypes.ZeroInt()
	}
	unbondings := mapper.GetUnbondingDelegationsByValidator(valAddr)
	for _, unbonding := range unbondings {
		// 仅对完成高度大于infractionHeight的解除委托做处理
		if unbonding.Height >= infractionHeight {
			// 惩罚总量不能大于maxSlash
			if !maxSlash.GT(btypes.ZeroInt()) {
				break
			}
			amountSlash := fraction.MulInt(unbonding.Amount).TruncateInt()
			if !maxSlash.GT(amountSlash) {
				amountSlash = maxSlash
			}
			if amountSlash.Equal(unbonding.Amount) {
				mapper.RemoveUnbondingDelegation(unbonding.CompleteHeight, unbonding.DelegatorAddr, unbonding.ValidatorAddr)
			} else {
				unbonding.Amount = unbonding.Amount.Sub(amountSlash)
				mapper.SetUnbondingDelegation(unbonding)
			}
			maxSlash = maxSlash.Sub(amountSlash)
		}
	}

	return maxSlash
}

// 从转委托中扣除惩罚
func (mapper *Mapper) SlashRedelegations(valAddr btypes.ValAddress, infractionHeight int64, fraction qtypes.Dec, maxSlash btypes.BigInt) btypes.BigInt {
	if !maxSlash.GT(btypes.ZeroInt()) {
		return btypes.ZeroInt()
	}
	redelegations := mapper.GetRedelegationsByFromValidator(valAddr)
	for _, redelegation := range redelegations {
		// 仅对完成高度大于infractionHeight的转委托做处理
		if redelegation.Height >= infractionHeight {
			// 惩罚总量不能大于maxSlash
			if maxSlash.Equal(btypes.ZeroInt()) {
				break
			}
			amountSlash := fraction.MulInt(redelegation.Amount).TruncateInt()
			if !maxSlash.GT(amountSlash) {
				amountSlash = maxSlash
			}
			if amountSlash.Equal(redelegation.Amount) {
				mapper.RemoveRedelegation(redelegation.CompleteHeight, redelegation.DelegatorAddr, redelegation.FromValidator)
			} else {
				redelegation.Amount = redelegation.Amount.Sub(amountSlash)
				mapper.SetRedelegation(redelegation)
			}
			maxSlash = maxSlash.Sub(amountSlash)
		}
	}

	return maxSlash
}
