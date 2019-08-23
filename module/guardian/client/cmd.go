package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		QueryGuardianCmd(cdc),
		QueryGuardiansCmd(cdc),
	)
}

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands(
		AddGuardianCmd(cdc),
		DeleteGuardianCmd(cdc),
		HaltCmd(cdc),
	)
}

var (
	flagAddress     = "address"
	flagCreator     = "creator"
	flagDescription = "description"

	flagDeletedBy = "deleted-by"
)
