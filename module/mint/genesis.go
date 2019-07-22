package mint

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.SetMintParams(data.Params)

	if data.FirstBlockTime > 0 {
		mapper.SetFirstBlockTime(data.FirstBlockTime)
	}

	if data.AppliedQOSAmount > 0 {
		mapper.SetAllTotalMintQOSAmount(data.AppliedQOSAmount)
	}

}

func ExportGenesis(ctx context.Context) types.GenesisState {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	firstBlockTime := mapper.GetFirstBlockTime()
	return types.GenesisState{
		Params:           mapper.GetMintParams(),
		FirstBlockTime:   firstBlockTime,
		AppliedQOSAmount: mapper.GetAllTotalMintQOSAmount(),
	}
}
