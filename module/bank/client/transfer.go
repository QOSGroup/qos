package client

import (
	"fmt"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qclitx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/module/bank/txs"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
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
			return qclitx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (btxs.ITx, error) {
				sendersStr := viper.GetString(flagSenders)
				senders, err := parseTransItem(ctx, sendersStr)
				if err != nil {
					return nil, err
				}

				receiversStr := viper.GetString(flagReceivers)
				receivers, err := parseTransItem(ctx, receiversStr)
				if err != nil {
					return nil, err
				}

				return txs.TxTransfer{
					Senders:   senders,
					Receivers: receivers,
				}, nil
			})
		},
	}

	cmd.Flags().String(flagSenders, "", "Senders, eg: Arya,10qos,100qstar. multiple users separated by ';' ")
	cmd.Flags().String(flagReceivers, "", "Receivers, eg: address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,10qos,100qstar. multiple users separated by ';'")
	cmd.MarkFlagRequired(flagSenders)
	cmd.MarkFlagRequired(flagReceivers)

	return cmd
}

// Parse flags from string
func parseTransItem(cliCtx context.CLIContext, str string) (types.TransItems, error) {
	items := make(types.TransItems, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		if ti == "" {
			continue
		}

		addrAndCoins := strings.Split(ti, ",")
		if len(addrAndCoins) < 2 {
			return nil, fmt.Errorf("`%s` not match rules", ti)
		}

		addr, err := qcliacc.GetAddrFromValue(cliCtx, addrAndCoins[0])
		if err != nil {
			return nil, err
		}
		qos, qscs, err := qtypes.ParseCoins(strings.Join(addrAndCoins[1:], ","))
		if err != nil {
			return nil, err
		}
		items = append(items, types.TransItem{
			Address: addr,
			QOS:     qos,
			QSCs:    qscs,
		})
	}

	return items, nil
}
