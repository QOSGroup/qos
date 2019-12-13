package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/module/qcp/txs"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCustomMaxGasCommands([]*cobra.Command{
		InitQCPCmd(cdc),
	}, []int64{
		txs.GasForCreateQCP + bctypes.DefaultMaxGas,
	})
}

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands()
}
