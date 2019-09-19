package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/types"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/qsc"
	qtypes "github.com/QOSGroup/qos/types"
)

// 创建授权
type TxCreateApprove struct {
	types.Approve
}

var _ txs.ITx = (*TxCreateApprove)(nil)

// 参数校验
func (tx TxCreateApprove) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须不存在
	mapper := mapper.GetMapper(ctx)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if exists {
		return types.ErrApproveExists()
	}

	return nil
}

// 执行交易并返回交易结果，QOS公链执行交易crossTxQcps为空
func (tx TxCreateApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 初始化账户信息
	bankMapper := bank.GetMapper(ctx)
	fromAcc := bankMapper.GetAccount(tx.From)
	if fromAcc == nil {
		fromAcc = bankMapper.NewAccountWithAddress(tx.From)
		bankMapper.SetAccount(fromAcc)
	}
	toAcc := bankMapper.GetAccount(tx.To)
	if toAcc == nil {
		toAcc = bankMapper.NewAccountWithAddress(tx.To)
		bankMapper.SetAccount(toAcc)
	}

	// 创建授权
	mapper := mapper.GetMapper(ctx)
	mapper.SaveApprove(tx.Approve)

	// 交易事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeCreateApprove),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：授权账号
func (tx TxCreateApprove) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.From}
}

// Gas
func (tx TxCreateApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：授权账号
func (tx TxCreateApprove) GetGasPayer() btypes.AccAddress {
	return tx.From
}

// 签名字节
func (tx TxCreateApprove) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return ret
}

// 增加授权
type TxIncreaseApprove struct {
	types.Approve
}

var _ txs.ITx = (*TxIncreaseApprove)(nil)

// 参数校验
func (tx TxIncreaseApprove) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须存在
	mapper := mapper.GetMapper(ctx)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists()
	}

	return nil
}

// 执行交易并返回交易结果，QOS公链执行交易crossTxQcps为空
func (tx TxIncreaseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := mapper.GetMapper(ctx)
	approve, _ := mapper.GetApprove(tx.From, tx.To)
	approve = approve.Plus(tx.QOS, tx.QSCs)

	// 保存更新
	mapper.SaveApprove(approve)

	// 交易事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeIncreaseApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeIncreaseApprove),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：授权账号
func (tx TxIncreaseApprove) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.From}
}

// Gas
func (tx TxIncreaseApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：授权账号
func (tx TxIncreaseApprove) GetGasPayer() btypes.AccAddress {
	return tx.From
}

// 签名字节
func (tx TxIncreaseApprove) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return ret
}

// 减少授权
type TxDecreaseApprove struct {
	types.Approve
}

var _ txs.ITx = (*TxDecreaseApprove)(nil)

// 参数校验
func (tx TxDecreaseApprove) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须存在
	mapper := mapper.GetMapper(ctx)
	approve, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists()
	}

	// 校验减少预授权后数值
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return types.ErrApproveNotEnough()
	}

	return nil
}

// 执行交易并返回交易结果，QOS公链执行交易crossTxQcps为空
func (tx TxDecreaseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := mapper.GetMapper(ctx)
	approve, _ := mapper.GetApprove(tx.From, tx.To)
	approve = approve.Minus(tx.QOS, tx.QSCs)

	// 保存更新
	if approve.IsPositive() {
		mapper.SaveApprove(approve)
	} else {
		mapper.DeleteApprove(approve.From, approve.To)
	}

	// 交易事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeDecreaseApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeDecreaseApprove),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：授权账号
func (tx TxDecreaseApprove) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.From}
}

// Gas
func (tx TxDecreaseApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：授权账号
func (tx TxDecreaseApprove) GetGasPayer() btypes.AccAddress {
	return tx.From
}

// 签名字节
func (tx TxDecreaseApprove) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return ret
}

// 使用授权
type TxUseApprove struct {
	types.Approve
}

var _ txs.ITx = (*TxUseApprove)(nil)

// 参数校验
func (tx TxUseApprove) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 校验授权信息
	approveMapper := mapper.GetMapper(ctx)
	approve, exists := approveMapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists()
	}
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return types.ErrApproveNotEnough()
	}

	// 校验授权用户状态
	from := bank.GetAccount(ctx, tx.From)
	if tx.IsGT(from.QOS, from.QSCs) {
		return types.ErrFromAccountCoinsNotEnough()
	}

	return nil
}

// 执行交易并返回交易结果，QOS公链执行交易crossTxQcps为空
func (tx TxUseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	accountMapper := bank.GetMapper(ctx)
	from := accountMapper.GetAccount(tx.From).(*qtypes.QOSAccount)
	to := accountMapper.GetAccount(tx.To).(*qtypes.QOSAccount)

	approveMapper := mapper.GetMapper(ctx)
	approve, _ := approveMapper.GetApprove(tx.From, tx.To)

	// 更新授权用户状态
	from.MustMinus(tx.QOS, tx.QSCs)
	accountMapper.SetAccount(from)

	// 更新被授权账户
	to.MustPlus(tx.QOS, tx.QSCs)
	accountMapper.SetAccount(to)

	// 保存更新，若授权已经使用完，删除预授权记录
	approve = approve.Minus(tx.QOS, tx.QSCs)
	if approve.IsPositive() {
		approveMapper.SaveApprove(approve)
	} else {
		approveMapper.DeleteApprove(approve.From, approve.To)
	}

	// 交易事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeUseApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeUseApprove),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：被授权账户
func (tx TxUseApprove) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.To}
}

// Gas
func (tx TxUseApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账户
func (tx TxUseApprove) GetGasPayer() btypes.AccAddress {
	return tx.To
}

// 签名字节
func (tx TxUseApprove) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return ret
}

// 取消授权 Tx
type TxCancelApprove struct {
	From btypes.AccAddress `json:"from"` // 授权账号
	To   btypes.AccAddress `json:"to"`   // 被授权账号
}

var _ txs.ITx = (*TxCancelApprove)(nil)

// 参数校验
func (tx TxCancelApprove) ValidateData(ctx context.Context) error {
	// 授权账户不能为空
	if len(tx.From) == 0 {
		return types.ErrInvalidInput("from address is empty")
	}
	// 被授权账户不能为空
	if len(tx.To) == 0 {
		return types.ErrInvalidInput("to address is empty")
	}

	// 授权必须存在
	mapper := mapper.GetMapper(ctx)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists()
	}

	return nil
}

// 执行交易并返回交易结果，QOS公链执行交易crossTxQcps为空
func (tx TxCancelApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper.GetMapper(ctx).DeleteApprove(tx.From, tx.To)

	// 交易事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCancelApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeCancelApprove),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：被授权账号
func (tx TxCancelApprove) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.From}
}

// Gas
func (tx TxCancelApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账号
func (tx TxCancelApprove) GetGasPayer() btypes.AccAddress {
	return tx.From
}

// 签名字节
func (tx TxCancelApprove) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)
	return ret
}

// 基础数据校验
func ValidateData(ctx context.Context, msg types.Approve) error {
	if err := msg.Valid(); err != nil {
		return err
	}

	// from账户必须存在
	from := bank.GetAccount(ctx, msg.From)
	if from == nil {
		return types.ErrFromAccountNotExists()
	}

	if msg.QSCs.Len() > 0 {
		mapper := qsc.GetMapper(ctx)
		for _, val := range msg.QSCs {
			// 币种存在
			if !mapper.Exists(val.Name) {
				return types.ErrQSCNotExists()
			}
		}
	}

	return nil
}
