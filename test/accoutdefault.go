package test

import (
	baccount "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type KeyPV struct {
	Name       string         `json:"name"`
	AddrBech32 string         `json:"addrbech32"`
	Address    cmn.HexBytes   `json:"address"`
	PubKey     crypto.PubKey  `json:"pub_key"`
	PrivKey    crypto.PrivKey `json:"priv_key"`
}

type AccAndPrivkey struct {
	PrivKey crypto.PrivKey     `json:"priv_key"`
	Acc     account.QOSAccount `json:"acc"`
}

func InitKeys(cdc *go_amino.Codec) (accret []*AccAndPrivkey) {
	jsbyte := []byte(`[
{
 "name": "creator",
 "addrbech32": "address1auug9tjmkm00w36savxjywmj0sjccaam3pvjfu",
 "address": "EF3882AE5BB6DEF74750EB0D223B727C258C77BB",
 "pub_key": {
  "type": "tendermint/PubKeyEd25519",
  "value": "z25F0KGTCB/fkYItf7xXDSVB6fPHymQzRXRIn/HYO9w="
 },
 "priv_key": {
  "type": "tendermint/PrivKeyEd25519",
  "value": "rDwWppdGKFCv0wUxFqVID87GI/CFwLbL9p6EM6ug5brPbkXQoZMIH9+Rgi1/vFcNJUHp88fKZDNFdEif8dg73A=="
 }
}
]`)

	var accs KeyPV
	cdc.UnmarshalJSON(jsbyte, &accs)
	accret = []*AccAndPrivkey{}

	addr, _ := btypes.GetAddrFromBech32(accs.AddrBech32)
	acc := account.QOSAccount{
		baccount.BaseAccount{
			addr,
			accs.PubKey,
			0},
		btypes.NewInt(200000),
		[]*types.QSC{},
	}
	accret = append(accret, &AccAndPrivkey{accs.PrivKey, acc})

	return
}
