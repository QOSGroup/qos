package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func AccountInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		var msg string
		coins := btypes.BaseCoins{}
		var count int

		accounts := GetAccounts(ctx)
		for _, account := range accounts {
			coins = coins.Plus(append(btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, account.QOS)}, account.QSCs...))
			if account.QOS.LT(btypes.ZeroInt()) || !account.QSCs.IsNotNegative() {
				count++
				msg += fmt.Sprintf("account %s has a negative values %s %s\n",
					account.AccountAddress.String(), account.QOS.String(), account.QSCs.String())
			}
		}
		broken := count != 0

		return qtypes.FormatInvariant(module, "account",
			fmt.Sprintf("amount of negative accounts found %d\n%s", count, msg), coins, broken)
	}
}
