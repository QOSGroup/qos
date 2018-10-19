package types

import (
	"fmt"
	"github.com/QOSGroup/qbase/types"
	"sort"
	"strings"
)

// qscs for the specific regional chain
type QSC struct {
	Name   string       `json:"coin_name"`
	Amount types.BigInt `json:"amount"`
}

func NewQSC(name string, amount types.BigInt) *QSC {
	return &QSC{
		name,
		amount,
	}
}

// getter of qsc name
func (qsc *QSC) GetName() string {
	return qsc.Name
}

// getter of qsc name
func (qsc *QSC) GetAmount() types.BigInt {
	return qsc.Amount
}

// setter of qsc amount
func (qsc *QSC) SetAmount(amount types.BigInt) {
	qsc.Amount = amount
}

func NewInt64QSC(denom string, amount int64) *QSC {
	return NewQSC(denom, types.NewInt(amount))
}

// String provides a human-readable representation of a qsc
func (qsc QSC) String() string {
	return fmt.Sprintf("%v%v", qsc.Amount, qsc.Name)
}

// SameDenomAs returns true if the two qscs are the same denom
func (qsc QSC) SameDenomAs(other QSC) bool {
	return (qsc.Name == other.Name)
}

// IsZero returns if this represents no money
func (qsc QSC) IsZero() bool {
	return qsc.Amount.IsZero()
}

// IsGTE returns true if they are the same type and the receiver is
// an equal or greater value
func (qsc QSC) IsGTE(other QSC) bool {
	return qsc.SameDenomAs(other) && (!qsc.Amount.LT(other.Amount))
}

// IsLT returns true if they are the same type and the receiver is
// a smaller value
func (qsc QSC) IsLT(other QSC) bool {
	return !qsc.IsGTE(other)
}

// IsEqual returns true if the two sets of QSCS have the same value
func (qsc QSC) IsEqual(other QSC) bool {
	return qsc.SameDenomAs(other) && (qsc.Amount.Equal(other.Amount))
}

// IsPositive returns true if qsc amount is positive
func (qsc QSC) IsPositive() bool {
	return (qsc.Amount.Sign() == 1)
}

// IsNotNegative returns true if qsc amount is not negative
func (qsc QSC) IsNotNegative() bool {
	return (qsc.Amount.Sign() != -1)
}

// Adds amounts of two qscs with same denom
func (qsc QSC) Plus(coinB QSC) QSC {
	if !qsc.SameDenomAs(coinB) {
		return qsc
	}
	return QSC{qsc.Name, qsc.Amount.Add(coinB.Amount)}
}

// Subtracts amounts of two qscs with same denom
func (qsc QSC) Minus(coinB QSC) QSC {
	if !qsc.SameDenomAs(coinB) {
		return qsc
	}
	return QSC{qsc.Name, qsc.Amount.Sub(coinB.Amount)}
}

// QSCS is a set of Qsc, one per currency
type QSCS []QSC

func (qscs QSCS) String() string {
	if len(qscs) == 0 {
		return ""
	}

	out := ""
	for _, qsc := range qscs {
		out += fmt.Sprintf("%v,", qsc.String())
	}
	return out[:len(out)-1]
}

// IsValid asserts the QSCS are sorted, and don't have 0 amounts
func (qscs QSCS) IsValid() bool {
	switch len(qscs) {
	case 0:
		return true
	case 1:
		return !qscs[0].IsZero()
	default:
		lowDenom := qscs[0].Name
		for _, qsc := range qscs[1:] {
			if qsc.Name <= lowDenom {
				return false
			}
			if qsc.IsZero() {
				return false
			}
			// we compare each qsc against the last denom
			lowDenom = qsc.Name
		}
		return true
	}
}

// Plus combines two sets of qscs
// CONTRACT: Plus will never return QSCS where one Qsc has a 0 amount.
func (qscs QSCS) Plus(coinsB QSCS) QSCS {
	sum := ([]QSC)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(qscs), len(coinsB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, coinsB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, qscs[indexA:]...)
		}
		coinA, coinB := qscs[indexA], coinsB[indexB]
		switch strings.Compare(coinA.Name, coinB.Name) {
		case -1:
			sum = append(sum, coinA)
			indexA++
		case 0:
			if coinA.Amount.Add(coinB.Amount).IsZero() {
				// ignore 0 sum qsc type
			} else {
				sum = append(sum, coinA.Plus(coinB))
			}
			indexA++
			indexB++
		case 1:
			sum = append(sum, coinB)
			indexB++
		}
	}
}

