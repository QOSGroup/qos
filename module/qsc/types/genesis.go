package types

import (
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	RootPubKey crypto.PubKey `json:"ca_root_pub_key"`
	QSCs       []Info        `json:"qscs"`
}

func NewGenesisState(pubKey crypto.PubKey, qscs []Info) GenesisState {
	return GenesisState{
		RootPubKey: pubKey,
		QSCs:       qscs,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	return nil
}
