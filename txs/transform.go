package txs

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qos/types"
)

// 功能：转账
type TxTransform struct {
	Senders   []AddrTrans `json:"senders"`   //转账人
	Receivers []AddrTrans `json:"receivers"` //收款人
}

// 功能：转账/收款人信息
type AddrTrans struct {
	Address btypes.Address `json:"address"` //账户地址
	Amount  btypes.BigInt  `json:"amount"`  //金额
	QscName string         `json:"qscname"` //qsc名称
}

// 功能：检测合法性
// 规则：
//    	1,senders/receivers不能为空
//		2,sender花费qsc数量不能超过持有数
//		3,sender必须存在，否则无人支付
//		4,receiver若不存在，创建一个账户
//		5,send总数 == receive总数
func (tx TxTransform) ValidateData(ctx context.Context) bool {
	if len(tx.Senders) == 0 || len(tx.Receivers) == 0 {
		return false
	}

	for _, sd := range tx.Senders {
		if !(CheckAddr(sd.Address) && btypes.BigInt.GT(sd.Amount, btypes.ZeroInt()) && btypes.CheckQscName(sd.QscName)) {
			return false
		}
	}

	for _, rv := range tx.Receivers {
		if !(CheckAddr(rv.Address) && btypes.BigInt.GT(rv.Amount, btypes.ZeroInt()) && btypes.CheckQscName(rv.QscName)) {
			return false
		}
	}

	return true
}

// 功能：执行transaction
// 备注：Gas的逻辑判断及扣除应该在外层操作，不放到Tx执行体内
// todo: 某一个用户转账出错，如何处理（跳过？全部还原？）
func (tx TxTransform) Exec(ctx context.Context) (ret btypes.Result, crossTxQcps *btxs.TxQcp) {
	var sdCoin int64 = 0
	var rvCoin int64 = 0
	for _, sd := range tx.Senders {
		sender := GetAccount(ctx, sd.Address)
		qsc := sender.GetQSC(sd.QscName)
		if sd.Amount.GT(qsc.GetAmount()) {
			ret.Code = btypes.ToABCICode(btypes.CodespaceRoot, btypes.CodeInternal) //todo: code?
			ret.Log = "error: " + sd.Address.String() +  " havn't enought money"
			return
		}
		sdCoin += sd.Amount.Int64()
	}

	for _, rv := range tx.Receivers {
		receiver := GetAccount(ctx, rv.Address)
		if receiver == nil {
			CreateAccount(ctx, rv.Address)
		}
		rvCoin += rv.Amount.Int64()
	}

	if sdCoin != rvCoin {
		ret.Code = btypes.ToABCICode(btypes.CodespaceRoot, btypes.CodeInternal) //todo: code?
		ret.Log = "error: coins from senders not equal receivers"
		return
	}

	for _, sd := range tx.Senders {
		sender := GetAccount(ctx, sd.Address)
		qsc := sender.GetQSC(sd.QscName).GetAmount().Sub(sd.Amount)
		sender.SetQSC(types.NewQSC(sd.QscName, qsc))
	}
	for _, rv := range tx.Receivers {
		receiver := GetAccount(ctx, rv.Address)
		qsc := receiver.GetQSC(rv.QscName).GetAmount().Add(rv.Amount)
		receiver.SetQSC(types.NewQSC(rv.QscName, qsc))
	}

	ret.Code = btypes.ABCICodeOK
	ret.Log = "resutl: done!"
	ret.GasUsed = tx.CalcGas().Int64() //type of gas used in tendermint is int64

	return
}

// 功能：返回签名者
func (tx TxTransform) GetSigner() (ret []btypes.Address) {
	ret = []btypes.Address{}
	for _, val := range tx.Senders {
		ret = append(ret, val.Address)
	}

	return
}

// 计算gas
// 基础价为10，每多一个sender/receiver，gas增1
func (tx TxTransform) CalcGas() btypes.BigInt {
	baseNum := 10
	gas := (int64)(baseNum + len(tx.Senders) + len(tx.Receivers) - 2)

	return btypes.NewInt(gas)
}

// 功能：返回gas付费人
// 算法：第一个sender付费
// 注：  返回[]commmand.Address,为方便以后扩展
func (tx TxTransform) GetGasPayer() (payer btypes.Address) {
	payer = tx.Senders[0].Address

	return
}

// 获取签名字段
func (tx TxTransform) GetSignData() (ret []byte) {
	for _, sd := range tx.Senders {
		ret = append(ret, sd.Address.Bytes()...)
		ret = append(ret, btypes.Int2Byte(sd.Amount.Int64())...)
		ret = append(ret, []byte(sd.QscName)...)
	}

	for _, rv := range tx.Receivers {
		ret = append(ret, rv.Address.Bytes()...)
		ret = append(ret, btypes.Int2Byte(rv.Amount.Int64())...)
		ret = append(ret, []byte(rv.QscName)...)
	}

	return
}

// 功能：构建Transform结构体
func NewTransform(senders *[]AddrTrans, receiver *[]AddrTrans) (rTx *TxTransform) {
	rTx = &TxTransform{
		*senders,
		*receiver,
	}

	return
}

// 功能：检查 commmon.Address 的合法性
// todo: types.Address的其他规则需在此处检测
func CheckAddr(addr btypes.Address) (ret bool) {
	ret = true
	if addr.Empty() {
		ret = false
	}

	return
}
