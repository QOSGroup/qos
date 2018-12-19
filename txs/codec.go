package txs

import (
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/txs/qcp"
	"github.com/QOSGroup/qos/txs/qsc"
	"github.com/QOSGroup/qos/txs/staking"
	"github.com/QOSGroup/qos/txs/transfer"
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec) {
	approve.RegisterCodec(cdc)
	qsc.RegisterCodec(cdc)
	transfer.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	qcp.RegisterCodec(cdc)
}
