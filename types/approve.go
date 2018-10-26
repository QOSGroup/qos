package types

import (
	"bytes"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"sort"
	"strings"
)

// 授权 Common 结构
type Approve struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
	QOS  btypes.BigInt  `json:"qos"`  // QOS
	QSCs QSCs           `json:"qscs"` // QSCs
}

func NewApprove(from btypes.Address, to btypes.Address, qos *btypes.BigInt, qscs QSCs) Approve {
	if qos == nil {
		val := btypes.NewInt(0)
		qos = &val
	}
	if qscs == nil {
		qscs = QSCs{}
	}
	return Approve{
		From: from,
		To:   to,
		QOS:  *qos,
		QSCs: qscs,
	}
}

// 基础数据校验
// 1.From，To不为空
// 2.QOS、QscList内币值大于0
// 3.QscList内币种不能重复，不能为qos(大小写不敏感)
func (tx Approve) ValidateData(ctx context.Context) bool {
	if tx.From == nil || tx.To == nil || !tx.IsPositive() {
		return false
	}

	m := make(map[string]bool)
	for _, val := range tx.QSCs {
		if strings.ToLower(val.Name) == "qos" {
			return false
		}
		if _, ok := m[val.Name]; !ok {
			m[val.Name] = true
		} else {
			return false
		}
	}

	return true
}

// 签名账号：授权账号，使用授权签名者：被授权账号
func (tx Approve) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx Approve) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer 授权账号，使用授权：被授权账号
func (tx Approve) GetGasPayer() btypes.Address {
	return tx.From
}

// 签名字节
func (tx Approve) GetSignData() (ret []byte) {
	tx.QOS = ZeroNilBigInt(tx.QOS)

	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)
	ret = append(ret, tx.QOS.String()...)
	for _, coin := range tx.QSCs {
		ret = append(ret, []byte(coin.Name)...)
		ret = append(ret, []byte(coin.Amount.String())...)
	}

	return ret
}

// 通用方法
//-----------------------------------------------------------------

// 是否为正值
func (tx Approve) IsPositive() bool {
	if tx.QOS.IsNil() || tx.QOS.IsZero() {
		return tx.QSCs.IsPositive()
	} else if tx.QOS.GT(btypes.NewInt(0)) {
		return tx.IsNotNegative()
	} else {
		return false
	}
}

// 是否为非负值
func (tx Approve) IsNotNegative() bool {
	tx.QOS = ZeroNilBigInt(tx.QOS)

	if tx.QOS.LT(btypes.NewInt(0)) {
		return false
	}

	return tx.QSCs.IsNotNegative()
}

// 返回相反值
func (tx Approve) Negative() (a Approve) {
	a = NewApprove(tx.From, tx.To, nil, nil)
	a.QOS = tx.QOS.Neg()
	a.QSCs = tx.QSCs.Negative()

	return a
}

// Plus
func (tx Approve) Plus(qos btypes.BigInt, qscs QSCs) (a Approve) {
	qos = ZeroNilBigInt(qos)
	a = NewApprove(tx.From, tx.To, nil, nil)
	a.QOS = tx.QOS.Add(qos)
	a.QSCs = tx.QSCs.Plus(qscs)

	return a
}

// Minus
func (tx Approve) Minus(qos btypes.BigInt, qscs QSCs) (a Approve) {
	tx.QOS = ZeroNilBigInt(tx.QOS)
	qos = ZeroNilBigInt(qos)
	a = NewApprove(tx.From, tx.To, nil, nil)
	a.QOS = tx.QOS.Add(qos.Neg())
	a.QSCs = tx.QSCs.Minus(qscs)

	return a
}

// 是否大于等于
func (tx Approve) IsGTE(qos btypes.BigInt, qscs QSCs) bool {
	tx.QOS = ZeroNilBigInt(tx.QOS)
	qos = ZeroNilBigInt(qos)

	if tx.QOS.LT(qos) {
		return false
	}

	return tx.QSCs.IsGTE(qscs)
}

// 是否大于
func (tx Approve) IsGT(qos btypes.BigInt, qscs QSCs) bool {
	tx.QOS = ZeroNilBigInt(tx.QOS)
	qos = ZeroNilBigInt(qos)

	if tx.QOS.LT(qos) {
		return false
	} else if tx.QOS.Equal(qos) {
		return !tx.QSCs.IsLT(qscs) && !tx.QSCs.IsEqual(qscs)
	} else {
		return qscs.IsNotNegative()
	}
}

// 重写Equals
func (tx Approve) Equals(approve Approve) bool {
	return tx.String() == approve.String()
}

// 输出字符串
func (tx Approve) String() string {
	tx.QOS = ZeroNilBigInt(tx.QOS)

	var buf bytes.Buffer
	buf.WriteString("from:" + tx.From.String() + " ")
	buf.WriteString("to:" + tx.To.String() + " ")
	buf.WriteString("qos:" + tx.QOS.String() + " ")
	names := make([]string, 0, len(tx.QSCs))
	m1 := make(map[string]btypes.BigInt)
	for _, val := range tx.QSCs {
		names = append(names, val.Name)
		m1[val.Name] = val.Amount
	}
	sort.Strings(names)
	for _, name := range names {
		buf.WriteString(name + ":")
		buf.WriteString(m1[name].String() + " ")
	}
	return buf.String()
}

//-----------------------------------------------------------------

// 取消授权 结构
type ApproveCancel struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
}

// 基础数据校验
// 1.From，To不为空
func (tx ApproveCancel) ValidateData(ctx context.Context) bool {
	if tx.From == nil || tx.To == nil {
		return false
	}
	return true
}

// 签名账号：被授权账号
func (tx ApproveCancel) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx ApproveCancel) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账号
func (tx ApproveCancel) GetGasPayer() btypes.Address {
	return tx.From
}

// 签名字节
func (tx ApproveCancel) GetSignData() (ret []byte) {
	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)

	return ret
}