// Negative returns a set of qscs with all amount negative
func (qscs QSCS) Negative() QSCS {
	res := make([]QSC, 0, len(qscs))
	for _, qsc := range qscs {
		res = append(res, QSC{
			Name:   qsc.Name,
			Amount: qsc.Amount.Neg(),
		})
	}
	return res
}

// Minus subtracts a set of qscs from another (adds the inverse)
func (qscs QSCS) Minus(coinsB QSCS) QSCS {
	return qscs.Plus(coinsB.Negative())
}

// IsGTE returns True iff qscs is NonNegative(), and for every
// currency in coinsB, the currency is present at an equal or greater
// amount in coinsB
func (qscs QSCS) IsGTE(coinsB QSCS) bool {
	diff := qscs.Minus(coinsB)
	if len(diff) == 0 {
		return true
	}
	return diff.IsNotNegative()
}

// IsLT returns True iff every currency in qscs, the currency is
// present at a smaller amount in qscs
func (qscs QSCS) IsLT(coinsB QSCS) bool {
	return !qscs.IsGTE(coinsB)
}

// IsZero returns true if there are no qscs
// or all qscs are zero.
func (qscs QSCS) IsZero() bool {
	for _, qsc := range qscs {
		if !qsc.IsZero() {
			return false
		}
	}
	return true
}

// IsEqual returns true if the two sets of QSCS have the same value
func (qscs QSCS) IsEqual(coinsB QSCS) bool {
	if len(qscs) != len(coinsB) {
		return false
	}
	for i := 0; i < len(qscs); i++ {
		if qscs[i].Name != coinsB[i].Name || !qscs[i].Amount.Equal(coinsB[i].Amount) {
			return false
		}
	}
	return true
}

// IsPositive returns true if there is at least one qsc, and all
// currencies have a positive value
func (qscs QSCS) IsPositive() bool {
	if len(qscs) == 0 {
		return false
	}
	for _, qsc := range qscs {
		if !qsc.IsPositive() {
			return false
		}
	}
	return true
}

// IsNotNegative returns true if there is no currency with a negative value
// (even no qscs is true here)
func (qscs QSCS) IsNotNegative() bool {
	if len(qscs) == 0 {
		return true
	}
	for _, qsc := range qscs {
		if !qsc.IsNotNegative() {
			return false
		}
	}
	return true
}

// Returns the amount of a denom from qscs
func (qscs QSCS) AmountOf(denom string) types.BigInt {
	switch len(qscs) {
	case 0:
		return types.ZeroInt()
	case 1:
		qsc := qscs[0]
		if qsc.Name == denom {
			return qsc.Amount
		}
		return types.ZeroInt()
	default:
		midIdx := len(qscs) / 2 // 2:1, 3:1, 4:2
		qsc := qscs[midIdx]
		if denom < qsc.Name {
			return qscs[:midIdx].AmountOf(denom)
		} else if denom == qsc.Name {
			return qsc.Amount
		} else {
			return qscs[midIdx+1:].AmountOf(denom)
		}
	}
}

//----------------------------------------
// Sort interface

//nolint
func (qscs QSCS) Len() int           { return len(qscs) }
func (qscs QSCS) Less(i, j int) bool { return qscs[i].Name < qscs[j].Name }
func (qscs QSCS) Swap(i, j int)      { qscs[i], qscs[j] = qscs[j], qscs[i] }

var _ sort.Interface = QSCS{}

// Sort is a helper function to sort the set of qscs inplace
func (qscs QSCS) Sort() QSCS {
	sort.Sort(qscs)
	return qscs
}

//---------------------------------------
// QOSAccount
func (qscs QSCS) Qos() types.BigInt {
	for _, qsc := range qscs {
		if qsc.Name == "qos" {
			return qsc.Amount
		}
	}
	return types.NewInt(0)
}

// 返回account对应的QSC列表，不包含QOS
func (qscs QSCS) QscList() []*QSC {
	var list []*QSC
	for _, qsc := range qscs {
		if qsc.Name == "qos" {
			continue
		}
		list = append(list, &qsc)
	}
	return list
}
