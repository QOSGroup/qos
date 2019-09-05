package approve

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/types"
)

// 初始化创世状态
func InitGenesis(ctx context.Context, data types.GenesisState) {
	approveMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	for _, approve := range data.Approves {
		// 正常情况不会存在非正预授权
		if approve.IsPositive() {
			approveMapper.SaveApprove(approve)
		}
	}
}

// 导出状态数据
func ExportGenesis(ctx context.Context) types.GenesisState {
	approveMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	return types.NewGenesisState(approveMapper.GetApproves())
}
