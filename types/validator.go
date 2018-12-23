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
	Revoke InactiveCode = iota
	MissVoteBlock
	MaxValidator
)

type Validator struct {
	Name            string         `json:"name"`
	Owner           btypes.Address `json:"owner"`
	ValidatorPubKey crypto.PubKey  `json:"validatorPubkey"`
	BondTokens      uint64         `json:"bondTokens"` //不能超过int64最大值
	Description     string         `json:"description"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactiveCode"`
	InactiveTime   time.Time    `json:"inactiveTime"`
	InactiveHeight uint64       `json:"inactiveHeight"`

	BondHeight uint64 `json:"bondHeight"`
}

func (val Validator) ToABCIValidator() (abciVal abci.Validator) {
	abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ValidatorPubKey)
	abciVal.Power = int64(val.BondTokens)
	abciVal.Address = val.ValidatorPubKey.Address()
	return
}

func (val Validator) IsActive() bool {
	return val.Status == Active
}
