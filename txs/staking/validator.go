package staking

import (
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
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
	if len(tx.Name) == 0 {
		return errors.New("Name is empty")
	}

	if tx.PubKey == nil {
		return errors.New("PubKey is empty")
	}

	if len(tx.Owner) == 0 {
		return errors.New("Owner is empty")
	}

	if tx.BondTokens <= 0 {
		return errors.New("BondToken lte zero")
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	owner := accountMapper.GetAccount(tx.Owner)
	if nil == owner {
		return errors.New("Owner not exists")
	}
	ownerAccount := owner.(*account.QOSAccount)
	if !ownerAccount.EnoughOfQOS(btypes.NewInt(int64(tx.BondTokens))) {
		return errors.New("Owner has no enough token")
	}

	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	if mapper.Exists(tx.PubKey.Address().Bytes()) {
		return errors.New("validator already exists")
	}
	if mapper.ExistsWithOwner(tx.Owner) {
		return fmt.Errorf("owner %s already bind a validator", tx.Owner)
	}

	return nil
}

func (tx *TxCreateValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	accMapper := baseabci.GetAccountMapper(ctx)
	// 扣除owner等量QOS
	owner := accMapper.GetAccount(tx.Owner).(*account.QOSAccount)
	owner.MustMinusQOS(btypes.NewInt(int64(tx.BondTokens)))
	accMapper.SetAccount(owner)

	validator := types.Validator{
		Name:            tx.Name,
		Owner:           tx.Owner,
		ValidatorPubKey: tx.PubKey,
		BondTokens:      tx.BondTokens,
		Description:     tx.Description,
		Status:          types.Active,
		BondHeight:      uint64(ctx.BlockHeight()),
	}
	validatorMapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	validatorMapper.CreateValidator(validator)

	return btypes.Result{Code: btypes.ABCICodeOK}, nil
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
		return errors.New("Owner is empty")
	}

	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return errors.New("validator not exists")
	}
	if validator.Status != types.Active {
		return errors.New("validator not active")
	}

	return nil
}

func (tx *TxRevokeValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeInternal)}, nil
	}
	mapper.MakeValidatorInactive(validator.ValidatorPubKey.Address().Bytes(), uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.Revoke)

	return btypes.Result{Code: btypes.ABCICodeOK}, nil
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
		return errors.New("Owner is empty")
	}

	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return errors.New("validator not exists")
	}
	if validator.Status == types.Active {
		return errors.New("validator is active")
	}

	return nil
}

func (tx *TxActiveValidator) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	validator, exists := mapper.GetValidatorByOwner(tx.Owner)
	if !exists {
		return btypes.Result{Code: btypes.ABCICodeType(btypes.CodeInternal)}, nil
	}
	mapper.MakeValidatorActive(validator.ValidatorPubKey.Address().Bytes())

	voteInfoMapper := ctx.Mapper(VoteInfoMapperName).(*VoteInfoMapper)
	voteInfo := types.NewValidatorVoteInfo(validator.BondHeight+1, 0, 0)
	voteInfoMapper.ResetValidatorVoteInfo(validator.ValidatorPubKey.Address().Bytes(), voteInfo)

	return btypes.Result{Code: btypes.ABCICodeOK}, nil
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
