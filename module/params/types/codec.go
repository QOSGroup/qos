package types

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*ParamSet)(nil), nil)
}
