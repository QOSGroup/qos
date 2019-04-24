package types

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/keys"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = go_amino.NewCodec()

// 包初始化，注册codec
func init() {
	cryptoAmino.RegisterAmino(cdc)
	cdc.RegisterInterface((*account.Account)(nil), nil)
	keys.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

// 为包内定义结构注册codec
func RegisterCodec(cdc *go_amino.Codec) {
	cdc.RegisterConcrete(&QOSAccount{}, "qos/types/QOSAccount", nil)
	cdc.RegisterConcrete(&Fraction{}, "qos/types/Fraction", nil)
	cdc.RegisterConcrete(&Dec{}, "qos/types/Dec", nil)
}
