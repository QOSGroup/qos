package txs

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	account.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	approve.RegisterCodec(cdc)
	cdc.RegisterConcrete(&CA{}, "qos/txs/ca", nil)
	cdc.RegisterConcrete(&TxCreateQSC{}, "qos/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&TxIssueQsc{}, "qos/txs/TxIssueQsc", nil)
	cdc.RegisterConcrete(&TransferTx{}, "qos/txs/TransferTx", nil)
}
