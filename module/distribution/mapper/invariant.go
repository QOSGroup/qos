package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/types"
	qtypes "github.com/QOSGroup/qos/types"
)

// 社区费池数据检查
func FeePoolInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, GetMapper(ctx).GetCommunityFeePool())}
		var broken bool
		// 非负值
		if !coins.IsNotNegative() {
			broken = true
		}

		return qtypes.FormatInvariant(module, "fee-pool",
			fmt.Sprintf("fee pool %s\n", coins.String()), coins, broken)
	}
}

// 待分发奖励检查
func PreDistributionInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		coins := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, GetMapper(ctx).GetPreDistributionQOS())}
		var broken bool
		// 非负值
		if !coins.IsNotNegative() {
			broken = true
		}

		return qtypes.FormatInvariant(module, "pre-distribution",
			fmt.Sprintf("pre distribution %s\n", coins.String()), coins, broken)
	}
}

// 验证节点共享费池检查
func ValidatorFeePoolInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		var msg string
		coins := btypes.BaseCoins{}
		var count int

		// 非负值
		GetMapper(ctx).IteratorValidatorEcoFeePools(func(validatorAddr btypes.ValAddress, pool types.ValidatorEcoFeePool) {

			if pool.ProposerTotalRewardFee.LT(btypes.ZeroInt()) {
				count++
				msg += fmt.Sprintf("validator %s has a negative fee pool proposerTotalRewardFee value %s\n",
					validatorAddr.String(), pool.ProposerTotalRewardFee.String())
			}
			if pool.CommissionTotalRewardFee.LT(btypes.ZeroInt()) {
				count++
				msg += fmt.Sprintf("validator %s has a negative fee pool commissionTotalRewardFee value %s\n",
					validatorAddr.String(), pool.CommissionTotalRewardFee.String())
			}
			if pool.PreDistributeTotalRewardFee.LT(btypes.ZeroInt()) {
				count++
				msg += fmt.Sprintf("validator %s has a negative fee pool preDistributeTotalRewardFee value %s\n",
					validatorAddr.String(), pool.PreDistributeTotalRewardFee.String())
			}
			tokens := btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, pool.PreDistributeRemainTotalFee)}
			coins = coins.Plus(tokens)
			if !tokens.IsNotNegative() {
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
