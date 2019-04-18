package gov

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
)

type TxVote struct {
	ProposalID uint64            `json:"proposal_id"` // ID of the proposal
	Voter      btypes.Address    `json:"voter"`       //  address of the voter
	Option     gtypes.VoteOption `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewTxVote(proposalID uint64, voter btypes.Address, option gtypes.VoteOption) *TxVote {
	return &TxVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

var _ txs.ITx = (*TxVote)(nil)

func (tx TxVote) ValidateData(ctx context.Context) error {
	if len(tx.Voter) == 0 {
		return ErrInvalidInput("depositor is empty")
	}

	if !gtypes.ValidVoteOption(tx.Option) {
		return ErrInvalidInput("invalid voting option")
	}

	proposal, ok := GetGovMapper(ctx).GetProposal(tx.ProposalID)
	if !ok {
		return ErrUnknownProposal(tx.ProposalID)
	}

	if proposal.Status != gtypes.StatusVotingPeriod {
		return ErrWrongProposalStatus(tx.ProposalID)
	}

	return nil
}

func (tx TxVote) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	err := GetGovMapper(ctx).AddVote(tx.ProposalID, tx.Voter, tx.Option)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	result.Tags = btypes.NewTags(btypes.TagAction, TagActionVoteProposal,
		TagProposalID, tx.ProposalID,
		TagVoter, tx.Voter.String())

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
	ret = append(ret, types.Uint64ToBigEndian(tx.ProposalID)...)
	ret = append(ret, tx.Voter...)
	ret = append(ret, byte(tx.Option))

	return
}
