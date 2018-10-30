package txs

import (
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTxIssueQsc(t *testing.T) {
	var bankeraddr = []byte("address1zcjduepqdrq67yachuqs7w87syaq62x999r5dddtwehtamhymnpmxwew9j0srz0a36")
	txIssue := NewTxIssueQsc("qsc10", types.NewInt(10000), bankeraddr)
	require.NotNil(t, txIssue)

	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)

	txIssue.ValidateData(ctx)
	//require.Equal(t, isvalid, true)

	//result, _ := txIssue.Exec(ctx)
	//require.Equal(t, result.Code, types.ABCICodeOK)

	return
}
