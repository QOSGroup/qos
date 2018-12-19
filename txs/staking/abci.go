package staking

import (
	"github.com/QOSGroup/qbase/context"

	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

}

func EndBlocker(ctx context.Context) (res abci.ResponseEndBlock) {

	return
}
