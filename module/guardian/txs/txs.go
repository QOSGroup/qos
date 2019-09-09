package txs

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	"github.com/QOSGroup/qos/module/guardian/types"
)

const (
	// 描述信息长度限制
	MaxDescriptionLen = 1000
)

// 创建系统账户
type TxAddGuardian struct {
	Description string            `json:"description"` // 描述信息
	Address     btypes.AccAddress `json:"address"`     // 账户地址
	Creator     btypes.AccAddress `json:"creator"`     // 创建账户地址
}

func NewTxAddGuardian(description string, address, creator btypes.AccAddress) *TxAddGuardian {
	return &TxAddGuardian{
		Description: description,
		Address:     address,
		Creator:     creator,
	}
}

var _ txs.ITx = (*TxAddGuardian)(nil)

// 数据校验
func (tx TxAddGuardian) ValidateData(ctx context.Context) error {
	// 描述信息不能太长
	if len(tx.Description) > MaxDescriptionLen {
		return types.ErrInvalidInput("description is too long")
	}

	// 账户地址不能为空
	if len(tx.Address) == 0 {
		return types.ErrInvalidInput("address is empty")
	}

	// 创建者账户地址不能为空
	if len(tx.Creator) == 0 {
		return types.ErrInvalidInput("creator is empty")
	}

	// 系统账户不存在
	mapper := mapper.GetMapper(ctx)
	if _, exists := mapper.GetGuardian(tx.Address); exists {
		return types.ErrGuardianAlreadyExists()
	}

	// 创建账户必须存在且类型是Genesis
	guardian, exists := mapper.GetGuardian(tx.Creator)
	if !exists || guardian.GuardianType != types.Genesis {
		return types.ErrInvalidCreator()
	}

	return nil
}

// 交易执行
func (tx TxAddGuardian) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 保存系统账户
	mapper.GetMapper(ctx).AddGuardian(*types.NewGuardian(tx.Description, types.Ordinary, tx.Address, tx.Creator))

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeAddGuardian,
			btypes.NewAttribute(types.AttributeKeyCreator, tx.Creator.String()),
			btypes.NewAttribute(types.AttributeKeyGuardian, tx.Address.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeAddGuardian),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户：创建账户
func (tx TxAddGuardian) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Creator}
}

// 交易费：0
func (tx TxAddGuardian) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// 交易费支付账户：创建账户
func (tx TxAddGuardian) GetGasPayer() btypes.AccAddress {
	return tx.Creator
}

// 签名字节
func (tx TxAddGuardian) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 删除系统账户
type TxDeleteGuardian struct {
	Address   btypes.AccAddress `json:"address"`    // this guardian's address
	DeletedBy btypes.AccAddress `json:"deleted_by"` // address that initiated the AddGuardian tx
}

func NewTxDeleteGuardian(address, deletedBy btypes.AccAddress) *TxDeleteGuardian {
	return &TxDeleteGuardian{
		Address:   address,
		DeletedBy: deletedBy,
	}
}

var _ txs.ITx = (*TxDeleteGuardian)(nil)

// 数据校验
func (tx TxDeleteGuardian) ValidateData(ctx context.Context) error {
	// 账户地址不能为空
	if len(tx.Address) == 0 {
		return types.ErrInvalidInput("address is empty")
	}

	// 操作账户不能为空
	if len(tx.DeletedBy) == 0 {
		return types.ErrInvalidInput("deleted_by is empty")
	}

	mapper := mapper.GetMapper(ctx)

	// 系统账户必须存在
	guardian, exists := mapper.GetGuardian(tx.Address)
	if !exists {
		return types.ErrUnKnownGuardian()
	}

	// 账户类型必须是Ordinary
	if guardian.GuardianType == types.Genesis {
		return types.ErrInvalidInput("can not delete genesis guardian")
	}

	// 操作账户必须存在且类型是Genesis
	deletedBy, exists := mapper.GetGuardian(tx.DeletedBy)
	if !exists || deletedBy.GuardianType != types.Genesis {
		return types.ErrInvalidCreator()
	}

	return nil
}

// 执行
func (tx TxDeleteGuardian) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 删除特权账户
	mapper.GetMapper(ctx).DeleteGuardian(tx.Address)

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeDeleteGuardian,
			btypes.NewAttribute(types.AttributeKeyDeleteBy, tx.DeletedBy.String()),
			btypes.NewAttribute(types.AttributeKeyGuardian, tx.Address.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeDeleteGuardian),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户：DeletedBy
func (tx TxDeleteGuardian) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.DeletedBy}
}

// 交易费：0
func (tx TxDeleteGuardian) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer: DeletedBy
func (tx TxDeleteGuardian) GetGasPayer() btypes.AccAddress {
	return tx.DeletedBy
}

// 签名字节
func (tx TxDeleteGuardian) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 停网
type TxHaltNetwork struct {
	Guardian btypes.AccAddress `json:"guardian"` // guardian's address
	Reason   string            `json:"reason"`   // reason for halting the network
}

func NewTxHaltNetwork(address btypes.AccAddress, reason string) *TxHaltNetwork {
	return &TxHaltNetwork{
		Guardian: address,
		Reason:   reason,
	}
}

var _ txs.ITx = (*TxHaltNetwork)(nil)

// 数据校验
func (tx TxHaltNetwork) ValidateData(ctx context.Context) error {
	// 操作账户不能为空
	if len(tx.Guardian) == 0 {
		return types.ErrInvalidInput("guardian is empty")
	}

	// 操作原因不能为空且不能大于MaxDescriptionLen
	if len(tx.Reason) == 0 {
		return types.ErrInvalidInput("reason is empty")
	}
	if len(tx.Reason) > MaxDescriptionLen {
		return types.ErrInvalidInput("reason is too long")
	}

	mapper := mapper.GetMapper(ctx)
	guardian, exists := mapper.GetGuardian(tx.Guardian)

	// 操作账户必须存在
	if !exists {
		return types.ErrUnKnownGuardian()
	}

	// 操作账户类型必须是Genesis
	if guardian.GuardianType != types.Genesis {
		return types.ErrInvalidInput("can not halt the network")
	}

	return nil
}

// 执行停网
func (tx TxHaltNetwork) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 设置停网标志
	mapper.GetMapper(ctx).SetHalt(tx.Reason)

	// 发送事件
	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeHaltNetwork,
			btypes.NewAttribute(types.AttributeKeyGuardian, tx.Guardian.String()),
			btypes.NewAttribute(types.AttributeKeyReason, tx.Reason),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeHaltNetwork),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户：操作账户
func (tx TxHaltNetwork) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Guardian}
}

// 交易费：0
func (tx TxHaltNetwork) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer：操作账户
func (tx TxHaltNetwork) GetGasPayer() btypes.AccAddress {
	return tx.Guardian
}

// 签名字节
func (tx TxHaltNetwork) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}
