package transfer

import (
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/txs/transfer"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"strings"
)

const (
	flagSenders   = "senders"
	flagReceivers = "receivers"
)

func TransferCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer QOS and QSCs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				sendersStr := viper.GetString(flagSenders)
				_, senders, err := parseSenderTransItem(&ctx, sendersStr)
				if err != nil {
					return nil, err
				}
				receiversStr := viper.GetString(flagReceivers)
				receivers, err := parseReceiverTransItem(receiversStr)
				if err != nil {
					return nil, err
				}
				return transfer.TxTransfer{
					Senders:   senders,
					Receivers: receivers,
				}, nil
			})
		},
	}

	cmd.Flags().String(flagSenders, "", "Senders, eg: Arya,10qos,100qstar")
	cmd.Flags().String(flagReceivers, "", "Receivers, eg: address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,10qos,100qstar")
	cmd.MarkFlagRequired(flagSenders)
	cmd.MarkFlagRequired(flagReceivers)

	return cmd
}

// Parse SenderTransItems from string
func parseSenderTransItem(cliCtx *context.CLIContext, str string) ([]string, []transfer.TransItem, error) {
	names := make([]string, 0)
	items := make([]transfer.TransItem, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		index := strings.Index(ti, ",")
		name := ti[:index]
		names = append(names, name)
		info, err := keys.GetKeyInfo(*cliCtx, name)
		if err != nil {
			return nil, nil, err
		}
		qos, qscs, err := types.ParseCoins(ti[index+1:])
		if err != nil {
			return nil, nil, err
		}
		items = append(items, transfer.TransItem{
			Address: info.GetAddress(),
			QOS:     qos,
			QSCs:    qscs,
		})
	}

	return names, items, nil
}

// Parse ReceiverTransItems from string
func parseReceiverTransItem(str string) ([]transfer.TransItem, error) {
	items := make([]transfer.TransItem, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		index := strings.Index(ti, ",")
		addr := ti[:index]
		address, err := btypes.GetAddrFromBech32(addr)
		if err != nil {
			return nil, err
		}
		qos, qscs, err := types.ParseCoins(ti[index+1:])
		if err != nil {
			return nil, err
		}
		items = append(items, transfer.TransItem{
			Address: address,
			QOS:     qos,
			QSCs:    qscs,
		})
	}

	return items, nil
}
