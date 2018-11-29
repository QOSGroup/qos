package transfer

import (
	"fmt"
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
				senders, err := parseTransItem(ctx, sendersStr)
				if err != nil {
					return nil, err
				}

				receiversStr := viper.GetString(flagReceivers)
				receivers, err := parseTransItem(ctx, receiversStr)
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

	cmd.Flags().String(flagSenders, "", "Senders, eg: Arya,10qos,100qstar. multiple users separated by `;` ")
	cmd.Flags().String(flagReceivers, "", "Receivers, eg: address1vkl6nc6eedkxwjr5rsy2s5jr7qfqm487wu95w7,10qos,100qstar. multiple users separated by `;`")
	cmd.MarkFlagRequired(flagSenders)
	cmd.MarkFlagRequired(flagReceivers)

	return cmd
}

// Parse flags from string
func parseTransItem(cliCtx context.CLIContext, str string) ([]transfer.TransItem, error) {
	items := make([]transfer.TransItem, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		if ti == "" {
			continue
		}

		addrAndCoins := strings.Split(ti , ",")
		if len(addrAndCoins) < 2 {
			return nil , fmt.Errorf("`%s` not match rules", ti)
		}

		addr, err := getAddress(cliCtx, addrAndCoins[0])
		if err != nil {
			return nil, err
		}
		qos, qscs, err := types.ParseCoins(strings.Join(addrAndCoins[1:],","))
		if err != nil {
			return nil, err
		}
		items = append(items, transfer.TransItem{
			Address: addr,
			QOS:     qos,
			QSCs:    qscs,
		})
	}

	return items, nil
}


func getAddress(ctx context.CLIContext , value string) (btypes.Address , error) {
		address, err := btypes.GetAddrFromBech32(value)
		if err == nil {
			return address , nil
		}

		info, err := keys.GetKeyInfo(ctx, value)
		if err != nil {
			return nil, err
		}

		return info.GetAddress(),nil
}
