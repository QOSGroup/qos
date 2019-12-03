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

var GasForUnbond = int64(0.18*qtypes.UnitQOS) * qtypes.UnitQOSGas // 0.18 QOS

// 委托Tx
type TxCreateDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人
	Amount        btypes.BigInt     `json:"amount"`         // 委托QOS数量
	IsCompound    bool              `json:"is_compound"`    // 定期收益是否复投
}

var _ txs.ITx = (*TxCreateDelegation)(nil)

// 数据验证
func (tx *TxCreateDelegation) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 {
		return types.ErrInvalidInput("delegator is empty")
	}
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput("validator address is empty")
	}
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return types.ErrInvalidInput("amount must be a positive")
	}

	if _, err := validateValidator(ctx, tx.ValidatorAddr, false, types.Active, btypes.AccAddress{}, false); err != nil {
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
	sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens().Add(tx.Amount))

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateDelegation,
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Delegator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeCreateDelegation),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户，delegator
func (tx *TxCreateDelegation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

// Tx Gas, 0
func (tx *TxCreateDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, delegator
func (tx *TxCreateDelegation) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.Delegator)
}

// 签名字节
func (tx *TxCreateDelegation) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)
	return
}

// 修改单复利Tx
type TxModifyCompound struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证者
	IsCompound    bool              `json:"is_compound"`    // 周期收益是否复投: 收益发放周期内多次修改,仅最后一次生效
}

var _ txs.ITx = (*TxModifyCompound)(nil)

// 数据校验
func (tx *TxModifyCompound) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 {
		return types.ErrInvalidInput("delegator is empty")
	}
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput("validator address is empty")
	}

	validator, err := validateValidator(ctx, tx.ValidatorAddr, false, 0, btypes.AccAddress{}, false)
	if nil != err {
		return err
	}

	info, err := validateDelegator(ctx, validator.GetValidatorAddress(), tx.Delegator, false, btypes.ZeroInt())
	if err != nil {
		return err
	}

	if info.IsCompound == tx.IsCompound {
		return types.ErrInvalidInput("delegator's compound not change")
	}

	return nil
}

// 交易执行
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeModifyCompound),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账户，delegator
func (tx *TxModifyCompound) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

// Tx Gas, 0
func (tx *TxModifyCompound) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, delegator
func (tx *TxModifyCompound) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.Delegator)
}

// 签名字节
func (tx *TxModifyCompound) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)
	return
}

// 解除委托Tx
type TxUnbondDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // 委托人
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证者
	UnbondAmount  btypes.BigInt     `json:"unbond_amount"`  // unbond数量
	UnbondAll     bool              `json:"unbond_all"`     // 是否全部解绑, 为true时覆盖UnbondAmount
}

var _ txs.ITx = (*TxUnbondDelegation)(nil)

// 数据校验
func (tx *TxUnbondDelegation) ValidateData(ctx context.Context) error {

	if !tx.UnbondAll && !tx.UnbondAmount.GT(btypes.ZeroInt()) {
		return errors.New("unbond QOS amount must be positive")
	}

	validator, err := validateValidator(ctx, tx.ValidatorAddr, false, types.Active, btypes.AccAddress{}, false)
	if nil != err {
		return err
	}

	if !tx.UnbondAll && (validator.BondTokens.LT(tx.UnbondAmount)) {
		return types.ErrInvalidInput("validator does't have enough tokens")
	}

	checkAmount := btypes.ZeroInt()

	if !tx.UnbondAll {
		checkAmount = tx.UnbondAmount
	}

	if _, err = validateDelegator(ctx, validator.GetValidatorAddress(), tx.Delegator, !tx.UnbondAll, checkAmount); err != nil {
		return err
	}

	return nil
}

