package bank

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx context.Context, req abci.RequestEndBlock) {
	// 存在数据检查请求时向Event中添加EventTypeInvariantCheck事件
	if NeedInvariantCheck(ctx) {
		ctx.EventManager().EmitEvent(btypes.NewEvent(qtypes.EventTypeInvariantCheck))
	}

	return
}
