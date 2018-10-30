package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/mapper"
	baccount "github.com/QOSGroup/qbase/account"
)

// 功能：发币 对应的Tx结构
type TxIssueQsc struct {
	QscName string        `json:"qscName"` //发币账户名
	Amount  types.BigInt  `json:"amount"`  //金额
	Banker  types.Address `json:"banker"`  //banker地址
}

// 功能：检测TxIssuQsc结构体字段是否合法
func (tx *TxIssueQsc) ValidateData(ctx context.Context) bool {
	if tx.Amount.LT(types.NewInt(0)) || !types.CheckQscName(tx.QscName) {
		return false
	}

	acc := GetAccount(ctx, tx.Banker)
	mainmapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.MainMapper)
	qscinfo := mainmapper.GetQsc(tx.QscName)
	if qscinfo == nil {
		return false
	}

	return acc.GetPubicKey().Equals(qscinfo.PubkeyBank)
}

// 功能：tx执行
// 发币过程：banker向自己的账户发币
func (tx *TxIssueQsc) Exec(ctx context.Context) (ret types.Result, crossTxQcps *btxs.TxQcp) {
	banker := GetAccount(ctx, tx.Banker)
	if &banker == nil {
		ret.Log = "result: Can't find Bulanker"
		ret = types.ErrInternal(ret.Log).Result()
		return
	}

	err := banker.SetQOS(banker.GetQOS().Add(tx.Amount))
	if err != nil {
		ret.Log = "result: set qos to banker error!"
		ret = types.ErrInternal(ret.Log).Result()
		return
	}
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		ret.Log = "result: Get mapper error"
		ret = types.ErrInternal(ret.Log).Result()
		return
	}
	mapper.SetAccount(banker)

	ret.Code = types.ABCICodeOK
	ret.Log += fmt.Sprintf("result: Done! Send to banker qos(%d)", tx.Amount.Int64())
	ret.GasUsed = tx.CalcGas().Int64()

	return
}

// 功能：签名者
func (tx *TxIssueQsc) GetSigner() (singer []types.Address) {
	singer = append(singer, tx.Banker)
	return
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
	ret = append(ret, []byte(tx.Banker)...)

	return
}

// 构建 TxIsssueQsc 结构体
func NewTxIssueQsc(qsc string, amount types.BigInt, banker types.Address) (rTx *TxIssueQsc) {
	rTx = &TxIssueQsc{
		qsc,
		amount,
		banker,
	}

	return
}
