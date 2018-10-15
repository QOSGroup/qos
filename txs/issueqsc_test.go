package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/types"
	"testing"
)

func TestNewTxIssueQsc(t *testing.T) {
	txIssue := NewTxIssueQsc("qsc10", types.NewInt(10000))
	if txIssue == nil {
		t.Error("new TxIssueQsc error!")
		return
	}

	if !txIssue.ValidateData() {
		t.Error("TxIssueQsc Invalidate error")
		return
	}

	key := store.NewKVStoreKey(t.Name())
	ctx := defaultContext(key)
	result := txIssue.Exec(ctx)
	if result.Code != types.ABCICodeOK {
		fmt.Printf("TxIssueQsc execute Error: %d", result.Code)
		return
	}

	fmt.Print("TxIssueQsc excute successful!")
	return
}
