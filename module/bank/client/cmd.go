package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qos/module/bank/txs"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		QueryLockAccountCmd(cdc),
	)
}

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCustomMaxGasCommands([]*cobra.Command{TransferCmd(cdc), InvariantCheckCmd(cdc)}, []int64{
		txs.GasForTransfer + bctypes.DefaultMaxGas,
		txs.GasForInvariantCheck + bctypes.DefaultMaxGas,
	})
}
