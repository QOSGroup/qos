package client

import (
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/bank/txs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

const (
	flagSender = "sender"
)

func InvariantCheckCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invariant-check",
		Short: "submit invariant checking tx to proof an invariant broken and halt the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				addr, err := qcliacc.GetAddrFromValue(ctx, viper.GetString(flagSender))
				if err != nil {
					return nil, err
				}

				return txs.TxInvariantCheck{
					Sender: addr,
				}, nil
			})
		},
	}

	cmd.Flags().String(flagSender, "", "Sender's keybase name or address")
	cmd.MarkFlagRequired(flagSender)

	return cmd
}
