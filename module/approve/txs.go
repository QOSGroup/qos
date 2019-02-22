package approve

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	approvetypes "github.com/QOSGroup/qos/module/approve/types"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/types"
)

// 创建授权
type TxCreateApprove struct {
	approvetypes.Approve
}

func (tx TxCreateApprove) ValidateData(ctx context.Context) error {
	err := validateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须不存在
	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if exists {
		return ErrApproveExists(DefaultCodeSpace, "")
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
		fromAcc = accountMapper.NewAccountWithAddress(tx.From).(*types.QOSAccount)
		accountMapper.SetAccount(fromAcc)
	}
	toAcc := accountMapper.GetAccount(tx.To)
	if toAcc == nil {
		toAcc = accountMapper.NewAccountWithAddress(tx.To).(*types.QOSAccount)
		accountMapper.SetAccount(toAcc)
	}

	// 创建授权
	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	mapper.SaveApprove(tx.Approve)

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
	approvetypes.Approve
}

func (tx TxIncreaseApprove) ValidateData(ctx context.Context) error {
	err := validateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须存在
	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return ErrApproveNotExists(DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxIncreaseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approve, _ := mapper.GetApprove(tx.From, tx.To)
	approve = approve.Plus(tx.QOS, tx.QSCs)

	// 保存更新
	mapper.SaveApprove(approve)

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
	approvetypes.Approve
}

func (tx TxDecreaseApprove) ValidateData(ctx context.Context) error {
	err := validateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 授权必须存在
	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approve, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return ErrApproveNotExists(DefaultCodeSpace, "")
	}

	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return ErrApproveNotEnough(DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxDecreaseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approve, _ := mapper.GetApprove(tx.From, tx.To)
	approve = approve.Minus(tx.QOS, tx.QSCs)

	// 保存更新
	mapper.SaveApprove(approve)

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
	approvetypes.Approve
}

func (tx TxUseApprove) ValidateData(ctx context.Context) error {
	err := validateData(ctx, tx.Approve)
	if err != nil {
		return err
	}

	// 校验授权信息
	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approve, exisit := approveMapper.GetApprove(tx.From, tx.To)
	if !exisit {
		return ErrApproveNotExists(DefaultCodeSpace, "")
	}
	if !approve.IsGTE(tx.QOS, tx.QSCs) {
		return ErrApproveNotEnough(DefaultCodeSpace, "")
	}

	// 校验授权用户状态
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	iAcc := accountMapper.GetAccount(tx.From)
	if iAcc == nil {
		return ErrFromAccountNotExists(DefaultCodeSpace, "")
	}
	from := iAcc.(*types.QOSAccount)
	if tx.IsGT(from.QOS, from.QSCs) {
		return ErrFromAccountCoinsNotEnough(DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxUseApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	from := accountMapper.GetAccount(tx.From).(*types.QOSAccount)
	to := accountMapper.GetAccount(tx.To).(*types.QOSAccount)

	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	approve, _ := approveMapper.GetApprove(tx.From, tx.To)

	// 更新授权用户状态
	from.MustMinus(tx.QOS, tx.QSCs)
	accountMapper.SetAccount(from)

	// 更新被授权账户
	to.MustPlus(tx.QOS, tx.QSCs)
	accountMapper.SetAccount(to)

	// 保存更新
	approveMapper.SaveApprove(approve.Minus(tx.QOS, tx.QSCs))

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
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	// 授权是否存在
	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	_, exists := mapper.GetApprove(tx.From, tx.To)
	if !exists {
		return ErrApproveNotExists(DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxCancelApprove) Exec(ctx context.Context) (result btypes.Result, crossTxQcps *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	mapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	mapper.DeleteApprove(tx.From, tx.To)

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
func validateData(ctx context.Context, msg approvetypes.Approve) error {
	if valid, err := msg.IsValid(); !valid {
		return ErrInvalidInput(DefaultCodeSpace, err.Error())
	}

	if msg.QSCs.Len() > 0 {
		mapper := ctx.Mapper(qsc.QSCMapperName).(*qsc.QSCMapper)
		for _, val := range msg.QSCs {
			// 币种存在
			if !mapper.Exists(val.Name) {
				return ErrQSCNotExists(DefaultCodeSpace, "")
			}
		}
	}

	return nil
}
