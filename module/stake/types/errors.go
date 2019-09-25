package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// stake errors reserve 500 ~ 599.
const (
	DefaultCodeSpace btypes.CodespaceType = "stake"

	CodeInvalidInput                     btypes.CodeType = 501 // 输入有误
	CodeOwnerNoEnoughToken               btypes.CodeType = 503 // Owner账户Tokens不足
	CodeValidatorExists                  btypes.CodeType = 504 // Validator已存在
	CodeConsensusHasValidator            btypes.CodeType = 505 // Owner已绑定有Validator
	CodeValidatorNotExists               btypes.CodeType = 506 // Validator不存在
	CodeErrCommissionNegative            btypes.CodeType = 510 // Negative commission
	CodeErrCommissionHuge                btypes.CodeType = 511 // Commission too large
	CodeErrCommissionGTMaxRate           btypes.CodeType = 512 // Commission cannot be more than the max rate
	CodeErrCommissionChangeRateNegative  btypes.CodeType = 513 // Commission change rate must be positive
	CodeErrCommissionChangeRateGTMaxRate btypes.CodeType = 514 // Commission change rate cannot be more than the max rate
	CodeErrCommissionUpdateTime          btypes.CodeType = 515 // Commission cannot be changed more than once in 24h
	CodeErrCommissionGTMaxChangeRate     btypes.CodeType = 516 // Commission cannot be changed more than max change rate
	CoderErrOwnerNotMatch                btypes.CodeType = 517 // validator owner不匹配
)

func ErrInvalidInput(msg string) btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeInvalidInput, msg)
}

func ErrOwnerNoEnoughToken() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeOwnerNoEnoughToken, "owner has no enough token")
}

func ErrValidatorExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeValidatorExists, "validator exists")
}

func ErrConsensusHasValidator() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeConsensusHasValidator, "consensus pubkey already bind a validator")
}

func ErrValidatorNotExists() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeValidatorNotExists, "validator not exists")
}

func ErrCommissionNegative() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionNegative, "commission must be positive")
}

func ErrCommissionHuge() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionHuge, "commission cannot be more than 100%")
}

func ErrCommissionGTMaxRate() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionGTMaxRate, "commission cannot be more than the max rate")
}

func ErrCommissionChangeRateNegative() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionChangeRateNegative, "commission change rate must be positive")
}

func ErrCommissionChangeRateGTMaxRate() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionChangeRateGTMaxRate, "commission change rate cannot be more than the max rate")
}

func ErrCommissionUpdateTime() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionUpdateTime, "commission cannot be changed more than once in 24h")
}

func ErrCommissionGTMaxChangeRate() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CodeErrCommissionGTMaxChangeRate, "commission cannot be changed more than max change rate")
}

func ErrOwnerNotMatch() btypes.Error {
	return btypes.NewError(DefaultCodeSpace, CoderErrOwnerNotMatch, "validator owner not match")
}
