package txs

import (
	"fmt"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/tendermint/tendermint/crypto/ed25519"

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

// 创建验证节点Tx
type TxCreateValidator struct {
	Owner       btypes.AccAddress      `json:"owner"`        // 操作者, self delegator
	ConsPubKey  crypto.PubKey          `json:"cons_pub_key"` // validator公钥
	BondTokens  btypes.BigInt          `json:"bond_tokens"`  // 绑定Token数量
	IsCompound  bool                   `json:"is_compound"`  // 周期收益是否复投
	Description types.Description      `json:"description"`  // 描述信息
	Commission  types.CommissionRates  `json:"commission"`   // 佣金比例
	Delegations []types.DelegationInfo `json:"delegations"`  // 初始委托，仅在iniChainer中执行有效
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(operator btypes.AccAddress, bech32ConPubKey crypto.PubKey, bondTokens btypes.BigInt, isCompound bool, description types.Description, commission types.CommissionRates, delegations []types.DelegationInfo) *TxCreateValidator {
	return &TxCreateValidator{
		Owner:       operator,
		ConsPubKey:  bech32ConPubKey,
		BondTokens:  bondTokens,
		IsCompound:  isCompound,
		Description: description,
		Commission:  commission,
		Delegations: delegations,
	}
}

// 基础数据校验
func (tx *TxCreateValidator) ValidateInputs() (err error) {
	if len(tx.Description.Moniker) == 0 {
		return types.ErrInvalidInput("moniker is empty")
	}
	if len(tx.Description.Moniker) > MaxNameLen {
		return types.ErrInvalidInput("moniker is too long")
	}
	if tx.ConsPubKey == nil {
		return types.ErrInvalidInput("cons_pub_key is empty")
	}
	if len(tx.Description.Logo) > MaxLinkLen {
		return types.ErrInvalidInput("logo is too long")
	}
	if len(tx.Description.Website) > MaxLinkLen {
		return types.ErrInvalidInput("website is too long")
	}
	if len(tx.Description.Details) > MaxDescriptionLen {
		return types.ErrInvalidInput("details is too long")
	}
	if len(tx.Owner) == 0 {
		return types.ErrInvalidInput("operator is empty")
	}
	if !tx.BondTokens.GT(btypes.ZeroInt()) {
		return types.ErrInvalidInput("bond_tokens must be positive")
	}

	// 验证创世文件中包含委托信息的验证节点创建交易
	if len(tx.Delegations) != 0 {
		totalDelegation := btypes.ZeroInt()
		for _, delegation := range tx.Delegations {
			totalDelegation = totalDelegation.Add(delegation.Amount)
		}
		if !totalDelegation.Equal(tx.BondTokens) {
			return types.ErrInvalidInput("validator bondTokens must equal sum(amount) of delegations")
		}
	}

	// 佣金参数校验
	err = tx.Commission.Validate()
	if err != nil {
		return
	}

	return nil
}

// 数据校验
func (tx *TxCreateValidator) ValidateData(ctx context.Context) (err error) {
	// 基础数据校验
	err = tx.ValidateInputs()
	if err != nil {
		return err
	}

	// 验证创世文件中包含委托信息的验证节点创建交易
	if ctx.BlockHeader().Height == 0 && len(tx.Delegations) != 0 {
		for _, delegation := range tx.Delegations {
			err = validateQOSAccount(ctx, delegation.DelegatorAddr, delegation.Amount)
			if nil != err {
				return err
			}
		}
	} else {
		err = validateQOSAccount(ctx, tx.Owner, tx.BondTokens)
		if nil != err {
			return err
		}
	}

	mapper := mapper.GetMapper(ctx)
	valAddr := btypes.ValAddress(tx.Owner)
	// 验证节点已存在校验
	if mapper.Exists(valAddr) {
		return types.ErrValidatorExists()
	}

	consAddr := btypes.ConsAddress(tx.ConsPubKey.Address())
	// 共识地址已存在校验
	if mapper.ExistsWithConsensusAddr(consAddr) {
		return types.ErrConsensusHasValidator()
	}

	return nil
}

// 交易执行
func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}
	valAddr := btypes.ValAddress(ed25519.GenPrivKey().PubKey().Address())

	validator := types.Validator{
		OperatorAddress: valAddr,
		Owner:           tx.Owner,
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
	sm := mapper.GetMapper(ctx)
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

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateValidator,
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyOwner, tx.Owner.String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeCreateValidator),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户，operator
func (tx *TxCreateValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

// Tx Gas, 1.8QOS
func (tx *TxCreateValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForCreateValidator)
}

// Gas payer, operator
func (tx *TxCreateValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

// 签名字节
func (tx *TxCreateValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

// 修改验证节点基础信息Tx
type TxModifyValidator struct {
	Owner          btypes.AccAddress `json:"owner"`           // 验证人Owner地址
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`  // 验证人地址
	Description    types.Description `json:"description"`     // 描述信息
	CommissionRate *qtypes.Dec       `json:"commission_rate"` // 佣金比例
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

// 基础数据校验
func (tx *TxModifyValidator) ValidateInputs() (err error) {
	if len(tx.Owner) == 0 {
		return types.ErrInvalidInput("owner is empty")
	}
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput("validator address is empty")
	}
	if len(tx.Description.Moniker) > MaxNameLen {
		return types.ErrInvalidInput("moniker is too long")
	}
	if len(tx.Description.Logo) > MaxLinkLen {
		return types.ErrInvalidInput("logo url is too long")
	}
	if len(tx.Description.Website) > MaxLinkLen {
		return types.ErrInvalidInput("website is too long")
	}
	if len(tx.Description.Details) > MaxDescriptionLen {
		return types.ErrInvalidInput("details is too long")
	}

	return
}

// 数据校验
func (tx *TxModifyValidator) ValidateData(ctx context.Context) (err error) {
	// 校验基础数据
	err = tx.ValidateInputs()
	if err != nil {
		return err
	}

	// 校验验证节点信息
	validator, err := validateValidator(ctx, tx.ValidatorAddr, false, types.Active, tx.Owner, true)
	if err != nil {
		return err
	}

	// valid commission rate
	if tx.CommissionRate != nil {
		err = validator.Commission.ValidateNewRate(*tx.CommissionRate, ctx.BlockHeader().Time.UTC())
	}

	return
}

// 交易执行
func (tx *TxModifyValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	// 获取验证节点信息
	validatorMapper := mapper.GetMapper(ctx)
	validator, exists := validatorMapper.GetValidator(tx.ValidatorAddr)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// 更新描述信息
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

	// 更新佣金比例信息
	if tx.CommissionRate != nil {
		validator.Commission = types.NewCommissionWithTime(*tx.CommissionRate, validator.Commission.MaxRate, validator.Commission.MaxChangeRate, ctx.BlockHeader().Time.UTC())
	}

	// 更新验证节点信息
	validatorMapper.Set(types.BuildValidatorKey(validator.GetValidatorAddress()), validator)

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeModifyValidator,
			btypes.NewAttribute(types.AttributeKeyOwner, tx.ValidatorAddr.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeModifyValidator),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户
func (tx *TxModifyValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

// Tx Gas, 0.18QOS
func (tx *TxModifyValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForModifyValidator)
}

// Gas payer,
func (tx *TxModifyValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

// 签名字节
func (tx *TxModifyValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

// 撤销验证节点Tx
type TxRevokeValidator struct {
	Owner         btypes.AccAddress `json:"owner"`          // 验证人Owner地址
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人地址
}

var _ txs.ITx = (*TxRevokeValidator)(nil)

func NewRevokeValidatorTx(owner btypes.AccAddress, validatorAddr btypes.ValAddress) *TxRevokeValidator {
	return &TxRevokeValidator{
		Owner:         owner,
		ValidatorAddr: validatorAddr,
	}
}

// 数据校验
func (tx *TxRevokeValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Owner) == 0 {
		return types.ErrInvalidInput("owner is empty")
	}
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput("validator address is empty")
	}

	_, err = validateValidator(ctx, tx.ValidatorAddr, true, types.Active, tx.Owner, true)
	if nil != err {
		return err
	}

	return nil
}

// 交易执行
func (tx *TxRevokeValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	// 获取验证节点信息
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	validator, exists := mapper.GetValidator(tx.ValidatorAddr)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}
	valAddr := validator.GetValidatorAddress()

	// 更新验证节点状态, active -> inactive
	mapper.MakeValidatorInactive(valAddr, ctx.BlockHeight(), ctx.BlockHeader().Time.UTC(), types.Revoke)

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeRevokeValidator,
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyOwner, validator.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeRevokeValidator),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户，owenr
func (tx *TxRevokeValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

// Tx Gas, 18QOS
func (tx *TxRevokeValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForRevokeValidator)
}

// Gas payer, owenr
func (tx *TxRevokeValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

// 签名字节
func (tx *TxRevokeValidator) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 激活验证节点Tx
type TxActiveValidator struct {
	Owner         btypes.AccAddress `json:"owner"`          // 验证人Owner地址
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // 验证人地址
	BondTokens    btypes.BigInt     `json:"bond_tokens"`    // 绑定Token数量
}

var _ txs.ITx = (*TxActiveValidator)(nil)

func NewActiveValidatorTx(owner btypes.AccAddress, validatorAddr btypes.ValAddress, bondTokens btypes.BigInt) *TxActiveValidator {
	return &TxActiveValidator{
		Owner:         owner,
		ValidatorAddr: validatorAddr,
		BondTokens:    bondTokens,
	}
}

// 数据验证
func (tx *TxActiveValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Owner) == 0 {
		return types.ErrInvalidInput("owner is empty")
	}
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput("validator address is empty")
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

// 交易执行
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeActiveValidator),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户，owner
func (tx *TxActiveValidator) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Owner}
}

// Tx Gas, 0
func (tx *TxActiveValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, owner
func (tx *TxActiveValidator) GetGasPayer() btypes.AccAddress {
	return tx.Owner
}

// 签名字节
func (tx *TxActiveValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

// 验证账户余额
func validateQOSAccount(ctx context.Context, addr btypes.AccAddress, toPay btypes.BigInt) error {
	accountMapper := bank.GetMapper(ctx)
	acc := accountMapper.GetAccount(addr)
	if toPay.GT(btypes.ZeroInt()) {
		if acc != nil {
			qosAccount := acc.(*qtypes.QOSAccount)
			if !qosAccount.EnoughOfQOS(toPay) {
				return types.ErrOwnerNoEnoughToken()
			}
		} else {
			return types.ErrOwnerNoEnoughToken()
		}
	}
	return nil
}

// 验证验证节点状态， 返回验证节点信息
func validateValidator(ctx context.Context, validatorAddr btypes.ValAddress, checkStatus bool, expectedStatus int8, owner btypes.AccAddress, checkOwner bool) (validator types.Validator, err error) {
	valMapper := mapper.GetMapper(ctx)
	validator, exists := valMapper.GetValidator(validatorAddr)
	if !exists {
		return validator, types.ErrValidatorNotExists()
	}

	if checkOwner && !validator.Owner.Equals(owner) {
		return validator, types.ErrOwnerNotMatch()
	}

	if checkStatus {
		if validator.Status != expectedStatus {
			return validator, fmt.Errorf("validator status not match. except: %d, actual:%d", expectedStatus, validator.Status)
		}
	}
	return validator, nil
}
