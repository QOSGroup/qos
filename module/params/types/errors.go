package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// Params errors reserve 700 ~ 799.
const (
	DefaultCodeSpace btypes.CodespaceType = "params"

	CodeInvalidParam btypes.CodeType = 701 // invalid param
)

func ErrInvalidParam(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidParam, msg)
}
