package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// Guardian errors reserve 600 ~ 699.
const (
	DefaultCodeSpace btypes.CodespaceType = "guardian"

	CodeInvalidInput          btypes.CodeType = 601 // invalid input
	CodeInvalidCreator        btypes.CodeType = 602 // invalid creator
	CodeUnKnownGuardian       btypes.CodeType = 603 // unknown guardian
	CodeGuardianAlreadyExists btypes.CodeType = 604 // guardian already exists
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrInvalidCreator() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidCreator, "invalid creator")
}

func ErrUnKnownGuardian() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeUnKnownGuardian, "unknown guardian")
}

func ErrGuardianAlreadyExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeGuardianAlreadyExists, "guardian already exists")
}
