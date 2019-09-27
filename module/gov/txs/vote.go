package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/types"
)

// 提议投票
type TxVote struct {
	ProposalID int64             `json:"proposal_id"` // ID of the proposal
	Voter      btypes.AccAddress `json:"voter"`       //  address of the voter
	Option     types.VoteOption  `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewTxVote(proposalID int64, voter btypes.AccAddress, option types.VoteOption) *TxVote {
	return &TxVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

var _ txs.ITx = (*TxVote)(nil)

// 数据校验
func (tx TxVote) ValidateInputs() error {
	// 投票账户存在
	if len(tx.Voter) == 0 {
		return types.ErrInvalidInput("voter is empty")
	}
	// 投票有效
	if !types.ValidVoteOption(tx.Option) {
		return types.ErrInvalidInput("invalid voting option")
	}

	return nil
}

// 数据校验
func (tx TxVote) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := tx.ValidateInputs()
	if err != nil {
		return err
	}

	// 提议存在，且处于投票期
	proposal, ok := mapper.GetMapper(ctx).GetProposal(tx.ProposalID)
	if !ok {
		return types.ErrUnknownProposal(tx.ProposalID)
	}
	if proposal.Status != types.StatusVotingPeriod {
		return types.ErrWrongProposalStatus(tx.ProposalID)
	}

	return nil
}

// 交易执行
func (tx TxVote) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 保存投票
	err := mapper.GetMapper(ctx).AddVote(tx.ProposalID, tx.Voter, tx.Option)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeVoteProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", tx.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyVoter, tx.Voter.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeVoteProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账户, Voter
func (tx TxVote) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Voter}
}

// Tx gas, 0
func (tx TxVote) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, Voter
func (tx TxVote) GetGasPayer() btypes.AccAddress {
	return tx.Voter
}

// 签名字节
func (tx TxVote) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}
