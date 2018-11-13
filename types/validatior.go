package types

import "github.com/tendermint/tendermint/crypto"

type Validator struct {
	Name        string        `json:"name""`
	PubKey      crypto.PubKey `json:"pub_key"`
	VotingPower int64         `json:"voting_power"`
	Height      int64         `json:"height"`
}

func NewValidator(name string, pubKey crypto.PubKey, votingPower int64, height int64) Validator {
	return Validator{
		Name:        name,
		PubKey:      pubKey,
		VotingPower: votingPower,
		Height:      height,
	}
}