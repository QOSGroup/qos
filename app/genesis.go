package app

import "github.com/QOSGroup/qos/types"

func NewDefaultGenesisState() types.GenesisState {
	return ModuleBasics.DefaultGenesis()
}
