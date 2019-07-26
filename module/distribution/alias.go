package distribution

import (
	"github.com/QOSGroup/qos/module/distribution/mapper"
	"github.com/QOSGroup/qos/module/distribution/types"
)

var (
	ModuleName          = types.Distribution
	RegisterCodec       = types.RegisterCodec
	DefaultGenesisState = types.DefaultGenesisState

	MapperName = types.MapperName
	NewMapper  = mapper.NewMapper
	GetMapper  = mapper.GetMapper
	Query      = mapper.Query

	BuildDelegatorEarningStartInfoKey = types.BuildDelegatorEarningStartInfoKey

	ParamsSpace = types.ParamSpace

	NewStakingHooks = mapper.NewStakingHooks
)

type (
	GenesisState = types.GenesisState

	Params = types.Params

	DelegatorEarningsStartInfo = types.DelegatorEarningsStartInfo

	StakingHooks = mapper.StakingHooks
)
