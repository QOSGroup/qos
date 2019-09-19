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

			path := types.BuildQueryValidatorPeriodInfoCustomQueryPath(validator)
			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result mapper.ValidatorPeriodInfoQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
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

			path := types.BuildQueryDelegatorIncomeInfoCustomQueryPath(delegator, validator)
			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result mapper.DelegatorIncomeInfoQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagValidator, "", "validator's address")
	cmd.Flags().String(flagDelegator, "", "delegator account address")

	cmd.MarkFlagRequired(flagDelegator)
	cmd.MarkFlagRequired(flagValidator)
	return cmd
}

func queryCommunityFeePoolCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-fee-pool",
		Short: "Query community fee pool",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.Query(fmt.Sprintf("/store/%s/key", types.MapperName), types.BuildCommunityFeePoolKey())
			if err != nil {
				return err
			}

			var result btypes.BigInt
			cdc.MustUnmarshalBinaryBare(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}
