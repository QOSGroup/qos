package guardian

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {
	if mapper.GetMapper(ctx).NeedHalt(uint64(ctx.BlockHeight())) {
		panic("HALT THE NETWORK")
	}
}
