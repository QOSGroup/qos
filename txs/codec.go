package txs

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
	cdc.RegisterConcrete(&TxCreateQSC{}, "qos/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&TxIssueQsc{}, "qos/txs/TxIssueQsc", nil)
	cdc.RegisterConcrete(&TransferTx{}, "qos/txs/TransferTx", nil)
	cdc.RegisterConcrete(&ApproveCreateTx{}, "qos/txs/ApproveCreateTx", nil)
	cdc.RegisterConcrete(&ApproveIncreaseTx{}, "qos/txs/ApproveIncreaseTx", nil)
	cdc.RegisterConcrete(&ApproveDecreaseTx{}, "qos/txs/ApproveDecreaseTx", nil)
	cdc.RegisterConcrete(&ApproveUseTx{}, "qos/txs/ApproveUseTx", nil)
	cdc.RegisterConcrete(&ApproveCancelTx{}, "qos/txs/ApproveCancelTx", nil)
}
