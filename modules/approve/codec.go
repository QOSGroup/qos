package approve

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	types.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxCreateApprove{}, "qos/txs/TxCreateApprove", nil)
	cdc.RegisterConcrete(&TxIncreaseApprove{}, "qos/txs/TxIncreaseApprove", nil)
	cdc.RegisterConcrete(&TxDecreaseApprove{}, "qos/txs/TxDecreaseApprove", nil)
	cdc.RegisterConcrete(&TxUseApprove{}, "qos/txs/TxUseApprove", nil)
	cdc.RegisterConcrete(&TxCancelApprove{}, "qos/txs/TxCancelApprove", nil)
}
