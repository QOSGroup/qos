package types

import (
	"encoding/binary"
	"fmt"
	"time"

	btypes "github.com/QOSGroup/qbase/types"
)

const (
	ValidatorMapperName = "validator"

	MintMapperName      = "mint"
	MintParamsKey       = "mintparams"

	DelegationMapperName = "delegation"

)

var (
	//keys see docs/spec/staking.md
	validatorKey            = []byte{0x01} // 保存Validator信息. key: ValidatorAddress
	validatorByOwnerKey     = []byte{0x02} // 保存Owner与Validator的映射关系. key: OwnerAddress, value : ValidatorAddress
	validatorByInactiveKey  = []byte{0x03} // 保存处于`inactive`状态的Validator. key: ValidatorInactiveTime + ValidatorAddress
	validatorByVotePowerKey = []byte{0x04} // 按VotePower排序的Validator地址,不包含`pending`状态的Validator. key: VotePower + ValidatorAddress

	DelegationByDelValKey            	= []byte{0x31}	// key: delegator add + validator owner add, value: delegationInfo
	DelegationByValDelKey				= []byte{0x32}	// key: validator owner add + delegator add, value: nil
	DelegatorUnbondingQOSatHeightKey 	= []byte{0x41}	// key: height + delegator add, value: the amount of qos going to be unbonded on this height

	currentValidatorAddressKey = []byte("currentValidatorAddressKey")

	// params
	stakeParamsKey = []byte("stake_params")
)

func BuildValidatorStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", ValidatorMapperName))
}

func BuildCurrentValidatorAddressKey() []byte {
	return currentValidatorAddressKey
}

func BuildValidatorKey(valAddress btypes.Address) []byte {
	return append(validatorKey, valAddress...)
}

func BulidValidatorPrefixKey() []byte {
	return validatorKey
}

func BuildOwnerWithValidatorKey(ownerAddress btypes.Address) []byte {

	lenz := 1 + len(ownerAddress)
	bz := make([]byte, lenz)

	copy(bz[0:1], validatorByOwnerKey)
	copy(bz[1:len(ownerAddress)+1], ownerAddress)

	return bz
}

func BuildInactiveValidatorKeyByTime(inactiveTime time.Time, valAddress btypes.Address) []byte {
	return BuildInactiveValidatorKey(uint64(inactiveTime.UTC().Unix()), valAddress)
}

func BuildInactiveValidatorKey(sec uint64, valAddress btypes.Address) []byte {
	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, sec)

	copy(bz[0:1], validatorByInactiveKey)
	copy(bz[1:9], secBytes)
	copy(bz[9:len(valAddress)+9], valAddress)

	return bz
}

func GetValidatorByInactiveKey() []byte {
	return validatorByInactiveKey
}

func GetValidatorByVotePowerKey() []byte {
	return validatorByVotePowerKey
}

func BuildValidatorByVotePower(votePower uint64, valAddress btypes.Address) []byte {
	lenz := 1 + 8 + len(valAddress)
	bz := make([]byte, lenz)

	votePowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(votePowerBytes, votePower)

	copy(bz[0:1], validatorByVotePowerKey)
	copy(bz[1:9], votePowerBytes)
	copy(bz[9:len(valAddress)+9], valAddress)

	return bz
}

func BuildStakeParamsKey() []byte {
	return stakeParamsKey
}

func BuildMintParamsKey() []byte {
	return []byte(MintParamsKey)
}

func BuildDelegationByDelValKey(delAdd btypes.Address, valAdd btypes.Address) []byte{
	bz := append(DelegationByDelValKey, delAdd...)
	return append(bz, valAdd...)
}

func BuildDelegationByValDelKey(valAdd btypes.Address, delAdd btypes.Address) []byte{
	bz := append(DelegationByValDelKey, valAdd...)
	return append(bz, delAdd...)
}

func BuildUnbondingDelegationByHeightDelKey(height uint64, delAdd btypes.Address) []byte{
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	bz := append(DelegatorUnbondingQOSatHeightKey, heightBytes...)
	return append(bz, delAdd...)
}
