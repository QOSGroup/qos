package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func ApproveInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		var msg string
		var count int

		approves := GetMapper(ctx).GetApproves()
		for _, approve := range approves {
			if approve.QOS.LT(btypes.ZeroInt()) || !approve.QSCs.IsNotNegative() {
				count++
				msg += fmt.Sprintf("approve from %s to %s has a negative values %s\n",
					approve.From.String(), approve.To.String(),
					approve.String())
			}
		}
		broken := count != 0

		return qtypes.FormatInvariant(module, "approve",
			fmt.Sprintf("amount of negative approve found %d\n%s", count, msg), btypes.BaseCoins{}, broken)
	}
}
