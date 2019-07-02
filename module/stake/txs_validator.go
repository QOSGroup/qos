package stake

import (
	"fmt"

	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MaxNameLen        = 300
	MaxLinkLen        = 255
	MaxDescriptionLen = 1000
)

type TxCreateValidator struct {
	Owner       btypes.Address       //操作者, self delegator
	PubKey      crypto.PubKey        //validator公钥
	BondTokens  uint64               //绑定Token数量
	IsCompound  bool                 //周期收益是否复投
	Description ecotypes.Description //描述信息
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(owner btypes.Address, pubKey crypto.PubKey, bondTokens uint64, isCompound bool, description ecotypes.Description) *TxCreateValidator {
	return &TxCreateValidator{
		Owner:       owner,
		PubKey:      pubKey,
		BondTokens:  bondTokens,
		IsCompound:  isCompound,
		Description: description,
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
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	err = validateQOSAccount(ctx, tx.Owner, tx.BondTokens)
	if nil != err {
		return err
	}

	mapper := ecomapper.GetValidatorMapper(ctx)
	if mapper.Exists(tx.PubKey.Address().Bytes()) {
		return ErrValidatorExists(DefaultCodeSpace, "")
	}
	if mapper.ExistsWithOwner(tx.Owner) {
		return ErrOwnerHasValidator(DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	result = btypes.Result{Code: btypes.CodeOK}

	err := eco.DecrAccountQOS(ctx, tx.Owner, btypes.NewInt(int64(tx.BondTokens)))
	if err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	validator := ecotypes.Validator{
		Owner:           tx.Owner,
		ValidatorPubKey: tx.PubKey,
		BondTokens:      tx.BondTokens,
		Description:     tx.Description,
		Status:          ecotypes.Active,
		MinPeriod:       uint64(0),
		BondHeight:      uint64(ctx.BlockHeight()),
	}

	valAddr := validator.GetValidatorAddress()
	delegatorAddr := validator.Owner

	//初始化validator self delegate 数据
	delegationMapper := ecomapper.GetDelegationMapper(ctx)
	delegationInfo := ecotypes.NewDelegationInfo(delegatorAddr, valAddr, tx.BondTokens, tx.IsCompound)
	delegationMapper.SetDelegationInfo(delegationInfo)

	//初始化validator distribution数据
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	distributionMapper.InitValidatorPeriodSummaryInfo(valAddr)
	distributionMapper.InitDelegatorIncomeInfo(ctx, valAddr, delegatorAddr, tx.BondTokens, validator.BondHeight)

	validatorMapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	validatorMapper.CreateValidator(validator)

	result.Tags = btypes.NewTags(btypes.TagAction, TagActionCreateValidator,
		TagValidator, valAddr.String(),
		TagOwner, tx.Owner.String(),
		TagDelegator, tx.Owner.String())

	return
}

func (tx *TxCreateValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxCreateValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxCreateValidator) GetSignData() (ret []byte) {
	return cdc.MustMarshalJSON(*tx)
}

type TxModifyValidator struct {
	Owner       btypes.Address       //节点所有账户
	Description ecotypes.Description //描述信息
}

var _ txs.ITx = (*TxModifyValidator)(nil)

func NewModifyValidatorTx(owner btypes.Address, description ecotypes.Description) *TxModifyValidator {
	return &TxModifyValidator{
		Owner:       owner,
		Description: description,
	}
}

func (tx *TxModifyValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Description.Moniker) > MaxNameLen ||
		len(tx.Description.Logo) > MaxLinkLen ||
		len(tx.Description.Website) > MaxLinkLen ||
		len(tx.Description.Details) > MaxDescriptionLen ||
		len(tx.Owner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	mapper := ecomapper.GetValidatorMapper(ctx)
	if !mapper.ExistsWithOwner(tx.Owner) {
		return ErrOwnerHasValidator(DefaultCodeSpace, fmt.Sprintf("%s has no validator", tx.Owner.String()))
	}

	return nil
}

func (tx *TxModifyValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	result = btypes.Result{Code: btypes.CodeOK}

	validatorMapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
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
	validatorMapper.Set(ecotypes.BuildValidatorKey(validator.GetValidatorAddress()), validator)

	result.Tags = btypes.NewTags(btypes.TagAction, TagActionModifyValidator,
		TagOwner, tx.Owner.String(),
		TagDelegator, tx.Owner.String())

	return
}

func (tx *TxModifyValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxModifyValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxModifyValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxModifyValidator) GetSignData() (ret []byte) {
	return cdc.MustMarshalJSON(*tx)
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
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	_, err = validateValidator(ctx, tx.Owner, true, ecotypes.Active, true)
	if nil != err {
		return err
	}

	return nil
}

func (tx *TxRevokeValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	mapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	valAddr := validator.GetValidatorAddress()
	mapper.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), ecotypes.Revoke)

	result.Tags = btypes.NewTags(btypes.TagAction, TagActionRevokeValidator,
		TagValidator, valAddr.String(),
		TagOwner, tx.Owner.String())

	return
}

func (tx *TxRevokeValidator) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Owner}
}

