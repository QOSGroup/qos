package mint

import (
	"github.com/QOSGroup/qbase/context"
	minttypes "github.com/QOSGroup/qos/module/mint/types"
)

type GenesisState struct {
	Params minttypes.Params `json:"params"`
}

func NewGenesisState(params minttypes.Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(minttypes.DefaultParams())
}

func InitGenesis(ctx context.Context, data GenesisState) {
	mintMapper := ctx.Mapper(MintMapperName).(*MintMapper)
	mintMapper.SetParams(data.Params)
}

func ExportGenesis(ctx context.Context) GenesisState {
	mintMapper := ctx.Mapper(MintMapperName).(*MintMapper)
	return GenesisState{
		mintMapper.GetParams(),
	}
}
