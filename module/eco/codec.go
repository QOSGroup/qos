package eco

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/module/eco/types"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterConcrete(&types.DistributionParams{}, "eco/types/DistributionParams", nil)
	cdc.RegisterConcrete(&types.StakeParams{}, "eco/types/StakeParams", nil)
	cdc.RegisterConcrete(&types.MintParams{}, "eco/types/MintParams", nil)
	cdc.RegisterConcrete(&types.InflationPhrase{}, "eco/types/InflationPhrase", nil)

	cdc.RegisterConcrete(&types.Validator{}, "eco/types/Validator", nil)
	cdc.RegisterConcrete(&types.DelegationInfo{}, "eco/types/DelegationInfo", nil)
	cdc.RegisterConcrete(&types.DelegatorEarningsStartInfo{}, "eco/types/DelegatorEarningsStartInfo", nil)
	cdc.RegisterConcrete(&types.ValidatorCurrentPeriodSummary{}, "eco/types/ValidatorCurrentPeriodSummary", nil)
	cdc.RegisterConcrete(&types.ValidatorVoteInfo{}, "eco/types/ValidatorVoteInfo", nil)
}
