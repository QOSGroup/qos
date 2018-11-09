package txs

import (
	"bytes"
	"errors"
	"fmt"
	baccount "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/types"
)

// 功能：发币 对应的Tx结构
type TxIssueQsc struct {
	QscName string         `json:"qscName"` //币名
	Amount  btypes.BigInt  `json:"amount"`  //金额
	Banker  btypes.Address `json:"banker"`  //banker地址
}

// 功能：检测TxIssuQsc结构体字段是否合法
func (tx *TxIssueQsc) ValidateData(ctx context.Context) error {
	if tx.Amount.LT(btypes.NewInt(0)) || !btypes.CheckQscName(tx.QscName) {
		return errors.New("QscName or Amount not valid")
	}

	acc := GetAccount(ctx, tx.Banker)
	mainmapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.MainMapper)
	qscinfo := mainmapper.GetQsc(tx.QscName)
	if qscinfo == nil {
		return errors.New("QscName not exists")
	}

	if !bytes.Equal(acc.GetAddress(), qscinfo.BankAddr) {
		return errors.New("Banker Address not match")
	}

	if qscinfo.Qscname != tx.QscName {
		errs := fmt.Sprintf("Banker expect coin(%s) but get (%s)", qscinfo.Qscname, tx.QscName)
		return errors.New(errs)
	}

	return nil
}

// 功能：tx执行
// 发币过程：banker向自己的账户发币
func (tx *TxIssueQsc) Exec(ctx context.Context) (ret btypes.Result, crossTxQcps *btxs.TxQcp) {
	banker := GetAccount(ctx, tx.Banker)
	if &banker == nil {
		ret.Log = "result: Can't find banker"
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}

	var qsccoin *types.QSC
	qsccoin = banker.GetQSC(tx.QscName)
	if qsccoin == nil {
		qsccoin = types.NewQSC(tx.QscName, tx.Amount)
	}else{
		qsccoin.PlusByAmount(tx.Amount)
	}

	err := banker.SetQSC(qsccoin)
	if err != nil {
		ret.Log = "result: set qos to banker error!"
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}

	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		ret.Log = "result: Get mapper error"
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}
	mapper.SetAccount(banker)

	ret.Code = btypes.ABCICodeOK
	ret.Log += fmt.Sprintf("result: Done! Send to banker qos(%d)", tx.Amount.Int64())
	ret.GasUsed = tx.CalcGas().Int64()

	return
}

// 功能：签名者
func (tx *TxIssueQsc) GetSigner() (singer []btypes.Address) {
	singer = append(singer, tx.Banker)
	return
}

// 计算gas
// 此处设置gas = 0
func (tx *TxIssueQsc) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// gas付费人
// no gas, no payer
func (tx *TxIssueQsc) GetGasPayer() btypes.Address {
	return btypes.Address{}
}

// 获取签名字段
func (tx *TxIssueQsc) GetSignData() (ret []byte) {
	ret = append(ret, []byte(tx.QscName)...)
	ret = append(ret, btypes.Int2Byte(tx.Amount.Int64())...)
	ret = append(ret, []byte(tx.Banker)...)

	return
}

// 构建 TxIsssueQsc 结构体
func NewTxIssueQsc(qsc string, amount btypes.BigInt, banker btypes.Address) (rTx *TxIssueQsc) {
	rTx = &TxIssueQsc{
		qsc,
		amount,
		banker,
	}

	return
}
