package types

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
	CodeErrCommissionNegative            btypes.CodeType = 510 // Negative commission
	CodeErrCommissionHuge                btypes.CodeType = 511 // Validator处于非激活状态时收益非法
	CodeErrCommissionGTMaxRate           btypes.CodeType = 512 // Validator处于非激活状态时收益非法
	CodeErrCommissionChangeRateNegative  btypes.CodeType = 513 // Validator处于非激活状态时收益非法
	CodeErrCommissionChangeRateGTMaxRate btypes.CodeType = 514 // Validator处于非激活状态时收益非法
	CodeErrCommissionUpdateTime          btypes.CodeType = 515 // Validator处于非激活状态时收益非法
	CodeErrCommissionGTMaxChangeRate     btypes.CodeType = 516 // Validator处于非激活状态时收益非法
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
	case CodeErrCommissionNegative:
		return "commission must be positive"
	case CodeErrCommissionHuge:
		return "commission cannot be more than 100%"
	case CodeErrCommissionGTMaxRate:
		return "commission cannot be more than the max rate"
	case CodeErrCommissionChangeRateNegative:
		return "commission change rate must be positive"
	case CodeErrCommissionChangeRateGTMaxRate:
		return "commission change rate cannot be more than the max rate"
	case CodeErrCommissionUpdateTime:
		return "commission cannot be changed more than once in 24h"
	case CodeErrCommissionGTMaxChangeRate:
		return "commission cannot be changed more than max change rate"
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

func ErrCommissionNegative(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionNegative, "")
}

func ErrCommissionHuge(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionHuge, "")
}

func ErrCommissionGTMaxRate(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionGTMaxRate, "")
}

func ErrCommissionChangeRateNegative(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionChangeRateNegative, "")
}

func ErrCommissionChangeRateGTMaxRate(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionChangeRateGTMaxRate, "")
}

func ErrCommissionUpdateTime(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionUpdateTime, "")
}

func ErrCommissionGTMaxChangeRate(codeSpace btypes.CodespaceType) btypes.Error {
	return newError(codeSpace, CodeErrCommissionGTMaxChangeRate, "")
}
