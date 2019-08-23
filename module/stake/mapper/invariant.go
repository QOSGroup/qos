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

		tokens := uint64(0)
		sm.IterateUnbondingDelegations(func(unbondings []types.UnbondingDelegationInfo) {
			for _, unbonding := range unbondings {
				tokens += unbonding.Amount
			}
		})

		return qtypes.FormatInvariant(module, "unbonding",
			fmt.Sprintf("total unbond tokens %d\n", tokens), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, btypes.NewInt(int64(tokens)))}, false)
	}
}

func RedelegationInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		sm := GetMapper(ctx)

		tokens := uint64(0)
		sm.IterateRedelegationsInfo(func(redelegations []types.RedelegationInfo) {
			for _, redelegation := range redelegations {
				tokens += redelegation.Amount
			}
		})

		return qtypes.FormatInvariant(module, "redelegation",
			fmt.Sprintf("total redelegation tokens %d\n", tokens), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, btypes.NewInt(int64(tokens)))}, false)
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

		valTokens := uint64(0)
		for _, validator := range validators {
			valTokens += validator.BondTokens
			delTokens := uint64(0)
			delegations := sm.GetDelegationsByValidator(validator.GetValidatorAddress())
			for _, delegation := range delegations {
				delTokens += delegation.Amount
			}
			if validator.BondTokens != delTokens {
				count++
				msg += fmt.Sprintf("validator %s bond tokens %d not equals its delegations %d\n",
					validator.GetValidatorAddress().String(), validator.BondTokens, delTokens)
			}
		}
		broken := count != 0

		return qtypes.FormatInvariant(module, "delegation",
			fmt.Sprintf("validator delegations not equals found %d\n%s", count, msg), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, btypes.NewInt(int64(valTokens)))}, broken)
	}
}
