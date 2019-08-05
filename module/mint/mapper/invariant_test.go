package mapper

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTotalAppliedInvariant(t *testing.T) {
	ctx := defaultMintContext()

	mm := GetMapper(ctx)
	mm.SetAllTotalMintQOSAmount(100)

	_, coins, broken := TotalAppliedInvariant("mint")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(types.QOSCoinName), btypes.NewInt(-100))
}
