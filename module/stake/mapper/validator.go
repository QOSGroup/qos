package mapper

import (
	"encoding/binary"
	"github.com/QOSGroup/qos/module/params"
	"time"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
)

func (mapper *Mapper) CreateValidator(validator types.Validator) {
	valAddr := validator.GetValidatorAddress()
	mapper.Set(types.BuildValidatorKey(valAddr), validator)
	mapper.Set(types.BuildOwnerWithValidatorKey(validator.Owner), valAddr)
	mapper.Set(types.BuildValidatorByVotePower(validator.BondTokens, valAddr), true)
}

func (mapper *Mapper) ChangeValidatorBondTokens(validator types.Validator, updatedTokens uint64) {
	valAddr := validator.GetValidatorAddress()
	mapper.Del(types.BuildValidatorByVotePower(validator.BondTokens, valAddr))
	validator.BondTokens = updatedTokens
	mapper.CreateValidator(validator)
}

func (mapper *Mapper) Exists(valAddress btypes.Address) bool {
	return mapper.Get(types.BuildValidatorKey(valAddress), &(types.Validator{}))
}

func (mapper *Mapper) ExistsWithOwner(owner btypes.Address) bool {
	return mapper.Get(types.BuildOwnerWithValidatorKey(owner), &(btypes.Address{}))
}

func (mapper *Mapper) GetValidator(valAddress btypes.Address) (validator types.Validator, exists bool) {
	validatorKey := types.BuildValidatorKey(valAddress)
	exists = mapper.Get(validatorKey, &validator)
	return
}

func (mapper *Mapper) MakeValidatorInactive(valAddress btypes.Address, inactiveHeight uint64, inactiveTime time.Time, code types.InactiveCode) {
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

	validatorVotePowerKey := types.BuildValidatorByVotePower(validator.BondTokens, valAddress)
	mapper.Del(validatorVotePowerKey)
}

func (mapper *Mapper) KickValidator(valAddress btypes.Address) (validator types.Validator, ok bool) {
	validator, exists := mapper.GetValidator(valAddress)
	if !exists {
		return validator, false
	}
	mapper.Del(types.BuildValidatorKey(valAddress))
	mapper.Del(types.BuildOwnerWithValidatorKey(validator.Owner))
	mapper.Del(types.BuildInactiveValidatorKeyByTime(validator.InactiveTime, valAddress))
	mapper.Del(types.BuildValidatorByVotePower(validator.BondTokens, valAddress))

	return validator, true
}

func (mapper *Mapper) IteratorInactiveValidator(fromSecond, endSecond uint64) store.Iterator {

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, fromSecond)
	startKey := append(types.GetValidatorByInactiveKey(), secBytes...)

	binary.BigEndian.PutUint64(secBytes, endSecond)
	endKey := append(types.GetValidatorByInactiveKey(), secBytes...)

	return mapper.GetStore().Iterator(startKey, endKey)
}

func (mapper *Mapper) IteratorInactiveValidatorByTime(fromTime, endTime time.Time) store.Iterator {
	return mapper.IteratorInactiveValidator(uint64(fromTime.UTC().Unix()), uint64(endTime.UTC().Unix()))
}

func (mapper *Mapper) IteratorValidatorByVoterPower(ascending bool) store.Iterator {
	if ascending {
		return btypes.KVStorePrefixIterator(mapper.GetStore(), types.GetValidatorByVotePowerKey())
	}
	return btypes.KVStoreReversePrefixIterator(mapper.GetStore(), types.GetValidatorByVotePowerKey())
}

func (mapper *Mapper) GetActiveValidatorSet(ascending bool) (validators []btypes.Address) {
	iterator := mapper.IteratorValidatorByVoterPower(ascending)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])
		if _, exists := mapper.GetValidator(valAddr); exists {
			validators = append(validators, valAddr)
		}
	}

	return validators
}

func (mapper *Mapper) MakeValidatorActive(valAddress btypes.Address, addTokens uint64) {
	validator, exists := mapper.GetValidator(valAddress)
	if !exists {
		return
	}
	mapper.Del(types.BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()))
	validator.Status = types.Active
	validator.BondTokens += addTokens

	mapper.Set(types.BuildValidatorKey(validator.ValidatorPubKey.Address().Bytes()), validator)
	mapper.Del(types.BuildInactiveValidatorKey(uint64(validator.InactiveTime.UTC().Unix()), valAddress))
	mapper.Set(types.BuildValidatorByVotePower(validator.BondTokens, validator.ValidatorPubKey.Address().Bytes()), 1)
}

func (mapper *Mapper) GetValidatorByOwner(owner btypes.Address) (validator types.Validator, exists bool) {
	var valAddress btypes.Address
	exists = mapper.Get(types.BuildOwnerWithValidatorKey(owner), &valAddress)
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

	iter := btypes.KVStorePrefixIterator(mapper.GetStore(), types.BulidValidatorPrefixKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		mapper.DecodeObject(iter.Value(), &validator)
		fn(validator)
	}
}
