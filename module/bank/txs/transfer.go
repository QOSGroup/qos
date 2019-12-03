package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
)

const (
	// 转账发送+接收账户总数限制
	MaxTransLen = 500
)

// 转账交易Gas
var GasForTransfer = int64(0.018*qtypes.UnitQOS) * qtypes.UnitQOSGas // 0.018 QOS

type TxTransfer struct {
	Senders   types.TransItems `json:"senders"`   // 发送集合
	Receivers types.TransItems `json:"receivers"` // 接收集合
}

// 基础数据校验
func (tx TxTransfer) ValidateInputs() error {
	// 发送/接收集合基础数据校验
	if err := tx.Senders.Valid(); err != nil {
		return err
	}
	if err := tx.Receivers.Valid(); err != nil {
		return err
	}
	if err := tx.Senders.Match(tx.Receivers); err != nil {
		return err
	}
	if len(tx.Senders)+len(tx.Receivers) > MaxTransLen {
		return types.ErrInvalidInput(fmt.Sprintf("len(senders) + len(receivers) must lte %d", MaxTransLen))
	}

	return nil
}

// 数据校验
func (tx TxTransfer) ValidateData(ctx context.Context) error {
	// 发送/接收集合基础数据校验
	if err := tx.ValidateInputs(); err != nil {
		return err
	}

	// 发送账户余额校验
	accountMapper := mapper.GetMapper(ctx)
	for _, sender := range tx.Senders {
		a := accountMapper.GetAccount(sender.Address)
		if a == nil {
			return types.ErrSenderAccountNotExists()
		}
		acc := a.(*qtypes.QOSAccount)
		if !acc.EnoughOf(sender.QOS, sender.QSCs) {
			return types.ErrSenderAccountCoinsNotEnough()
		}
	}

	return nil
}

// 转账
func (tx TxTransfer) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	accountMapper := mapper.GetMapper(ctx)

	// 消息类型事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeTransfer),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	// 处理发送账户及事件
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

	// 处理接收账户及事件
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

// 所有Senders参与签名
func (tx TxTransfer) GetSigner() []btypes.AccAddress {
	addrs := make([]btypes.AccAddress, 0)
	for _, sender := range tx.Senders {
		addrs = append(addrs, sender.Address)
	}

	return addrs
}

// Gas
func (tx TxTransfer) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForTransfer)
}

// 第一个Sender支付gas费
func (tx TxTransfer) GetGasPayer() btypes.AccAddress {
	return tx.Senders[0].Address
}

// 签名字节
func (tx TxTransfer) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return ret
}
