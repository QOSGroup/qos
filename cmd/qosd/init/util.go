package init

import (
	"fmt"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"io/ioutil"
	"strings"

	baccount "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
)

func loadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc tmtypes.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

// Parse accounts from string
// address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'
func ParseAccounts(str string) ([]*account.QOSAccount, error) {
	accounts := make([]*account.QOSAccount, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		if ti == "" {
			continue
		}

		addrAndCoins := strings.Split(ti, ",")
		if len(addrAndCoins) < 2 {
			return nil, fmt.Errorf("`%s` not match rules", ti)
		}

		addr, err := btypes.GetAddrFromBech32(addrAndCoins[0])
		if err != nil {
			return nil, err
		}
		qos, qscs, err := types.ParseCoins(strings.Join(addrAndCoins[1:], ","))
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account.QOSAccount{
			BaseAccount: baccount.BaseAccount{
				AccountAddress: addr,
				Nonce:          0,
			},
			QOS:  qos,
			QSCs: qscs,
		})
	}

	return accounts, nil
}
