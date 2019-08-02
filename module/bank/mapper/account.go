package mapper

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func GetMapper(ctx context.Context) *account.AccountMapper {
	return baseabci.GetAccountMapper(ctx)
}

func GetAccount(ctx context.Context, addr btypes.Address) *qtypes.QOSAccount {
	account := baseabci.GetAccountMapper(ctx).GetAccount(addr)
	if account != nil {
		return account.(*qtypes.QOSAccount)
	} else {
		return nil
	}
}

func GetAccounts(ctx context.Context) []*qtypes.QOSAccount {
	accounts := []*qtypes.QOSAccount{}
	baseabci.GetAccountMapper(ctx).IterateAccounts(func(acc account.Account) (stop bool) {
		accounts = append(accounts, acc.(*qtypes.QOSAccount))
		return false
	})

	return accounts
}
