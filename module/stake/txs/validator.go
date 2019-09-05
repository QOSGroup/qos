package txs

import (
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MaxNameLen        = 300
	MaxLinkLen        = 255
	MaxDescriptionLen = 1000

	GasForCreateValidator = int64(1.8*qtypes.QOSUnit) * qtypes.GasPerUnitCost  // 1.8 QOS
	GasForModifyValidator = int64(0.18*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 0.18 QOS
	GasForRevokeValidator = int64(18*qtypes.QOSUnit) * qtypes.GasPerUnitCost   // 18 QOS
)

type TxCreateValidator struct {
	Operator    btypes.AccAddress      //操作者, self delegator
	ConsPubKey  crypto.PubKey          //validator公钥
	BondTokens  btypes.BigInt          //绑定Token数量
	IsCompound  bool                   //周期收益是否复投
	Description types.Description      //描述信息
	Commission  types.CommissionRates  //佣金比例
	Delegations []types.DelegationInfo //初始委托，仅在iniChainer中执行有效
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(operator btypes.AccAddress, bech32ConPubKey crypto.PubKey, bondTokens btypes.BigInt, isCompound bool, description types.Description, commission types.CommissionRates, delegations []types.DelegationInfo) *TxCreateValidator {
	return &TxCreateValidator{
		Operator:    operator,
		ConsPubKey:  bech32ConPubKey,
		BondTokens:  bondTokens,
		IsCompound:  isCompound,
		Description: description,
		Commission:  commission,
		Delegations: delegations,
	}
}

func (tx *TxCreateValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Description.Moniker) == 0 ||
		len(tx.Description.Moniker) > MaxNameLen ||
		tx.ConsPubKey == nil ||
		len(tx.Description.Logo) > MaxLinkLen ||
		len(tx.Description.Website) > MaxLinkLen ||
		len(tx.Description.Details) > MaxDescriptionLen ||
		len(tx.Operator) == 0 ||
		!tx.BondTokens.GT(btypes.ZeroInt()) {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	if ctx.BlockHeader().Height == 0 && len(tx.Delegations) != 0 {
		totalDelegation := btypes.ZeroInt()
		for _, delegation := range tx.Delegations {
			totalDelegation = totalDelegation.Add(delegation.Amount)
			err = validateQOSAccount(ctx, delegation.DelegatorAddr, delegation.Amount)
			if nil != err {
				return err
			}
		}
		if !totalDelegation.Equal(tx.BondTokens) {
			return types.ErrInvalidInput(types.DefaultCodeSpace, "validator bondTokens must equal sum(amount) of delegations")
		}
	} else {
		err = validateQOSAccount(ctx, tx.Operator, tx.BondTokens)
		if nil != err {
			return err
		}
	}

	err = tx.Commission.Validate()
	if err != nil {
		return
	}

	mapper := mapper.GetMapper(ctx)
	valAddr := btypes.ValAddress(tx.Operator)
	if mapper.Exists(valAddr) {
		return types.ErrValidatorExists(types.DefaultCodeSpace, "")
	}

	consAddr := btypes.ConsAddress(tx.ConsPubKey.Address())
	if mapper.ExistsWithConsensusAddr(consAddr) {
		return types.ErrConsensusHasValidator(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	result = btypes.Result{Code: btypes.CodeOK}
	valAddr := btypes.ValAddress(tx.Operator)

	validator := types.Validator{
		OperatorAddress: valAddr,
		Owner:           tx.Operator,
		ConsPubKey:      tx.ConsPubKey,
		BondTokens:      tx.BondTokens,
		Description:     tx.Description,
		Status:          types.Active,
		MinPeriod:       int64(0),
		BondHeight:      ctx.BlockHeight(),
		Commission:      types.NewCommissionWithTime(tx.Commission.Rate, tx.Commission.MaxRate, tx.Commission.MaxChangeRate, ctx.BlockHeader().Time.UTC()),
	}

	delegatorAddr := validator.Owner

	// 创建validator
	sm := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	sm.CreateValidator(validator)

	sm.AfterValidatorCreated(ctx, valAddr)

	// 初始化delegations
	if ctx.BlockHeader().Height == 0 && len(tx.Delegations) != 0 {
		for _, delegation := range tx.Delegations {
			sm.Delegate(ctx, delegation, false)
		}
	} else {
		delegationInfo := types.NewDelegationInfo(delegatorAddr, valAddr, tx.BondTokens, tx.IsCompound)
		sm.Delegate(ctx, delegationInfo, false)

	}

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateValidator,
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyOwner, tx.Operator.String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Operator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxCreateValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Operator}
}

func (tx *TxCreateValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForCreateValidator)
}

func (tx *TxCreateValidator) GetGasPayer() btypes.AccAddress {
	return tx.Operator
}

func (tx *TxCreateValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

type TxModifyValidator struct {
	Owner          btypes.AccAddress //验证人Owner地址
	ValidatorAddr  btypes.ValAddress //验证人地址
	Description    types.Description //描述信息
	CommissionRate *qtypes.Dec       //佣金比例
}

var _ txs.ITx = (*TxModifyValidator)(nil)

func NewModifyValidatorTx(owner btypes.AccAddress, validatorAddr btypes.ValAddress, description types.Description, commissionRate *qtypes.Dec) *TxModifyValidator {
	return &TxModifyValidator{
		Owner:          owner,
		ValidatorAddr:  validatorAddr,
		Description:    description,
		CommissionRate: commissionRate,
	}
}

func (tx *TxModifyValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Description.Moniker) > MaxNameLen ||
		len(tx.Description.Logo) > MaxLinkLen ||
		len(tx.Description.Website) > MaxLinkLen ||
		len(tx.Description.Details) > MaxDescriptionLen ||
		len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	validator, err := validateValidator(ctx, tx.ValidatorAddr, false, 0, tx.Owner, true)
	if err != nil {
		return err
	}

	// valid commission rate
	if tx.CommissionRate != nil {
		err = validator.Commission.ValidateNewRate(*tx.CommissionRate, ctx.BlockHeader().Time.UTC())
	}

	return
}

func (tx *TxModifyValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	result = btypes.Result{Code: btypes.CodeOK}

	validatorMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	validator, exists := validatorMapper.GetValidator(tx.ValidatorAddr)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	description := validator.Description
	if len(tx.Description.Moniker) != 0 {
		description.Moniker = tx.Description.Moniker
	}
	if len(tx.Description.Logo) != 0 {
		description.Logo = tx.Description.Logo
	}
	if len(tx.Description.Website) != 0 {
		description.Website = tx.Description.Website
	}
	if len(tx.Description.Details) != 0 {
		description.Details = tx.Description.Details
	}
	validator.Description = description

	if tx.CommissionRate != nil {
		validator.Commission = types.NewCommissionWithTime(*tx.CommissionRate, validator.Commission.MaxRate, validator.Commission.MaxChangeRate, ctx.BlockHeader().Time.UTC())
	}

	validatorMapper.Set(types.BuildValidatorKey(validator.GetValidatorAddress()), validator)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeModifyValidator,
			btypes.NewAttribute(types.AttributeKeyOwner, tx.ValidatorAddr.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxModifyValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

func (tx *TxModifyValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForModifyValidator)
}

func (tx *TxModifyValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

func (tx *TxModifyValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

type TxRevokeValidator struct {
	Owner         btypes.AccAddress //验证人Owner地址
	ValidatorAddr btypes.ValAddress //验证人地址
}

var _ txs.ITx = (*TxRevokeValidator)(nil)

func NewRevokeValidatorTx(owner btypes.AccAddress, validatorAddr btypes.ValAddress) *TxRevokeValidator {
	return &TxRevokeValidator{
		Owner:         owner,
		ValidatorAddr: validatorAddr,
	}
}

func (tx *TxRevokeValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	_, err = validateValidator(ctx, tx.ValidatorAddr, true, types.Active, tx.Owner, true)
	if nil != err {
		return err
	}

	return nil
}

func (tx *TxRevokeValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	validator, exists := mapper.GetValidator(tx.ValidatorAddr)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	valAddr := validator.GetValidatorAddress()
	mapper.MakeValidatorInactive(valAddr, ctx.BlockHeight(), ctx.BlockHeader().Time.UTC(), types.Revoke)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeRevokeValidator,
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyOwner, validator.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxRevokeValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

func (tx *TxRevokeValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForRevokeValidator)
}

func (tx *TxRevokeValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

func (tx *TxRevokeValidator) GetSignData() (ret []byte) {
	ret = append(ret, tx.ValidatorAddr...)
	ret = append(ret, tx.Owner...)
	return
}

type TxActiveValidator struct {
	Owner         btypes.AccAddress //验证人Owner地址
	ValidatorAddr btypes.ValAddress //验证人地址
	BondTokens    btypes.BigInt     //绑定Token数量
}

var _ txs.ITx = (*TxActiveValidator)(nil)

func NewActiveValidatorTx(owner btypes.AccAddress, validatorAddr btypes.ValAddress, bondTokens btypes.BigInt) *TxActiveValidator {
	return &TxActiveValidator{
		Owner:         owner,
		ValidatorAddr: validatorAddr,
		BondTokens:    bondTokens,
	}
}

func (tx *TxActiveValidator) ValidateData(ctx context.Context) (err error) {

	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	err = validateQOSAccount(ctx, btypes.AccAddress(tx.ValidatorAddr), tx.BondTokens)
	if nil != err {
		return err
	}

	_, err = validateValidator(ctx, tx.ValidatorAddr, true, types.Inactive, tx.Owner, true)
	if nil != err {
		return err
	}

	return nil
}

func (tx *TxActiveValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	// 准备 Mapper
	stakeMapper := mapper.GetMapper(ctx)
	accountMapper := baseabci.GetAccountMapper(ctx)

	// 获取 Owner 对应的 Validator
	validator, exists := stakeMapper.GetValidator(tx.ValidatorAddr)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// 激活 Validator
	validatorAddr := validator.GetValidatorAddress()
	stakeMapper.MakeValidatorActive(validatorAddr, tx.BondTokens)

	// 重置 ValidatorVoteInfo
	voteInfo := types.NewValidatorVoteInfo(validator.BondHeight+1, 0, 0)
	stakeMapper.ResetValidatorVoteInfo(validator.GetValidatorAddress(), voteInfo)

	// 获取 Owner 对应的 Delegator 账户
	delegatorAddr := validator.Owner
	delegator := accountMapper.GetAccount(delegatorAddr).(*qtypes.QOSAccount)

	// 当增加委托的tokens大于0时
	if tx.BondTokens.GT(btypes.ZeroInt()) {

		// 从 delegator 账户, 扣去增加的自委托token数量
		delegator.MustMinusQOS(tx.BondTokens)
		accountMapper.SetAccount(delegator)

		// 获取 delegationInfo
		delegationInfo, exists := stakeMapper.GetDelegationInfo(delegatorAddr, validatorAddr)
		if !exists {
			return btypes.Result{Code: btypes.CodeInternal}, nil
		}

		// 修改 delegationInfo 中的token amount, 并保存.
		delegationInfo.Amount = delegationInfo.Amount.Add(tx.BondTokens)
		stakeMapper.BeforeDelegationModified(ctx, validatorAddr, delegatorAddr, delegationInfo.Amount)
		stakeMapper.SetDelegationInfo(delegationInfo)
	}

	// 设置Events
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeActiveValidator,
			btypes.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
			btypes.NewAttribute(types.AttributeKeyOwner, validator.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxActiveValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

func (tx *TxActiveValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxActiveValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

func (tx *TxActiveValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

func validateQOSAccount(ctx context.Context, addr btypes.AccAddress, toPay btypes.BigInt) error {
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	acc := accountMapper.GetAccount(addr)
	if toPay.GT(btypes.ZeroInt()) {
		if acc != nil {
			qosAccount := acc.(*qtypes.QOSAccount)
			if !qosAccount.EnoughOfQOS(toPay) {
				return types.ErrOwnerNoEnoughToken(types.DefaultCodeSpace, "No enough QOS in account: "+addr.String())
			}
		} else {
			return types.ErrOwnerNoEnoughToken(types.DefaultCodeSpace, "account not exists: "+addr.String())
		}
	}
	return nil
}

func validateValidator(ctx context.Context, validatorAddr btypes.ValAddress, checkStatus bool, expectedStatus int8, owner btypes.AccAddress, checkOwner bool) (validator types.Validator, err error) {
	valMapper := mapper.GetMapper(ctx)
	validator, exists := valMapper.GetValidator(validatorAddr)
	if !exists {
		return validator, types.ErrValidatorNotExists(types.DefaultCodeSpace, "Validator not exists.")
	}

	if checkOwner && !validator.Owner.Equals(owner) {
		return validator, types.ErrOwnerNotMatch(types.DefaultCodeSpace, fmt.Sprintf("Owner:%s does not have right operate validator:%s", owner, validatorAddr))
	}

	if checkStatus {
		if validator.Status != expectedStatus {
			return validator, fmt.Errorf("Validator status not match. except: %d, actual:%d.", expectedStatus, validator.Status)
		}
	}
	return validator, nil
}
