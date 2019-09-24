package types

import (
	"github.com/tendermint/tendermint/crypto"
)

// 创世状态
type GenesisState struct {
	RootPubKey crypto.PubKey `json:"ca_root_pub_key"` // kepler根证书公钥
	QCPs       []QCPInfo     `json:"qcps""`           // 初始联盟链信息
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
