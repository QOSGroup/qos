package transfer

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/account"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	account.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {

	cdc.RegisterConcrete(&TransferTx{}, "qos/txs/TransferTx", nil)
}
