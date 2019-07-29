package init

import (
	"fmt"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	btypes "github.com/QOSGroup/qbase/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	ttypes "github.com/tendermint/tendermint/types"
)

var (
	chainId  string
	compound bool

	nValidators   int
	outputDir     string
	nodeDirPrefix string

	populatePersistentPeers bool
	hostnamePrefix          string
	startingIPAddress       string

	qcpRootCA string
	qscRootCA string
	accounts  string

	guardianAddresses string
)

const (
	nodeDirPerm  = 0755
	nodeFilePerm = 0644

	validatorBondTokens   = 1000
	validatorOwnerInitQOS = 1000000
	validatorOperatorFile = "priv_validator_owner.json"
)

func TestnetFileCmd(ctx *server.Context, cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a QOS testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:

	qosd testnet --chain-id=qostest --v=4 --o=./output --starting-ip-address=192.168.1.2 --genesis-accounts=address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config

			// accounts
			genesisAccounts := make([]*types.QOSAccount, 0)
			var err error
			if accounts != "" {
				genesisAccounts, err = types.ParseAccounts(accounts, viper.GetString(flagClientHome))
				if err != nil {
					return err
				}
			}

			// root ca
			var qcpPubKey crypto.PubKey
			if qcpRootCA != "" {
				err := cdc.UnmarshalJSON(cmn.MustReadFile(qcpRootCA), &qcpPubKey)
				if err != nil {
					return err
				}
			}
			var qscPubKey crypto.PubKey
			if qscRootCA != "" {
				err := cdc.UnmarshalJSON(cmn.MustReadFile(qscRootCA), &qscPubKey)
				if err != nil {
					return err
				}
			}

			// chainId
			if len(chainId) == 0 {
				chainId = "test-chain-" + cmn.RandStr(6)
			}

			// validators
			genTxDir := filepath.Join(outputDir, "gentxs")
			var nodeDirs []string
			var nodeIDs []string
			for i := 0; i < nValidators; i++ {
				nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
				nodeDir := filepath.Join(outputDir, nodeDirName)
				nodeDirs = append(nodeDirs, nodeDir)
				config.SetRoot(nodeDir)
				config.Moniker = nodeDirName

				err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				err = os.MkdirAll(filepath.Join(nodeDir, "data"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				nodeID, valPubKey, err := server.InitializeNodeValidatorFiles(config)
				nodeIDs = append(nodeIDs, nodeID)
				if err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}

				// create gentx file
				owner := ed25519.GenPrivKey()
				desc := stake.Description{Moniker: nodeDirName}
				txCreateValidator := stake.NewCreateValidatorTx(btypes.Address(owner.PubKey().Address()), valPubKey, validatorBondTokens, compound, desc)
				txStd := txs.NewTxStd(txCreateValidator, chainId, btypes.NewInt(1000000))
				sig, err := owner.Sign(txStd.BuildSignatureBytes(1, ""))
				if err != nil {
					return err
				}
				txStd.Signature = append(txStd.Signature, txs.Signature{
					Pubkey:    owner.PubKey(),
					Signature: sig,
					Nonce:     1,
				})
				writeSignedGenTx(cdc, genTxDir, nodeID, hostnameOrIP(i), txStd)

				genesisAccounts = append(genesisAccounts, types.NewQOSAccount(owner.PubKey().Address().Bytes(), btypes.NewInt(validatorOwnerInitQOS), nil))

				// write private key of validator owner
				ownerFile := filepath.Join(nodeDir, "config", validatorOperatorFile)
				ownerBz, _ := cdc.MarshalJSON(owner)
				cmn.MustWriteFile(ownerFile, ownerBz, nodeFilePerm)

				// write config file
				config.P2P.AddrBookStrict = false
				cfg.WriteConfigFile(filepath.Join(nodeDirs[i], "config", "config.toml"), config)
			}

			// guardians
			var guardians []guardian.Guardian
			if len(guardianAddresses) != 0 {
				addressArr := strings.Split(guardianAddresses, ",")
				for _, address := range addressArr {
					addr, err := btypes.GetAddrFromBech32(address)
					if err != nil {
						return err
					}
					guardians = append(guardians, *guardian.NewGuardian("genesis guardian", guardian.Genesis, addr, nil))
				}
			}
			guardianState := guardian.NewGenesisState(guardians)
			err = guardian.ValidateGenesis(guardianState)
			if err != nil {
				return err
			}

			appliedQOSAmount := btypes.ZeroInt()
			for _, account := range genesisAccounts {
				appliedQOSAmount = appliedQOSAmount.Add(account.QOS)
			}
			appState := app.GenesisState{
				Accounts:         genesisAccounts,
				MintData:         mint.DefaultGenesisState(),
				StakeData:        stake.DefaultGenesisState(),
				QCPData:          qcp.NewGenesisState(qcpPubKey, nil),
				QSCData:          qsc.NewGenesisState(qscPubKey, nil),
				DistributionData: distribution.DefaultGenesisState(),
				GovData:          gov.DefaultGenesisState(),
				GuardianData:     guardianState,
			}
			appState.MintData.AppliedQOSAmount = uint64(appliedQOSAmount.Int64())

			rawState, _ := cdc.MarshalJSON(appState)
			genDoc := &ttypes.GenesisDoc{
				ChainID:         chainId,
				GenesisTime:     time.Now(),
				ConsensusParams: defaultConsensusParams(),
				AppState:        rawState,
			}

			// collect gentxs, write genesis files and update config files
			for i := 0; i < nValidators; i++ {
				if err := genDoc.SaveAs(filepath.Join(nodeDirs[i], config.Genesis)); err != nil {
					_ = os.RemoveAll(outputDir)
					return err
				}
				config.SetRoot(nodeDirs[i])
				err = updateGenesisStateFromGenTxs(config, cdc, nodeIDs[i], genTxDir)
				if err != nil {
					return err
				}
			}

			fmt.Printf("Successfully initialized %v node directories\n", nValidators)
			return nil
		},
	}

	cmd.Flags().IntVar(&nValidators, "v", 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().StringVar(&outputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().StringVar(&nodeDirPrefix, "node-dir-prefix", "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().StringVar(&hostnamePrefix, "hostname-prefix", "node",
		"Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")
	cmd.Flags().StringVar(&startingIPAddress, "starting-ip-address", "",
		"Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:26656, ID1@192.168.0.2:26656, ...)")
	cmd.Flags().StringVar(&accounts, "genesis-accounts", "",
		"Add genesis accounts to genesis.json, eg: address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'")
	cmd.Flags().StringVar(&qcpRootCA, "qsc-root-ca", "", "Config pubKey of root CA for QCP")
	cmd.Flags().StringVar(&qscRootCA, "qcp-root-ca", "", "Config pubKey of root CA for QSC")
	cmd.Flags().StringVar(&chainId, "chain-id", "", "Chain ID")
	cmd.Flags().BoolVar(&compound, "compound", true, "whether the validator's income is calculated as compound interest, default: true")
	cmd.Flags().StringVar(&guardianAddresses, "guardians", "", "addresses for guardian. Multiple addresses separated by ','")
	cmd.Flags().String(flagClientHome, types.DefaultCLIHome, "directory for keybase")

	return cmd
}

func hostnameOrIP(i int) string {
	if startingIPAddress != "" {
		ip := net.ParseIP(startingIPAddress)
		ip = ip.To4()
		if ip == nil {
			fmt.Printf("%v: non ipv4 address\n", startingIPAddress)
			os.Exit(1)
		}

		for j := 0; j < i; j++ {
			ip[3]++
		}
		return ip.String()
	}

	return fmt.Sprintf("%s%d", hostnamePrefix, i)
}
