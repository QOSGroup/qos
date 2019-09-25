package types

import (
	"github.com/tendermint/tendermint/crypto"
)

// 创世状态
type GenesisState struct {
	RootPubKey crypto.PubKey `json:"ca_root_pub_key"` // kepler根证书公钥
	QSCs       []QSCInfo     `json:"qscs"`            // 代币信息
}

func NewGenesisState(pubKey crypto.PubKey, qscs []QSCInfo) GenesisState {
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
