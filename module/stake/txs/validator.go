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

	GasForCreateValidator = uint64(1.8*qtypes.QOSUnit) * qtypes.GasPerUnitCost  // 1.8 QOS
	GasForModifyValidator = uint64(0.18*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 0.18 QOS
	GasForRevokeValidator = uint64(18*qtypes.QOSUnit) * qtypes.GasPerUnitCost   // 18 QOS
)

type TxCreateValidator struct {
	Operator    btypes.AccAddress     //操作者, self delegator
	ConsPubKey  crypto.PubKey         //validator公钥
	BondTokens  uint64                //绑定Token数量
	IsCompound  bool                  //周期收益是否复投
	Description types.Description     //描述信息
	Commission  types.CommissionRates //佣金比例
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(operator btypes.AccAddress, bech32ConPubKey crypto.PubKey, bondTokens uint64, isCompound bool, description types.Description, commission types.CommissionRates) *TxCreateValidator {
	return &TxCreateValidator{
		Operator:    operator,
		ConsPubKey:  bech32ConPubKey,
		BondTokens:  bondTokens,
		IsCompound:  isCompound,
		Description: description,
		Commission:  commission,
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
		tx.BondTokens <= 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	err = tx.Commission.Validate()
	if err != nil {
		return
	}

	err = validateQOSAccount(ctx, tx.Operator, tx.BondTokens)
	if nil != err {
		return err
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
		MinPeriod:       uint64(0),
		BondHeight:      uint64(ctx.BlockHeight()),
		Commission:      types.NewCommissionWithTime(tx.Commission.Rate, tx.Commission.MaxRate, tx.Commission.MaxChangeRate, ctx.BlockHeader().Time.UTC()),
	}

	delegatorAddr := validator.Owner

	// 创建validator
	sm := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	sm.CreateValidator(validator)

	sm.AfterValidatorCreated(ctx, valAddr)

	// 初始化self-delegation
	delegationInfo := types.NewDelegationInfo(delegatorAddr, valAddr, tx.BondTokens, tx.IsCompound)
	sm.Delegate(ctx, delegationInfo, false)

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
	return btypes.NewInt(int64(GasForCreateValidator))
}

func (tx *TxCreateValidator) GetGasPayer() btypes.AccAddress {
	return tx.Operator
}

func (tx *TxCreateValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

type TxModifyValidator struct {
	ValidatorAddr  btypes.ValAddress //节点所有账户
	Description    types.Description //描述信息
	CommissionRate *qtypes.Dec       //佣金比例
}

var _ txs.ITx = (*TxModifyValidator)(nil)

func NewModifyValidatorTx(validatorAddr btypes.ValAddress, description types.Description, commissionRate *qtypes.Dec) *TxModifyValidator {
	return &TxModifyValidator{
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

	mapper := mapper.GetMapper(ctx)

	validator, exists := mapper.GetValidator(tx.ValidatorAddr)
	if !exists {
		return types.ErrOwnerNotExists(types.DefaultCodeSpace, fmt.Sprintf("%s has no validator", tx.ValidatorAddr.String()))
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
	return []btypes.AccAddress{btypes.AccAddress(tx.ValidatorAddr)}
}

func (tx *TxModifyValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForModifyValidator))
}

func (tx *TxModifyValidator) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.ValidatorAddr)
}

func (tx *TxModifyValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

type TxRevokeValidator struct {
	ValidatorAddr btypes.ValAddress //操作者
}

var _ txs.ITx = (*TxRevokeValidator)(nil)

func NewRevokeValidatorTx(validatorAddr btypes.ValAddress) *TxRevokeValidator {
	return &TxRevokeValidator{
		ValidatorAddr: validatorAddr,
	}
}

func (tx *TxRevokeValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.ValidatorAddr) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	_, err = validateValidator(ctx, tx.ValidatorAddr, true, types.Active, true)
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
	mapper.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.Revoke)

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
	return []btypes.AccAddress{btypes.AccAddress(tx.ValidatorAddr)}
}

func (tx *TxRevokeValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForRevokeValidator))
}

func (tx *TxRevokeValidator) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.ValidatorAddr)
}

func (tx *TxRevokeValidator) GetSignData() (ret []byte) {
	ret = append(ret, tx.ValidatorAddr...)

	return
}

type TxActiveValidator struct {
	ValidatorAddr btypes.ValAddress //操作者
	BondTokens    uint64            //绑定Token数量
}

var _ txs.ITx = (*TxActiveValidator)(nil)

func NewActiveValidatorTx(validatorAddr btypes.ValAddress, bondTokens uint64) *TxActiveValidator {
	return &TxActiveValidator{
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

	_, err = validateValidator(ctx, tx.ValidatorAddr, true, types.Inactive, true)
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
	if tx.BondTokens > 0 {

		// 从 delegator 账户, 扣去增加的自委托token数量
		delegator.MustMinusQOS(btypes.NewInt(int64(tx.BondTokens)))
		accountMapper.SetAccount(delegator)

		// 获取 delegationInfo
		delegationInfo, exists := stakeMapper.GetDelegationInfo(delegatorAddr, validatorAddr)
		if !exists {
			return btypes.Result{Code: btypes.CodeInternal}, nil
		}

		// 修改 delegationInfo 中的token amount, 并保存.
		delegationInfo.Amount += tx.BondTokens
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
	return []btypes.AccAddress{btypes.AccAddress(tx.ValidatorAddr)}
}

func (tx *TxActiveValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxActiveValidator) GetGasPayer() btypes.AccAddress {
	return btypes.AccAddress(tx.ValidatorAddr)
}

func (tx *TxActiveValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

func validateQOSAccount(ctx context.Context, addr btypes.AccAddress, toPay uint64) error {
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	acc := accountMapper.GetAccount(addr)

	qtypes.AssertUint64NotOverflow(toPay)

	if toPay > 0 {
		if acc != nil {
			qosAccount := acc.(*qtypes.QOSAccount)
			toPay := btypes.NewInt(int64(toPay))

			//溢出校验
			if toPay.LT(btypes.ZeroInt()) {
				return types.ErrInvalidInput(types.DefaultCodeSpace, "Bind tokens is lt zero: "+addr.String())
			}

			if !qosAccount.EnoughOfQOS(toPay) {
				return types.ErrOwnerNoEnoughToken(types.DefaultCodeSpace, "No enough QOS in account: "+addr.String())
			}
		} else {
			return types.ErrOwnerNoEnoughToken(types.DefaultCodeSpace, "account not exists: "+addr.String())
		}
	}
	return nil
}

func validateValidator(ctx context.Context, validatorAddr btypes.ValAddress, checkStatus bool, expectedStatus int8, checkJail bool) (validator types.Validator, err error) {
	valMapper := mapper.GetMapper(ctx)
	validator, exists := valMapper.GetValidator(validatorAddr)
	if !exists {
		return validator, types.ErrValidatorNotExists(types.DefaultCodeSpace, validatorAddr.String()+" does't have validator.")
	}
	if checkStatus {
		if validator.Status != expectedStatus {
			return validator, fmt.Errorf("validator status not match. except: %d, actual:%d", expectedStatus, validator.Status)
		}
	}
	if checkJail {
		// TODO: block jailed validator
	}
	return validator, nil
}
