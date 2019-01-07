package app

import (
	"github.com/QOSGroup/qos/account"
	staketypes "github.com/QOSGroup/qos/modules/stake/types"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

// QOS初始状态
type GenesisState struct {
	CAPubKey   crypto.PubKey          `json:"ca_pub_key"`
	Accounts   []*account.QOSAccount  `json:"accounts"`
	Validators []staketypes.Validator `json:"validators"`

	SPOConfig   types.SPOConfig   `json:"spo_config"`
	StakeConfig types.StakeConfig `json:"stake_config"`
}

func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		SPOConfig:   types.DefaultSPOConfig(),
		StakeConfig: types.DefaultStakeConfig(),
	}
}
