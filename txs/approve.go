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
type TxApproveCreate struct {
	*types.Approve
}

func (tx *TxApproveCreate) ValidateData(ctx context.Context) bool {
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

func (tx *TxApproveCreate) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	fromAcc := accountMapper.GetAccount(tx.To)
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
type TxApproveIncrease struct {
	*types.Approve
}

func (tx *TxApproveIncrease) ValidateData(ctx context.Context) bool {
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

func (tx *TxApproveIncrease) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
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
	approve.Coins = approve.Coins.Plus(tx.Coins)
	mapper.SaveApprove(&approve)

	return
}

// 减少授权
type TxApproveDecrease struct {
	*types.Approve
}

func (tx *TxApproveDecrease) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 授权必须存在
	mapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	approve, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return false
	}

	if !approve.Coins.IsGTE(tx.Coins) {
		return false
	}

	return true
}

func (tx *TxApproveDecrease) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
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
	if !approve.Coins.IsGTE(tx.Coins) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 保存更新
	approve.Coins = approve.Coins.Minus(tx.Coins)
	mapper.SaveApprove(&approve)

	return
}

// 使用授权
type TxApproveUse struct {
	*types.Approve
}

func (tx *TxApproveUse) ValidateData(ctx context.Context) bool {
	if !tx.Approve.ValidateData(ctx) {
		return false
	}

	// 校验授权信息
	approveMapper := ctx.Mapper(mapper.ApproveMapperName).(*mapper.ApproveMapper)
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		return false
	}
	if !approve.Coins.IsGTE(tx.Coins) {
		return false
	}

	// 校验授权用户状态
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	iAcc := accountMapper.GetAccount(tx.From)
	if iAcc == nil {
		return false
	}
	from := iAcc.(*account.QOSAccount)
	if !from.Coins().IsGTE(tx.Coins) {
		return false
	}

	return true
}

func (tx *TxApproveUse) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
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
	if !approve.Coins.IsGTE(tx.Coins) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 校验授权用户状态
	if !from.Coins().IsGTE(tx.Coins) {
		result.Code = btypes.ABCICodeType(btypes.CodeInternal)
		return
	}

	// 更新授权用户状态
	fromQscs := from.Coins().Minus(tx.Coins)
	from.Qos = fromQscs.Qos()
	from.QscList = fromQscs.QscList()
	accountMapper.SetAccount(from)

	// 更新被授权账户
	toList := to.Coins().Plus(tx.Coins)
	to.QscList = toList.QscList()
	to.Qos = toList.Qos()
	accountMapper.SetAccount(to)
	// 保存更新
	approve.Coins = approve.Coins.Minus(tx.Coins)
	approveMapper.SaveApprove(&approve)

	return
}

func (tx *TxApproveUse) GetSigner() []btypes.Address {
	return []btypes.Address{tx.To}
}

func (tx *TxApproveUse) GetGasPayer() btypes.Address {
	return tx.To
}

// 取消授权 Tx
type TxApproveCancel struct {
	*types.ApproveCancel
}

func (tx *TxApproveCancel) ValidateData(ctx context.Context) bool {
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

func (tx *TxApproveCancel) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
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
