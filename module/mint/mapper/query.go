package mapper

import (
	"fmt"
	"runtime/debug"
)

import (
	"errors"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
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
	Total   = "total"
	Applied = "applied"
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
	case Total:
		data, e = queryTotal(ctx)
	case Applied:
		data, e = queryApplied(ctx)
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
	mintMapper := GetMapper(ctx)
	phrases := mintMapper.MustGetInflationPhrases()

	return mintMapper.GetCodec().MarshalJSON(phrases)
}

func queryTotal(ctx context.Context) ([]byte, error) {
	mintMapper := GetMapper(ctx)
	total := mintMapper.GetTotalQOSAmount()

	return mintMapper.GetCodec().MarshalJSON(total)
}

func queryApplied(ctx context.Context) ([]byte, error) {
	mintMapper := GetMapper(ctx)
	applied := mintMapper.GetAllTotalMintQOSAmount()

	return mintMapper.GetCodec().MarshalJSON(applied)
}

//nolint
func BuildQueryPhrasesPath() string {
	return fmt.Sprintf("custom/%s/%s", Mint, Phrases)
}

func BuildQueryTotalPath() string {
	return fmt.Sprintf("custom/%s/%s", Mint, Total)
}

func BuildQueryAppliedPath() string {
	return fmt.Sprintf("custom/%s/%s", Mint, Applied)
}
