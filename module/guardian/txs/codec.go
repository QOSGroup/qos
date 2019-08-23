package txs

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	types.RegisterCodec(Cdc)
	RegisterCodec(Cdc)
	cert.RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxAddGuardian{}, "guardian/txs/TxAddGuardian", nil)
	cdc.RegisterConcrete(&TxDeleteGuardian{}, "guardian/txs/TxDeleteGuardian", nil)
	cdc.RegisterConcrete(&TxHaltNetwork{}, "guardian/txs/TxHaltNetwork", nil)
}
