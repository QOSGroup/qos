package types

import (
	"fmt"
	qtypes "github.com/QOSGroup/qos/types"
)

type GenesisState struct {
	Accounts []*qtypes.QOSAccount `json:"accounts"`
}

func NewGenesisState(accounts []*qtypes.QOSAccount) GenesisState {
	return GenesisState{
		accounts,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{[]*qtypes.QOSAccount{}}
}

func ValidateGenesis(gs GenesisState) error {
	addrMap := make(map[string]bool, len(gs.Accounts))
	for i := 0; i < len(gs.Accounts); i++ {
		acc := gs.Accounts[i]
		strAddr := string(acc.AccountAddress)
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("duplicate account in genesis state: Address %v", acc.AccountAddress)
		}
		addrMap[strAddr] = true
	}
	return nil
}
