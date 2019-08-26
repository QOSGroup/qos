package mint

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.SetInflationPhrases(data.InflationPhrases)

	if data.FirstBlockTime > 0 {
		mapper.SetFirstBlockTime(data.FirstBlockTime)
	}

	if data.AppliedQOSAmount > 0 {
		mapper.SetAllTotalMintQOSAmount(data.AppliedQOSAmount)
	}

	if data.TotalQOSAmount > 0 {
		mapper.SetTotalQOSAmount(data.TotalQOSAmount)
	}

}

func ExportGenesis(ctx context.Context) types.GenesisState {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	firstBlockTime := mapper.GetFirstBlockTime()
	return types.GenesisState{
		InflationPhrases: mapper.MustGetInflationPhrases(),
		FirstBlockTime:   firstBlockTime,
		AppliedQOSAmount: mapper.GetAllTotalMintQOSAmount(),
		TotalQOSAmount:   mapper.GetTotalQOSAmount(),
	}
}
