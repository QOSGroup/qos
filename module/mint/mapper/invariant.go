package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func TotalAppliedInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, btypes.NewInt(int64(GetMapper(ctx).GetAllTotalMintQOSAmount())))}
		var broken bool
		if !coins.IsNotNegative() {
			broken = true
		}

		return qtypes.FormatInvariant(module, "total-applied",
			fmt.Sprintf("total applied QOS %s\n", coins.String()), coins.Negative(), broken)
	}
}
