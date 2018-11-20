package transfer

import (
	"fmt"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	btx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sendersStr := viper.GetString(flagSenders)
			names, senders, err := parseSenderTransItem(&cliCtx, sendersStr)
			if err != nil {
				return err
			}
			receiversStr := viper.GetString(flagReceivers)
			receivers, err := parseReceiverTransItem(receiversStr)
			if err != nil {
				return err
			}
			transferTx := transfer.TxTransfer{
				Senders:   senders,
				Receivers: receivers,
			}

			chainId, err := types.GetDefaultChainId()
			if err != nil {
				return nil
			}

			stdTx := btxs.NewTxStd(&transferTx, chainId, btypes.ZeroInt())
			for _, name := range names {
				info, err := keys.GetKeyInfo(cliCtx, name)
				if err != nil {
					return err
				}
				account, err := cliacc.GetAccount(cliCtx, info.GetAddress())
				if err != nil {
					return err
				}
				stdTx, err = btx.SignStdTx(cliCtx, name, account.GetNonce()+1, stdTx)
				if err != nil {
					return err
				}
			}

			result, err := cliCtx.BroadcastTx(cdc.MustMarshalBinaryBare(stdTx))

			msg, _ := cdc.MarshalJSON(result)
			fmt.Println(string(msg))

			return err
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
