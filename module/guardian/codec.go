package guardian

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	types.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cert.RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxAddGuardian{}, "guardian/txs/TxAddGuardian", nil)
	cdc.RegisterConcrete(&TxDeleteGuardian{}, "guardian/txs/TxDeleteGuardian", nil)
}
