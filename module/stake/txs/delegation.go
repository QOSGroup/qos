package txs

import (
	"errors"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
)

const GasForUnbond = uint64(0.18*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 0.18 QOS

type TxCreateDelegation struct {
	Delegator     btypes.AccAddress //委托人
	ValidatorAddr btypes.ValAddress // 验证人
	Amount        uint64            //委托QOS数量
	IsCompound    bool              //定期收益是否复投
}

var _ txs.ITx = (*TxCreateDelegation)(nil)

func (tx *TxCreateDelegation) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 || len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	if tx.Amount <= 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "Delegation amount must be a positive integer.")
	}

	if _, err := validateValidator(ctx, tx.ValidatorAddr, false, types.Active, true, btypes.AccAddress{}, false); err != nil {
		return err
	}

	if err := validateQOSAccount(ctx, tx.Delegator, tx.Amount); err != nil {
		return err
	}

	return nil
}

//创建或新增委托
func (tx *TxCreateDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	sm := mapper.GetMapper(ctx)
	validator, _ := sm.GetValidator(tx.ValidatorAddr)

	// delegation
	info := types.NewDelegationInfo(tx.Delegator, validator.GetValidatorAddress(), tx.Amount, tx.IsCompound)
	sm.Delegate(ctx, info, false)

	// update validator
	sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens()+tx.Amount)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateDelegation,
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Delegator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxCreateDelegation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

func (tx *TxCreateDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateDelegation) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.Delegator)
}

func (tx *TxCreateDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorAddr...)
	ret = append(ret, btypes.Int2Byte(int64(tx.Amount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}

type TxModifyCompound struct {
	Delegator     btypes.AccAddress //委托人
	ValidatorAddr btypes.ValAddress //验证者
	IsCompound    bool              //周期收益是否复投: 收益发放周期内多次修改,仅最后一次生效
}

var _ txs.ITx = (*TxModifyCompound)(nil)

func (tx *TxModifyCompound) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 || len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	// TODO:是否允许validator为inactive/jailed时修改
	validator, err := validateValidator(ctx, tx.ValidatorAddr, false, 0, true, btypes.AccAddress{}, false)
	if nil != err {
		return err
	}

	info, err := validateDelegator(ctx, validator.GetValidatorAddress(), tx.Delegator, false, 0)
	if err != nil {
		return err
	}

	if info.IsCompound == tx.IsCompound {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "delegator's compound not change")
	}

	return nil
}

//修改收益单复利
func (tx *TxModifyCompound) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	sm := mapper.GetMapper(ctx)

	validator, _ := sm.GetValidator(tx.ValidatorAddr)
	info, _ := sm.GetDelegationInfo(tx.Delegator, validator.GetValidatorAddress())

	info.IsCompound = tx.IsCompound
	sm.SetDelegationInfo(info)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeModifyCompound,
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Delegator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxModifyCompound) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

func (tx *TxModifyCompound) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxModifyCompound) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.Delegator)
}

func (tx *TxModifyCompound) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorAddr...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}

type TxUnbondDelegation struct {
	Delegator     btypes.AccAddress //委托人
	ValidatorAddr btypes.ValAddress //验证者
	UnbondAmount  uint64            //unbond数量
	IsUnbondAll   bool              //是否全部解绑, 为true时覆盖UnbondAmount
}

var _ txs.ITx = (*TxUnbondDelegation)(nil)

func (tx *TxUnbondDelegation) ValidateData(ctx context.Context) error {

	if !tx.IsUnbondAll && tx.UnbondAmount == 0 {
		return errors.New("unbond QOS amount is zero")
	}

	validator, err := validateValidator(ctx, tx.ValidatorAddr, false, types.Active, true, btypes.AccAddress{}, false)
	if nil != err {
		return err
	}

	if !tx.IsUnbondAll && (validator.BondTokens < tx.UnbondAmount) {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "validator does't have enough tokens")
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
	result = btypes.Result{Code: btypes.CodeOK}

	sm := mapper.GetMapper(ctx)
	validator, _ := sm.GetValidator(tx.ValidatorAddr)
	delegation, _ := sm.GetDelegationInfo(tx.Delegator, validator.GetValidatorAddress())

	if tx.IsUnbondAll {
		tx.UnbondAmount = delegation.Amount
	}

	if delegation.Amount < tx.UnbondAmount || validator.GetBondTokens() < tx.UnbondAmount {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// unBond delegation tokens
	sm.UnbondTokens(ctx, delegation, tx.UnbondAmount)

	// update validator
	sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens()-tx.UnbondAmount)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeUnbondDelegation,
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Delegator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxUnbondDelegation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

func (tx *TxUnbondDelegation) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForUnbond))
}

