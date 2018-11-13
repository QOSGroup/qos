package validator

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
)

// 授权 Common 结构
type CreateValidatorTx struct {
	Name   string        `json:"name""`
	PubKey crypto.PubKey `json:"pub_key"`
}

func NewCreateValidatorTx(name string, pubkey crypto.PubKey) CreateValidatorTx {
	return CreateValidatorTx{name, pubkey}
}

func (tx CreateValidatorTx) ValidateData(ctx context.Context) error {
	if len(tx.Name) == 0 {
		return errors.New("Name is empty")
	}

	if tx.PubKey == nil {
		return errors.New("PubKey is empty")
	}

	// Validator必须不存在
	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	if mapper.Exists(tx.PubKey.Address().Bytes()) {
		return errors.New("Validator already exists")
	}

	return nil
}

func (tx CreateValidatorTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	mapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)
	// TODO votingPower 公式
	validator := types.NewValidator(tx.Name, tx.PubKey, 10, ctx.BlockHeight())
	mapper.SaveValidator(validator)
	mapper.SetUpdated(false)
	return btypes.Result{Code: btypes.ABCICodeOK,}, nil
}

func (tx CreateValidatorTx) GetSigner() []btypes.Address {
	return nil
}

func (tx CreateValidatorTx) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

func (tx CreateValidatorTx) GetGasPayer() btypes.Address {
	return tx.PubKey.Address().Bytes()
}

func (tx CreateValidatorTx) GetSignData() (ret []byte) {
	ret = append(ret, tx.Name...)
	ret = append(ret, tx.PubKey.Bytes()...)
	return
}
