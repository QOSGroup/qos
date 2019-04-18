package gov

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
)

type TxDeposit struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Depositor  btypes.Address `json:"depositor"`   // Address of the depositor
	Amount     uint64         `json:"amount"`      // Percent of QOS to add to the proposal's deposit
}

func NewTxDeposit(proposalID uint64, depositor btypes.Address, amount uint64) *TxDeposit {
	return &TxDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
}

var _ txs.ITx = (*TxDeposit)(nil)

func (tx TxDeposit) ValidateData(ctx context.Context) error {
	if len(tx.Depositor) == 0 {
		return ErrInvalidInput("depositor is empty")
	}

	if tx.Amount == 0 {
		return ErrInvalidInput("amount of deposit is zero")
	}

	proposal, ok := GetGovMapper(ctx).GetProposal(tx.ProposalID)
	if !ok {
		return ErrUnknownProposal(tx.ProposalID)
	}

	if (proposal.Status != gtypes.StatusDepositPeriod) && (proposal.Status != gtypes.StatusVotingPeriod) {
		return ErrFinishedProposal(tx.ProposalID)
	}

	accountMapper := baseabci.GetAccountMapper(ctx)
	account := accountMapper.GetAccount(tx.Depositor).(*types.QOSAccount)
	if !account.EnoughOfQOS(btypes.NewInt(int64(tx.Amount))) {
		return ErrInvalidInput("depositor has no enough qos")
	}

	return nil
}

func (tx TxDeposit) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	err, _ := GetGovMapper(ctx).AddDeposit(ctx, tx.ProposalID, tx.Depositor, tx.Amount)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	result.Tags = btypes.NewTags(btypes.TagAction, TagActionDepositProposal,
		TagProposalID, tx.ProposalID,
		TagDepositor, tx.Depositor.String())

	return
}

func (tx TxDeposit) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Depositor}
}

func (tx TxDeposit) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxDeposit) GetGasPayer() btypes.Address {
	return tx.Depositor
}

func (tx TxDeposit) GetSignData() (ret []byte) {
	ret = append(ret, types.Uint64ToBigEndian(tx.ProposalID)...)
	ret = append(ret, tx.Depositor...)
	ret = append(ret, types.Uint64ToBigEndian(tx.Amount)...)

	return
}
