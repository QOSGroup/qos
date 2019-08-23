package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func TotalAppliedInvariant(module string) qtypes.Invariant {
	// always return negative coins
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		mm := GetMapper(ctx)
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, btypes.NewInt(int64(mm.GetAllTotalMintQOSAmount())))}
		var broken bool
		if !coins.IsNotNegative() {
			broken = true
		}

		phrases, _ := mm.GetInflationPhrases()
		phraseApplied := uint64(0)
		phraseTotal := uint64(0)
		for _, phrase := range phrases {
			phraseApplied += phrase.AppliedAmount
			phraseTotal += phrase.TotalAmount
		}

		if mm.GetAllTotalMintQOSAmount()-phraseApplied+phraseTotal != mm.GetTotalQOSAmount() {
			broken = true
		}

		return qtypes.FormatInvariant(module, "total-applied",
			fmt.Sprintf("total applied QOS %s\n", coins.String()), coins.Negative(), broken)
	}
}
