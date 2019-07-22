package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/types"
	qtypes "github.com/QOSGroup/qos/types"
)

type TxVote struct {
	ProposalID uint64           `json:"proposal_id"` // ID of the proposal
	Voter      btypes.Address   `json:"voter"`       //  address of the voter
	Option     types.VoteOption `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewTxVote(proposalID uint64, voter btypes.Address, option types.VoteOption) *TxVote {
	return &TxVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

var _ txs.ITx = (*TxVote)(nil)

func (tx TxVote) ValidateData(ctx context.Context) error {
	if len(tx.Voter) == 0 {
		return types.ErrInvalidInput("depositor is empty")
	}

	if !types.ValidVoteOption(tx.Option) {
		return types.ErrInvalidInput("invalid voting option")
	}

	proposal, ok := mapper.GetMapper(ctx).GetProposal(tx.ProposalID)
	if !ok {
		return types.ErrUnknownProposal(tx.ProposalID)
	}

	if proposal.Status != types.StatusVotingPeriod {
		return types.ErrWrongProposalStatus(tx.ProposalID)
	}

	return nil
}

func (tx TxVote) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	err := mapper.GetMapper(ctx).AddVote(tx.ProposalID, tx.Voter, tx.Option)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeVoteProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", tx.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyVoter, tx.Voter.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxVote) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Voter}
}

func (tx TxVote) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxVote) GetGasPayer() btypes.Address {
	return tx.Voter
}

func (tx TxVote) GetSignData() (ret []byte) {
	ret = append(ret, qtypes.Uint64ToBigEndian(tx.ProposalID)...)
	ret = append(ret, tx.Voter...)
	ret = append(ret, byte(tx.Option))

	return
}
