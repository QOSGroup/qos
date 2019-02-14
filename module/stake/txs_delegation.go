package stake

import (
	"errors"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco"
	"github.com/QOSGroup/qos/module/eco/mapper"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
)

type TxCreateDelegation struct {
	Delegator      btypes.Address //委托人
	ValidatorOwner btypes.Address //验证者Owner
	Amount         uint64         //委托QOS数量
	IsCompound     bool           //定期收益是否复投
}

var _ txs.ITx = (*TxCreateDelegation)(nil)

func (tx *TxCreateDelegation) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 || len(tx.ValidatorOwner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	if tx.Amount == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Delegation amount must be a positive integer.")
	}

	if _, err := validateValidator(ctx, tx.ValidatorOwner, true, staketypes.Active, true); err != nil {
		return err
	}

	if err := validateQOSAccount(ctx, tx.Delegator, tx.Amount); err != nil {
		return err
	}

	return nil
}

//创建或新增委托
func (tx *TxCreateDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	e := eco.GetEco(ctx)

	validator, _ := e.ValidatorMapper.GetValidatorByOwner(tx.ValidatorOwner)

	if err := e.DelegateValidator(validator, tx.Delegator, tx.Amount, tx.IsCompound, true); err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	return btypes.Result{Code: btypes.CodeOK}, nil
}

func (tx *TxCreateDelegation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxCreateDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateDelegation) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Delegator)
}

func (tx *TxCreateDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorOwner...)
	ret = append(ret, btypes.Int2Byte(int64(tx.Amount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}

type TxModifyCompound struct {
	Delegator      btypes.Address //委托人
	ValidatorOwner btypes.Address //验证者Owner
	IsCompound     bool           //周期收益是否复投: 收益发放周期内多次修改,仅最后一次生效
}

var _ txs.ITx = (*TxModifyCompound)(nil)

func (tx *TxModifyCompound) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 || len(tx.ValidatorOwner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	// TODO:是否允许validator为inactive/jailed时修改
	validator, err := validateValidator(ctx, tx.ValidatorOwner, true, staketypes.Active, true)
	if nil != err {
		return err
	}

	info, err := validateDelegator(ctx, validator.GetValidatorAddress(), tx.Delegator, false, 0)
	if err != nil {
		return err
	}

	if info.IsCompound == tx.IsCompound {
		return ErrInvalidInput(DefaultCodeSpace, "delegator's compound not change")
	}

	return nil
}

//修改收益单复利
func (tx *TxModifyCompound) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	e := eco.GetEco(ctx)

	validator, _ := e.ValidatorMapper.GetValidatorByOwner(tx.ValidatorOwner)
	info, _ := e.DelegationMapper.GetDelegationInfo(tx.Delegator, validator.GetValidatorAddress())

	info.IsCompound = tx.IsCompound
	e.DelegationMapper.SetDelegationInfo(info)

	return btypes.Result{Code: btypes.CodeOK}, nil
}

func (tx *TxModifyCompound) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxModifyCompound) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxModifyCompound) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Delegator)
}

func (tx *TxModifyCompound) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorOwner...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}

type TxUnbondDelegation struct {
	Delegator      btypes.Address //委托人
	ValidatorOwner btypes.Address //验证者Owner
	UnbondAmount   uint64         //unbond数量
	IsUnbondAll    bool           //是否全部解绑, 为true时覆盖UnbondAmount
}

var _ txs.ITx = (*TxUnbondDelegation)(nil)

func (tx *TxUnbondDelegation) ValidateData(ctx context.Context) error {

	if !tx.IsUnbondAll && tx.UnbondAmount == 0 {
		return errors.New("unbond QOS amount is zero")
	}

	validator, err := validateValidator(ctx, tx.ValidatorOwner, false, staketypes.Active, true)
	if nil != err {
		return err
	}

	if !tx.IsUnbondAll && (validator.BondTokens < tx.UnbondAmount) {
		return ErrInvalidInput(DefaultCodeSpace, "validator does't have enough tokens")
	}

	isCheckAmount := !tx.IsUnbondAll
	checkAmount := uint64(0)

	if isCheckAmount {
		checkAmount = tx.UnbondAmount
	}

	if _, err = validateDelegator(ctx, validator.GetValidatorAddress(), tx.Delegator, isCheckAmount, checkAmount); err != nil {
		return err
	}

	return nil
}

