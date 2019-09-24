package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
)

// 解除委托数据检查
func UnbondingInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		sm := GetMapper(ctx)

		tokens := btypes.ZeroInt()
		var msg string
		var count int

		sm.IterateUnbondingDelegations(func(unbondings []types.UnbondingDelegationInfo) {
			for _, unbonding := range unbondings {
				// 不能存在负值
				if unbonding.Amount.LT(btypes.ZeroInt()) {
					count++
					msg += fmt.Sprintf("unbond token not positive, validator:%s, delegator:%s, token:%s\n",
						unbonding.ValidatorAddr.String(), unbonding.DelegatorAddr.String(), unbonding.Amount)
				}
				tokens = tokens.Add(unbonding.Amount)
			}
		})

		broken := count != 0

		return qtypes.FormatInvariant(module, "unbonding",
			fmt.Sprintf("unbonding negetive found %d\n", count), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, tokens)}, broken)
	}
}

// 转委托数据检查
func RedelegationInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		sm := GetMapper(ctx)

		tokens := btypes.ZeroInt()
		var msg string
		var count int

		sm.IterateRedelegationsInfo(func(redelegations []types.RedelegationInfo) {
			for _, redelegation := range redelegations {
				// 不能存在负值
				if redelegation.Amount.LT(btypes.ZeroInt()) {
					count++
					msg += fmt.Sprintf("redelegation token not positive, from-validator:%s, to-validator:%s, delegator:%s, token:%s\n",
						redelegation.FromValidator.String(), redelegation.ToValidator.String(), redelegation.DelegatorAddr.String(), redelegation.Amount)
				}
				tokens = tokens.Add(redelegation.Amount)
			}
		})

		broken := count != 0

		return qtypes.FormatInvariant(module, "redelegation",
			fmt.Sprintf("redelegation negetive found %d\n", count), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, tokens)}, broken)
	}
}

// 委托数据检查
func DelegationInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		sm := GetMapper(ctx)
		var msg string
		var count int

		var validators []types.Validator
		sm.IterateValidators(func(validator types.Validator) {
			validators = append(validators, validator)
		})

		valTokens := btypes.ZeroInt()
		for _, validator := range validators {
			valTokens = valTokens.Add(validator.BondTokens)
			delTokens := btypes.ZeroInt()
			delegations := sm.GetDelegationsByValidator(validator.GetValidatorAddress())
			for _, delegation := range delegations {
				delTokens = delTokens.Add(delegation.Amount)
			}
			if !validator.BondTokens.Equal(delTokens) {
				count++
				msg += fmt.Sprintf("validator %s bond tokens %s not equals its delegations %s\n",
					validator.GetValidatorAddress().String(), validator.BondTokens, delTokens.String())
			}
		}
		broken := count != 0

		return qtypes.FormatInvariant(module, "delegation",
			fmt.Sprintf("validator delegations not equals found %d\n%s", count, msg), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, valTokens)}, broken)
	}
}
