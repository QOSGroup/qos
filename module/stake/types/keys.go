package types

import (
	"encoding/binary"
	"fmt"
	"time"

	btypes "github.com/QOSGroup/qbase/types"
)

const (
	AddrLen = 20

	MapperName = "validator"

	//------query-------
	Stake         = "stake"
	Delegation    = "delegation"
	Delegations   = "delegations"
	Owner         = "owner"
	Delegator     = "delegator"
	Unbondings    = "Unbondings"
	Redelegations = "Redelegations"
)

var (
	//keys see docs/spec/staking.md
	validatorKey            = []byte{0x01} // 保存Validator信息. key: ValidatorAddress
	validatorByOwnerKey     = []byte{0x02} // 保存Owner与Validator的映射关系. key: OwnerAddress, value : ValidatorAddress
	validatorByInactiveKey  = []byte{0x03} // 保存处于`inactive`状态的Validator. key: ValidatorInactiveTime + ValidatorAddress
	validatorByVotePowerKey = []byte{0x04} // 按VotePower排序的Validator地址,不包含`pending`状态的Validator. key: VotePower + ValidatorAddress

	//keys see docs/spec/staking.md
	validatorVoteInfoKey         = []byte{0x11} // 保存Validator在窗口的统计信息
	validatorVoteInfoInWindowKey = []byte{0x12} // 保存Validator在指定窗口签名信息

	DelegationByDelValKey = []byte{0x31} // key: delegator add + validator owner add, value: delegationInfo
	DelegationByValDelKey = []byte{0x32} // key: validator owner add + delegator add, value: nil

	UnbondingHeightDelegatorKey = []byte{0x41} // key: height + delegator add, value: the amount of qos going to be unbonded on this height
	UnbondingDelegatorHeightKey = []byte{0x42} // key: delegator + height add, value: nil

	RedelegationHeightDelegatorKey = []byte{0x51} // key: height + delegator add, value: redelegations going to be complete on this height
	RedelegationDelegatorHeightKey = []byte{0x52} // key: delegator + height add, value: nil

	currentValidatorsAddressKey = []byte("currentValidatorsAddressKey")
)

func BuildStakeStoreQueryPath() []byte {
	return []byte(fmt.Sprintf("/store/%s/key", MapperName))
}

func BuildCurrentValidatorsAddressKey() []byte {
	return currentValidatorsAddressKey
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

func GetValidatorVoteInfoInWindowKey() []byte {
	return validatorVoteInfoInWindowKey
}

func GetValidatorVoteInfoKey() []byte {
	return validatorVoteInfoKey
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

func BuildDelegationByDelValKey(delAdd btypes.Address, valAdd btypes.Address) []byte {
	bz := append(DelegationByDelValKey, delAdd...)
	return append(bz, valAdd...)
}

func BuildDelegationByValDelKey(valAdd btypes.Address, delAdd btypes.Address) []byte {
	bz := append(DelegationByValDelKey, valAdd...)
	return append(bz, delAdd...)
}

func GetDelegationValDelKeyAddress(key []byte) (valAddr btypes.Address, deleAddr btypes.Address) {
	if len(key) != 1+2*AddrLen {
		panic("invalid DelegationValDelKey length")
	}

	valAddr = key[1 : 1+AddrLen]
	deleAddr = key[1+AddrLen:]
	return
}

func BuildValidatorVoteInfoKey(valAddress btypes.Address) []byte {
	return append(validatorVoteInfoKey, valAddress...)
}

func BuildValidatorVoteInfoInWindowPrefixKey(valAddress btypes.Address) []byte {
	return append(validatorVoteInfoInWindowKey, valAddress...)
}

func GetValidatorVoteInfoAddr(key []byte) btypes.Address {
	return btypes.Address(key[1:])
}

func BuildValidatorVoteInfoInWindowKey(index uint64, valAddress btypes.Address) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, index)

	bz := append(validatorVoteInfoInWindowKey, valAddress...)
	bz = append(bz, b...)

	return bz
}

