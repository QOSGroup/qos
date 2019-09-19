package mapper

import (
	"encoding/binary"
	"time"

	"github.com/QOSGroup/qos/module/params"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
)

func (mapper *Mapper) CreateValidator(validator types.Validator) {
	valAddr := validator.GetValidatorAddress()
	consAddr := validator.ConsAddress()
	mapper.Set(types.BuildValidatorKey(valAddr), validator)
	mapper.Set(types.BuildValidatorByConsensusKey(consAddr), valAddr)
	mapper.Set(types.BuildValidatorByVotePower(validator.ConsensusPower(), valAddr), true)

}

func (mapper *Mapper) ChangeValidatorBondTokens(validator types.Validator, updatedTokens btypes.BigInt) {
	valAddr := validator.GetValidatorAddress()
	mapper.Del(types.BuildValidatorByVotePower(validator.ConsensusPower(), valAddr))
	validator.BondTokens = updatedTokens
	mapper.CreateValidator(validator)
}

func (mapper *Mapper) Exists(valAddress btypes.ValAddress) bool {
	return mapper.Get(types.BuildValidatorKey(valAddress), &(types.Validator{}))
}

func (mapper *Mapper) ExistsWithConsensusAddr(consensusAddr btypes.ConsAddress) bool {
	return mapper.Get(types.BuildValidatorByConsensusKey(consensusAddr), &(btypes.ValAddress{}))
}

func (mapper *Mapper) GetValidator(valAddress btypes.ValAddress) (validator types.Validator, exists bool) {
	validatorKey := types.BuildValidatorKey(valAddress)
	exists = mapper.Get(validatorKey, &validator)
	return
}

func (mapper *Mapper) MakeValidatorInactive(valAddress btypes.ValAddress, inactiveHeight int64, inactiveTime time.Time, code types.InactiveCode) {
	validator, exists := mapper.GetValidator(valAddress)
	if !exists {
		return
	}
	validator.Status = types.Inactive
	validator.InactiveCode = code
	validator.InactiveHeight = inactiveHeight
	validator.InactiveTime = inactiveTime.UTC()
	mapper.Set(types.BuildValidatorKey(valAddress), validator)

	validatorInactiveKey := types.BuildInactiveValidatorKeyByTime(inactiveTime, valAddress)
	mapper.Set(validatorInactiveKey, inactiveTime.UTC().Unix())

	validatorVotePowerKey := types.BuildValidatorByVotePower(validator.ConsensusPower(), valAddress)
	mapper.Del(validatorVotePowerKey)
}

func (mapper *Mapper) KickValidator(valAddress btypes.ValAddress) (validator types.Validator, ok bool) {
	validator, exists := mapper.GetValidator(valAddress)
	if !exists {
		return validator, false
	}
	mapper.Del(types.BuildValidatorKey(valAddress))
	mapper.Del(types.BuildValidatorByConsensusKey(validator.ConsAddress()))
	mapper.Del(types.BuildInactiveValidatorKeyByTime(validator.InactiveTime, valAddress))
	mapper.Del(types.BuildValidatorByVotePower(validator.ConsensusPower(), valAddress))

	return validator, true
}

func (mapper *Mapper) IteratorInactiveValidator(fromSecond, endSecond int64) store.Iterator {

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, uint64(fromSecond))
	startKey := append(types.GetValidatorByInactiveKey(), secBytes...)

	binary.BigEndian.PutUint64(secBytes, uint64(endSecond))
	endKey := append(types.GetValidatorByInactiveKey(), secBytes...)

	return mapper.GetStore().Iterator(startKey, endKey)
}

func (mapper *Mapper) IteratorInactiveValidatorByTime(fromTime, endTime time.Time) store.Iterator {
	return mapper.IteratorInactiveValidator(fromTime.UTC().Unix(), endTime.UTC().Unix())
}

func (mapper *Mapper) IteratorValidatorByVoterPower(ascending bool) store.Iterator {
	if ascending {
		return btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetValidatorByVotePowerKey())
	}
	return btypes.KVStoreReversePrefixIterator(mapper.GetStore(), types.GetValidatorByVotePowerKey())
}

func (mapper *Mapper) GetActiveValidatorSet(ascending bool) (validators []btypes.ValAddress) {
	iterator := mapper.IteratorValidatorByVoterPower(ascending)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.ValAddress(key[9:])
		if _, exists := mapper.GetValidator(valAddr); exists {
			validators = append(validators, valAddr)
		}
	}

	return validators
}

func (mapper *Mapper) MakeValidatorActive(valAddress btypes.ValAddress, addTokens btypes.BigInt) {
	validator, exists := mapper.GetValidator(valAddress)
	if !exists {
		return
	}
	mapper.Del(types.BuildValidatorByVotePower(validator.ConsensusPower(), validator.GetValidatorAddress()))
	bondTokens := validator.BondTokens.Add(addTokens)

	validator.Status = types.Active
	validator.BondTokens = bondTokens

	mapper.Set(types.BuildValidatorKey(validator.GetValidatorAddress()), validator)
	mapper.Del(types.BuildInactiveValidatorKey(validator.InactiveTime.UTC().Unix(), valAddress))
	mapper.Set(types.BuildValidatorByVotePower(validator.ConsensusPower(), validator.GetValidatorAddress()), 1)
}

func (mapper *Mapper) GetValidatorByConsensusAddr(consensusAddr btypes.ConsAddress) (validator types.Validator, exists bool) {
	var valAddress btypes.ValAddress
	exists = mapper.Get(types.BuildValidatorByConsensusKey(consensusAddr), &valAddress)
	if !exists {
		return validator, false
	}

	return mapper.GetValidator(valAddress)
}

func (mapper *Mapper) SetParams(ctx context.Context, p types.Params) {
	params.GetMapper(ctx).SetParamSet(&p)
}

func (mapper *Mapper) GetParams(ctx context.Context) types.Params {
	p := types.Params{}
	params.GetMapper(ctx).GetParamSet(&p)
	return p
}

//-------------------------genesis export

func (mapper *Mapper) IterateValidators(fn func(types.Validator)) {

	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BuildValidatorPrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		mapper.DecodeObject(iter.Value(), &validator)
		fn(validator)
	}
}
