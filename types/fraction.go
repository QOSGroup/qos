package types

import (
	"math/big"

	btypes "github.com/QOSGroup/qbase/types"
)

type Fraction struct {
	Numer   btypes.BigInt `json:"numer"`
	Denomin btypes.BigInt `json:"denomin"`
}

func NewFraction(numer, denomin int64) Fraction {
	if denomin == int64(0) {
		panic("Denomin cannot be zero")
	}

	return Fraction{
		Numer:   btypes.NewInt(numer),
		Denomin: btypes.NewInt(denomin),
	}
}

func ZeroFraction() Fraction {
	return NewFraction(int64(0), int64(1))
}

func OneFraction() Fraction {
	return NewFraction(int64(1), int64(1))
}

func (frac Fraction) Add(f1 Fraction) Fraction {
	numer := btypes.BigInt{}
	denomin := btypes.BigInt{}
	if frac.Denomin.Equal(f1.Denomin) {
		denomin = frac.Denomin
		numer = frac.Numer.Add(f1.Numer)
	} else {
		denomin = frac.Denomin.Mul(f1.Denomin)
		numer = (frac.Numer.Mul(f1.Denomin)).Add(frac.Denomin.Mul(f1.Numer))
	}

	return Fraction{
		Numer:   numer,
		Denomin: denomin,
	}
}

func (frac Fraction) Sub(f1 Fraction) Fraction {

	numer := btypes.BigInt{}
	denomin := btypes.BigInt{}

	if frac.Denomin.Equal(f1.Denomin) {
		denomin = frac.Denomin
		numer = frac.Numer.Sub(f1.Numer)
	} else {
		denomin = frac.Denomin.Mul(f1.Denomin)
		numer = (frac.Numer.Mul(f1.Denomin)).Sub(frac.Denomin.Mul(f1.Numer))
	}

	return Fraction{
		Numer:   numer,
		Denomin: denomin,
	}
}

func (frac Fraction) GCD() Fraction {
	n := new(big.Int).Abs(frac.Numer.BigInt())
	d := new(big.Int).Abs(frac.Denomin.BigInt())
	gcd := new(big.Int).GCD(nil, nil, n, d)

	numer := new(big.Int).Div(n, gcd)
	denomin := new(big.Int).Div(d, gcd)

	return Fraction{
		Numer:   btypes.NewIntFromBigInt(numer),
		Denomin: btypes.NewIntFromBigInt(denomin),
	}
}

func (frac Fraction) MultiInt64(t int64) btypes.BigInt {
	return frac.Numer.Mul(btypes.NewInt(t)).Div(frac.Denomin)
}

func (frac Fraction) Equal(f1 Fraction) bool {
	frac = frac.GCD()
	f1 = f1.GCD()

	return frac.Numer.Equal(f1.Numer) && frac.Denomin.Equal(f1.Denomin)
}
