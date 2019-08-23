package types

import (
	"time"

	"github.com/tendermint/tendermint/crypto"

	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type InactiveCode int8

const (
	//Active 可获得挖矿奖励状态
	Active int8 = iota

	//Inactive
	Inactive

	//Inactive Code
	Revoke        InactiveCode = iota // 2
	MissVoteBlock                     // 3
	MaxValidator                      // 4
	DoubleSign                        // 5
)

// Description - description fields for a validator
type Description struct {
	Moniker string `json:"moniker"` // name
	Logo    string `json:"logo"`    // optional logo link
	Website string `json:"website"` // optional website link
	Details string `json:"details"` // optional details
}

type Validator struct {
	Owner           btypes.Address `json:"owner"`
	ValidatorPubKey crypto.PubKey  `json:"pub_key"`
	BondTokens      uint64         `json:"bond_tokens"` //不能超过int64最大值
	Description     Description    `json:"description"`
	Commission      Commission     `json:"commission"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactive_code"`
	InactiveTime   time.Time    `json:"inactive_time"`
	InactiveHeight uint64       `json:"inactive_height"`

	MinPeriod  uint64 `json:"min_period"`
	BondHeight uint64 `json:"bond_height"`
}

func (val Validator) GetValidatorAddress() btypes.Address {
	return btypes.Address(val.ValidatorPubKey.Address())
}

func (val Validator) ToABCIValidator() (abciVal abci.Validator) {
	// abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ValidatorPubKey)
	abciVal.Power = int64(val.BondTokens)
	abciVal.Address = val.ValidatorPubKey.Address()
	return
}

func (val Validator) ToABCIValidatorUpdate(isRemoved bool) (abciVal abci.ValidatorUpdate) {
	abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ValidatorPubKey)
	if isRemoved {
		abciVal.Power = int64(0)
	} else {
		abciVal.Power = int64(val.BondTokens)
	}
	return
}

func (val Validator) IsActive() bool {
	return val.Status == Active
}

func (val Validator) GetBondTokens() uint64 {
	return val.BondTokens
}

func (val Validator) GetValidatorPubKey() crypto.PubKey {
	return val.ValidatorPubKey
}

func (val Validator) GetOwner() btypes.Address {
	return val.Owner
}
