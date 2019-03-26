package gov

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
)

const (
	MaxTitleLen       = 200
	MaxDescriptionLen = 1000
)

type TxProposal struct {
	Title          string              `json:"title"`           //  Title of the proposal
	Description    string              `json:"description"`     //  Description of the proposal
	ProposalType   gtypes.ProposalType `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       btypes.Address      `json:"proposer"`        //  Address of the proposer
	InitialDeposit uint64              `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
}

var _ txs.ITx = (*TxProposal)(nil)

func (tx TxProposal) ValidateData(ctx context.Context) error {
	if len(tx.Title) == 0 || len(tx.Title) > MaxTitleLen {
		return ErrInvalidInput("invalid title")
	}
	if len(tx.Description) == 0 || len(tx.Description) > MaxDescriptionLen {
		return ErrInvalidInput("invalid description")
	}
	if !gtypes.ValidProposalType(tx.ProposalType) {
		return ErrInvalidInput("unknown proposal type")
	}

	govMapper := GetGovMapper(ctx)
	if tx.InitialDeposit < govMapper.GetDepositParams().MinDeposit/3 {
		return ErrInvalidInput("initial deposit is too small")
	}
	return nil
}

func (tx TxProposal) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := GetGovMapper(ctx)

	textContent := gtypes.NewTextProposal(tx.Title, tx.Description, tx.InitialDeposit)
	_, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	return
}

func (tx TxProposal) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Proposer}
}

func (tx TxProposal) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxProposal) GetGasPayer() btypes.Address {
	return tx.Proposer
}

func (tx TxProposal) GetSignData() (ret []byte) {
	ret = append(ret, tx.Title...)
	ret = append(ret, tx.Description...)
	ret = append(ret, byte(tx.ProposalType))
	ret = append(ret, tx.Proposer...)
	ret = append(ret, types.Uint64ToBigEndian(tx.InitialDeposit)...)

	return
}
