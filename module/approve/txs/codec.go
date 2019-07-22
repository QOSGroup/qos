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
	cdc.RegisterConcrete(&TxCreateApprove{}, "approve/txs/TxCreateApprove", nil)
	cdc.RegisterConcrete(&TxIncreaseApprove{}, "approve/txs/TxIncreaseApprove", nil)
	cdc.RegisterConcrete(&TxDecreaseApprove{}, "approve/txs/TxDecreaseApprove", nil)
	cdc.RegisterConcrete(&TxUseApprove{}, "approve/txs/TxUseApprove", nil)
	cdc.RegisterConcrete(&TxCancelApprove{}, "approve/txs/TxCancelApprove", nil)
}
