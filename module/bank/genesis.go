package bank

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	if len(data.Accounts) == 0 {
		return
	}
	accountMapper := baseabci.GetAccountMapper(ctx)
	for _, acc := range data.Accounts {
		accountMapper.SetAccount(acc)
	}
}

func ExportGenesis(ctx context.Context) types.GenesisState {
	return types.NewGenesisState(mapper.GetAccounts(ctx))
}
