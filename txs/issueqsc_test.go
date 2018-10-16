package txs

import (
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTxIssueQsc(t *testing.T) {
	txIssue := NewTxIssueQsc("qsc10", types.NewInt(10000))
	require.NotNil(t, txIssue)

	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	isvalid := txIssue.ValidateData(ctx)
	require.Equal(t, isvalid, true)

	result,_ := txIssue.Exec(ctx)
	require.Equal(t, result.Code, types.ABCICodeOK)

	return
}
