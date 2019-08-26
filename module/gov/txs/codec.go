package txs

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/module/gov/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	qtypes.RegisterCodec(Cdc)
	RegisterCodec(Cdc)
	cert.RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*types.ProposalContent)(nil), nil)
	cdc.RegisterConcrete(&types.TextProposal{}, "gov/TextProposal", nil)
	cdc.RegisterConcrete(&types.TaxUsageProposal{}, "gov/TaxUsageProposal", nil)
	cdc.RegisterConcrete(&types.ParameterProposal{}, "gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&types.ModifyInflationProposal{}, "gov/ModifyInflationProposal", nil)
	cdc.RegisterConcrete(&types.SoftwareUpgradeProposal{}, "gov/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&TxParameterChange{}, "gov/txs/TxParameterChange", nil)
	cdc.RegisterConcrete(&TxTaxUsage{}, "gov/txs/TxTaxUsage", nil)
	cdc.RegisterConcrete(&TxModifyInflation{}, "gov/txs/TxModifyInflation", nil)
	cdc.RegisterConcrete(&TxSoftwareUpgrade{}, "gov/txs/TxSoftwareUpgrade", nil)
	cdc.RegisterConcrete(&TxProposal{}, "gov/txs/TxProposal", nil)
	cdc.RegisterConcrete(&TxDeposit{}, "gov/txs/TxDeposit", nil)
	cdc.RegisterConcrete(&TxVote{}, "gov/txs/TxVote", nil)
	cdc.RegisterConcrete(&types.Params{}, "gov/params", nil)
}
