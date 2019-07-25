package txs

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	types.RegisterCodec(Cdc)
	RegisterCodec(Cdc)
	cert.RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxCreateQSC{}, "qsc/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&TxIssueQSC{}, "qsc/txs/TxIssueQSC", nil)
}
