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
},
{
 "name": "banker",
 "addrbech32": "address1yyyyshlpl3rs27f69hdr6r2xyt9pqqdkvtdqxj",
 "address": "2108485FE1FC4705793A2DDA3D0D4622CA1001B6",
 "pub_key": {
  "type": "tendermint/PubKeyEd25519",
  "value": "8jzqwFn1U88f6ExtOMuPHFeyv3qHCKHu8BQdvZlorIc="
 },
 "priv_key": {
  "type": "tendermint/PrivKeyEd25519",
  "value": "xs3y4D3kOfNjzL2U5iGPOGUJgU1hCXzJofP0wrDu3ofyPOrAWfVTzx/oTG04y48cV7K/eocIoe7wFB29mWishw=="
 }
},
{
 "name": "acc1",
 "addrbech32": "address1zsqzn6wdecyar6c6nzem3e8qss2ws95csr8d0r",
 "address": "140029E9CDCE09D1EB1A98B3B8E4E08414E81698",
 "pub_key": {
  "type": "tendermint/PubKeyEd25519",
  "value": "DcCiMiGEAbPVWcRsg/RUuSeF72Z7Z1Gk7aBrBJe/t4Y="
 },
 "priv_key": {
  "type": "tendermint/PrivKeyEd25519",
  "value": "vAeIlHuWjvz/JmyGcB46ZHfCZdXCYuRogqxDgjYUM5wNwKIyIYQBs9VZxGyD9FS5J4XvZntnUaTtoGsEl7+3hg=="
 }
},
{
 "name": "acc2",
 "addrbech32": "address12as5uhdpf2y9zjkurx2l6dz8g98qkgryc4x355",
 "address": "57614E5DA14A88514ADC1995FD3447414E0B2064",
 "pub_key": {
  "type": "tendermint/PubKeyEd25519",
  "value": "va8Kjc8UCZUD4efbWaW0tRHSE+kxOTei+9rHvQVHvYs="
 },
 "priv_key": {
  "type": "tendermint/PrivKeyEd25519",
  "value": "31PlT2p6UICjV63dG7Nh3Mh9W0b+7FAEU+KOAxyNbZ29rwqNzxQJlQPh59tZpbS1EdIT6TE5N6L72se9BUe9iw=="
 }
},
{
 "name": "acc3",
 "addrbech32": "address1y9r4pjjnvkmpvw46de8tmwunw4nx4qnz2ax5ux",
 "address": "214750CA5365B6163ABA6E4EBDBB9375666A8262",
 "pub_key": {
  "type": "tendermint/PubKeyEd25519",
  "value": "O8YzKrGr4XNfdMQhs2bl3I2qVGoYu2O79gw+Ic1QgLE="
 },
 "priv_key": {
  "type": "tendermint/PrivKeyEd25519",
  "value": "9QkouVPl29N2v1lBO1+azUDqm38fAgs6d3Xo8DcnCus7xjMqsavhc190xCGzZuXcjapUahi7Y7v2DD4hzVCAsQ=="
 }
}
]`)

	var accs []KeyPV
	cdc.UnmarshalJSON(jsbyte, &accs)
	accret = []*AccAndPrivkey{}

	for _, ks := range accs {
		var qos int64 = 0

		switch ks.Name {
		case "creator":
			qos = 2000000
			break
		case "banker":
			qos = 10
			break
		default:
			qos = 0
		}

		addr, _ := btypes.GetAddrFromBech32(ks.AddrBech32)
		acc := account.QOSAccount{
			baccount.BaseAccount{
				addr,
				ks.PubKey,
				0},
			btypes.NewInt(qos),
			[]*types.QSC{},
		}
		accret = append(accret, &AccAndPrivkey{ks.PrivKey, acc})
	}

	return
}
