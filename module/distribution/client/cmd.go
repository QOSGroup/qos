package distribution

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		queryValidatorPeriodCommand(cdc),
		queryDelegatorIncomeInfoCommand(cdc),
	)
}
