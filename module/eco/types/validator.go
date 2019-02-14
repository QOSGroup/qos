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
)

type Validator struct {
	Name            string         `json:"name"`
	Owner           btypes.Address `json:"owner"`
	ValidatorPubKey crypto.PubKey  `json:"pub_key"`
	BondTokens      uint64         `json:"bond_tokens"` //不能超过int64最大值
	Description     string         `json:"description"`

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
