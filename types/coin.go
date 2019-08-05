package types

import (
	"github.com/QOSGroup/qbase/types"
)

const (
	QOSCoinName = "QOS"
	Qos         = int64(1e8)
	BlockReward = int64(0 * Qos)
)

type QSC = types.BaseCoin

func NewQSC(name string, amount types.BigInt) *QSC {
	return &QSC{
		name, amount,
	}
}

type QSCs = types.BaseCoins
