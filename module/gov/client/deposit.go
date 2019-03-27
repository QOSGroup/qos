package gov

import (
	"errors"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/gov"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

func DepositCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				proposalID := viper.GetInt64(flagProposalID)
				if proposalID <= 0 {
					return nil, errors.New("proposal-id must be positive")
				}

				proposer, err := qcliacc.GetAddrFromFlag(ctx, flagDepositor)
				if err != nil {
					return nil, err
				}

				deposit := viper.GetInt64(flagAmount)
				if deposit <= 0 {
					return nil, errors.New("deposit must be positive")
				}

				return gov.NewTxDeposit(uint64(proposalID), proposer, uint64(deposit)), nil
			})
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.Flags().String(flagDepositor, "", "Depositor")
	cmd.Flags().Uint64(flagAmount, 0, "Percent of QOS for deposit")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagDepositor)
	cmd.MarkFlagRequired(flagAmount)

	return cmd
}

func QueryDepositCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "Query deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryPath := "store/governance/key"

			proposalID := viper.GetInt64(flagProposalID)
			if proposalID <= 0 {
				return errors.New("proposal-id must be positive")
			}

			depositor, err := qcliacc.GetAddrFromFlag(cliCtx, flagDepositor)
			if err != nil {
				return err
			}

			output, err := cliCtx.Query(queryPath, gov.KeyDeposit(uint64(proposalID), depositor))
			if err != nil {
				return err
			}

			if output == nil {
				return errors.New("deposit does not exist")
			}

			deposit := types.Deposit{}
			cdc.MustUnmarshalBinaryBare(output, &deposit)

			return cliCtx.PrintResult(deposit)
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.Flags().String(flagDepositor, "", "Depositor")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagDepositor)

	return cmd
}

func QueryDepositsCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits",
		Short: "Query deposit list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			proposalID := viper.GetInt64(flagProposalID)
			if proposalID <= 0 {
				return errors.New("proposal-id must be positive")
			}

			result, err := node.ABCIQuery("store/governance/subspace", gov.KeyDepositsSubspace(uint64(proposalID)))

			if err != nil {
				return err
			}

			if len(result.Response.Value) == 0 {
				return errors.New("no deposit")
			}

			var deposits []types.Deposit
			var vKVPair []store.KVPair
			cdc.UnmarshalBinaryLengthPrefixed(result.Response.Value, &vKVPair)
			for _, kv := range vKVPair {
				var deposit types.Deposit
				cdc.UnmarshalBinaryBare(kv.Value, &deposit)
				deposits = append(deposits, deposit)
			}

			return cliCtx.PrintResult(deposits)
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.MarkFlagRequired(flagProposalID)

	return cmd
}
