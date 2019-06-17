package mint

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/mapper"
	abci "github.com/tendermint/tendermint/abci/types"
	"runtime/debug"
)

/*
	custom query : /custom/gov/$queryPath

	$queryPath:
		1. /proposals
		2. /proposal/:pid
		3. /votes/:pid
		4. /vote/:pid/:addr
		5. /deposit/:pid/:addr
		6. /deposits/:pid
		7. /tally/:pid
		8. /params
*/

//nolint
const (
	Mint    = "mint"
	Phrases = "phrases"
)

//nolint
func Query(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err btypes.Error) {

	defer func() {
		if r := recover(); r != nil {
			err = btypes.ErrInternal(string(debug.Stack()))
			return
		}
	}()

	var data []byte
	var e error

	switch route[0] {
	case Phrases:
		data, e = queryPhrases(ctx)
	default:
		data = nil
		e = errors.New("not found match path")
	}

	if e != nil {
		return nil, btypes.ErrInternal(e.Error())
	}

	return data, nil
}

func queryPhrases(ctx context.Context) ([]byte, error) {
	mintMapper := mapper.GetMintMapper(ctx)
	phrases := mintMapper.GetMintParams().Phrases

	return mintMapper.GetCodec().MarshalJSON(phrases)
}

//nolint
func BuildQueryProposalPath() string {
	return fmt.Sprintf("custom/%s/%s", Mint, Phrases)
}