func (tx *TxUnbondDelegation) GetGasPayer() btypes.AccAddress {
	return tx.Delegator
}

func (tx *TxUnbondDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorAddr...)
	ret = append(ret, btypes.Int2Byte(int64(tx.UnbondAmount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsUnbondAll)...)
	return
}

type TxCreateReDelegation struct {
	Delegator         btypes.AccAddress //委托人
	FromValidatorAddr btypes.ValAddress //原委托验证人
	ToValidatorAddr   btypes.ValAddress //现委托验证人
	Amount            uint64            //委托数量
	IsRedelegateAll   bool              //
	IsCompound        bool              //
}

var _ txs.ITx = (*TxCreateReDelegation)(nil)

func (tx *TxCreateReDelegation) ValidateData(ctx context.Context) error {

	if !tx.IsRedelegateAll && tx.Amount == 0 {
		return errors.New("redelegate QOS amount is zero")
	}

	//1. 校验fromValidator是否存在
	validator, err := validateValidator(ctx, tx.FromValidatorAddr, false, 0, true, btypes.AccAddress{}, false)
	if err != nil {
		return err
	}

	//2. 校验toValidator是否存在 <del>且 状态为active</del>
	_, err = validateValidator(ctx, tx.ToValidatorAddr, false, 0, true, btypes.AccAddress{}, false)
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
	result = btypes.Result{Code: btypes.CodeOK}

	sm := mapper.GetMapper(ctx)

	fromValidator, _ := sm.GetValidator(tx.FromValidatorAddr)
	toValidator, _ := sm.GetValidator(tx.ToValidatorAddr)
	delegation, _ := sm.GetDelegationInfo(tx.Delegator, fromValidator.GetValidatorAddress())

	if tx.IsRedelegateAll {
		tx.Amount = delegation.Amount
	}

	qtypes.AssertUint64NotOverflow(tx.Amount)
	qtypes.AssertUint64NotOverflow(toValidator.GetBondTokens() + tx.Amount)

	if fromValidator.GetBondTokens() < tx.Amount {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// redelegate
	redelegateHeight := uint64(sm.GetParams(ctx).DelegatorRedelegationActiveHeight) + uint64(ctx.BlockHeight())
	sm.ReDelegate(ctx, delegation, types.NewRedelegateInfo(delegation.DelegatorAddr, fromValidator.GetValidatorAddress(), toValidator.GetValidatorAddress(), tx.Amount, uint64(ctx.BlockHeight()), redelegateHeight, tx.IsCompound))

	// update validator
	sm.ChangeValidatorBondTokens(fromValidator, fromValidator.GetBondTokens()-tx.Amount)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateReDelegation,
			btypes.NewAttribute(types.AttributeKeyValidator, fromValidator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyNewValidator, toValidator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Delegator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return

}

func (tx *TxCreateReDelegation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

func (tx *TxCreateReDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateReDelegation) GetGasPayer() btypes.AccAddress {
	return tx.Delegator
}

func (tx *TxCreateReDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.FromValidatorAddr...)
	ret = append(ret, tx.ToValidatorAddr...)
	ret = append(ret, btypes.Int2Byte(int64(tx.Amount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	ret = append(ret, btypes.Bool2Byte(tx.IsRedelegateAll)...)
	return
}

func validateDelegator(ctx context.Context, valAddr btypes.ValAddress, deleAddr btypes.AccAddress, checkAmount bool, maxAmount uint64) (types.DelegationInfo, error) {

	sm := mapper.GetMapper(ctx)
	info, exists := sm.GetDelegationInfo(deleAddr, valAddr)
	if !exists {
		return info, types.ErrInvalidInput(types.DefaultCodeSpace, "delegator not delegate the owner's validator")
	}

	if checkAmount {
		if info.Amount < maxAmount {
			return info, types.ErrInvalidInput(types.DefaultCodeSpace, "delegator does't have enough amount of QOS")
		}
	}

	return info, nil
}
