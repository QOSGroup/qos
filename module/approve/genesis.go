package approve

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	approveMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	for _, approve := range data.Approves {
		if approve.IsPositive() {
			approveMapper.SaveApprove(approve)
		}
	}
}

func ExportGenesis(ctx context.Context) types.GenesisState {
	approveMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	return types.NewGenesisState(approveMapper.GetApproves())
}
