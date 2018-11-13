package txs

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/txs/validator"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	account.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	approve.RegisterCodec(cdc)
	CrtRegisterCodec(cdc)
	cdc.RegisterConcrete(&Issuer{}, "qos/txs/Issuer", nil)
	cdc.RegisterConcrete(&TxCreateQSC{}, "qos/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&TxIssueQsc{}, "qos/txs/TxIssueQsc", nil)
	cdc.RegisterConcrete(&TransferTx{}, "qos/txs/TransferTx", nil)
	validator.RegisterCodec(cdc)
}
