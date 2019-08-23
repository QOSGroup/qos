package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/mapper"
	"github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
)

const GasForInvariantCheck = uint64(200000*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 200000QOS

// 全网账户数据检查
type TxInvariantCheck struct {
	Sender btypes.Address
}

func (tx TxInvariantCheck) ValidateData(ctx context.Context) error {
	if len(tx.Sender) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "Sender empty")
	}

	return nil
}

func (tx TxInvariantCheck) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	// 清空检查事件
	mapper.ClearInvariantCheck(ctx)
	// 设置检查时间
	mapper.SetInvariantCheck(ctx)
	result.Events = result.Events.AppendEvent(btypes.NewEvent(types.EventTypeInvariantCheck,
		btypes.NewAttribute(types.AttributeKeySender, tx.Sender.String()),
		btypes.NewAttribute(types.AttributeKeyHeight, string(ctx.BlockHeight())),
	))
	return btypes.Result{Code: btypes.CodeOK, Events: result.Events}, nil
}

func (tx TxInvariantCheck) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Sender}
}

func (tx TxInvariantCheck) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForInvariantCheck))
}

func (tx TxInvariantCheck) GetGasPayer() btypes.Address {
	return tx.Sender
}

func (tx TxInvariantCheck) GetSignData() (ret []byte) {
	ret = append(ret, tx.Sender...)

	return ret
}
