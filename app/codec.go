package app

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/go-amino"
)

func MakeCodec() *amino.Codec {
	cdc := baseabci.MakeQBaseCodec()
	RegisterCodec(cdc)
	return cdc
}

func RegisterCodec(cdc *amino.Codec) {
	txs.RegisterCodec(cdc)
	account.RegisterCodec(cdc)
}