//unbond delegator tokens
func (tx *TxUnbondDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	e := eco.GetEco(ctx)

	validator, _ := e.ValidatorMapper.GetValidatorByOwner(tx.ValidatorOwner)
	if err := e.UnbondValidator(validator, tx.Delegator, tx.IsUnbondAll, tx.UnbondAmount, false); err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	return btypes.Result{Code: btypes.CodeOK}, nil
}

func (tx *TxUnbondDelegation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxUnbondDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxUnbondDelegation) GetGasPayer() btypes.Address {
	return tx.Delegator
}

func (tx *TxUnbondDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorOwner...)
	ret = append(ret, btypes.Int2Byte(int64(tx.UnbondAmount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsUnbondAll)...)
	return
}

type TxCreateReDelegation struct {
	Delegator          btypes.Address //委托人
	FromValidatorOwner btypes.Address //原委托验证人Owner
	ToValidatorOwner   btypes.Address //现委托验证人Owner
	Amount             uint64         //委托数量
	IsRedelegateAll    bool           //
	IsCompound         bool           //
}

var _ txs.ITx = (*TxCreateReDelegation)(nil)

func (tx *TxCreateReDelegation) ValidateData(ctx context.Context) error {

	if !tx.IsRedelegateAll && tx.Amount == 0 {
		return errors.New("redelegate QOS amount is zero")
	}

	//1. 校验fromValidator是否存在
	validator, err := validateValidator(ctx, tx.FromValidatorOwner, false, 0, true)
	if err != nil {
		return err
	}

	//2. 校验toValidator是否存在 且 状态为active
	_, err = validateValidator(ctx, tx.ToValidatorOwner, true, staketypes.Active, true)
	if err != nil {
		return err
	}

	//3. 校验当前用户是否委托了fromValidator
	_, err = validateDelegator(ctx, validator.GetValidatorAddress(), tx.Delegator, true, tx.Amount)
	if err != nil {
		return err
	}

	return nil
}

//delegate from one to another
func (tx *TxCreateReDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	e := eco.GetEco(ctx)

	fromValidator, _ := e.ValidatorMapper.GetValidatorByOwner(tx.FromValidatorOwner)
	toValidator, _ := e.ValidatorMapper.GetValidatorByOwner(tx.ToValidatorOwner)
	info, _ := e.DelegationMapper.GetDelegationInfo(tx.Delegator, fromValidator.GetValidatorAddress())

	reDelegateAmount := tx.Amount
	if tx.IsRedelegateAll {
		reDelegateAmount = info.Amount
	}

	if err := e.UnbondValidator(fromValidator, tx.Delegator, false, reDelegateAmount, true); err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	if err := e.DelegateValidator(toValidator, tx.Delegator, reDelegateAmount, tx.IsCompound, false); err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	return btypes.Result{Code: btypes.CodeOK}, nil

}

func (tx *TxCreateReDelegation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxCreateReDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateReDelegation) GetGasPayer() btypes.Address {
	return tx.Delegator
}

func (tx *TxCreateReDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.FromValidatorOwner...)
	ret = append(ret, tx.ToValidatorOwner...)
	ret = append(ret, btypes.Int2Byte(int64(tx.Amount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	ret = append(ret, btypes.Bool2Byte(tx.IsRedelegateAll)...)
	return
}

func validateDelegator(ctx context.Context, valAddr, deleAddr btypes.Address, checkAmount bool, maxAmount uint64) (staketypes.DelegationInfo, error) {

	delegationMapper := mapper.GetDelegationMapper(ctx)
	info, exsits := delegationMapper.GetDelegationInfo(deleAddr, valAddr)
	if !exsits {
		return info, ErrInvalidInput(DefaultCodeSpace, "delegator not delegate the owner's validator")
	}

	if checkAmount {
		if info.Amount < maxAmount {
			return info, ErrInvalidInput(DefaultCodeSpace, "delegator does't have enough amount of QOS")
		}
	}

	return info, nil
}
