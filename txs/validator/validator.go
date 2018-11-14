package validator

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
)

type CreateValidatorTx struct {
	Name       string         `json:"name""`
	ConsPubKey crypto.PubKey  `json:"cons_pubkey"` //参与共识的validator公钥
	Operator   btypes.Address `json:"operator"`    //QOS操作账户地址
}

var _ txs.ITx = (*CreateValidatorTx)(nil)

func NewCreateValidatorTx(name string, consPubKey crypto.PubKey, operator btypes.Address) *CreateValidatorTx {
	return &CreateValidatorTx{
		Name:       name,
		ConsPubKey: consPubKey,
		Operator:   operator,
	}
}

func (tx *CreateValidatorTx) ValidateData(ctx context.Context) error {
	if len(tx.Name) == 0 {
		return errors.New("Name is empty")
	}

	if tx.ConsPubKey == nil {
		return errors.New("ConsPubKey is empty")
	}

	if tx.Operator == nil {
		return errors.New("Operator is empty")
	}

	// Validator必须不存在
	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	if mapper.Exists(tx.ConsPubKey.Address().Bytes()) {
		return errors.New("Validator already exists")
	}

	return nil
}

func (tx *CreateValidatorTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	accMapper := baseabci.GetAccountMapper(ctx)

	// TODO votingPower 公式
	validator := types.NewValidator(tx.Name, tx.ConsPubKey, tx.Operator, 10, ctx.BlockHeight())
	mapper.SaveValidator(validator)
	mapper.SetValidatorChanged()

	acc := accMapper.GetAccount(tx.Operator)
	if acc.GetPubicKey() == nil {
		acc = accMapper.NewAccountWithAddress(tx.Operator)
		accMapper.SetAccount(acc)
	}

	return btypes.Result{Code: btypes.ABCICodeOK}, nil
}

func (tx *CreateValidatorTx) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Operator}
}

func (tx *CreateValidatorTx) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *CreateValidatorTx) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Operator)
}

func (tx *CreateValidatorTx) GetSignData() (ret []byte) {
	ret = append(ret, tx.Name...)
	ret = append(ret, tx.ConsPubKey.Bytes()...)
	ret = append(ret, tx.Operator...)
	return
}
