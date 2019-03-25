package gov

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	types.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cert.RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&TxProposal{}, "gov/TxProposal", nil)
	cdc.RegisterConcrete(&TxDeposit{}, "gov/TxDeposit", nil)
	cdc.RegisterConcrete(&TxVote{}, "gov/TxVote", nil)
}
