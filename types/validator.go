package types

import (
	"github.com/tendermint/tendermint/crypto"

	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Validator struct {
	Name        string         `json:"name""`
	ConsPubKey  crypto.PubKey  `json:"cons_pubkey"`
	Operator    btypes.Address `json:"operator"`
	VotingPower int64          `json:"voting_power"`
	Height      int64          `json:"height"`
}

func NewValidator(name string, consPubKey crypto.PubKey, operator btypes.Address, votingPower int64, height int64) Validator {
	return Validator{
		Name:        name,
		ConsPubKey:  consPubKey,
		Operator:    operator,
		VotingPower: votingPower,
		Height:      height,
	}
}

func (val Validator) ToABCIValidator() (abciVal abci.Validator) {
	abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ConsPubKey)
	abciVal.Power = val.VotingPower
	abciVal.Address = val.ConsPubKey.Address()
	return
}
