package types

import (
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParams_Validate(t *testing.T) {
	p := Params{}
	cases := []struct {
		key      string
		value    string
		expected interface{}
	}{
		{"min_deposit", "1", uint64(1)},
		{"min_deposit", "0.5", nil},
		{"max_deposit_period", "1", time.Duration(1)},
		{"max_deposit_period", "aa", nil},
		{"voting_period", "1", time.Duration(1)},
		{"voting_period", "lk", nil},
		{"quorum", "0.5", types.MustNewDecFromStr("0.5")},
		{"quorum", "www", nil},
		{"threshold", "0.5", types.MustNewDecFromStr("0.5")},
		{"threshold", "", nil},
		{"veto", "0.5", types.MustNewDecFromStr("0.5")},
		{"veto", "ee", nil},
		{"penalty", "0.5", types.MustNewDecFromStr("0.5")},
		{"penalty", "6h", nil},
	}

	for tcIndex, tc := range cases {
		v, _ := p.Validate(tc.key, tc.value)
		require.Equal(t, tc.expected, v, "tc #%d", tcIndex)
	}
}

func TestParams_GetParamSpace(t *testing.T) {
	params := Params{
		MinDeposit:       10,
		MaxDepositPeriod: DefaultPeriod,
		VotingPeriod:     DefaultPeriod,
		Quorum:           types.NewDecWithPrec(334, 3),
		Threshold:        types.NewDecWithPrec(5, 1),
		Veto:             types.NewDecWithPrec(334, 3),
		Penalty:          types.ZeroDec(),
	}

	require.Equal(t, ParamSpace, params.GetParamSpace())
}
