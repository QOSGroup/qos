package mapper

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"runtime/debug"
)

const (
	QSC  = "qsc"
	Qsc  = "qsc"
	Qscs = "qscs"
)

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
	case Qsc:
		data, e = queryQsc(ctx, route[1])
	case Qscs:
		data, e = queryQSCs(ctx)
	default:
		data = nil
		e = errors.New("not found match path")
	}

	if e != nil {
		return nil, btypes.ErrInternal(e.Error())
	}

	return data, nil
}

func queryQsc(ctx context.Context, qscName string) ([]byte, error) {
	mapper := GetMapper(ctx)
	qsc, exists := mapper.GetQsc(qscName)
	if !exists {
		return nil, fmt.Errorf("qsc %d not exists", qscName)
	}

	return mapper.GetCodec().MarshalJSON(qsc)
}

func queryQSCs(ctx context.Context) ([]byte, error) {
	mapper := GetMapper(ctx)
	qscs := mapper.GetQSCs()

	return mapper.GetCodec().MarshalJSON(qscs)
}

func BuildQueryQSCPath(qscName string) string {
	return fmt.Sprintf("custom/%s/%s/%s", QSC, Qsc, qscName)
}

func BuildQueryQSCsPath() string {
	return fmt.Sprintf("custom/%s/%s", QSC, Qscs)
}
