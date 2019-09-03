package init

import (
	"errors"
	"fmt"

	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"time"

	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
	tmtypes "github.com/tendermint/tendermint/types"
)

func AddGenesisAccount(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-accounts [accounts]",
		Short: "Add genesis accounts to genesis.json",
		Long: `add-genesis-accounts [accounts] will add [accounts] into app_state.
Multiple accounts separated by ';'.

Example:

	qosd add-genesis-accounts "qosacc128xejc2f9lp7pjm56j8sahcvkxan9lqfqp6azy,1000000qos;qosacc12kjmpgyg23l7axhzzne33jmd0r9y083w6mpa33,100000qos"
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			accounts, err := types.ParseAccounts(args[0], viper.GetString(flagClientHome))

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			var appState types.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			var bankState bank.GenesisState
			cdc.MustUnmarshalJSON(appState[bank.ModuleName], &bankState)

			var mintState mint.GenesisState
			cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintState)

			for _, v := range bankState.Accounts {
				for _, acc := range accounts {
					if acc.AccountAddress.Equals(v.GetAddress()) {
						return fmt.Errorf("addr: %s has already exists", v.AccountAddress.String())
					}
				}
			}

			bankState.Accounts = append(bankState.Accounts, accounts...)
			for _, acc := range accounts {
				mintState.AppliedQOSAmount = mintState.AppliedQOSAmount.Add(acc.QOS)
			}

			appState[bank.ModuleName] = cdc.MustMarshalJSON(bankState)
			appState[mint.ModuleName] = cdc.MustMarshalJSON(mintState)

			rawMessage, _ := cdc.MarshalJSON(appState)
			genDoc.AppState = rawMessage

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(config.GenesisFile())
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "directory for node's data and config files")
	cmd.Flags().String(flagClientHome, types.DefaultCLIHome, "directory for keybase")

	return cmd
}

func AddLockAccount(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-lock-account",
		Short: "Add lock account to genesis.json",
		Long: `add-lock-account will add locked account into app_state.

Example:

	qosd add-lock-account --receiver qosacc1lly0audg7yem8jt77x2jc6wtrh7v96hgve8fh8 --total-amount 10000000000000 --released-amount 1000000000000 --release-time '2023-10-20T00:00:00Z' --release-interval 30 --release-times 10"
	`,
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			receiverStr := viper.GetString(flagReceiver)
			if len(receiverStr) == 0 {
				return errors.New("empty receiver")
			}
			receiver, err := btypes.AccAddressFromBech32(receiverStr)
			if err != nil {
				return errors.New("invalid receiver address")
			}

			totalAmount, err := types.GetIntFromFlag(flagTotalAmount, false)
			if err != nil {
				return err
			}
			releasedAmount, err := types.GetIntFromFlag(flagReleasedAmount, false)
			if err != nil {
				return err
			}
			if !totalAmount.GT(releasedAmount) {
				return errors.New("released-amount must lt total-amount")
			}
			releaseInterval := viper.GetInt64(flagReleaseInterval)
			if releaseInterval <= 0 {
				return errors.New("release-interval must be positive")
			}
			releaseTimes := viper.GetInt64(flagReleaseTimes)
			if releaseTimes <= 0 {
				return errors.New("release-times must be positive")
			}
			releaseTime := viper.GetTime(flagReleaseTime).UTC()
			if releaseTime.Before(time.Now().UTC()) {
				return errors.New("release-time must after now")
			}

			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}

			var appState types.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			var bankState bank.GenesisState
			cdc.MustUnmarshalJSON(appState[bank.ModuleName], &bankState)
			if bankState.LockInfo != nil {
				return errors.New("lock account already set")
			}

			var mintState mint.GenesisState
			cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintState)

			lockedAddress := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			lockAccount := types.NewQOSAccount(lockedAddress, totalAmount.Sub(releasedAmount), nil)
			lockInfo := bank.NewLockInfo(lockedAddress, receiver, totalAmount, releasedAmount, releaseTime, releaseInterval, releaseTimes)
			bankState.Accounts = append(bankState.Accounts, lockAccount)
			bankState.LockInfo = &lockInfo

			mintState.AppliedQOSAmount = mintState.AppliedQOSAmount.Add(lockAccount.QOS)

			appState[bank.ModuleName] = cdc.MustMarshalJSON(bankState)
			appState[mint.ModuleName] = cdc.MustMarshalJSON(mintState)

			rawMessage, _ := cdc.MarshalJSON(appState)
			genDoc.AppState = rawMessage

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(config.GenesisFile())
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "directory for node's data and config files")
	cmd.Flags().String(flagClientHome, types.DefaultCLIHome, "directory for keybase")
	cmd.Flags().String(flagReceiver, "", "keybase name or address to hold the released QOS")
	cmd.Flags().String(flagTotalAmount, "0", "total QOS amount locked")
	cmd.Flags().String(flagReleasedAmount, "0", "total released QOS amount")
	cmd.Flags().String(flagReleaseTime, "0", "first release time(UTC)")
	cmd.Flags().String(flagReleaseInterval, "0", "release interval, in days")
	cmd.Flags().String(flagReleaseTimes, "0", "release times")

	cmd.MarkFlagRequired(flagReceiver)
	cmd.MarkFlagRequired(flagTotalAmount)
	cmd.MarkFlagRequired(flagReleasedAmount)
	cmd.MarkFlagRequired(flagReleaseTime)
	cmd.MarkFlagRequired(flagReleaseInterval)
	cmd.MarkFlagRequired(flagReleaseTimes)

	return cmd
}
