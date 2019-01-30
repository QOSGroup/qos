package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

type Fraction struct {
	Value Dec `json:"value"`
}

func NewFraction(numer, denomin int64) Fraction {
	if denomin == int64(0) {
		panic("Denomin cannot be zero")
	}

	return Fraction{
		Value: NewDec(numer).Quo(NewDec(denomin)),
	}
}

func NewFractionFromBigInt(numer, denomin btypes.BigInt) Fraction {
	if denomin.NilToZero().IsZero() {
		panic("Denomin cannot be zero")
	}
	return Fraction{
		Value: NewDecFromInt(numer).Quo(NewDecFromInt(denomin)),
	}
}

func ZeroFraction() Fraction {
	return Fraction{
		Value: ZeroDec(),
	}
}

func OneFraction() Fraction {
	return Fraction{
		Value: OneDec(),
	}
}

func (frac Fraction) Add(f1 Fraction) Fraction {

	value := frac.Value.Add(f1.Value)

	return Fraction{
		Value: value,
	}
}

func (frac Fraction) Sub(f1 Fraction) Fraction {
	value := frac.Value.Sub(f1.Value)
	return Fraction{
		Value: value,
	}
}

func (frac Fraction) Mul(f1 Fraction) Fraction {
	value := frac.Value.Mul(f1.Value)
	return Fraction{
		Value: value,
	}
}

func (frac Fraction) MultiInt64(t int64) btypes.BigInt {
	return frac.Value.MulInt(btypes.NewInt(t)).TruncateInt()
}

func (frac Fraction) MultiBigInt(t btypes.BigInt) btypes.BigInt {
	return frac.Value.MulInt(t).TruncateInt()
}

func (frac Fraction) Equal(f1 Fraction) bool {
	return frac.Value.Equal(f1.Value)
}
