package modules

import (
	"github.com/QOSGroup/qos/modules/approve"
	"github.com/QOSGroup/qos/modules/qcp"
	"github.com/QOSGroup/qos/modules/qsc"
	"github.com/QOSGroup/qos/modules/stake"
	"github.com/QOSGroup/qos/modules/transfer"
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec) {
	approve.RegisterCodec(cdc)
	qsc.RegisterCodec(cdc)
	transfer.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	qcp.RegisterCodec(cdc)
}
