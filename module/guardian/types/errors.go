package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// QCP errors reserve 700 ~ 799.
const (
	DefaultCodeSpace btypes.CodespaceType = "guardian"

	CodeInvalidInput          btypes.CodeType = 601 // invalid input
	CodeInvalidCreator        btypes.CodeType = 602 // invalid creator
	CodeUnKnownGuardian       btypes.CodeType = 603 // unknown guardian
	CodeGuardianAlreadyExists btypes.CodeType = 604 // guardian already exists
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
	case CodeInvalidCreator:
		return "invalid creator"
	case CodeUnKnownGuardian:
		return "unknown guardian"
	case CodeGuardianAlreadyExists:
		return "guardian already exists"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrInvalidCreator(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidCreator, msg)
}

func ErrUnKnownGuardian(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeUnKnownGuardian, msg)
}

func ErrGuardianAlreadyExists(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeGuardianAlreadyExists, msg)
}
