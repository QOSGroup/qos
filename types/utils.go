package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

// BigInt nil值转换成0值
func ZeroNilBigInt(val btypes.BigInt) btypes.BigInt {
	if val.IsNil() || val.Sign() == 0 {
		return btypes.ZeroInt()
	}
	return val
}
