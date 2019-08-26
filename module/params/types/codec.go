package types

import (
	"github.com/QOSGroup/qbase/baseabci"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*qtypes.ParamSet)(nil), nil)
}
