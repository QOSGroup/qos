package types

import (
	btypes "github.com/QOSGroup/qbase/types"
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
		{"normal_min_deposit", "1", btypes.NewInt(1)},
		{"normal_min_deposit", "0.5", nil},
		{"normal_max_deposit_period", "1", time.Duration(1)},
		{"normal_max_deposit_period", "aa", nil},
		{"normal_voting_period", "1", time.Duration(1)},
		{"normal_voting_period", "lk", nil},
		{"normal_quorum", "0.5", types.MustNewDecFromStr("0.5")},
		{"normal_quorum", "www", nil},
		{"normal_threshold", "0.5", types.MustNewDecFromStr("0.5")},
		{"normal_threshold", "", nil},
		{"normal_veto", "0.5", types.MustNewDecFromStr("0.5")},
		{"normal_veto", "ee", nil},
		{"normal_penalty", "0.5", types.MustNewDecFromStr("0.5")},
		{"normal_penalty", "6h", nil},
	}

	for tcIndex, tc := range cases {
		v, _ := p.ValidateKeyValue(tc.key, tc.value)
		require.Equal(t, tc.expected, v, "tc #%d", tcIndex)
	}
}

func TestParams_GetParamSpace(t *testing.T) {
	params := Params{
		NormalMinDeposit:       btypes.NewInt(10),
		NormalMaxDepositPeriod: DefaultDepositPeriod,
		NormalVotingPeriod:     DefaultVotingPeriod,
		NormalQuorum:           types.NewDecWithPrec(334, 3),
		NormalThreshold:        types.NewDecWithPrec(5, 1),
		NormalVeto:             types.NewDecWithPrec(334, 3),
		NormalPenalty:          types.ZeroDec(),
	}

	require.Equal(t, ParamSpace, params.GetParamSpace())
}
