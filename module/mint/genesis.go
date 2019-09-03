package mint

import (
	"github.com/QOSGroup/qbase/context"
	qtypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/mint/mapper"
	"github.com/QOSGroup/qos/module/mint/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.SetInflationPhrases(data.InflationPhrases)

	if data.FirstBlockTime > 0 {
		mapper.SetFirstBlockTime(data.FirstBlockTime)
	}

	if data.AppliedQOSAmount.GT(qtypes.ZeroInt()) {
		mapper.SetAllTotalMintQOSAmount(data.AppliedQOSAmount)
	}

	if data.TotalQOSAmount.GT(qtypes.ZeroInt()) {
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
