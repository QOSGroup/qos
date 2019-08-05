package client

import (
	"errors"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	gtxs "github.com/QOSGroup/qos/module/gov/txs"
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

				return gtxs.NewTxVote(uint64(proposalID), voter, option), nil
			})
		},
	}

	cmd.Flags().Uint64(flagProposalID, 0, "Proposal ID")
	cmd.Flags().String(flagVoter, "", "Voter")
	cmd.Flags().String(flagVoteOption, "", "Vote option, possible values: Yes、Abstain、No、NoWithVeto")
	cmd.MarkFlagRequired(flagProposalID)
	cmd.MarkFlagRequired(flagVoter)
	cmd.MarkFlagRequired(flagVoteOption)

	return cmd
}
