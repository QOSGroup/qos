package client

import (
	"errors"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	gtxs "github.com/QOSGroup/qos/module/gov/txs"
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

				return gtxs.NewTxDeposit(uint64(proposalID), proposer, uint64(deposit)), nil
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
