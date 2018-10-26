package main

import (
	"flag"
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/app"
	"github.com/tendermint/tendermint/rpc/client"
)

// 查询账户
//-addr=address1k0m8ucnqug974maa6g36zw7g2wvfd4sug6uxay
func main() {
	cdc := app.MakeCodec()

	addr := flag.String("addr", "", "input account addr(bech32)")

	flag.Parse()

	http := client.NewHTTP("tcp://127.0.0.1:26657", "/websocket")

	if *addr == "" {
		panic("usage: -m=acc -addr=xxx")
	}
	address, _ := btypes.GetAddrFromBech32(*addr)
	key := bacc.AddressStoreKey(address)
	result, err := http.ABCIQuery("/store/acc/key", key)
	if err != nil {
		panic(err)
	}

	queryValueBz := result.Response.GetValue()
	var acc *account.QOSAccount
	cdc.UnmarshalBinaryBare(queryValueBz, &acc)

	json, _ := cdc.MarshalJSON(acc)
	fmt.Println(fmt.Sprintf("query addr is %s = %s", *addr, string(json)))
}
