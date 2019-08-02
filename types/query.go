package types

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QueryRegistry interface {
	RegisterQueryRoute(module string, query Querier)
}

type Querier = func(ctx context.Context, path []string, req abci.RequestQuery) (res []byte, err btypes.Error)
