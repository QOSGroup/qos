package stake

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxCreateValidator{}, "qos/txs/TxCreateValidator", nil)
	cdc.RegisterConcrete(&TxRevokeValidator{}, "qos/txs/TxRevokeValidator", nil)
	cdc.RegisterConcrete(&TxActiveValidator{}, "qos/txs/TxActiveValidator", nil)
}
