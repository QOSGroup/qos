package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/module/qsc/txs"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		QueryQSCCmd(cdc),
		QueryQSCsCmd(cdc),
	)
}

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCustomMaxGasCommands([]*cobra.Command{
		CreateQSCCmd(cdc),
		IssueQSCCmd(cdc),
	}, []int64{
		txs.GasForCreateQSC + bctypes.DefaultMaxGas,
		txs.GasForIssueQSC + bctypes.DefaultMaxGas,
	})

}
