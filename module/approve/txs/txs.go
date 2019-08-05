package txs

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/types"
	"github.com/QOSGroup/qos/module/qsc"
	qtypes "github.com/QOSGroup/qos/types"
)

// 创建授权
type TxCreateApprove struct {
	types.Approve
}

func (tx TxCreateApprove) ValidateData(ctx context.Context) error {
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须不存在
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if exists {
		return types.ErrApproveExists(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxCreateApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	fromAcc := accountMapper.GetAccount(tx.From)
	if fromAcc == nil {
		fromAcc = accountMapper.NewAccountWithAddress(tx.From).(*qtypes.QOSAccount)
		accountMapper.SetAccount(fromAcc)
	}
	toAcc := accountMapper.GetAccount(tx.To)
	if toAcc == nil {
		toAcc = accountMapper.NewAccountWithAddress(tx.To).(*qtypes.QOSAccount)
		accountMapper.SetAccount(toAcc)
	}

	// 创建授权
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.SaveApprove(tx.Approve)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：授权账号
func (tx TxCreateApprove) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx TxCreateApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：授权账号
func (tx TxCreateApprove) GetGasPayer() btypes.Address {
	return tx.From
}

// 增加授权
type TxIncreaseApprove struct {
	types.Approve
}

func (tx TxIncreaseApprove) ValidateData(ctx context.Context) error {
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须存在
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxIncreaseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	approve, _ := mapper.GetApprove(tx.From, tx.To)
	approve = approve.Plus(tx.QOS, tx.QSCs)

	// 保存更新
	mapper.SaveApprove(approve)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeIncreaseApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：授权账号
func (tx TxIncreaseApprove) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx TxIncreaseApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：授权账号
func (tx TxIncreaseApprove) GetGasPayer() btypes.Address {
	return tx.From
}

// 减少授权
type TxDecreaseApprove struct {
	types.Approve
}

func (tx TxDecreaseApprove) ValidateData(ctx context.Context) error {
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须存在
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	approve, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists(types.DefaultCodeSpace, "")
	}

	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return types.ErrApproveNotEnough(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxDecreaseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	approve, _ := mapper.GetApprove(tx.From, tx.To)
	approve = approve.Minus(tx.QOS, tx.QSCs)

	// 保存更新
	mapper.SaveApprove(approve)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeDecreaseApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：授权账号
func (tx TxDecreaseApprove) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx TxDecreaseApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：授权账号
func (tx TxDecreaseApprove) GetGasPayer() btypes.Address {
	return tx.From
}

// 使用授权
type TxUseApprove struct {
	types.Approve
}

func (tx TxUseApprove) ValidateData(ctx context.Context) error {
	err := ValidateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 校验授权信息
	approveMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		return types.ErrApproveNotExists(types.DefaultCodeSpace, "")
	}
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return types.ErrApproveNotEnough(types.DefaultCodeSpace, "")
	}

	// 校验授权用户状态
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	iAcc := accountMapper.GetAccount(tx.From)
	if iAcc == nil {
		return types.ErrFromAccountNotExists(types.DefaultCodeSpace, "")
	}
	from := iAcc.(*qtypes.QOSAccount)
	if tx.IsGT(from.QOS, from.QSCs) {
		return types.ErrFromAccountCoinsNotEnough(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxUseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	from := accountMapper.GetAccount(tx.From).(*qtypes.QOSAccount)
	to := accountMapper.GetAccount(tx.To).(*qtypes.QOSAccount)

	approveMapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	approve, _ := approveMapper.GetApprove(tx.From, tx.To)

	// 更新授权用户状态
	from.MustMinus(tx.QOS, tx.QSCs)
	accountMapper.SetAccount(from)

	// 更新被授权账户
	to.MustPlus(tx.QOS, tx.QSCs)
	accountMapper.SetAccount(to)

	// 保存更新
	approveMapper.SaveApprove(approve.Minus(tx.QOS, tx.QSCs))

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeUseApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：被授权账户
func (tx TxUseApprove) GetSigner() []btypes.Address {
	return []btypes.Address{tx.To}
}

// Gas TODO
func (tx TxUseApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账户
func (tx TxUseApprove) GetGasPayer() btypes.Address {
	return tx.To
}

// 取消授权 Tx
type TxCancelApprove struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
}

func (tx TxCancelApprove) ValidateData(ctx context.Context) error {
	if tx.From == nil || tx.To == nil {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	// 授权是否存在
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return types.ErrApproveNotExists(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxCancelApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.DeleteApprove(tx.From, tx.To)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCancelApprove,
			btypes.NewAttribute(types.AttributeKeyApproveFrom, tx.From.String()),
			btypes.NewAttribute(types.AttributeKeyApproveTo, tx.To.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

// 签名账号：被授权账号
func (tx TxCancelApprove) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx TxCancelApprove) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账号
func (tx TxCancelApprove) GetGasPayer() btypes.Address {
	return tx.From
}

// 签名字节
func (tx TxCancelApprove) GetSignData() (ret []byte) {
	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)

	return ret
}

// 基础数据校验
func ValidateData(ctx context.Context, msg types.Approve) error {
	if valid, err := msg.IsValid(); !valid {
		return types.ErrInvalidInput(types.DefaultCodeSpace, err.Error())
	}

	if msg.QSCs.Len() > 0 {
		mapper := qsc.GetMapper(ctx)
		for _, val := range msg.QSCs {
			// 币种存在
			if !mapper.Exists(val.Name) {
				return types.ErrQSCNotExists(types.DefaultCodeSpace, "")
			}
		}
	}

	return nil
}
