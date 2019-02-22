package transfer

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	transfertypes "github.com/QOSGroup/qos/module/transfer/types"
	"github.com/QOSGroup/qos/types"
)

type TxTransfer struct {
	Senders   transfertypes.TransItems `json:"senders"`   // 发送集合
	Receivers transfertypes.TransItems `json:"receivers"` // 接收集合
}

// 数据校验
func (tx TxTransfer) ValidateData(ctx context.Context) error {
	if valid, err := tx.Senders.IsValid(); !valid {
		return ErrInvalidInput(DefaultCodeSpace, err.Error())
	}
	if valid, err := tx.Receivers.IsValid(); !valid {
		return ErrInvalidInput(DefaultCodeSpace, err.Error())
	}

	if valid, err := tx.Senders.Match(tx.Receivers); !valid {
		return ErrInvalidInput(DefaultCodeSpace, err.Error())
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	for _, sender := range tx.Senders {
		a := accountMapper.GetAccount(sender.Address)
		if a == nil {
			return ErrSenderAccountNotExists(DefaultCodeSpace, "")
		}
		acc := a.(*types.QOSAccount)
		if !acc.EnoughOf(sender.QOS, sender.QSCs) {
			return ErrSenderAccountCoinsNotEnough(DefaultCodeSpace, "")
		}
	}

	return nil
}

// 转账
func (tx TxTransfer) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)

	for _, sender := range tx.Senders {
		acc := accountMapper.GetAccount(sender.Address).(*types.QOSAccount)
		acc.MustMinus(sender.QOS, sender.QSCs)
		accountMapper.SetAccount(acc)
	}
	for _, receiver := range tx.Receivers {
		a := accountMapper.GetAccount(receiver.Address)
		var acc *types.QOSAccount
		if a != nil {
			acc = a.(*types.QOSAccount)
		} else {
			acc = types.NewQOSAccountWithAddress(receiver.Address)
		}
		acc.MustPlus(receiver.QOS, receiver.QSCs)
		accountMapper.SetAccount(acc)
	}

	return btypes.Result{Code: btypes.CodeOK}, nil
}

// 所有Senders
func (tx TxTransfer) GetSigner() []btypes.Address {
	addrs := make([]btypes.Address, 0)
	for _, sender := range tx.Senders {
		addrs = append(addrs, sender.Address)
	}

	return addrs
}

// Gas TODO
func (tx TxTransfer) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Senders[0]
func (tx TxTransfer) GetGasPayer() btypes.Address {
	return tx.Senders[0].Address
}

// 签名字节
func (tx TxTransfer) GetSignData() (ret []byte) {
	for _, sender := range tx.Senders {
		ret = append(ret, sender.Address...)
		ret = append(ret, (sender.QOS.NilToZero()).String()...)
		ret = append(ret, sender.QSCs.String()...)
	}
	for _, receiver := range tx.Receivers {
		ret = append(ret, receiver.Address...)
		ret = append(ret, (receiver.QOS.NilToZero()).String()...)
		ret = append(ret, receiver.QSCs.String()...)
	}

	return ret
}
