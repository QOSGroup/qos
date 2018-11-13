package qsc

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
	CrtRegisterCodec(cdc)
	cdc.RegisterConcrete(&TxCreateQSC{}, "qos/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&TxIssueQsc{}, "qos/txs/TxIssueQsc", nil)
}
