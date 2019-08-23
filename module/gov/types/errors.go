package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
)

// QCP errors reserve 600 ~ 699.
const (
	DefaultCodeSpace btypes.CodespaceType = "gov"

	CodeInvalidInput        btypes.CodeType = 601 // invalid input
	CodeInvalidGenesis      btypes.CodeType = 602 // invalid genesis
	CodeUnknownProposal     btypes.CodeType = 603 // unknown proposal
	CodeInactiveProposal    btypes.CodeType = 604 // inactive proposal
	CodeInvalidVote         btypes.CodeType = 605 // invalid vote
	CodeFinishedProposal    btypes.CodeType = 606 // finished proposal
	CodeWrongProposalStatus btypes.CodeType = 607 // wrong status of proposal
)

func msgOrDefaultMsg(msg string, code btypes.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codeSpace btypes.CodespaceType, code btypes.CodeType, msg string) btypes.Error {
	msg = msgOrDefaultMsg(msg, code)
	return btypes.NewError(codeSpace, code, msg)
}

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code btypes.CodeType) string {
	switch code {
	case CodeInvalidInput:
		return "invalid tx msg"
	case CodeInvalidGenesis:
		return "invalid genesis"
	case CodeUnknownProposal:
		return "unknown proposal"
	case CodeInactiveProposal:
		return "inactive proposal"
	case CodeInvalidVote:
		return "invalida vote"
	case CodeFinishedProposal:
		return "finished proposal"
	case CodeWrongProposalStatus:
		return "wrong status of proposal"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrInvalidGenesis(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidGenesis, msg)
}

func ErrUnknownProposal(proposalID uint64) btypes.Error {
	return newError(DefaultCodeSpace, CodeUnknownProposal, fmt.Sprintf("unknown proposal %v", proposalID))
}

func ErrInactiveProposal(proposalID uint64) btypes.Error {
	return newError(DefaultCodeSpace, CodeInactiveProposal, fmt.Sprintf("inactive proposal %v", proposalID))
}

func ErrInvalidVote(voteOption VoteOption) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidVote, fmt.Sprintf("'%v' is not a valid voting option", voteOption))
}

func ErrFinishedProposal(proposalID uint64) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidVote, fmt.Sprintf("'%v' already finished", proposalID))
}

func ErrWrongProposalStatus(proposalID uint64) btypes.Error {
	return newError(DefaultCodeSpace, CodeWrongProposalStatus, fmt.Sprintf("wrong status of proposal %v", proposalID))
}
