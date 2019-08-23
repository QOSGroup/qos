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
	Owner       btypes.Address        //操作者, self delegator
	PubKey      crypto.PubKey         //validator公钥
	BondTokens  uint64                //绑定Token数量
	IsCompound  bool                  //周期收益是否复投
	Description types.Description     //描述信息
	Commission  types.CommissionRates //佣金比例
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(owner btypes.Address, pubKey crypto.PubKey, bondTokens uint64, isCompound bool, description types.Description, commission types.CommissionRates) *TxCreateValidator {
	return &TxCreateValidator{
		Owner:       owner,
		PubKey:      pubKey,
		BondTokens:  bondTokens,
		IsCompound:  isCompound,
		Description: description,
		Commission:  commission,
	}
}

func (tx *TxCreateValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Description.Moniker) == 0 ||
		len(tx.Description.Moniker) > MaxNameLen ||
		tx.PubKey == nil ||
		len(tx.Description.Logo) > MaxLinkLen ||
		len(tx.Description.Website) > MaxLinkLen ||
		len(tx.Description.Details) > MaxDescriptionLen ||
		len(tx.Owner) == 0 ||
		tx.BondTokens == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	err = tx.Commission.Validate()
	if err != nil {
		return
	}

	err = validateQOSAccount(ctx, tx.Owner, tx.BondTokens)
	if nil != err {
		return err
	}

	mapper := mapper.GetMapper(ctx)
	if mapper.Exists(tx.PubKey.Address().Bytes()) {
		return types.ErrValidatorExists(types.DefaultCodeSpace, "")
	}
	if mapper.ExistsWithOwner(tx.Owner) {
		return types.ErrOwnerHasValidator(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	result = btypes.Result{Code: btypes.CodeOK}

	validator := types.Validator{
		Owner:           tx.Owner,
		ValidatorPubKey: tx.PubKey,
		BondTokens:      tx.BondTokens,
		Description:     tx.Description,
		Status:          types.Active,
		MinPeriod:       uint64(0),
		BondHeight:      uint64(ctx.BlockHeight()),
		Commission:      types.NewCommissionWithTime(tx.Commission.Rate, tx.Commission.MaxRate, tx.Commission.MaxChangeRate, ctx.BlockHeader().Time.UTC()),
	}

	valAddr := validator.GetValidatorAddress()
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
			btypes.NewAttribute(types.AttributeKeyOwner, tx.Owner.String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxCreateValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxCreateValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForCreateValidator))
}

func (tx *TxCreateValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxCreateValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

type TxModifyValidator struct {
	Owner          btypes.Address    //节点所有账户
	Description    types.Description //描述信息
	CommissionRate *qtypes.Dec       //佣金比例
}

var _ txs.ITx = (*TxModifyValidator)(nil)

func NewModifyValidatorTx(owner btypes.Address, description types.Description, commissionRate *qtypes.Dec) *TxModifyValidator {
	return &TxModifyValidator{
		Owner:          owner,
		Description:    description,
		CommissionRate: commissionRate,
	}
}

func (tx *TxModifyValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Description.Moniker) > MaxNameLen ||
		len(tx.Description.Logo) > MaxLinkLen ||
		len(tx.Description.Website) > MaxLinkLen ||
		len(tx.Description.Details) > MaxDescriptionLen ||
		len(tx.Owner) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	mapper := mapper.GetMapper(ctx)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return types.ErrOwnerHasValidator(types.DefaultCodeSpace, fmt.Sprintf("%s has no validator", tx.Owner.String()))
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
	validator, exists := validatorMapper.GetValidatorByOwner(tx.Owner)
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
			btypes.NewAttribute(types.AttributeKeyOwner, tx.Owner.String()),
			btypes.NewAttribute(types.AttributeKeyDelegator, tx.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxModifyValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxModifyValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForModifyValidator))
}

func (tx *TxModifyValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxModifyValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

type TxRevokeValidator struct {
	Owner btypes.Address //操作者
}

var _ txs.ITx = (*TxRevokeValidator)(nil)

func NewRevokeValidatorTx(owner btypes.Address) *TxRevokeValidator {
	return &TxRevokeValidator{
		Owner: owner,
	}
}

func (tx *TxRevokeValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Owner) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	_, err = validateValidator(ctx, tx.Owner, true, types.Active, true)
	if nil != err {
		return err
	}

	return nil
}

func (tx *TxRevokeValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	valAddr := validator.GetValidatorAddress()
	mapper.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.Revoke)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeRevokeValidator,
			btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			btypes.NewAttribute(types.AttributeKeyOwner, tx.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxRevokeValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxRevokeValidator) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForRevokeValidator))
}

func (tx *TxRevokeValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxRevokeValidator) GetSignData() (ret []byte) {
	ret = append(ret, tx.Owner...)

	return
}

type TxActiveValidator struct {
	Owner      btypes.Address //操作者
	BondTokens uint64         //绑定Token数量
}

var _ txs.ITx = (*TxActiveValidator)(nil)

func NewActiveValidatorTx(owner btypes.Address, bondTokens uint64) *TxActiveValidator {
	return &TxActiveValidator{
		Owner:      owner,
		BondTokens: bondTokens,
	}
}

func (tx *TxActiveValidator) ValidateData(ctx context.Context) (err error) {

	if len(tx.Owner) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	err = validateQOSAccount(ctx, tx.Owner, tx.BondTokens)
	if nil != err {
		return err
	}

	_, err = validateValidator(ctx, tx.Owner, true, types.Inactive, true)
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
	validator, exists := stakeMapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	// 激活 Validator
	validatorAddr := validator.GetValidatorAddress()
	stakeMapper.MakeValidatorActive(validatorAddr, tx.BondTokens)

	// 重置 ValidatorVoteInfo
	voteInfo := types.NewValidatorVoteInfo(validator.BondHeight+1, 0, 0)
	stakeMapper.ResetValidatorVoteInfo(validator.ValidatorPubKey.Address().Bytes(), voteInfo)

	// 获取 Owner 对应的 Delegator 账户
	delegatorAddr := tx.Owner
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
			btypes.NewAttribute(types.AttributeKeyOwner, tx.Owner.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx *TxActiveValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxActiveValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxActiveValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxActiveValidator) GetSignData() (ret []byte) {
	return Cdc.MustMarshalJSON(*tx)
}

func validateQOSAccount(ctx context.Context, addr btypes.Address, toPay uint64) error {
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	acc := accountMapper.GetAccount(addr)

	if toPay > 0 {
		if acc != nil {
			qosAccount := acc.(*qtypes.QOSAccount)
			if !qosAccount.EnoughOfQOS(btypes.NewInt(int64(toPay))) {
				return types.ErrOwnerNoEnoughToken(types.DefaultCodeSpace, "No enough QOS in account: "+addr.String())
			}
		} else {
			return types.ErrOwnerNoEnoughToken(types.DefaultCodeSpace, "account not exists: "+addr.String())
		}
	}
	return nil
}

func validateValidator(ctx context.Context, ownerAddr btypes.Address, checkStatus bool, expectedStatus int8, checkJail bool) (validator types.Validator, err error) {
	valMapper := mapper.GetMapper(ctx)
	validator, exists := valMapper.GetValidatorByOwner(ownerAddr)
	if !exists {
		return validator, types.ErrValidatorNotExists(types.DefaultCodeSpace, ownerAddr.String()+" does't have validator.")
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
