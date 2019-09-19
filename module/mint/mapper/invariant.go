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
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, mm.GetAllTotalMintQOSAmount())}
		var broken bool
		if !coins.IsNotNegative() {
			broken = true
		}

		phrases, _ := mm.GetInflationPhrases()
		phraseApplied := btypes.ZeroInt()
		phraseTotal := btypes.ZeroInt()
		for _, phrase := range phrases {
			phraseApplied = phraseApplied.Add(phrase.AppliedAmount)
			phraseTotal = phraseTotal.Add(phrase.TotalAmount)
		}

		if !(mm.GetAllTotalMintQOSAmount().Sub(phraseApplied).Add(phraseTotal)).Equal(mm.GetTotalQOSAmount()) {
			broken = true
		}

		return qtypes.FormatInvariant(module, "total-applied",
			fmt.Sprintf("total applied QOS %s\n", coins.String()), coins.Negative(), broken)
	}
}
