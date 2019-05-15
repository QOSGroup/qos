package eco

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/module/eco/types"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&types.DistributionParams{}, "distribution/params", nil)
	cdc.RegisterConcrete(&types.StakeParams{}, "stake/params", nil)
	cdc.RegisterConcrete(&types.MintParams{}, "mint/params", nil)
}