// 交易执行
func (tx *TxUnbondDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	sm := mapper.GetMapper(ctx)
	validator, _ := sm.GetValidator(tx.ValidatorAddr)
	delegation, _ := sm.GetDelegationInfo(tx.Delegator, validator.GetValidatorAddress())

	if tx.UnbondAll {
		tx.UnbondAmount = delegation.Amount
	}

	if delegation.Amount.LT(tx.UnbondAmount) || validator.GetBondTokens().LT(tx.UnbondAmount) {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// unBond delegation tokens
	sm.UnbondTokens(ctx, delegation, tx.UnbondAmount)

	// update validator
	sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens().Sub(tx.UnbondAmount))

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeUnbondDelegation,
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Delegator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeUnbondDelegation),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户，delegator
func (tx *TxUnbondDelegation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

// Tx Gas, 0.18QOS
func (tx *TxUnbondDelegation) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForUnbond)
}

// Gas payer, delegator
func (tx *TxUnbondDelegation) GetGasPayer() btypes.AccAddress {
	return tx.Delegator
}

// 签名字节
func (tx *TxUnbondDelegation) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 转委托Tx
type TxCreateReDelegation struct {
	Delegator         btypes.AccAddress `json:"delegator"`           // 委托人
	FromValidatorAddr btypes.ValAddress `json:"from_validator_addr"` // 原委托验证人
	ToValidatorAddr   btypes.ValAddress `json:"to_validator_addr"`   // 现委托验证人
	Amount            btypes.BigInt     `json:"amount"`              // 委托数量
	RedelegateAll     bool              `json:"redelegate_all"`      // 转委托所有
	Compound          bool              `json:"compound"`            // 复投
}

var _ txs.ITx = (*TxCreateReDelegation)(nil)

// 数据校验
func (tx *TxCreateReDelegation) ValidateData(ctx context.Context) error {

	if !tx.RedelegateAll && !tx.Amount.GT(btypes.ZeroInt()) {
		return errors.New("redelegate QOS amount is zero")
	}

	//1. 校验fromValidator是否存在
	validator, err := validateValidator(ctx, tx.FromValidatorAddr, false, types.Active, btypes.AccAddress{}, false)
	if err != nil {
		return err
	}

	//2. 校验toValidator是否存在 <del>且 状态为active</del>
	_, err = validateValidator(ctx, tx.ToValidatorAddr, true, types.Active, btypes.AccAddress{}, false)
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

// 交易执行
func (tx *TxCreateReDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	sm := mapper.GetMapper(ctx)

	fromValidator, _ := sm.GetValidator(tx.FromValidatorAddr)
	toValidator, _ := sm.GetValidator(tx.ToValidatorAddr)
	delegation, _ := sm.GetDelegationInfo(tx.Delegator, fromValidator.GetValidatorAddress())

	if tx.RedelegateAll {
		tx.Amount = delegation.Amount
	}

	if fromValidator.GetBondTokens().LT(tx.Amount) {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// redelegate
	redelegateHeight := sm.GetParams(ctx).DelegatorUnbondFrozenHeight + ctx.BlockHeight()
	sm.ReDelegate(ctx, delegation, types.NewRedelegateInfo(delegation.DelegatorAddr, fromValidator.GetValidatorAddress(), toValidator.GetValidatorAddress(), tx.Amount, ctx.BlockHeight(), redelegateHeight, tx.Compound))

	// update validator
	sm.ChangeValidatorBondTokens(fromValidator, fromValidator.GetBondTokens().Sub(tx.Amount))

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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeCreateReDelegation),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return

}

// 签名账户，delegator
func (tx *TxCreateReDelegation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Delegator}
}

// Tx Gas, 0
func (tx *TxCreateReDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, delegator
func (tx *TxCreateReDelegation) GetGasPayer() btypes.AccAddress {
	return tx.Delegator
}

// 签名字节
func (tx *TxCreateReDelegation) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)
	return
}

// 验证委托状态
func validateDelegator(ctx context.Context, valAddr btypes.ValAddress, deleAddr btypes.AccAddress, checkAmount bool, maxAmount btypes.BigInt) (types.DelegationInfo, error) {

	sm := mapper.GetMapper(ctx)
	info, exists := sm.GetDelegationInfo(deleAddr, valAddr)
	if !exists {
		return info, types.ErrInvalidInput("delegator not delegate the owner's validator")
	}

	if checkAmount {
		if info.Amount.LT(maxAmount) {
			return info, types.ErrInvalidInput("delegator does't have enough amount of QOS")
		}
	}

	return info, nil
}
