package client

import (
	"fmt"
	"github.com/QOSGroup/qbase/client"
	"github.com/QOSGroup/qbase/client/context"
	btypes "github.com/QOSGroup/qbase/types"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"regexp"
	"strconv"
	"strings"
)

// Get default chain-id.
func GetDefaultChainId() (string, error) {
	chainId := viper.GetString(client.FlagChainID)
	if len(chainId) == 0 {
		chainId, err := btypes.GetChainID(app.DefaultNodeHome)
		if err != nil {
			return "", errors.New("use --chain-id flag to specify your chain")
		}
		return chainId, nil
	}
	return chainId, nil
}

// Parse QOS and QSCs from string
// str example : 100qos,100qstar
func ParseCoins(str string) (btypes.BigInt, types.QSCs, error) {
	reDnm := `[[:alpha:]][[:alnum:]]{2,15}`
	reAmt := `[[:digit:]]+`
	reSpc := `[[:space:]]*`
	reCoin := regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))

	arr := strings.Split(str, ",")
	qos := btypes.ZeroInt()
	qscs := types.QSCs{}
	for _, q := range arr {
		coin := reCoin.FindStringSubmatch(q)
		coin[2] = strings.TrimSpace(coin[2])
		amount, err := strconv.ParseInt(strings.TrimSpace(coin[1]), 10, 64)
		if err != nil {
			return btypes.ZeroInt(), nil, err
		}
		if strings.ToLower(coin[2]) == "qos" {
			qos = btypes.NewInt(amount)
		} else {
			qscs = append(qscs, &types.QSC{
				coin[2],
				btypes.NewInt(amount),
			})
		}

	}

	return qos, qscs, nil
}

type Accsign struct {
	Privkey crypto.PrivKey `json:"privkey"`
	Nonce   int64          `json:"nonce"`
}

func GenTxStd(cdc *amino.Codec, itx btxs.ITx, chainid string, maxgas int64, accsigner []*Accsign) (txstd *btxs.TxStd) {
	if len(accsigner) == 0 {
		return nil
	}

	tx := btxs.NewTxStd(itx, chainid, btypes.NewInt(maxgas))

	for _, acc := range accsigner {
		signdata, _ := tx.SignTx(acc.Privkey, acc.Nonce)
		tx.Signature = append(tx.Signature, btxs.Signature{
			Pubkey:    acc.Privkey.PubKey(),
			Signature: signdata,
			Nonce:     acc.Nonce,
		})
	}

	return tx
}

func BroadcastTxStd(ctx *context.CLIContext, cdc *amino.Codec, txstd *btxs.TxStd) error {
	tx, err := cdc.MarshalBinaryBare(txstd)
	if err != nil {
		return err
	}

	result, err := ctx.BroadcastTxSync(tx)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("tx result:  %v", result))

	return nil
}

func ParseJsonPrivkey(cdc *amino.Codec, jsprivkey string) (addr []byte, privkey ed25519.PrivKeyEd25519) {
	privstr := fmt.Sprintf(` {
 			 	"type": "tendermint/PrivKeyEd25519",
 			 	"value": "%s"
 			}`, jsprivkey)
	err := cdc.UnmarshalJSON([]byte(privstr), &privkey)
	if err != nil {
		panic("parse json privkey error!")
	}

	addr = privkey.PubKey().Address()

	return
}

func Cmdcheck(mode string, option...string) bool {
	switch mode {
	case "txcreateqsc":
		//*privkeybank, *pathbank, *pathqsc
		for idx,v := range option {
			switch idx {
			case 0:
				if len(v) < 64 {
					fmt.Print("invalide private key!")
					return false
				}
			case 1,2:
				if len(v) < 6 {
					fmt.Print("invalide path!")
					return false
				}
			}
		}
	default:
		fmt.Printf("cmd (%s) not support!", mode)
		return false
	}

	return true
}
