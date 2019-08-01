package bank

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
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
	accounts := []*qtypes.QOSAccount{}
	appendAccount := func(acc account.Account) (stop bool) {
		accounts = append(accounts, acc.(*qtypes.QOSAccount))
		return false
	}
	baseabci.GetAccountMapper(ctx).IterateAccounts(appendAccount)

	return types.NewGenesisState(accounts)
}
