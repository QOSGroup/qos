package gov

import (
	"errors"
	"fmt"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/gov"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"strconv"
)

func ProposalCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				title := viper.GetString(flagTitle)
				description := viper.GetString(flagDescription)

				proposalType, err := gtypes.ProposalTypeFromString(viper.GetString(flagProposalType))
				if err != nil {
					return nil, err
				}

				proposer, err := qcliacc.GetAddrFromFlag(ctx, flagProposer)
				if err != nil {
					return nil, err
				}

				deposit := viper.GetInt64(flagDeposit)
				if deposit <= 0 {
					return nil, errors.New("deposit must be positive")
				}

				switch proposalType {
				case gtypes.ProposalTypeText:
					return gov.NewTxProposal(title, description, proposalType, proposer, uint64(deposit)), nil
				case gtypes.ProposalTypeTaxUsage:
					destAddress, err := qcliacc.GetAddrFromFlag(ctx, flagDestAddress)
					if err != nil {
						return nil, err
					}
					percent := viper.GetFloat64(flagPercent)
					if percent <= 0 {
						return nil, errors.New("deposit must be positive")
					}
					ps, _ := types.NewDecFromStr(fmt.Sprintf("%f", percent))
					return gov.NewTxTaxUsage(title, description, proposalType, proposer, uint64(deposit), destAddress, ps), nil
				}

				return nil, errors.New("unknown proposal-type")
			})
		},
	}

	cmd.Flags().String(flagTitle, "", "Proposal title")
	cmd.Flags().String(flagDescription, "", "Proposal description")
	cmd.Flags().String(flagProposalType, gtypes.ProposalTypeText.String(), "")
	cmd.Flags().String(flagProposer, "", "Proposer who submit the proposal")
	cmd.Flags().Uint64(flagDeposit, 0, "Initial deposit paid by proposer. Must be strictly positive")
	cmd.Flags().String(flagDestAddress, "", "Address to receive QOS")
	cmd.Flags().Float64(flagPercent, 0, "Percent of QOS in fee pool send to dest-address")
	cmd.MarkFlagRequired(flagTitle)
	cmd.MarkFlagRequired(flagDescription)
	cmd.MarkFlagRequired(flagProposalType)
	cmd.MarkFlagRequired(flagProposer)
	cmd.MarkFlagRequired(flagDeposit)

	return cmd
}

func QueryProposalCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Short: "Query proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryPath := "store/governance/key"

			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil || pID <= 0 {
				return errors.New("invalid proposal-id")
			}

			output, err := cliCtx.Query(queryPath, gov.KeyProposal(uint64(pID)))
			if err != nil {
				return err
			}

			if output == nil {
				return errors.New("unknown proposal")
			}

			proposal := gtypes.Proposal{}
			cdc.MustUnmarshalBinaryBare(output, &proposal)

			return cliCtx.PrintResult(proposal)
		},
	}

	return cmd
}

func QueryProposalsCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposal list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			result, err := node.ABCIQuery("store/governance/subspace", gov.KeyProposalSubspace())

			if err != nil {
				return err
			}

			if len(result.Response.Value) == 0 {
				return errors.New("no proposal")
			}

			var proposals []gtypes.ProposalContent
			var vKVPair []store.KVPair
			cdc.UnmarshalBinaryLengthPrefixed(result.Response.Value, &vKVPair)
			for _, kv := range vKVPair {
				var proposal gtypes.Proposal
				cdc.UnmarshalBinaryBare(kv.Value, &proposal)
				proposals = append(proposals, proposal)
			}

			return cliCtx.PrintResult(proposals)
		},
	}

	return cmd
}
