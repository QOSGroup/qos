package types

import (
	"fmt"
	"math/big"
	"testing"

	btypes "github.com/QOSGroup/qbase/types"
)

func TestFraction(t *testing.T) {

	t1 := NewFraction(int64(1), int64(3))
	t2 := NewFraction(int64(2), int64(3))

	t1 = t1.Add(t2)
	fmt.Println(t1)

	t3 := NewFraction(int64(21323), int64(324234234))
	t1 = t2.Add(t3)
	fmt.Println(t1)

	a := btypes.NewInt(int64(2134234))
	b := btypes.NewInt(int64(4564564564))

	fmt.Println(a.Div(b))

	c := new(big.Float).SetInt(a.BigInt())
	d := new(big.Float).SetInt(b.BigInt())

	fmt.Println(c.Quo(c, d))

}
