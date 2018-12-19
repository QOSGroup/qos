package staking

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {

}
