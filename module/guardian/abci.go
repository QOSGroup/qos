package guardian

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {
	if reason, exists := mapper.GetMapper(ctx).GetHalt(); exists {
		// 停止网络
		panic(fmt.Sprintf("HALT THE NETWORK FOR: %s", reason))
	}
}
