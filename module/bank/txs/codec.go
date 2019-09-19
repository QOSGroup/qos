package txs

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	types.RegisterCodec(Cdc)
	RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {

	cdc.RegisterConcrete(&TxTransfer{}, "bank/txs/TxTransfer", nil)
	cdc.RegisterConcrete(&TxInvariantCheck{}, "bank/txs/TxInvariantCheck", nil)
}
