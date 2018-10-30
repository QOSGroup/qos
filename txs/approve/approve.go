package approve

import (
	"bytes"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"sort"
	"strings"
)

// 授权 Common 结构
type Approve struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
	QOS  btypes.BigInt  `json:"qos"`  // QOS
	QSCs types.QSCs     `json:"qscs"` // QSCs
}

func NewApprove(from btypes.Address, to btypes.Address, qos btypes.BigInt, qscs types.QSCs) Approve {
	if qscs == nil {
		qscs = types.QSCs{}
	}
	return Approve{
		From: from,
		To:   to,
		QOS:  qos.NilToZero(),
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
	tx.QOS = tx.QOS.NilToZero()

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
	tx.QOS = tx.QOS.NilToZero()

	if tx.QOS.LT(btypes.NewInt(0)) {
		return false
	}

	return tx.QSCs.IsNotNegative()
}

// 返回相反值
func (tx Approve) Negative() (a Approve) {
	a = NewApprove(tx.From, tx.To, tx.QOS.Neg(), tx.QSCs.Negative())

	return a
}

// Plus
func (tx Approve) Plus(qos btypes.BigInt, qscs types.QSCs) (a Approve) {
	a = NewApprove(tx.From, tx.To, tx.QOS.Add(qos.NilToZero()), tx.QSCs.Plus(qscs))

	return a
}

// Minus
func (tx Approve) Minus(qos btypes.BigInt, qscs types.QSCs) (a Approve) {
	tx.QOS = tx.QOS.NilToZero()
	qos = qos.NilToZero()
	a = NewApprove(tx.From, tx.To, tx.QOS.Add(qos.Neg()), tx.QSCs.Minus(qscs))

	return a
}

// 是否大于等于
func (tx Approve) IsGTE(qos btypes.BigInt, qscs types.QSCs) bool {
	tx.QOS = tx.QOS.NilToZero()
	qos = qos.NilToZero()

	if tx.QOS.LT(qos) {
		return false
	}

	return tx.QSCs.IsGTE(qscs)
}

// 是否大于
func (tx Approve) IsGT(qos btypes.BigInt, qscs types.QSCs) bool {
	tx.QOS = tx.QOS.NilToZero()
	qos = qos.NilToZero()

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
	tx.QOS = tx.QOS.NilToZero()

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

// 创建授权
type ApproveCreateTx struct {
	Approve
}

func (tx ApproveCreateTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须不存在
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if exists {
		return false
	}

	return true
}

func (tx ApproveCreateTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	fromAcc := accountMapper.GetAccount(tx.From)
	if fromAcc == nil {
		fromAcc = accountMapper.NewAccountWithAddress(tx.From).(*account.QOSAccount)
		accountMapper.SetAccount(fromAcc)
	}
	toAcc := accountMapper.GetAccount(tx.To)
	if toAcc == nil {
		toAcc = accountMapper.NewAccountWithAddress(tx.To).(*account.QOSAccount)
		accountMapper.SetAccount(toAcc)
	}

	// 创建授权
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	mapper.SaveApprove(tx.Approve)

	return
}

// 增加授权
type ApproveIncreaseTx struct {
	Approve
}

func (tx ApproveIncreaseTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须存在
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return false
	}

	return true
}

func (tx ApproveIncreaseTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)

	// 校验授权信息
	approve, exisit := mapper.GetApprove(tx.From, tx.To)
	if !exisit {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 保存更新
	approve = approve.Plus(tx.QOS, tx.QSCs)
	mapper.SaveApprove(approve)

	return
}

// 减少授权
type ApproveDecreaseTx struct {
	Approve
}

func (tx ApproveDecreaseTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须存在
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	approve, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return false
	}

	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return false
	}

	return true
}

func (tx ApproveDecreaseTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)

	// 校验授权信息
	approve, exisit := mapper.GetApprove(tx.From, tx.To)
	if !exisit {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 保存更新
	approve = approve.Minus(tx.QOS, tx.QSCs)
	mapper.SaveApprove(approve)

	return
}

// 使用授权
type ApproveUseTx struct {
	Approve
}

func (tx ApproveUseTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 校验授权信息
	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		return false
	}
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return false
	}

	// 校验授权用户状态
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	iAcc := accountMapper.GetAccount(tx.From)
	if iAcc == nil {
		return false
	}
	from := iAcc.(*account.QOSAccount)
	if tx.IsGT(from.QOS, from.QSCs) {
		return false
	}

	return true
}

func (tx ApproveUseTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	from := accountMapper.GetAccount(tx.From).(*account.QOSAccount)
	to := accountMapper.GetAccount(tx.To).(*account.QOSAccount)

	approveMapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)

	// 校验授权信息
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 校验授权用户状态
	if tx.IsGT(from.QOS, from.QSCs) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 更新授权用户状态

	fromQscs := tx.Negative().Plus(from.QOS, from.QSCs)
	from.QOS = fromQscs.QOS
	from.QSCs = fromQscs.QSCs
	accountMapper.SetAccount(from)

	// 更新被授权账户
	toList := tx.Plus(to.QOS, to.QSCs)
	to.QOS = toList.QOS
	to.QSCs = toList.QSCs
	accountMapper.SetAccount(to)
	// 保存更新
	approveMapper.SaveApprove(approve.Minus(tx.QOS, tx.QSCs))

	return
}

func (tx ApproveUseTx) GetSigner() []btypes.Address {
	return []btypes.Address{tx.To}
}

func (tx ApproveUseTx) GetGasPayer() btypes.Address {
	return tx.To
}

// 取消授权 Tx
type ApproveCancelTx struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
}

func (tx ApproveCancelTx) ValidateData(ctx context.Context) bool {
	if tx.From == nil || tx.To == nil {
		return false
	}

	// 授权是否存在
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return false
	}

	return true
}

func (tx ApproveCancelTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	// 授权是否存在
	mapper := ctx.Mapper(GetApproveMapperStoreKey()).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	err := mapper.DeleteApprove(tx.From, tx.To)
	if err != nil {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	return
}

// 签名账号：被授权账号
func (tx ApproveCancelTx) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx ApproveCancelTx) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账号
func (tx ApproveCancelTx) GetGasPayer() btypes.Address {
	return tx.From
}

// 签名字节
func (tx ApproveCancelTx) GetSignData() (ret []byte) {
	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)

	return ret
}
