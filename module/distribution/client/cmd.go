package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands()
}

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		queryValidatorPeriodCommand(cdc),
		queryDelegatorIncomeInfoCommand(cdc),
		queryCommunityFeePoolCommand(cdc),
	)
}
