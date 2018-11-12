package createqsc

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/client"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
)

const(
	flagExtrate = "extrate"
	flagPathqsc = "pathqsc"
	flagPathbank = "pathbank"
	flagQscchainid = "qscchainid"
	flagPrivkey = "privkey"
	flagQoschainid = "qoschainid"
	flagMaxgas = "maxgas"
	flagInitAccount = "initaccount"
	flagQscname = "qscname"
	flagNonce = "nonce"
)

func CreateQscCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createqsc",
		Short: "create qsc tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			extrate := viper.GetString(flagExtrate)
			pathqsc := viper.GetString(flagPathqsc)
			pathbank := viper.GetString(flagPathbank)
			qscchainid := viper.GetString(flagQscchainid)
			privkeycreator := viper.GetString(flagPrivkey)
			qoschainid := viper.GetString(flagQoschainid)
			maxgas := viper.GetInt64(flagMaxgas)
			binitaccount := viper.GetBool(flagInitAccount)
			nonce := viper.GetInt64(flagNonce)

			if !client.Cmdcheck("txcreateqsc", privkeycreator, pathbank, pathqsc) {
				return errors.New("invalidate param")
			}

			caQsc := txs.FetchCA(pathqsc)
			caBanker := txs.FetchCA(pathbank)
			acc := []txs.AddrCoin{}
			if binitaccount {
				accprivkeys := []string {
					"vAeIlHuWjvz/JmyGcB46ZHfCZdXCYuRogqxDgjYUM5wNwKIyIYQBs9VZxGyD9FS5J4XvZntnUaTtoGsEl7+3hg==",
					"31PlT2p6UICjV63dG7Nh3Mh9W0b+7FAEU+KOAxyNbZ29rwqNzxQJlQPh59tZpbS1EdIT6TE5N6L72se9BUe9iw==",
					"9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ==",
				}

				for _, v := range accprivkeys {
					addracc, _ := client.ParseJsonPrivkey(cdc, v)
					acc = append(acc, txs.AddrCoin{addracc, btypes.NewInt(int64(100))})
				}
			}

			addr, privkey := client.ParseJsonPrivkey(cdc, privkeycreator)

			return stdTxCreateQSC(&cliCtx, cdc, caQsc, caBanker, addr, privkey, &acc, nonce, extrate, "", qscchainid, qoschainid, maxgas)
		},
	}

	cmd.Flags().String(flagExtrate, "1:280.0000", "extrate: qos/qscxxx")
	cmd.Flags().String(flagPathqsc, "", "path of CA(qsc)")
	cmd.Flags().String(flagPathbank, "", "path of CA(banker)")
	cmd.Flags().String(flagQscchainid, "", "chainid of qsc, used in tx")
	cmd.Flags().String(flagPrivkey, "", "private key of creator")
	cmd.Flags().String(flagQoschainid, "qos", "chainid of qos, used in tx")
	cmd.Flags().Int64(flagMaxgas, 0, "maxgas for txstd")
	cmd.Flags().Int64(flagNonce, 0, "nonce")
	cmd.Flags().Bool(flagInitAccount, true, "maxgas for txstd")

	return cmd
}

func QueryQscInfoCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queryqscinfo",
		Short: "query qsc info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			qscname := viper.GetString(flagQscname)
			key := fmt.Sprintf("qsc/[%s]", qscname)
			result, err := cliCtx.Client.ABCIQuery("store/base/key", []byte(key))
			if err != nil {
				return err
			}

			var qcpinfo mapper.QscInfo
			queryValueBz := result.Response.GetValue()
			if len(queryValueBz) == 0 {
				return errors.New(fmt.Sprintf("Chain (%s) not exist!", qscname))
			}

			err = cdc.UnmarshalBinaryBare(queryValueBz, &qcpinfo)
			if err != nil {
				panic(err)
			}

			addr, err := bech32.ConvertAndEncode("address", qcpinfo.BankAddr)
			fmt.Printf("qscname:%s \nbankeraddr:%s \nchainid: %s\n", qcpinfo.Qscname, addr, qcpinfo.ChainID)

			return nil
		},
	}

	cmd.Flags().String(flagQscname, "", "name of qsc")

	return cmd
}

func stdTxCreateQSC(ctx *context.CLIContext, cdc *amino.Codec, caqsc *[]byte, cabank *[]byte,
	createaddr btypes.Address, privkey crypto.PrivKey, accs *[]txs.AddrCoin, nonce int64,
	extrate string, dsp string, qscchainid string, qoschainid string, maxgas int64) error {

	var accsigner []*client.Accsign
	accsigner = append(accsigner, &client.Accsign{
		privkey,
		 nonce,
	})

	tx := txs.NewCreateQsc(cdc, caqsc, cabank, qscchainid, createaddr,accs, extrate, dsp)
	if tx == nil {
		return errors.New("createqsc error!")
	}

	txstd := client.GenTxStd(cdc, tx, qoschainid, maxgas, accsigner)
	if txstd == nil {
		return errors.New("gentxstd error!")
	}

	return client.BroadcastTxStd(ctx, cdc, txstd)
}