package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"math"
)

func Assert(condition bool, err error) {
	if !condition {
		panic(err)
	}
}

func AssertInt64GreaterThanZero(num int64) {
	Assert(num >= 0, fmt.Errorf("%d is less than zero. ", num))
}

func AssertUint64NotOverflow(num uint64) {
	Assert(num < math.MaxInt64, fmt.Errorf("%d is overflow. it must less than math.MaxInt64", num))
}

func AssertBigIntNotOverflow(num btypes.BigInt) {
	Assert(num.Int64() < math.MaxInt64, fmt.Errorf("%d is overflow. it must less than math.MaxInt64", num.Int64()))
}
