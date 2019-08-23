package txs

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
)

const GasForTransfer = uint64(0.018*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 0.018 QOS

type TxTransfer struct {
	Senders   types.TransItems `json:"senders"`   // 发送集合
	Receivers types.TransItems `json:"receivers"` // 接收集合
}

// 数据校验
func (tx TxTransfer) ValidateData(ctx context.Context) error {
	if valid, err := tx.Senders.IsValid(); !valid {
		return types.ErrInvalidInput(types.DefaultCodeSpace, err.Error())
	}
	if valid, err := tx.Receivers.IsValid(); !valid {
		return types.ErrInvalidInput(types.DefaultCodeSpace, err.Error())
	}

	if valid, err := tx.Senders.Match(tx.Receivers); !valid {
		return types.ErrInvalidInput(types.DefaultCodeSpace, err.Error())
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	for _, sender := range tx.Senders {
		a := accountMapper.GetAccount(sender.Address)
		if a == nil {
			return types.ErrSenderAccountNotExists(types.DefaultCodeSpace, "")
		}
		acc := a.(*qtypes.QOSAccount)
		if !acc.EnoughOf(sender.QOS, sender.QSCs) {
			return types.ErrSenderAccountCoinsNotEnough(types.DefaultCodeSpace, "")
		}
	}

	return nil
}

// 转账
func (tx TxTransfer) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)

	result.Events = btypes.Events{
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	for _, sender := range tx.Senders {
		acc := accountMapper.GetAccount(sender.Address).(*qtypes.QOSAccount)
		acc.MustMinus(sender.QOS, sender.QSCs)
		accountMapper.SetAccount(acc)
		result.Events = result.Events.AppendEvent(btypes.NewEvent(types.EventTypeSend,
			btypes.NewAttribute(types.AttributeKeyAddress, sender.Address.String()),
			btypes.NewAttribute(types.AttributeKeyQOS, sender.QOS.String()),
			btypes.NewAttribute(types.AttributeKeyQSCs, sender.QSCs.String()),
		),
		)
	}
	for _, receiver := range tx.Receivers {
		a := accountMapper.GetAccount(receiver.Address)
		var acc *qtypes.QOSAccount
		if a != nil {
			acc = a.(*qtypes.QOSAccount)
		} else {
			acc = qtypes.NewQOSAccountWithAddress(receiver.Address)
		}
		acc.MustPlus(receiver.QOS, receiver.QSCs)
		accountMapper.SetAccount(acc)
		result.Events = result.Events.AppendEvent(btypes.NewEvent(types.EventTypeReceive,
			btypes.NewAttribute(types.AttributeKeyAddress, receiver.Address.String()),
			btypes.NewAttribute(types.AttributeKeyQOS, receiver.QOS.String()),
			btypes.NewAttribute(types.AttributeKeyQSCs, receiver.QSCs.String()),
		),
		)
	}

	return btypes.Result{Code: btypes.CodeOK, Events: result.Events}, nil
}

// 所有Senders
func (tx TxTransfer) GetSigner() []btypes.Address {
	addrs := make([]btypes.Address, 0)
	for _, sender := range tx.Senders {
		addrs = append(addrs, sender.Address)
	}

	return addrs
}

// Gas
func (tx TxTransfer) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForTransfer))
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
