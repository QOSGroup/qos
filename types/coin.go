package types

import (
	"github.com/QOSGroup/qbase/types"
)

const (
	QOSCoinName = "QOS"
)

type QSC = types.BaseCoin

func NewQSC(name string, amount types.BigInt) *QSC {
	return &QSC{
		name, amount,
	}
}

type QSCs = types.BaseCoins
