package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func QSCsInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		var msg string
		coins := btypes.BaseCoins{}
		var count int

		qscs := GetMapper(ctx).GetQSCs()
		for _, qsc := range qscs {
			coins = coins.Plus(btypes.BaseCoins{btypes.NewBaseCoin(qsc.Name, qsc.TotalAmount)})
			if qsc.TotalAmount.LT(btypes.ZeroInt()) {
				count++
				msg += fmt.Sprintf("qsc %s has a negative values %s \n",
					qsc.Name, qsc.TotalAmount.String())
			}
		}
		broken := count != 0

		return qtypes.FormatInvariant(module, "account",
			fmt.Sprintf("qsc of negative amount found %d\n%s", count, msg), coins.Negative(), broken)
	}
}
