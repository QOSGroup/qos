package qcp

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	RootPubKey crypto.PubKey `json:"ca_root_pub_key"`
}

func NewGenesisState(pubKey crypto.PubKey) GenesisState {
	return GenesisState{
		pubKey,
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	SetQCPRootCA(ctx, data.RootPubKey)
}
