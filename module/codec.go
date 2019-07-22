package module

import (
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/module/transfer"
	"github.com/tendermint/go-amino"
)

func RegisterCodec(cdc *amino.Codec) {
	approve.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	mint.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	qcp.RegisterCodec(cdc)
	qsc.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	transfer.RegisterCodec(cdc)
}
