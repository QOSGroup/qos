package transfer

import (
	"github.com/QOSGroup/qos/module/transfer/txs"
)

var (
	ModuleName    = "transfer"
	RegisterCodec = txs.RegisterCodec
)
