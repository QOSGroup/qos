package staking

import (
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/QOSGroup/qos/txs/staking"
	"github.com/QOSGroup/qos/types"
	go_amino "github.com/tendermint/go-amino"
)

const (
	flagActive = "active"
)

func queryValidatorInfoCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Query validator's info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			owner := viper.GetString(flagOwner)
			ownerAddress, err := btypes.GetAddrFromBech32(owner)
			if err != nil {
				return err
			}

			opts := buildQueryOptions()

			result, err := node.ABCIQueryWithOptions(string(staking.BuildValidatorStoreQueryPath()), staking.BuildOwnerWithValidatorKey(ownerAddress), opts)
			if err != nil {
				return err
			}

			valueBz := result.Response.GetValue()
			if len(valueBz) == 0 {
				return errors.New("owner does't have validator")
			}

			var address btypes.Address
			cdc.UnmarshalBinaryBare(valueBz, &address)

			key := staking.BuildValidatorKey(address)
			result, err = node.ABCIQueryWithOptions(string(staking.BuildValidatorStoreQueryPath()), key, opts)
			if err != nil {
				return err
			}

			valueBz = result.Response.GetValue()
			if len(valueBz) == 0 {
				return errors.New("response empty value")
			}

			var validator types.Validator
			cdc.UnmarshalBinaryBare(valueBz, &validator)

			return cliCtx.PrintResult(validator)
		},
	}

	cmd.Flags().String(flagOwner, "", "validator's owner address")
	cmd.MarkFlagRequired(flagOwner)

	return cmd
}

func queryAllValidatorsCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query all validators info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			opts := buildQueryOptions()

			subspace := "/store/validator/subspace"
			result, err := node.ABCIQueryWithOptions(subspace, staking.BulidValidatorPrefixKey(), opts)

			if err != nil {
				return err
			}

			valueBz := result.Response.GetValue()
			if len(valueBz) == 0 {
				return errors.New("response empty value")
			}

			var validators []types.Validator

			var vKVPair []store.KVPair
			cdc.UnmarshalBinary(valueBz, &vKVPair)
			for _, kv := range vKVPair {
				var validator types.Validator
				cdc.UnmarshalBinaryBare(kv.Value, &validator)
				validators = append(validators, validator)
			}

			cPrint(cliCtx, validators)

			return nil
		},
	}

	// cmd.Flags().Bool(flagActive, false, "only query active status validator")

	return cmd
}

func buildQueryOptions() client.ABCIQueryOptions {
	height := viper.GetInt64(bctypes.FlagHeight)
	if height <= 0 {
		height = 0
	}

	trust := viper.GetBool(bctypes.FlagTrustNode)

	return client.ABCIQueryOptions{
		Height:  height,
		Trusted: trust,
	}
}

func cPrint(cliCtx context.CLIContext, validators []types.Validator) {
	fmt.Println("validators: ")
	for _, validator := range validators {
		cliCtx.PrintResult(validator)
	}
}