func GetValidatorVoteInfoInWindowIndexAddr(key []byte) (uint64, btypes.Address) {
	addr := btypes.Address(key[1 : AddrLen+1])
	index := binary.LittleEndian.Uint64(key[AddrLen+1:])
	return index, addr
}

func BuildUnbondingHeightDelegatorKey(height uint64, delAdd btypes.Address) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	bz := append(UnbondingHeightDelegatorKey, heightBytes...)
	return append(bz, delAdd...)
}

func BuildUnbondingDelegationByHeightPrefix(height uint64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	return append(UnbondingHeightDelegatorKey, heightBytes...)
}

func GetUnbondingDelegationHeightAddress(key []byte) (height uint64, deleAddr btypes.Address) {

	if len(key) != (1 + 8 + AddrLen) {
		panic("invalid UnbondingHeightDelegatorKey length")
	}

	height = binary.BigEndian.Uint64(key[1:9])
	deleAddr = btypes.Address(key[9:])
	return
}

func GetUnbondingDelegationAddressHeight(key []byte) (deleAddr btypes.Address, height uint64) {

	if len(key) != (1 + 8 + AddrLen) {
		panic("invalid UnbondingDelegatorHeightKey length")
	}

	deleAddr = btypes.Address(key[1 : AddrLen+1])
	height = binary.BigEndian.Uint64(key[AddrLen+1:])
	return
}

func BuildUnbondingDelegatorHeightKey(delAddr btypes.Address, height uint64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	bz := append(UnbondingDelegatorHeightKey, delAddr...)
	return append(bz, heightBytes...)
}

func BuildUnbondingByDelegatorPrefix(delAddr btypes.Address) []byte {

	return append(UnbondingDelegatorHeightKey, delAddr...)
}

func BuildRedelegationHeightDelegatorKey(height uint64, delAdd btypes.Address) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	bz := append(RedelegationHeightDelegatorKey, heightBytes...)
	return append(bz, delAdd...)
}

func BuildRedelegationByHeightPrefix(height uint64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	return append(RedelegationHeightDelegatorKey, heightBytes...)
}

func GetRedelegationHeightAddress(key []byte) (height uint64, deleAddr btypes.Address) {

	if len(key) != (1 + 8 + AddrLen) {
		panic("invalid RedelegationHeightDelegatorKey length")
	}

	height = binary.BigEndian.Uint64(key[1:9])
	deleAddr = btypes.Address(key[9:])
	return
}

func GetRedelegationAddressHeight(key []byte) (deleAddr btypes.Address, height uint64) {

	if len(key) != (1 + 8 + AddrLen) {
		panic("invalid RedelegationDelegatorHeightKey length")
	}

	deleAddr = btypes.Address(key[1 : AddrLen+1])
	height = binary.BigEndian.Uint64(key[AddrLen+1:])
	return
}

func BuildRedelegationDelegatorHeightKey(delAddr btypes.Address, height uint64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, height)

	bz := append(RedelegationDelegatorHeightKey, delAddr...)
	return append(bz, heightBytes...)
}

func BuildRedelegationByDelegatorPrefix(delAddr btypes.Address) []byte {

	return append(RedelegationDelegatorHeightKey, delAddr...)
}

//-------------------------query path

func BuildGetDelegationCustomQueryPath(deleAddr, owner btypes.Address) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Stake, Delegation, deleAddr.String(), owner.String())
}

func BuildQueryDelegationsByOwnerCustomQueryPath(owner btypes.Address) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Stake, Delegations, Owner, owner.String())
}

func BuildQueryDelegationsByDelegatorCustomQueryPath(deleAddr btypes.Address) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", Stake, Delegations, Delegator, deleAddr.String())
}

func BuildQueryUnbondingsByDelegatorCustomQueryPath(deleAddr btypes.Address) string {
	return fmt.Sprintf("custom/%s/%s/%s", Stake, Unbondings, deleAddr.String())
}

func BuildQueryRedelegationsByDelegatorCustomQueryPath(deleAddr btypes.Address) string {
	return fmt.Sprintf("custom/%s/%s/%s", Stake, Redelegations, deleAddr.String())
}
