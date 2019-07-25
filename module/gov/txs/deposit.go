package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/types"
	qtypes "github.com/QOSGroup/qos/types"
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
		return types.ErrInvalidInput("depositor is empty")
	}

	if tx.Amount == 0 {
		return types.ErrInvalidInput("amount of deposit is zero")
	}

	proposal, ok := mapper.GetMapper(ctx).GetProposal(tx.ProposalID)
	if !ok {
		return types.ErrUnknownProposal(tx.ProposalID)
	}

	if (proposal.Status != types.StatusDepositPeriod) && (proposal.Status != types.StatusVotingPeriod) {
		return types.ErrFinishedProposal(tx.ProposalID)
	}

	accountMapper := baseabci.GetAccountMapper(ctx)
	account := accountMapper.GetAccount(tx.Depositor).(*qtypes.QOSAccount)
	if !account.EnoughOfQOS(btypes.NewInt(int64(tx.Amount))) {
		return types.ErrInvalidInput("depositor has no enough qos")
	}

	return nil
}

func (tx TxDeposit) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	err, _ := mapper.GetMapper(ctx).AddDeposit(ctx, tx.ProposalID, tx.Depositor, tx.Amount)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeDepositProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", tx.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Depositor.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

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
	ret = append(ret, qtypes.Uint64ToBigEndian(tx.ProposalID)...)
	ret = append(ret, tx.Depositor...)
	ret = append(ret, qtypes.Uint64ToBigEndian(tx.Amount)...)

	return
}
