package mint

import (
	"github.com/QOSGroup/qbase/context"
	minttypes "github.com/QOSGroup/qos/modules/mint/types"
)

type GenesisState struct {
	Params minttypes.Params `json:"params"` // inflation params
}

func NewGenesisState(params minttypes.Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: minttypes.DefaultParams(),
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	ctx.Mapper(MintMapperName).(*MintMapper).SetParams(data.Params)
}
