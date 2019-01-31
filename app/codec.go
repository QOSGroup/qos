package app

import (
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/module"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = go_amino.NewCodec()

// 包初始化，注册codec
func init() {
	cryptoAmino.RegisterAmino(cdc)
	cdc.RegisterInterface((*bacc.Account)(nil), nil)
	RegisterCodec(cdc)
}

func MakeCodec() *amino.Codec {
	cdc := baseabci.MakeQBaseCodec()
	RegisterCodec(cdc)
	return cdc
}

func RegisterCodec(cdc *amino.Codec) {
	noPaincRegisterInterface(cdc)
	module.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	cert.RegisterCodec(cdc)
}

func noPaincRegisterInterface(cdc *go_amino.Codec) {
	defer func() {
		if r := recover(); r != nil {
			//nothing
		}
	}()
	cdc.RegisterInterface((*interface{})(nil), nil)
}
