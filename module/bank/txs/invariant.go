package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
)

var GasForInvariantCheck = int64(200000*qtypes.UnitQOS) * qtypes.UnitQOSGas // 200000QOS

// 发起全网数据检查
type TxInvariantCheck struct {
	Sender btypes.AccAddress `json:"sender"` // 发送交易账户地址
}

// 数据校验
func (tx TxInvariantCheck) ValidateData(ctx context.Context) error {
	if len(tx.Sender) == 0 {
		return types.ErrInvalidInput("sender is empty")
	}

	return nil
}

// 交易执行
func (tx TxInvariantCheck) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	// 清空检查事件
	mapper.ClearInvariantCheck(ctx)

	// 设置检查事件
	mapper.SetInvariantCheck(ctx)

	// 交易事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeInvariantCheck),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
		btypes.NewEvent(types.EventTypeInvariantCheck,
			btypes.NewAttribute(types.AttributeKeySender, tx.Sender.String()),
			btypes.NewAttribute(types.AttributeKeyHeight, string(ctx.BlockHeight())),
		)}
	return btypes.Result{Code: btypes.CodeOK, Events: result.Events}, nil
}

// 签名账户
func (tx TxInvariantCheck) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Sender}
}

// gas费
func (tx TxInvariantCheck) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForInvariantCheck)
}

// gas payer
func (tx TxInvariantCheck) GetGasPayer() btypes.AccAddress {
	return tx.Sender
}

// 签名字节
func (tx TxInvariantCheck) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return ret
}
