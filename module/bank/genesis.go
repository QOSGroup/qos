package bank

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	if len(data.Accounts) == 0 {
		return
	}
	// 初始化账户信息
	bm := mapper.GetMapper(ctx)
	for _, acc := range data.Accounts {
		bm.SetAccount(acc)
	}

	// 初始化锁定账户信息
	if data.LockInfo != nil {
		mapper.SetLockInfo(ctx, *data.LockInfo)
	}
}

// 状态数据导出
func ExportGenesis(ctx context.Context) types.GenesisState {
	lockInfo, exists := mapper.GetLockInfo(ctx)
	if !exists {
		return types.NewGenesisState(mapper.GetAccounts(ctx), nil)
	}

	return types.NewGenesisState(mapper.GetAccounts(ctx), &lockInfo)
}
