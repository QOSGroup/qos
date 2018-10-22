package txs

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/account"
	"github.com/tendermint/go-amino"
)

var cdc = makeCodec()

func makeCodec() *amino.Codec {
	cdc := baseabci.MakeQBaseCodec()
	account.RegisterCodec(cdc)
	RegisterCodec(cdc)
	return cdc
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxCreateQSC{}, "qos/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&TxIssueQsc{}, "qos/txs/TxIssueQsc", nil)
	cdc.RegisterConcrete(&TxTransform{}, "qos/txs/TxTransform", nil)
	cdc.RegisterConcrete(&ApproveCreateTx{}, "qos/txs/ApproveCreateTx", nil)
	cdc.RegisterConcrete(&ApproveIncreaseTx{}, "qos/txs/ApproveIncreaseTx", nil)
	cdc.RegisterConcrete(&ApproveDecreaseTx{}, "qos/txs/ApproveDecreaseTx", nil)
	cdc.RegisterConcrete(&ApproveUseTx{}, "qos/txs/ApproveUseTx", nil)
	cdc.RegisterConcrete(&ApproveCancelTx{}, "qos/txs/ApproveCancelTx", nil)
}
