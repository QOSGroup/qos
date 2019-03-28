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

func VoteCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				proposalID := viper.GetInt64(flagProposalID)
				if proposalID <= 0 {
					return nil, errors.New("proposal-id must be positive")
				}

				voter, err := qcliacc.GetAddrFromFlag(ctx, flagVoter)
				if err != nil {
					return nil, err
				}

				option, err := types.VoteOptionFromString(viper.GetString(flagVoteOption))
				if err != nil {
					return nil, errors.New("invalid option")
				}

				return gov.NewTxVote(uint64(proposalID), voter, option), nil
			})
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.Flags().String(flagVoter, "", "Voter")
	cmd.Flags().String(flagVoteOption, "", "Vote option")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagVoter)
	cmd.MarkFlagRequired(flagVoteOption)

	return cmd
}

func QueryVoteCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "Query vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			queryPath := "store/governance/key"

			proposalID := viper.GetInt64(flagProposalID)
			if proposalID <= 0 {
				return errors.New("proposal-id must be positive")
			}

			voter, err := qcliacc.GetAddrFromFlag(cliCtx, flagVoter)
			if err != nil {
				return err
			}

			output, err := cliCtx.Query(queryPath, gov.KeyVote(uint64(proposalID), voter))
			if err != nil {
				return err
			}

			if output == nil {
				return errors.New("vote does not exist")
			}

			deposit := types.Deposit{}
			cdc.MustUnmarshalBinaryBare(output, &deposit)

			return cliCtx.PrintResult(deposit)
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.Flags().String(flagVoter, "", "Voter")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagVoter)

	return cmd
}

func QueryVotesCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes",
		Short: "Query vote list",
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

			result, err := node.ABCIQuery("store/governance/subspace", gov.KeyVotesSubspace(uint64(proposalID)))

			if err != nil {
				return err
			}

			if len(result.Response.Value) == 0 {
				return errors.New("no vote")
			}

			var votes []types.Vote
			var vKVPair []store.KVPair
			cdc.UnmarshalBinaryLengthPrefixed(result.Response.Value, &vKVPair)
			for _, kv := range vKVPair {
				var vote types.Vote
				cdc.UnmarshalBinaryBare(kv.Value, &vote)
				votes = append(votes, vote)
			}

			return cliCtx.PrintResult(votes)
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.MarkFlagRequired(flagProposalID)

	return cmd
}
