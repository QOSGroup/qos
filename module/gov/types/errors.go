package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
)

// Governance errors reserve 700 ~ 799.
const (
	DefaultCodeSpace btypes.CodespaceType = "gov"

	CodeInvalidInput        btypes.CodeType = 701 // invalid input
	CodeInvalidGenesis      btypes.CodeType = 702 // invalid genesis
	CodeUnknownProposal     btypes.CodeType = 703 // unknown proposal
	CodeInvalidVote         btypes.CodeType = 704 // invalid vote
	CodeWrongProposalStatus btypes.CodeType = 705 // wrong status of proposal
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrInvalidGenesis(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidGenesis, msg)
}

func ErrUnknownProposal(proposalID int64) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeUnknownProposal, fmt.Sprintf("unknown proposal %v", proposalID))
}

func ErrFinishedProposal(proposalID int64) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidVote, fmt.Sprintf("'%v' already finished", proposalID))
}

func ErrWrongProposalStatus(proposalID int64) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeWrongProposalStatus, fmt.Sprintf("wrong status of proposal %v", proposalID))
}
