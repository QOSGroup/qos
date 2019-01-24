package stake

import (
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
	MaxDescriptionLen = 1000
)

type TxCreateValidator struct {
	Name        string
	Owner       btypes.Address //操作者, 作为self delegator
	PubKey      crypto.PubKey  //validator公钥
	BondTokens  uint64         //绑定Token数量
	IsCompound  bool
	Description string
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(name string, owner btypes.Address, pubKey crypto.PubKey, bondTokens uint64, isCompound bool, description string) *TxCreateValidator {
	return &TxCreateValidator{
		Name:        name,
		Owner:       owner,
		PubKey:      pubKey,
		BondTokens:  bondTokens,
		IsCompound:  isCompound,
		Description: description,
	}
}

func (tx *TxCreateValidator) ValidateData(ctx context.Context) (err error) {
	if len(tx.Name) == 0 ||
		len(tx.Name) > MaxNameLen ||
		tx.PubKey == nil ||
		len(tx.Description) > MaxDescriptionLen ||
		len(tx.Owner) == 0 ||
		tx.BondTokens == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	err = validateQOSAccount(ctx, tx.Owner, tx.BondTokens)
	if nil != err {
		return err
	}

	mapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	if mapper.Exists(tx.PubKey.Address().Bytes()) {
		return ErrValidatorExists(DefaultCodeSpace, "")
	}
	if mapper.ExistsWithOwner(tx.Owner) {
		return ErrOwnerHasValidator(DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	err := eco.DecrAccountQOS(ctx, tx.Owner, btypes.NewInt(int64(tx.BondTokens)))
	if err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	validator := ecotypes.Validator{
		Name:            tx.Name,
		Owner:           tx.Owner,
		ValidatorPubKey: tx.PubKey,
		BondTokens:      tx.BondTokens,
		Description:     tx.Description,
		Status:          ecotypes.Active,
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
	distributionMapper.InitDelegatorIncomeInfo(valAddr, delegatorAddr, tx.BondTokens, validator.BondHeight)

	validatorMapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	validatorMapper.CreateValidator(validator)

	return btypes.Result{Code: btypes.CodeOK}, nil
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
	ret = append(ret, tx.Name...)
	ret = append(ret, tx.Owner...)
	ret = append(ret, tx.PubKey.Bytes()...)
	ret = append(ret, btypes.Int2Byte(int64(tx.BondTokens))...)
	ret = append(ret, tx.Description...)

	return
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
	mapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	valAddr := validator.GetValidatorAddress()
	delegatorAddr := tx.Owner
	mapper.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), ecotypes.Revoke)

	//更新owner对应的delegator的token数量
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	distributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, uint64(0), uint64(ctx.BlockHeight()))

	return btypes.Result{Code: btypes.CodeOK}, nil
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

	validator, err := validateValidator(ctx, tx.Owner, true, ecotypes.Inactive, true)
	if nil != err {
		return err
	}

	//处于inactive状态的validator,当前汇总收益应该为0
	valAddr := validator.GetValidatorAddress()
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	vcps, _ := distributionMapper.GetValidatorCurrentPeriodSummary(valAddr)

	if vcps.Fees.NilToZero().GT(btypes.ZeroInt()) {
		return ErrCodeValidatorInactiveIncome(DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxActiveValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(ecotypes.ValidatorMapperName).(*ecomapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}

	valAddr := validator.GetValidatorAddress()
	delegatorAddr := tx.Owner
	mapper.MakeValidatorActive(valAddr)

	voteInfoMapper := ctx.Mapper(ecomapper.VoteInfoMapperName).(*ecomapper.VoteInfoMapper)
	voteInfo := ecotypes.NewValidatorVoteInfo(validator.BondHeight+1, 0, 0)
	voteInfoMapper.ResetValidatorVoteInfo(validator.ValidatorPubKey.Address().Bytes(), voteInfo)

	//更新owner对应的delegator的bondtokens
	delegationMapper := ecomapper.GetDelegationMapper(ctx)
	info, _ := delegationMapper.GetDelegationInfo(delegatorAddr, valAddr)

	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	distributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, info.Amount, uint64(ctx.BlockHeight()))

	return btypes.Result{Code: btypes.CodeOK}, nil
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
		qosAccount := acc.(*types.QOSAccount)
		if !qosAccount.EnoughOfQOS(btypes.NewInt(int64(toPay))) {
			return ErrOwnerNoEnoughToken(DefaultCodeSpace, "No enough QOS in account: "+addr.String())
		}
	}
	return nil
}

func validateValidator(ctx context.Context, ownerAddr btypes.Address, checkActive bool, expectingActiveCode int8, checkJail bool) (validator ecotypes.Validator, err error) {
	valMapper := ecomapper.GetValidatorMapper(ctx)
	validator, exists := valMapper.GetValidatorByOwner(ownerAddr)
	if !exists {
		return validator, ErrValidatorNotExists(DefaultCodeSpace, ownerAddr.String()+" does't have validator.")
	}
	if checkActive {
		if expectingActiveCode == ecotypes.Inactive && validator.Status == ecotypes.Active {
			return validator, ErrValidatorIsActive(DefaultCodeSpace, "Owner's Validator "+ownerAddr.String()+" is active")
		}
		if expectingActiveCode == ecotypes.Active && validator.Status == ecotypes.Inactive {
			return validator, ErrValidatorIsActive(DefaultCodeSpace, "Owner's Validator "+ownerAddr.String()+" is inactive")
		}
	}
	if checkJail {
		// TODO: block jailed validator
	}
	return validator, nil
}
