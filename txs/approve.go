package txs

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/types"
)

// 创建授权
type ApproveCreateTx struct {
	types.Approve
}

func (tx ApproveCreateTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须不存在
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
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
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	mapper.SaveApprove(tx.Approve)

	return
}

// 增加授权
type ApproveIncreaseTx struct {
	types.Approve
}

func (tx ApproveIncreaseTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须存在
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
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
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)

	// 校验授权信息
	approve, exisit := mapper.GetApprove(tx.From, tx.To)
	if !exisit {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 保存更新
	approve = approve.Plus(tx.Approve)
	mapper.SaveApprove(approve)

	return
}

// 减少授权
type ApproveDecreaseTx struct {
	types.Approve
}

func (tx ApproveDecreaseTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须存在
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	approve, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return false
	}

	if !approve.IsGTE(tx.Approve) {
		return false
	}

	return true
}

func (tx ApproveDecreaseTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)

	// 校验授权信息
	approve, exisit := mapper.GetApprove(tx.From, tx.To)
	if !exisit {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}
	if !approve.IsGTE(tx.Approve) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 保存更新
	approve = approve.Minus(tx.Approve)
	mapper.SaveApprove(approve)

	return
}

// 使用授权
type ApproveUseTx struct {
	types.Approve
}

func (tx ApproveUseTx) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 校验授权信息
	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		return false
	}
	if !approve.IsGTE(tx.Approve) {
		return false
	}

	// 校验授权用户状态
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	iAcc := accountMapper.GetAccount(tx.From)
	if iAcc == nil {
		return false
	}
	from := iAcc.(*account.QOSAccount)
	if tx.IsGTAccount(from.Qos, from.QscList) {
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

	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)

	// 校验授权信息
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}
	if !approve.IsGTE(tx.Approve) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 校验授权用户状态
	if tx.IsGTAccount(from.Qos, from.QscList) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 更新授权用户状态

	fromQscs := tx.Negative().PlusAccount(from.Qos, from.QscList)
	from.Qos = fromQscs.Qos
	from.QscList = fromQscs.QscList
	accountMapper.SetAccount(from)

	// 更新被授权账户
	toList := tx.PlusAccount(to.Qos, to.QscList)
	to.Qos = toList.Qos
	to.QscList = toList.QscList
	accountMapper.SetAccount(to)
	// 保存更新
	approveMapper.SaveApprove(approve.Minus(tx.Approve))

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
	types.ApproveCancel
}

func (tx ApproveCancelTx) ValidateData(ctx context.Context) bool {
	if !tx.ApproveCancel.ValidateData(ctx) {
		return false
	}

	// 授权是否存在
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
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
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	err := mapper.DeleteApprove(tx.ApproveCancel)
	if err != nil {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	return
}
