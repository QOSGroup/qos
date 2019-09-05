package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func UnbondingInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		sm := GetMapper(ctx)

		tokens := btypes.ZeroInt()
		sm.IterateUnbondingDelegations(func(unbondings []types.UnbondingDelegationInfo) {
			for _, unbonding := range unbondings {
				tokens = tokens.Add(unbonding.Amount)
			}
		})

		return qtypes.FormatInvariant(module, "unbonding",
			fmt.Sprintf("total unbond tokens %d\n", tokens), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, tokens)}, false)
	}
}

func RedelegationInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		sm := GetMapper(ctx)

		tokens := btypes.ZeroInt()
		sm.IterateRedelegationsInfo(func(redelegations []types.RedelegationInfo) {
			for _, redelegation := range redelegations {
				tokens = tokens.Add(redelegation.Amount)
			}
		})

		return qtypes.FormatInvariant(module, "redelegation",
			fmt.Sprintf("total redelegation tokens %d\n", tokens), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, tokens)}, false)
	}
}

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
				msg += fmt.Sprintf("validator %s bond tokens %d not equals its delegations %d\n",
					validator.GetValidatorAddress().String(), validator.BondTokens, delTokens)
			}
		}
		broken := count != 0

		return qtypes.FormatInvariant(module, "delegation",
			fmt.Sprintf("validator delegations not equals found %d\n%s", count, msg), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, valTokens)}, broken)
	}
}
