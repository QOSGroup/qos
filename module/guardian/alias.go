package guardian

import (
	"github.com/QOSGroup/qos/module/guardian/mapper"
	"github.com/QOSGroup/qos/module/guardian/txs"
	"github.com/QOSGroup/qos/module/guardian/types"
)

var (
	ModuleName      = "guardian"
	RegisterCodec   = txs.RegisterCodec
	NewGenesisState = types.NewGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = mapper.MapperName
	GetMapper  = mapper.GetMapper
	NewMapper  = mapper.NewMapper

	NewGuardian = types.NewGuardian
	Genesis     = types.Genesis
)

type (
	GenesisState = types.GenesisState
	Guardian = types.Guardian
)
