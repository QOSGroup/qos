package gov

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/baseabci"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
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
	cdc.RegisterInterface((*gtypes.ProposalContent)(nil), nil)
	cdc.RegisterConcrete(&gtypes.TextProposal{}, "gov/TextProposal", nil)
	cdc.RegisterConcrete(&gtypes.TaxUsageProposal{}, "gov/TaxUsageProposal", nil)
	cdc.RegisterConcrete(&gtypes.ParameterProposal{}, "gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&TxParameterChange{}, "gov/txs/TxParameterChange", nil)
	cdc.RegisterConcrete(&TxTaxUsage{}, "gov/txs/TxTaxUsage", nil)
	cdc.RegisterConcrete(&TxProposal{}, "gov/txs/TxProposal", nil)
	cdc.RegisterConcrete(&TxDeposit{}, "gov/txs/TxDeposit", nil)
	cdc.RegisterConcrete(&TxVote{}, "gov/txs/TxVote", nil)
}
