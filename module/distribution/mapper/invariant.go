package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func FeepoolInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, GetMapper(ctx).GetCommunityFeePool())}
		var broken bool
		if !coins.IsNotNegative() {
			broken = true
		}

		return qtypes.FormatInvariant(module, "fee-pool",
			fmt.Sprintf("fee pool %s\n", coins.String()), coins, broken)
	}
}

func PreDistributionInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, GetMapper(ctx).GetPreDistributionQOS())}
		var broken bool
		if !coins.IsNotNegative() {
			broken = true
		}

		return qtypes.FormatInvariant(module, "pre-distribution",
			fmt.Sprintf("pre distribution %s\n", coins.String()), coins, broken)
	}
}

func ValidatorFeePoolInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		var msg string
		coins := btypes.BaseCoins{}
		var count int

		GetMapper(ctx).IteratorValidatorEcoFeePools(func(validatorAddr btypes.Address, pool types.ValidatorEcoFeePool) {
			tokens := coins.Plus(btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, pool.PreDistributeRemainTotalFee)})
			coins = append(coins, tokens...)
			if !coins.IsNotNegative() {
				count++
				msg += fmt.Sprintf("validator %s has a negative fee pool value %s\n",
					validatorAddr.String(), tokens.String())
			}
		})
		broken := count != 0

		return qtypes.FormatInvariant(module, "validator-fee-pool",
			fmt.Sprintf("amount of negative validator fee pool found %d\n%s", count, msg), coins, broken)
	}
}
