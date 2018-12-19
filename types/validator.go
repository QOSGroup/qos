package types

import (
	"time"

	"github.com/tendermint/tendermint/crypto"

	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	//Active 可获得挖矿奖励状态
	Active int8 = iota

	//InActive
	InActive
)

type Validator struct {
	Name            string         `json:"name"`
	Owner           btypes.Address `json:"owner"`
	ValidatorPubKey crypto.PubKey  `json:"validatorPubkey"`
	BondTokens      uint64         `json:"bondTokens"` //不能超过int64最大值
	Description     string         `json:"description"`

	Status         int8      `json:"status"`
	IsRevoke       bool      `json:"isRevoke"`
	InActiveTime   time.Time `json:"inActiveTime"`
	InActiveHeight uint64    `json:"inActiveHeight"`

	BondHeight uint64 `json:"bondHeight"`
}

func (val Validator) ToABCIValidator() (abciVal abci.Validator) {
	abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ValidatorPubKey)
	abciVal.Power = int64(val.BondTokens)
	abciVal.Address = val.ValidatorPubKey.Address()
	return
}
