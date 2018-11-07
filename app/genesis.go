package app

import (
	"encoding/json"
	"errors"
	"fmt"
	clikeys "github.com/QOSGroup/qbase/client/keys"
	bkeys "github.com/QOSGroup/qbase/keys"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/server/config"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/spf13/pflag"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
	"os"
	"path/filepath"
)

const (
	DefaultAccountName = "Arya"
	DefaultAccountPass = "12345678"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.qoscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.qosd")
)

// QOS初始状态
type GenesisState struct {
	CAPubKey crypto.PubKey         `json:"ca_pub_key"`
	Accounts []*account.QOSAccount `json:"accounts"`
}

func QOSAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(server.FlagName, "", "validator moniker, required")
	fsAppGenTx.String(server.FlagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(server.FlagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         QOSAppGenTx,
		AppGenState:      QOSAppGenState,
	}
}

type QOSGenTx struct {
	Addr btypes.Address `json:"addr"`
}

// Generate a genesis transaction
func QOSAppGenTx(cdc *amino.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	var addr btypes.Address
	var secret string
	addr, secret, err = GenerateCoinKey(cdc, genTxConfig.CliRoot)
	if err != nil {
		return
	}

	var bz []byte
	simpleGenTx := QOSGenTx{addr}
	bz, err = cdc.MarshalJSON(simpleGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	mm := map[string]string{"name": DefaultAccountName, "pass": DefaultAccountPass, "address": addr.String(), "secret": secret}
	bz, err = cdc.MarshalJSON(mm)
	if err != nil {
		return
	}
	cliPrint = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  10,
	}
	return
}

// app_state初始配置项生成
func QOSAppGenState(cdc *amino.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	if len(appGenTxs) != 1 {
		err = errors.New("must provide a single genesis transaction")
		return
	}

	var genTx QOSGenTx
	err = cdc.UnmarshalJSON(appGenTxs[0], &genTx)
	if err != nil {
		return
	}

	addr, _ := cdc.MarshalJSON(genTx.Addr)
	appState = json.RawMessage(fmt.Sprintf(`{
		"ca_pub_key": {
			"type": "tendermint/PubKeyEd25519",
			"value": "Py/hnnJJKXkWLAx/g+bMt9WDLGDLLNt0l4OXezIEuyE="
		},
		"accounts": [{
			"base_account": {
				"account_address": %s
			},
			"qos": "100000000",
			"qscs": [{
				"coin_name": "qstar",
				"amount": "100000000"
			}]
		}]
	}`, addr))
	return
}

// 默认地址
func GenerateCoinKey(cdc *amino.Codec, clientRoot string) (addr btypes.Address, mnemonic string, err error) {
	db, err := dbm.NewGoLevelDB(clikeys.KeyDBName, filepath.Join(clientRoot, "keys"))
	if err != nil {
		return btypes.Address([]byte{}), "", err
	}
	keybase := bkeys.New(db, cdc)

	info, secret, err := keybase.CreateEnMnemonic(DefaultAccountName, DefaultAccountPass)
	if err != nil {
		return btypes.Address([]byte{}), "", err
	}

	addr = btypes.Address(info.GetPubKey().Address())
	return addr, secret, nil
}
