package types

import (
	"fmt"
	"github.com/pkg/errors"
)

type GenesisState struct {
	Guardians []Guardian `json:"guardians"`
}

func (gs GenesisState) Equals(gs1 GenesisState) bool {
	if len(gs.Guardians) != len(gs1.Guardians) {
		return false
	}
	for _, g := range gs.Guardians {
		exists := false
		for _, g1 := range gs1.Guardians {
			if g.Equals(g1) {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}

	return true
}

func NewGenesisState(guardians []Guardian) GenesisState {
	return GenesisState{
		Guardians: guardians,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(gs GenesisState) error {
	addrs := map[string]bool{}
	for _, g := range gs.Guardians {
		if g.GuardianType != Genesis || (g.Creator != nil && len(g.Creator) != 0) {
			return errors.New("invalid genesis guardian")
		}
		if _, ok := addrs[string(g.Address.String())]; ok {
			return errors.New(fmt.Sprintf("duplicate guardian %s", string(g.Address.String())))
		}
		addrs[string(g.Address.String())] = true
	}

	return nil
}
