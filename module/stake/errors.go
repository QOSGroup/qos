package stake

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// stake errors reserve 500 ~ 599.
const (
	DefaultCodeSpace btypes.CodespaceType = "stake"

	CodeInvalidInput            btypes.CodeType = 501 // 输入有误
	CodeOwnerNotExists          btypes.CodeType = 502 // Owner账户不存在
	CodeOwnerNoEnoughToken      btypes.CodeType = 503 // Owner账户Tokens不足
	CodeValidatorExists         btypes.CodeType = 504 // Validator已存在
	CodeOwnerHasValidator       btypes.CodeType = 505 // Owner已绑定有Validator
	CodeValidatorNotExists      btypes.CodeType = 506 // Validator不存在
	CodeValidatorIsActive       btypes.CodeType = 507 // Validator处于激活状态
	CodeValidatorIsInactive     btypes.CodeType = 508 // Validator处于非激活状态
	CodeValidatorInactiveIncome btypes.CodeType = 509 // Validator处于非激活状态时收益非法
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
		return "invalid input"
	case CodeOwnerNotExists:
		return "owner not exists"
	case CodeOwnerNoEnoughToken:
		return "owner has no enough token"
	case CodeValidatorExists:
		return "validator exists"
	case CodeOwnerHasValidator:
		return "owner already bind a validator"
	case CodeValidatorNotExists:
		return "validator not exists"
	case CodeValidatorIsActive:
		return "validator is active"
	case CodeValidatorIsInactive:
		return "validator is inactive"
	case CodeValidatorInactiveIncome:
		return "vaidator in inactive and got fees"
	default:
		return btypes.CodeToDefaultMsg(code)
	}
}

func ErrInvalidInput(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeInvalidInput, msg)
}

func ErrOwnerNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeOwnerNotExists, msg)
}

func ErrOwnerNoEnoughToken(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeOwnerNoEnoughToken, msg)
}

func ErrValidatorExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeValidatorExists, msg)
}

func ErrOwnerHasValidator(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeOwnerHasValidator, msg)
}

func ErrValidatorNotExists(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeValidatorNotExists, msg)
}

func ErrValidatorIsActive(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeValidatorIsActive, msg)
}

func ErrValidatorIsInactive(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeValidatorIsInactive, msg)
}

func ErrCodeValidatorInactiveIncome(codeSpace btypes.CodespaceType, msg string) btypes.Error {
	return newError(codeSpace, CodeValidatorInactiveIncome, msg)
}
