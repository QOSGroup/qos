package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	"log"
)

//功能：发币 对应的Tx结构
type TxIssueQsc struct {
	QscName string       `json:"qscName"` //发币账户名
	Amount  types.BigInt `json:"amount"`  //金额
}

//功能：检测TxIssuQsc结构体字段是否合法
func (tx *TxIssueQsc) ValidateData() bool {
	ret := !types.BigInt.LT(tx.Amount, types.ZeroInt())
	return ret && types.CheckQsc(tx.QscName)
}

//功能：tx执行
//发币过程：banker向自己的账户发币
func (tx *TxIssueQsc) Exec(ctx context.Context) (ret types.Result) {
	if !tx.ValidateData() {
		ret.Code = types.ToABCICode(types.CodespaceRoot, types.CodeInternal) //todo: code?
		ret.Log = "result: Invalidate Data"
		return
	}

	banker := GetBanker(tx.QscName)
	if &banker == nil {
		log.Panic("Can't find banker when execute issueqsc transaction")
		ret.Code = types.ToABCICode(types.CodespaceRoot, types.CodeInternal) //todo: code?
		ret.Log = "result: Can't find Bulanker"
		return
	}

	banker.SetQOS(banker.GetQOS().Add(tx.Amount))
	ret.Code = types.ABCICodeOK
	ret.Log += fmt.Sprintf("result: Done! Send to banker qos(%d)", tx.Amount.Int64())
	ret.GasUsed = tx.CalcGas()
	return
}

//功能：签名者
//todo: 从store查询
func (tx *TxIssueQsc) GetSigner() types.Address {
	banker := GetBanker(tx.QscName)
	return banker.BaseAccount.GetAddress()
}

//计算gas
//此处设置gas = 0
func (tx *TxIssueQsc) CalcGas() int64 {
	return 0
}

//gas付费人
//no gas, no payer
func (tx *TxIssueQsc) GetGasPayer() types.Address {
	return types.Address{}
}

//获取签名字段
func (tx *TxIssueQsc) GetSignData() (ret []byte) {
	ret = append(ret, []byte(tx.QscName)...)
	ret = append(ret, types.Int2Byte(tx.Amount.Int64())...)
	return
}

//构建 TxIsssueQsc 结构体
func NewTxIssueQsc(qsc string, amount types.BigInt) (rTx *TxIssueQsc) {
	if !types.CheckQsc(qsc) {
		log.Panic("Invalidate qsc name!")
		return nil
	}

	if !types.BigInt.LT(amount, types.ZeroInt()) {
		log.Panic("Error: param amount <= 0")
		return nil
	}

	if GetBanker(qsc) == nil {
		log.Panic("No banker exist, pleese create qsc with banker's CA first")
		return nil
	}

	rTx = new(TxIssueQsc)
	rTx.QscName = qsc
	rTx.Amount = amount
	return
}
