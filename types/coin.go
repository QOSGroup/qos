package types

import (
	"github.com/QOSGroup/qbase/types"
)

type QSC = types.BaseCoin

func NewQSC(name string, amount types.BigInt) *QSC {
	return &QSC{
		name, amount,
	}
}

type QSCs = types.BaseCoins