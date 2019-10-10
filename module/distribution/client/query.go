package client

import (
	"fmt"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution/mapper"
	"github.com/QOSGroup/qos/module/distribution/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

const (
	flagDelegator = "delegator"
	flagValidator = "validator"
)

func queryValidatorPeriodCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-period [validator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query distribution validator period info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var validator btypes.ValAddress
			if o, err := qcliacc.GetValidatorAddrFromValue(cliCtx, args[0]); err == nil {
				validator = o
			}

			result, err := queryValidatorPeriods(cliCtx, validator)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryValidatorPeriods(cliCtx context.CLIContext, valAddr btypes.ValAddress) (mapper.ValidatorPeriodInfoQueryResult, error) {
	path := types.BuildQueryValidatorPeriodInfoCustomQueryPath(valAddr)
	res, err := cliCtx.Query(path, []byte(""))
	if err != nil {
		return mapper.ValidatorPeriodInfoQueryResult{}, err
	}

	var result mapper.ValidatorPeriodInfoQueryResult
	err = cliCtx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return mapper.ValidatorPeriodInfoQueryResult{}, err
	}

	return result, nil
}

func queryDelegatorIncomeInfoCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegator-income",
		Short: "Query distribution delegator income info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var validator btypes.ValAddress
			var delegator btypes.AccAddress

			if o, err := qcliacc.GetValidatorAddrFromFlag(cliCtx, flagValidator); err == nil {
				validator = o
			}

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagDelegator); err == nil {
				delegator = d
			}

			result, err := queryDelegatorIncomes(cliCtx, delegator, validator)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagValidator, "", "validator's address")
	cmd.Flags().String(flagDelegator, "", "delegator account address")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagValidator)
	return cmd
}

func queryDelegatorIncomes(cliCtx context.CLIContext, delegator btypes.AccAddress, validator btypes.ValAddress) (mapper.DelegatorIncomeInfoQueryResult, error) {
	path := types.BuildQueryDelegatorIncomeInfoCustomQueryPath(delegator, validator)
	res, err := cliCtx.Query(path, []byte(""))
	if err != nil {
		return mapper.DelegatorIncomeInfoQueryResult{}, err
	}

	var result mapper.DelegatorIncomeInfoQueryResult
	err = cliCtx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return mapper.DelegatorIncomeInfoQueryResult{}, err
	}

	return result, nil
}

func queryCommunityFeePoolCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-fee-pool",
		Short: "Query community fee pool",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			result, err := getCommunityFeePool(cliCtx)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func getCommunityFeePool(cliCtx context.CLIContext) (btypes.BigInt, error) {
	res, err := cliCtx.Query(fmt.Sprintf("/store/%s/key", types.MapperName), types.BuildCommunityFeePoolKey())
	if err != nil {
		return btypes.BigInt{}, err
	}

	var result btypes.BigInt
	cliCtx.Codec.MustUnmarshalBinaryBare(res, &result)

	return result, nil
}
