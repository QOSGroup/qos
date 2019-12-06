package types

import "math"

const (
	Decimal        = 4     //4位精度
	TotalQOSAmount = 100e8 //100亿
)

var (
	UnitQOS                     = math.Pow(10, Decimal)           // QOS unit
	DefaultTotalUnitQOSQuantity = int64(TotalQOSAmount * UnitQOS) // total QOS amount

	DefaultBlockInterval = int64(5)   // 5 seconds
	UnitQOSGas           = int64(100) // 1 QOS = 100 gas
)
