package types

import (
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	RootPubKey crypto.PubKey `json:"ca_root_pub_key"`
	QCPs       []QCPInfo     `json:"qcps""`
}

func NewGenesisState(pubKey crypto.PubKey, qcps []QCPInfo) GenesisState {
	return GenesisState{
		RootPubKey: pubKey,
		QCPs:       qcps,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	return nil
}
