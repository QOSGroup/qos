package issue

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	btypes "github.com/QOSGroup/qbase/types"
	cliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qos/client"
	"github.com/QOSGroup/qos/txs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
)

const(
	flagAmount = "amount"
	flagQscname = "qscname"
	flagPrivkey = "privkeybank"
	flagMaxgas = "maxgas"
	flagQoschainid = "qoschainid"
)

func IssueCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "issue tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount := viper.GetInt64(flagAmount)
			qscname := viper.GetString(flagQscname)
			privkeybank := viper.GetString(flagPrivkey)
			maxgas := viper.GetInt64(flagMaxgas)
			qoschainid := viper.GetString(flagQoschainid)

			addr, privkey := client.ParseJsonPrivkey(cdc, privkeybank)
			return stdTxIssue(&cliCtx, cdc, qscname, btypes.NewInt(amount), addr, privkey, qoschainid, maxgas)
		},
	}

	cmd.Flags().Int64(flagAmount, 100000, "coins send to banker")
	cmd.Flags().String(flagQscname, "qsc1", "qsc name")
	cmd.Flags().String(flagPrivkey, "", "private key of creator")
	cmd.Flags().String(flagQoschainid, "qos", "chainid of qos, used in tx")
	cmd.Flags().Int64(flagMaxgas, 0, "maxgas for txstd")

	return cmd
}

func stdTxIssue(ctx *context.CLIContext, cdc *amino.Codec, qscname string, amount btypes.BigInt,
	bankaddr btypes.Address,privkey crypto.PrivKey, chainid string, maxgas int64) error {
	tx := txs.NewTxIssueQsc(qscname, amount, bankaddr)
	if tx == nil {
		return errors.New("create txissue error")
	}

	acc,err := cliacc.GetAccount(*ctx, bankaddr)
	if err != nil {
		return errors.New("Get banker nonce error!")
	}

	accsigner := []*client.Accsign{}
	accsigner = append(accsigner, &client.Accsign{privkey, acc.GetNonce() + 1})
	txstd := client.GenTxStd(cdc, tx, chainid, maxgas, accsigner)

	return client.BroadcastTxStd(ctx, cdc, txstd)
}