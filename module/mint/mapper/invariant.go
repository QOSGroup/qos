package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

// 通胀数据检查
func TotalAppliedInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		mm := GetMapper(ctx)

		// 总流通QOS不能为负
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

		// 总流通QOS - 通胀规则已发行QOS + 通胀规则总通胀 = QOS计划发行总量
		if !(mm.GetAllTotalMintQOSAmount().Sub(phraseApplied).Add(phraseTotal)).Equal(mm.GetTotalQOSAmount()) {
			broken = true
		}

		// 在QOSApp中会作总币种币量校验，这里返回的 btypes.BaseCoins 数据取反
		return qtypes.FormatInvariant(module, "total-applied",
			fmt.Sprintf("total applied QOS %s\n", coins.String()), coins.Negative(), broken)
	}
}