func (tx *TxRevokeValidator) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxRevokeValidator) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Owner)
}

func (tx *TxRevokeValidator) GetSignData() (ret []byte) {
	ret = append(ret, tx.Owner...)

	return
}

type TxActiveValidator struct {
	Owner btypes.Address //操作者
}

var _ txs.ITx = (*TxActiveValidator)(nil)

func NewActiveValidatorTx(owner btypes.Address) *TxActiveValidator {
	return &TxActiveValidator{
		Owner: owner,
	}
}

func (tx *TxActiveValidator) ValidateData(ctx context.Context) (err error) {

	if len(tx.Owner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	_, err = validateValidator(ctx, tx.Owner, true, ecotypes.Inactive, true)
	if nil != err {
		return err
	}

	return nil
}

func (tx *TxActiveValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{Code: btypes.CodeOK}

	mapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	valAddr := validator.GetValidatorAddress()
	// delegatorAddr := tx.Owner
	mapper.MakeValidatorActive(valAddr)

	voteInfoMapper := ecomapper.GetVoteInfoMapper(ctx)
	voteInfo := ecotypes.NewValidatorVoteInfo(validator.BondHeight+1, 0, 0)
	voteInfoMapper.ResetValidatorVoteInfo(validator.ValidatorPubKey.Address().Bytes(), voteInfo)

	// 更新owner对应的delegator的bondtokens
	// delegationMapper := ecomapper.GetDelegationMapper(ctx)
	// info, _ := delegationMapper.GetDelegationInfo(delegatorAddr, valAddr)

	// distributionMapper := ecomapper.GetDistributionMapper(ctx)
	// distributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, info.Amount, uint64(ctx.BlockHeight()))

	result.Tags = btypes.NewTags(btypes.TagAction, TagActionActiveValidator,
		TagValidator, valAddr.String(),
		TagOwner, tx.Owner.String())

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
	ret = append(ret, tx.Owner...)

	return
}

func validateQOSAccount(ctx context.Context, addr btypes.Address, toPay uint64) error {
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	acc := accountMapper.GetAccount(addr)

	if toPay > 0 {
		if acc != nil {
			qosAccount := acc.(*types.QOSAccount)
			if !qosAccount.EnoughOfQOS(btypes.NewInt(int64(toPay))) {
				return ErrOwnerNoEnoughToken(DefaultCodeSpace, "No enough QOS in account: "+addr.String())
			}
		} else {
			return ErrOwnerNoEnoughToken(DefaultCodeSpace, "account not exists: "+addr.String())
		}
	}
	return nil
}

func validateValidator(ctx context.Context, ownerAddr btypes.Address, checkStatus bool, expectedStatus int8, checkJail bool) (validator ecotypes.Validator, err error) {
	valMapper := ecomapper.GetValidatorMapper(ctx)
	validator, exists := valMapper.GetValidatorByOwner(ownerAddr)
	if !exists {
		return validator, ErrValidatorNotExists(DefaultCodeSpace, ownerAddr.String()+" does't have validator.")
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
