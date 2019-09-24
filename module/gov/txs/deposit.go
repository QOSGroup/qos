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

// 提议质押
type TxDeposit struct {
	ProposalID int64             `json:"proposal_id"` // ID of the proposal
	Depositor  btypes.AccAddress `json:"depositor"`   // Address of the depositor
	Amount     btypes.BigInt     `json:"amount"`      // Percent of QOS to add to the proposal's deposit
}

func NewTxDeposit(proposalID int64, depositor btypes.AccAddress, amount btypes.BigInt) *TxDeposit {
	return &TxDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
}

var _ txs.ITx = (*TxDeposit)(nil)

// 数据校验
func (tx TxDeposit) ValidateData(ctx context.Context) error {
	// 质押账户不能为空
	if len(tx.Depositor) == 0 {
		return types.ErrInvalidInput("depositor is empty")
	}

	// 质押必须为正
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return types.ErrInvalidInput("amount of deposit must be more than zero")
	}
	// 提议存在，且处于质押期
	proposal, ok := mapper.GetMapper(ctx).GetProposal(tx.ProposalID)
	if !ok {
		return types.ErrUnknownProposal(tx.ProposalID)
	}
	if (proposal.Status != types.StatusDepositPeriod) && (proposal.Status != types.StatusVotingPeriod) {
		return types.ErrFinishedProposal(tx.ProposalID)
	}
	// 质押账户有足够的QOS质押
	accountMapper := baseabci.GetAccountMapper(ctx)
	account := accountMapper.GetAccount(tx.Depositor).(*qtypes.QOSAccount)
	if !account.EnoughOfQOS(tx.Amount) {
		return types.ErrInvalidInput("depositor has no enough qos")
	}

	return nil
}

// 交易执行
func (tx TxDeposit) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 保存质押
	err, _ := mapper.GetMapper(ctx).AddDeposit(ctx, tx.ProposalID, tx.Depositor, tx.Amount)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeDepositProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", tx.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Depositor.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeDepositProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户, Depositor
func (tx TxDeposit) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Depositor}
}

// Tx gas, 0
func (tx TxDeposit) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, Depositor
func (tx TxDeposit) GetGasPayer() btypes.AccAddress {
	return tx.Depositor
}

// 签名字节
func (tx TxDeposit) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}
