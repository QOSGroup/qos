package distribution

import (
	"github.com/QOSGroup/qos/module/distribution/client"
	"github.com/QOSGroup/qos/module/distribution/mapper"
	"github.com/QOSGroup/qos/module/distribution/types"
)

var (
	ModuleName     = types.Distribution
	Cdc            = types.Cdc
	RegisterCodec  = types.RegisterCodec
	DefaultGenesis = types.DefaultGenesisState
	ValidateGenesis = types.ValidateGenesis

	MapperName = types.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper
	Query      = mapper.Query

	BuildDelegatorEarningStartInfoKey = types.BuildDelegatorEarningStartInfoKey

	ParamsSpace = types.ParamSpace

	NewStakingHooks = mapper.NewStakingHooks

	QueryCommands = client.QueryCommands
	TxCommands    = client.TxCommands
)

type (
	GenesisState = types.GenesisState

	Params = types.Params

	DelegatorEarningsStartInfo = types.DelegatorEarningsStartInfo

	StakingHooks = mapper.StakingHooks
)
