package transfer

import (
	"fmt"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/keys"
	btx "github.com/QOSGroup/qbase/client/tx"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/client"
	"github.com/QOSGroup/qos/txs"
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
			names, senders := parseTransItem(&cliCtx, sendersStr)
			receiversStr := viper.GetString(flagReceivers)
			_, receivers := parseTransItem(&cliCtx, receiversStr)
			transferTx := txs.TransferTx{
				Senders:   senders,
				Receivers: receivers,
			}

			chainId, err := client.GetDefaultChainId()
			if err != nil {
				return nil
			}

			stdTx := btxs.NewTxStd(&transferTx, chainId, types.ZeroInt())
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
	cmd.Flags().String(flagReceivers, "", "Receivers, eg: Sansa,10qos,100qstar")

	return cmd
}

// Parse TransItems from string
// str example: Jia,100qos,100qstar;Liu,100qos,100qstar
func parseTransItem(cliCtx *context.CLIContext, str string) ([]string, []txs.TransItem) {
	names := make([]string, 0)
	items := make([]txs.TransItem, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		index := strings.Index(ti, ",")
		name := ti[:index]
		names = append(names, name)
		info, err := keys.GetKeyInfo(*cliCtx, name)
		if err != nil {
			panic(info)
		}
		qos, qscs, err := client.ParseCoins(ti[index+1:])
		if err != nil {
			panic(err)
		}
		items = append(items, txs.TransItem{
			Address: info.GetAddress(),
			QOS:     qos,
			QSCs:    qscs,
		})
	}

	return names, items
}
