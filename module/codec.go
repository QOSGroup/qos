package module

import (
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/eco"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/module/transfer"
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec) {
	approve.RegisterCodec(cdc)
	qsc.RegisterCodec(cdc)
	transfer.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	qcp.RegisterCodec(cdc)
	eco.RegisterCodec(cdc)
}
