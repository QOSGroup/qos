package types

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
)

// 授权 Common 结构
type Approve struct {
	From  btypes.Address `json:"from"`  // 授权账号
	To    btypes.Address `json:"to"`    // 被授权账号
	Coins QSCS           `json:"coins"` // 授权币种、币值
}

// 基础数据校验
// 1.From，To不为空
// 2.Coins内币种存在，币值大于0
func (tx *Approve) ValidateData(ctx context.Context) bool {
	if tx.From == nil || tx.To == nil || !tx.Coins.IsPositive() {
		return false
	}
	return true
}

// 签名账号：授权账号，使用授权签名者：被授权账号
func (tx *Approve) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx *Approve) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer 授权账号，使用授权：被授权账号
func (tx *Approve) GetGasPayer() btypes.Address {
	return tx.From
}

// 签名字节
func (tx *Approve) GetSignData() (ret []byte) {
	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)
	for _, coin := range tx.Coins {
		ret = append(ret, []byte(coin.Name)...)
		ret = append(ret, []byte(coin.Amount.String())...)
	}

	return ret
}

// 取消授权 结构
type ApproveCancel struct {
	From btypes.Address `json:"from"` // 授权账号
	To   btypes.Address `json:"to"`   // 被授权账号
}

// 基础数据校验
// 1.From，To不为空
func (tx *ApproveCancel) ValidateData(ctx context.Context) bool {
	if tx.From == nil || tx.To == nil {
		return false
	}
	return true
}

// 签名账号：被授权账号
func (tx *ApproveCancel) GetSigner() []btypes.Address {
	return []btypes.Address{tx.From}
}

// Gas TODO
func (tx *ApproveCancel) CalcGas() btypes.BigInt {
	return btypes.NewInt(0)
}

// Gas Payer：被授权账号
func (tx *ApproveCancel) GetGasPayer() btypes.Address {
	return tx.From
}

// 签名字节
func (tx *ApproveCancel) GetSignData() (ret []byte) {
	ret = append(ret, tx.From...)
	ret = append(ret, tx.To...)

	return ret
}
