package mapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvariantCheck(t *testing.T) {
	ctx := defaultContext()
	ctx.WithBlockHeight(20)

	SetInvariantCheck(ctx)
	assert.True(t, NeedInvariantCheck(ctx))

	ClearInvariantCheck(ctx)
	assert.False(t, NeedInvariantCheck(ctx))
}
