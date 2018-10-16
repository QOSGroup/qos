package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	types "github.com/QOSGroup/qbase/types"
	btxs "github.com/QOSGroup/qbase/txs"
	"log"
)

// 功能：发币 对应的Tx结构
type TxIssueQsc struct {
	QscName string       `json:"qscName"` //发币账户名
	Amount  types.BigInt `json:"amount"`  //金额
}

// 功能：检测TxIssuQsc结构体字段是否合法
func (tx *TxIssueQsc) ValidateData(ctx context.Context) bool {
	ret := !types.BigInt.LT(tx.Amount, types.ZeroInt())

	return ret && types.CheckQscName(tx.QscName)
}

// 功能：tx执行
// 发币过程：banker向自己的账户发币
func (tx *TxIssueQsc) Exec(ctx context.Context) (ret types.Result, crossTxQcps *btxs.TxQcp) {
	banker := GetBanker(tx.QscName)
	if &banker == nil {
		ret.Code = types.ToABCICode(types.CodespaceRoot, types.CodeInternal) //todo: code?
		ret.Log = "result: Can't find Bulanker"
		return
	}

	err := banker.SetQOS(banker.GetQOS().Add(tx.Amount))
	if err != nil {
		ret.Code = types.ToABCICode(types.CodespaceRoot, types.CodeInternal) //todo: code?
		ret.Log = "result: set qos to banker error!"
		return
	}

	ret.Code = types.ABCICodeOK
	ret.Log += fmt.Sprintf("result: Done! Send to banker qos(%d)", tx.Amount.Int64())
	ret.GasUsed = tx.CalcGas().Int64()

	return
}

// 功能：签名者
// todo: 从store查询
func (tx *TxIssueQsc) GetSigner() (singer []types.Address) {
	banker := GetBanker(tx.QscName)
	return append(singer,banker.BaseAccount.GetAddress())
}

// 计算gas
// 此处设置gas = 0
func (tx *TxIssueQsc) CalcGas() types.BigInt {
	return types.NewInt(0)
}

// gas付费人
// no gas, no payer
func (tx *TxIssueQsc) GetGasPayer() types.Address {
	return types.Address{}
}

// 获取签名字段
func (tx *TxIssueQsc) GetSignData() (ret []byte) {
	ret = append(ret, []byte(tx.QscName)...)
	ret = append(ret, types.Int2Byte(tx.Amount.Int64())...)

	return
}

// 构建 TxIsssueQsc 结构体
func NewTxIssueQsc(qsc string, amount types.BigInt) (rTx *TxIssueQsc) {
	if !types.CheckQscName(qsc) || types.BigInt.LT(amount, types.ZeroInt()) {
		return nil
	}

	if GetBanker(qsc) == nil {
		log.Panic("No banker exist, pleese create qsc with banker's CA first")
		return nil
	}

	rTx = &TxIssueQsc{
		qsc,
		amount,
	}

	return
}
