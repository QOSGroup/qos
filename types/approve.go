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
	From    btypes.Address `json:"from"` // 授权账号
	To      btypes.Address `json:"to"`   // 被授权账号
	Qos     btypes.BigInt  `json:"qos"`  // qos
	QscList []*QSC         `json:"qsc"`  // qscs
}

func NewApprove(from btypes.Address, to btypes.Address) Approve {
	return Approve{
		From:    from,
		To:      to,
		Qos:     btypes.NewInt(0),
		QscList: []*QSC{},
	}
}

// 基础数据校验
// 1.From，To不为空
// 2.Qos、QscList内币种存在，币值大于0
// 3.QscList内币种不能重复，不能为qos(大小写不敏感)
func (tx Approve) ValidateData(ctx context.Context) bool {
	if tx.From == nil || tx.To == nil || !tx.IsPositive() {
		return false
	}

	m := make(map[string]bool)
	for _, val := range tx.QscList {
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
	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)
	if &tx.Qos == nil {
		ret = append(ret, tx.Qos.String()...)
	}
	for _, coin := range tx.QscList {
		ret = append(ret, []byte(coin.Name)...)
		ret = append(ret, []byte(coin.Amount.String())...)
	}

	return ret
}

// 通用方法
//-----------------------------------------------------------------

// 是否为正值
func (tx Approve) IsPositive() bool {
	if tx.Qos.IsZero() {
		sum := 0
		for _, qsc := range tx.QscList {
			if qsc.Amount.LT(btypes.NewInt(0)) {
				return false
			}
			if qsc.Amount.IsZero() {
				sum ++
			}
		}
		if sum == len(tx.QscList) {
			return false
		}
		return true
	} else if tx.Qos.GT(btypes.NewInt(0)) {
		for _, qsc := range tx.QscList {
			if qsc.Amount.LT(btypes.NewInt(0)) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// 是否为非负值
func (tx Approve) IsNotNegative() bool {
	// Qos > 0
	if tx.Qos.LT(btypes.NewInt(0)) {
		return false
	}
	// Qsc > 0
	for _, qsc := range tx.QscList {
		if qsc.Amount.LT(btypes.NewInt(0)) {
			return false
		}
	}

	return true
}

// 返回相反值
func (tx Approve) Negative() (a Approve) {
	a = NewApprove(tx.From, tx.To)
	a.Qos = tx.Qos.Neg()
	for _, val := range tx.QscList {
		qsc := QSC{
			Name:   val.Name,
			Amount: val.Amount.Neg(),
		}
		a.QscList = append(a.QscList, &qsc)
	}

	return a
}

// Plus
func (tx Approve) Plus(approve Approve) (a Approve) {
	a = NewApprove(tx.From, tx.To)
	a.Qos = tx.Qos.Add(approve.Qos)

	m1 := make(map[string]btypes.BigInt)
	for _, val := range tx.QscList {
		m1[val.Name] = val.Amount
	}
	m2 := make(map[string]btypes.BigInt)
	for _, val := range approve.QscList {
		m2[val.Name] = val.Amount
	}
	for key, val := range m1 {
		if val2, ok := m2[key]; ok {
			m1[key] = val.Add(val2)
			delete(m2, key)
		}
	}
	for key, val := range m2 {
		m1[key] = val
	}

	for key, val := range m1 {
		if val.GT(btypes.NewInt(0)) { // 保存大于0的数据
			a.QscList = append(a.QscList, &QSC{
				Name:   key,
				Amount: val,
			})
		}
	}

	return a
}

// Plus Account
func (tx Approve) PlusAccount(Qos btypes.BigInt, QscList []*QSC) (a Approve) {
	accApprove := NewApprove(tx.From, tx.To)
	accApprove.Qos = Qos
	accApprove.QscList = QscList
	a = tx.Plus(accApprove)

	return a
}

// Minus
func (tx Approve) Minus(approve Approve) (a Approve) {
	a = tx.Plus(approve.Negative())

	return a
}

// Plus Account
func (tx Approve) MinusAccount(Qos btypes.BigInt, QscList []*QSC) (a Approve) {
	accApprove := NewApprove(tx.From, tx.To)
	accApprove.Qos = Qos
	accApprove.QscList = QscList
	a = tx.Plus(accApprove.Negative())

	return a
}

// 是否大于等于
func (tx Approve) IsGTE(approve Approve) bool {
	if tx.Qos.LT(approve.Qos) {
		return false
	}
	diff := tx.Minus(approve)
	if diff.Qos.IsZero() && len(diff.QscList) == 0 {
		return true
	}
	return diff.IsNotNegative()
}

// 是否大于等于account持有币种、币值
func (tx Approve) IsGTEAccount(Qos btypes.BigInt, QscList []*QSC) bool {
	accApprove := NewApprove(tx.From, tx.To)
	accApprove.Qos = Qos
	accApprove.QscList = QscList

	return tx.IsGTE(accApprove)
}

// 是否大于
func (tx Approve) IsGT(approve Approve) bool {
	if tx.Qos.LT(approve.Qos) {
		return false
	}
	diff := tx.Minus(approve)
	return diff.IsPositive()
}

// 是否大于account持有币种、币值
func (tx Approve) IsGTAccount(Qos btypes.BigInt, QscList []*QSC) bool {
	accApprove := NewApprove(tx.From, tx.To)
	accApprove.Qos = Qos
	accApprove.QscList = QscList

	return tx.IsGT(accApprove)
}

// 重写Equals
func (tx Approve) Equals(approve Approve) bool {
	return tx.String() == approve.String()
}

// 输出字符串
func (tx Approve) String() string {
	var buf bytes.Buffer
	buf.WriteString("from:" + tx.From.String() + " ")
	buf.WriteString("to:" + tx.To.String() + " ")
	buf.WriteString("qos:" + tx.Qos.String() + " ")
	names := make([]string, 0, len(tx.QscList))
	m1 := make(map[string]btypes.BigInt)
	for _, val := range tx.QscList {
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
