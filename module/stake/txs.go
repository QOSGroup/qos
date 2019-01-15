package stake

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	stakemapper "github.com/QOSGroup/qos/module/eco/mapper"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MaxNameLen        = 300
	MaxDescriptionLen = 1000
)

type TxCreateValidator struct {
	Name        string
	Owner       btypes.Address //操作者
	PubKey      crypto.PubKey  //validator公钥
	BondTokens  uint64         //绑定Token数量
	Description string
}

var _ txs.ITx = (*TxCreateValidator)(nil)

func NewCreateValidatorTx(name string, owner btypes.Address, pubKey crypto.PubKey, bondTokens uint64, description string) *TxCreateValidator {
	return &TxCreateValidator{
		Name:        name,
		Owner:       owner,
		PubKey:      pubKey,
		BondTokens:  bondTokens,
		Description: description,
	}
}

func (tx *TxCreateValidator) ValidateData(ctx context.Context) error {
	if len(tx.Name) == 0 ||
		len(tx.Name) > MaxNameLen ||
		tx.PubKey == nil ||
		len(tx.Description) > MaxDescriptionLen ||
		len(tx.Owner) == 0 ||
		tx.BondTokens == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	owner := accountMapper.GetAccount(tx.Owner)
	if nil == owner {
		return ErrOwnerNotExists(DefaultCodeSpace, "")
	}
	ownerAccount := owner.(*types.QOSAccount)
	if !ownerAccount.EnoughOfQOS(btypes.NewInt(int64(tx.BondTokens))) {
		return ErrOwnerNoEnoughToken(DefaultCodeSpace, "")
	}

	mapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	if mapper.Exists(tx.PubKey.Address().Bytes()) {
		return ErrValidatorExists(DefaultCodeSpace, "")
	}
	if mapper.ExistsWithOwner(tx.Owner) {
		return ErrOwnerHasValidator(DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	accMapper := baseabci.GetAccountMapper(ctx)
	// 扣除owner等量QOS
	owner := accMapper.GetAccount(tx.Owner).(*types.QOSAccount)
	owner.MustMinusQOS(btypes.NewInt(int64(tx.BondTokens)))
	accMapper.SetAccount(owner)

	validator := staketypes.Validator{
		Name:            tx.Name,
		Owner:           tx.Owner,
		ValidatorPubKey: tx.PubKey,
		BondTokens:      tx.BondTokens,
		Description:     tx.Description,
		Status:          staketypes.Active,
		BondHeight:      uint64(ctx.BlockHeight()),
	}
	validatorMapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
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

func (tx *TxRevokeValidator) ValidateData(ctx context.Context) error {

	if len(tx.Owner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	mapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return ErrValidatorNotExists(DefaultCodeSpace, "")
	}
	if validator.Status != staketypes.Active {
		return ErrValidatorIsInactive(DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxRevokeValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}
	mapper.MakeValidatorInactive(validator.ValidatorPubKey.Address().Bytes(), uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), staketypes.Revoke)

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

func (tx *TxActiveValidator) ValidateData(ctx context.Context) error {

	if len(tx.Owner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	mapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return ErrValidatorNotExists(DefaultCodeSpace, "")
	}
	if validator.Status == staketypes.Active {
		return ErrValidatorIsActive(DefaultCodeSpace, "")
	}

	return nil
}

func (tx *TxActiveValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal}, nil
	}
	mapper.MakeValidatorActive(validator.ValidatorPubKey.Address().Bytes())

	voteInfoMapper := ctx.Mapper(stakemapper.VoteInfoMapperName).(*stakemapper.VoteInfoMapper)
	voteInfo := staketypes.NewValidatorVoteInfo(validator.BondHeight+1, 0, 0)
	voteInfoMapper.ResetValidatorVoteInfo(validator.ValidatorPubKey.Address().Bytes(), voteInfo)

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


type TxCreateDelegation struct {
	Delegator btypes.Address
	Validator btypes.Address
	Amount uint64
	isCompound bool
}

var _ txs.ITx = (*TxCreateDelegation)(nil)

func (tx *TxCreateDelegation) ValidateData(ctx context.Context) error {

	if len(tx.Delegator) == 0 || len(tx.Validator) == 0{
		return ErrInvalidInput(DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	// TODO: 是否应该在tx里做这种检查
	if tx.Amount <= 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Delegation amount must be a positive integer.")
	}

	valMapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	validator, exists := valMapper.GetValidatorByOwner(tx.Validator)
	if !exists {
		return ErrValidatorNotExists(DefaultCodeSpace, "")
	}
	if validator.Status == staketypes.Active {
		return ErrValidatorIsActive(DefaultCodeSpace, "")
	}
	// TODO: block jailed validator

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	delAcc := accountMapper.GetAccount(tx.Delegator)
	if nil == delAcc {
		return ErrOwnerNotExists(DefaultCodeSpace, "Delegator's address doesn't exist.")
	}
	delegatorAccount := delAcc.(*types.QOSAccount)
	if !delegatorAccount.EnoughOfQOS(btypes.NewInt(int64(tx.Amount))) {
		return ErrOwnerNoEnoughToken(DefaultCodeSpace, "No enough QOS in delegator's account.")
	}

	valAcc := accountMapper.GetAccount(tx.Validator)
	if nil == valAcc {
		return ErrOwnerNotExists(DefaultCodeSpace, "Validator's address doesn't exist.")
	}
	return nil
}

func (tx *TxCreateDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(staketypes.ValidatorMapperName).(*stakemapper.ValidatorMapper)
	_, exists := mapper.GetValidatorByOwner(tx.Validator)
	if !exists {
		return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeInternal)}, nil
	}
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	delAcc := accountMapper.GetAccount(tx.Delegator)
	if nil == delAcc {
		return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeInternal)}, nil
	}
	delegatorAccount := delAcc.(*types.QOSAccount)
	if !delegatorAccount.EnoughOfQOS(btypes.NewInt(int64(tx.Amount))) {
		return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeInternal)}, nil
	}

	delMapper := stakemapper.GetDelegationMapper(ctx)
	delegationInfo, exists := delMapper.GetDelegationInfo(tx.Validator, tx.Delegator)
	if exists {
		delegationInfo.Amount += tx.Amount
		delegationInfo.IsCompound = tx.isCompound
		delMapper.SetDelegationInfo(delegationInfo)
		return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeOK)}, nil
	}
	delegationInfo = staketypes.NewDelegationInfo(tx.Delegator, tx.Validator, tx.Amount, tx.isCompound)
	delMapper.SetDelegationInfo(delegationInfo)

	return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeOK)}, nil
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
	return
}