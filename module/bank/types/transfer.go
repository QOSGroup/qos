package types

import (
	"fmt"

	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
)

type TransItem struct {
	Address btypes.AccAddress `json:"addr"` // 账户地址
	QOS     btypes.BigInt     `json:"qos"`  // QOS
	QSCs    types.QSCs        `json:"qscs"` // QSCs
}

// 数据校验
func (item TransItem) Valid() error {
	item.QOS = item.QOS.NilToZero()

	// QOS和QSCs必须都为正
	if item.QOS.IsZero() && item.QSCs.IsZero() {
		return ErrInvalidInput(fmt.Sprintf("QOS and QSCs in %s are zero", item.Address.String()))
	}
	if btypes.ZeroInt().GT(item.QOS) {
		return ErrInvalidInput(fmt.Sprintf("QOS in %s is lte zero", item.Address.String()))
	}
	if !item.QSCs.IsNotNegative() {
		return ErrInvalidInput(fmt.Sprintf("QSCs in %s is lt zero", item.Address.String()))
	}

	return nil
}

type TransItems []TransItem

// 是否为空
func (items TransItems) IsEmpty() bool {
	return len(items) == 0
}

// 数据校验
func (items TransItems) Valid() error {
	// 不能为空
	if items.IsEmpty() {
		return ErrInvalidInput("transItems empty")
	}

	smap := map[string]bool{}
	for _, item := range items {
		// 转账地址不能重复
		if _, ok := smap[item.Address.String()]; ok {
			return ErrInvalidInput(fmt.Sprintf("repeat address:%s", item.Address.String()))
		}

		// 明细校验
		err := item.Valid()
		if err != nil {
			return err
		}

		smap[item.Address.String()] = true
	}

	return nil
}

// 判断转账列表是否匹配
func (items TransItems) Match(itemsB TransItems) error {
	sumsqos := btypes.ZeroInt()
	sumsqscs := types.QSCs{}
	for _, item := range items {
		sumsqos = sumsqos.Add(item.QOS)
		sumsqscs = sumsqscs.Plus(item.QSCs)
	}
	sumrqos := btypes.ZeroInt()
	sumrqscs := types.QSCs{}
	for _, item := range itemsB {
		sumrqos = sumrqos.Add(item.QOS)
		sumrqscs = sumrqscs.Plus(item.QSCs)
	}

	// 转入转出相等
	if !sumsqos.Equal(sumrqos) || !sumsqscs.IsEqual(sumrqscs) {
		return ErrInvalidInput("QOS,QSCs not equal in Senders and Receivers")
	}

	return nil
}
