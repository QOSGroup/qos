package mint

import (
	"github.com/QOSGroup/qbase/context"
	mintmapper "github.com/QOSGroup/qos/module/eco/mapper"
	minttypes "github.com/QOSGroup/qos/module/eco/types"
)

type GenesisState struct {
	Params           minttypes.MintParams `json:"params"`
	FirstBlockTime   int64                `json:"first_block_time"` //UTC().UNIX()
	AppliedQOSAmount uint64               `json:"applied_qos_amount"`
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

	if data.FirstBlockTime > 0 {
		mintMapper.SetFirstBlockTime(data.FirstBlockTime)
	}

	if data.AppliedQOSAmount > 0 {
		mintMapper.SetAppliedQOSAmount(data.AppliedQOSAmount)
	}

}

func ExportGenesis(ctx context.Context) GenesisState {
	mintMapper := ctx.Mapper(minttypes.MintMapperName).(*mintmapper.MintMapper)
	firstBlockTime := mintMapper.GetFirstBlockTime()
	return GenesisState{
		Params:           mintMapper.GetMintParams(),
		FirstBlockTime:   firstBlockTime,
		AppliedQOSAmount: mintMapper.GetAppliedQOSAmount(),
	}
}
