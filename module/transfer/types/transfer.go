package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
)

type TransItem struct {
	Address btypes.Address `json:"addr"` // 账户地址
	QOS     btypes.BigInt  `json:"qos"`  // QOS
	QSCs    types.QSCs     `json:"qscs"` // QSCs
}

// valid QSC、QSCs
func (item TransItem) IsValid() (bool, error) {
	item.QOS = item.QOS.NilToZero()

	if item.QOS.IsZero() && item.QSCs.IsZero() {
		return false, errors.New(fmt.Sprintf("QOS and QSCs in %s are zero", item.Address.String()))
	}
	if btypes.ZeroInt().GT(item.QOS) {
		return false, errors.New(fmt.Sprintf("QOS in %s is lte zero", item.Address.String()))
	}
	if !item.QSCs.IsNotNegative() {
		return false, errors.New(fmt.Sprintf("QSCs in %s is lt zero", item.Address.String()))
	}
	return true, nil
}

type TransItems []TransItem

// empty check
func (items TransItems) IsEmpty() bool {
	return len(items) == 0
}

// not empty and QOS、QSCs
func (items TransItems) IsValid() (bool, error) {
	// not empty
	if items.IsEmpty() {
		return false, errors.New("TransItems empty")
	}

	smap := map[string]bool{}
	for _, item := range items {
		// no repeat address
		if _, ok := smap[item.Address.String()]; ok {
			return false, errors.New(fmt.Sprintf("repeat address:%s", item.Address.String()))
		}

		// valid QOS & QSCs
		valid, err := item.IsValid()
		if !valid {
			return valid, err
		}

		smap[item.Address.String()] = true
	}

	return true, nil
}

// total QOS、QSCs in items and itemsB are equal
func (items TransItems) Match(itemsB TransItems) (bool, error) {
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
		return false, errors.New("QOS、QSCs not equal in Senders and Receivers")
	}

	return true, nil
}
