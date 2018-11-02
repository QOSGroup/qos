package app

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/server/config"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
)

// QOS初始状态
type GenesisState struct {
	CAPubKey crypto.PubKey `json:"ca_pub_key"`
	Accounts []*account.QOSAccount `json:"accounts"`
}

func QOSAppInit() server.AppInit {
	return server.AppInit{
		AppGenTx:    QOSAppGenTx,
		AppGenState: QOSAppGenState,
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
	addr, _, secret, err = GenerateCoinKey()
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

	mm := map[string]string{"secret": secret}
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
func GenerateCoinKey() (addr btypes.Address, pubkey crypto.PubKey, secret string, err error) {
	//ed25519
	addr, _ = btypes.GetAddrFromBech32("address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay")
	secret = "0xa328891040ae9b773bcd30005235f99a8d62df03a89e4f690f9fa03abb1bf22715fc9ca05613f2d8061492e9f8149510b5b67d340d199ff24f34c85dbbbd7e0df780e9a6cc"
	priHex, _ := hex.DecodeString(secret[2:])
	var priKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(priHex, &priKey)
	pubkey = priKey.PubKey()
	return
}
