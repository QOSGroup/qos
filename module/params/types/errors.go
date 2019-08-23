package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// QCP errors reserve 600 ~ 699.
const (
	DefaultCodeSpace btypes.CodespaceType = "params"

	CodeInvalidParam btypes.CodeType = 701 // invalid param
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
	case CodeInvalidParam:
		return "invalid param"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidParam(msg string) btypes.Error {
	return newError(DefaultCodeSpace, CodeInvalidParam, msg)
}
