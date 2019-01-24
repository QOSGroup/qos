package mint

import (
	"github.com/QOSGroup/qbase/context"
	mintmapper "github.com/QOSGroup/qos/module/eco/mapper"
	minttypes "github.com/QOSGroup/qos/module/eco/types"
)

type GenesisState struct {
	Params minttypes.MintParams `json:"params"`
}

func NewGenesisState(params minttypes.MintParams) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(minttypes.DefaultMintParams())
}

func InitGenesis(ctx context.Context, data GenesisState) {
	mintMapper := ctx.Mapper(minttypes.MintMapperName).(*mintmapper.MintMapper)
	mintMapper.SetMintParams(data.Params)
}

func ExportGenesis(ctx context.Context) GenesisState {
	mintMapper := ctx.Mapper(minttypes.MintMapperName).(*mintmapper.MintMapper)
	return GenesisState{
		mintMapper.GetMintParams(),
	}
}
