package init

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/QOSGroup/qbase/client/context"
	clikeys "github.com/QOSGroup/qbase/client/keys"
	"github.com/QOSGroup/qbase/keys"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenTxCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gentx",
		Short: "Generate a genesis tx carrying a self delegation",
		Args:  cobra.NoArgs,
		Long: `This command is an alias of the 'gaiad tx create-validator' command'.
qosd gentx --moniker validatorName --owner ownerName --tokens 100
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			nodeID, _, err := server.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}
			genDoc, err := tmtypes.GenesisDocFromFile(config.GenesisFile())
			if err != nil {
				return err
			}
			genesisState := types.GenesisState{}
			if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
				return err
			}

			name := viper.GetString(flagMoniker)
			if len(name) == 0 {
				return errors.New("moniker is empty")
			}
			tokens, err := types.GetIntFromFlag(flagBondTokens, false)
			if err != nil {
				return err
			}
			logo := viper.GetString(flagLogo)
			website := viper.GetString(flagWebsite)
			details := viper.GetString(flagDetails)
			desc := stake.Description{
				name, logo, website, details,
			}

			commission, err := stake.BuildCommissionRates()
			if err != nil {
				return err
			}

			privValidator := privval.LoadOrGenFilePV(filepath.Join(viper.GetString(cli.HomeFlag), cfg.DefaultConfig().PrivValidatorKeyFile()),
				filepath.Join(viper.GetString(cli.HomeFlag), cfg.DefaultConfig().PrivValidatorKeyFile()))

			owner := viper.GetString(flagOwner)
			var info keys.Info
			if len(owner) == 0 {
				return errors.New("creator is empty")
			}
			clientHome := viper.GetString(flagClientHome)
			keybase, err := clikeys.GetKeyBaseFromDir(cliCtx, clientHome)
			if err != nil {
				return err
			}
			if strings.HasPrefix(owner, btypes.GetAddressConfig().GetBech32AccountAddrPrefix()) {
				addr, err := btypes.AccAddressFromBech32(owner)
				if err != nil {
					return err
				}
				info, err = keybase.GetByAddress(addr)
				if err != nil {
					return err
				}
			} else {
				info, err = keybase.Get(owner)
				if err != nil {
					return err
				}
			}

			isCompound := viper.GetBool(flagCompound)

			delegations := viper.GetString(flagDelegations)
			var delegationInfos []stake.DelegationInfo
			if len(delegations) != 0 {
				delegators, err := types.ParseAccounts(delegations, clientHome)
				if err != nil {
					return err
				}
				totalBondTokens := btypes.ZeroInt()
				delegatorMap := map[string]bool{}
				for _, delegator := range delegators {
					if _, ok := delegatorMap[delegator.AccountAddress.String()]; ok {
						return errors.New("duplicate delegator in delegations")
					}
					err = validGenesisAccount(cdc, genesisState, delegator.AccountAddress, delegator.QOS)
					if err != nil {
						return err
					}
					totalBondTokens = totalBondTokens.Add(delegator.QOS)
					delegatorMap[delegator.AccountAddress.String()] = true
					delegationInfos = append(delegationInfos, stake.NewDelegationInfo(delegator.AccountAddress, btypes.ValAddress(info.GetAddress()), delegator.QOS, isCompound))
				}

				if !totalBondTokens.Equal(tokens) {
					return errors.New("tokens must equal sum(amount) of delegations")
				}
			} else {
				err = validGenesisAccount(cdc, genesisState, info.GetAddress(), tokens)
				if err != nil {
					return err
				}
			}

			itx := stake.NewCreateValidatorTx(info.GetAddress(), privValidator.GetPubKey(), tokens, isCompound, desc, *commission, delegationInfos)
			if err != nil {
				return err
			}
			txStd := txs.NewTxStd(itx, genDoc.ChainID, btypes.NewInt(1000000))
			sigdata := txStd.BuildSignatureBytes(1, "")
			pass, err := clikeys.GetPassphrase(cliCtx, info.GetName())
			if err != nil {
				panic(fmt.Sprintf("Get %s Passphrase error: %s", info.GetAddress(), err.Error()))
			}
			sig, pubkey, err := keybase.Sign(info.GetName(), pass, sigdata)
			if err != nil {
				panic(err.Error())
			}
			txStd.Signature = append(txStd.Signature, txs.Signature{
				Pubkey:    pubkey,
				Signature: sig,
				Nonce:     1,
			})

			if err := writeSignedGenTx(cdc, filepath.Join(config.RootDir, "config", "gentx"), nodeID, viper.GetString(flagIP), txStd); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagMoniker, "", "name for validator")
	cmd.Flags().String(flagOwner, "", "keystore name or account address for validator's owner")
	cmd.Flags().String(flagBondTokens, "0", "bond tokens amount")
	cmd.Flags().Bool(flagCompound, false, "as a self-delegator, whether the income is calculated as compound interest")
	cmd.Flags().String(flagClientHome, types.DefaultCLIHome, "directory for keybase")
	cmd.Flags().String(flagIP, "127.0.0.1", "ip of your node")
	cmd.Flags().String(flagLogo, "", "The optional logo link")
	cmd.Flags().String(flagWebsite, "", "The validator's (optional) website")
	cmd.Flags().String(flagDetails, "", "The validator's (optional) details")
	cmd.Flags().String(flagCommissionRate, stake.DefaultCommissionRate, "The initial commission rate percentage")
	cmd.Flags().String(flagCommissionMaxRate, stake.DefaultCommissionMaxRate, "The maximum commission rate percentage")
	cmd.Flags().String(flagCommissionMaxChangeRate, stake.DefaultCommissionMaxChangeRate, "The maximum commission change rate percentage (per day)")
	cmd.Flags().String(flagDelegations, "", "init delegations, 'address1,10000QOS,address2,10000QOS'")

	cmd.MarkFlagRequired(flagMoniker)
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}

func validGenesisAccount(cdc *amino.Codec, genesisState types.GenesisState, address btypes.AccAddress, amount btypes.BigInt) error {
	accountIsInGenesis := false

	var bankState bank.GenesisState
	cdc.MustUnmarshalJSON(genesisState[bank.ModuleName], &bankState)
	for _, acc := range bankState.Accounts {
		if acc.AccountAddress.Equals(address) {

			if !acc.EnoughOfQOS(amount) {
				return fmt.Errorf(
					"account %v is in genesis, but it only has %v QOS available to stake, not %v",
					address.String(), acc.QOS, amount,
				)
			}
			accountIsInGenesis = true
			break
		}
	}

	if accountIsInGenesis {
		return nil
	}

	return fmt.Errorf("account %s is not in genesis accounts, you can use `qosd add-genesis-accounts` to add it", address)
}

func writeSignedGenTx(cdc *amino.Codec, genTxDir, nodeID, ip string, tx *txs.TxStd) error {
	if err := common.EnsureDir(genTxDir, 0700); err != nil {
		return err
	}
	genTx := filepath.Join(genTxDir, fmt.Sprintf("%s@%s.json", nodeID, ip))
	outputFile, err := os.OpenFile(genTx, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	json, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(outputFile, "%s\n", json)
	return err
}
