package params

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/module/params/types"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*types.ParamSet)(nil), nil)
}
