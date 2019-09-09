package client

import (
	"errors"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

// 查询锁定-释放账户信息
func QueryLockAccountCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock-account",
		Short: "query lock account",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			result, err := cliCtx.Client.ABCIQuery("store/acc/key", mapper.LockInfoKey)
			if err != nil {
				return err
			}

			if len(result.Response.GetValue()) == 0 {
				return errors.New("no lock account")
			}

			var info types.LockInfo
			err = cdc.UnmarshalBinaryBare(result.Response.GetValue(), &info)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(info)
		},
	}

	return cmd
}
