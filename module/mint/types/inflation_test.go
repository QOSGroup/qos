package types

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInflationPhrase_Equals(t *testing.T) {
	ip1 := InflationPhrase{time.Now().AddDate(-1, 0, 0), 1000, 0}
	ip2 := InflationPhrase{time.Now().AddDate(0, 1, 0), 0, 0}
	ip3 := InflationPhrase{time.Now().AddDate(1, 0, 0), 1000, 0}

	cases := []struct {
		phrases   InflationPhrases
		nePhrases InflationPhrases
		equals    bool
	}{
		{InflationPhrases{}, InflationPhrases{}, true},
		{InflationPhrases{ip1}, InflationPhrases{ip1}, true},
		{InflationPhrases{ip1}, InflationPhrases{ip2}, false},
		{InflationPhrases{ip1, ip2}, InflationPhrases{ip1, ip3}, false},
		{InflationPhrases{ip1, ip2}, InflationPhrases{ip1, ip2}, true},
		{InflationPhrases{ip1, ip2}, InflationPhrases{ip2, ip1}, true},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.phrases.Equals(tc.nePhrases), tc.equals, "tc #%d", tcIndex)
	}
}

func TestInflationPhrases_Valid(t *testing.T) {
	ip1 := InflationPhrase{time.Now().AddDate(-1, 0, 0), 1000, 0}
	ip2 := InflationPhrase{time.Now().AddDate(0, 1, 0), 0, 0}
	ip3 := InflationPhrase{time.Now().AddDate(1, 0, 0), 1000, 0}

	cases := []struct {
		inputPhrases InflationPhrases
		valid        bool
	}{
		{InflationPhrases{ip1}, true},
		{InflationPhrases{ip1, ip1}, false},
		{InflationPhrases{ip1, ip2}, false},
		{InflationPhrases{ip1, ip3}, true},
		{InflationPhrases{ip3, ip1}, true},
	}

	for tcIndex, tc := range cases {
		err := tc.inputPhrases.Valid()
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
	}
}

func TestInflationPhrases_GetPhrase(t *testing.T) {
	ip1 := InflationPhrase{time.Now().AddDate(-1, 0, 0), 1000, 0}
	ip2 := InflationPhrase{time.Now().AddDate(0, 1, 0), 1000, 0}
	ip3 := InflationPhrase{time.Now().AddDate(1, 0, 0), 1000, 0}

	cases := []struct {
		inputPhrases InflationPhrases
		inputTime    time.Time
		expected     *InflationPhrase
	}{
		{InflationPhrases{ip1}, time.Now().UTC(), nil},
		{InflationPhrases{ip1, ip2}, time.Now().UTC(), &ip2},
		{InflationPhrases{ip1, ip2, ip3}, time.Now().UTC(), &ip2},
		{InflationPhrases{ip1, ip2, ip3}, ip2.EndTime.UTC(), &ip3},
		{InflationPhrases{ip3, ip1, ip2}, ip2.EndTime.UTC(), &ip3},
	}

	for tcIndex, tc := range cases {
		res, exists := tc.inputPhrases.GetPhrase(tc.inputTime)
		require.Equal(t, res != nil, exists, "tc #%d", tcIndex)
		if exists {
			require.True(t, tc.expected.Equals(*res), "tc #%d", tcIndex)
		}
	}
}

func TestInflationPhrases_GetPrePhrase(t *testing.T) {
	ip1 := InflationPhrase{time.Now().AddDate(-1, 0, 0), 1000, 0}
	ip2 := InflationPhrase{time.Now().AddDate(0, 1, 0), 1000, 0}
	ip3 := InflationPhrase{time.Now().AddDate(1, 0, 0), 1000, 0}

	cases := []struct {
		inputPhrases InflationPhrases
		inputTime    time.Time
		expected     *InflationPhrase
	}{
		{InflationPhrases{ip1}, time.Now().UTC(), &ip1},
		{InflationPhrases{ip1, ip2}, ip1.EndTime.UTC(), &ip1},
		{InflationPhrases{ip1, ip2}, time.Now().UTC(), &ip1},
		{InflationPhrases{ip1, ip2, ip3}, ip2.EndTime.UTC(), &ip2},
		{InflationPhrases{ip1, ip2, ip3}, ip3.EndTime.UTC(), &ip3},
		{InflationPhrases{ip3, ip1, ip2}, ip3.EndTime.UTC(), &ip3},
		{InflationPhrases{ip1, ip2, ip3}, ip3.EndTime.UTC().AddDate(0, 0, -1), &ip2},
	}

	for tcIndex, tc := range cases {
		res, exists := tc.inputPhrases.GetPrePhrase(tc.inputTime)
		require.Equal(t, res != nil, exists, "tc #%d", tcIndex)
		if exists {
			require.True(t, tc.expected.Equals(*res), "tc #%d", tcIndex)
		}
	}
}

func TestInflationPhrases_ValidNewPhrases(t *testing.T) {
	ip1 := InflationPhrase{time.Now().AddDate(-1, 0, 0), 1000, 0}
	ip2 := InflationPhrase{time.Now().AddDate(0, 1, 0), 1000, 0}
	ip3 := InflationPhrase{time.Now().AddDate(0, 1, 0), 2000, 0}
	ip4 := InflationPhrase{time.Now().AddDate(1, 0, 0), 2000, 0}

	cases := []struct {
		phrases      InflationPhrases
		totalApplied uint64
		newTotal     uint64
		newPhrases   InflationPhrases
		valid        bool
	}{
		{InflationPhrases{ip1}, 1000, 1000, InflationPhrases{ip1}, false},
		{InflationPhrases{ip1, ip2}, 1000, 1000, InflationPhrases{ip1, ip3}, false},
		{InflationPhrases{ip1, ip2}, 1000, 1000, InflationPhrases{ip1, ip3}, false},
		{InflationPhrases{ip1, ip2}, 1000, 4000, InflationPhrases{ip1, ip4}, false},
		{InflationPhrases{ip1, ip2}, 1000, 5000, InflationPhrases{ip1, ip3, ip4}, false},
		{InflationPhrases{ip1, ip2}, 1000, 4000, InflationPhrases{ip1, ip2, ip4}, false},
		{InflationPhrases{ip1, ip2}, 1000, 5000, InflationPhrases{ip1, ip2, ip4}, true},
	}

	for tcIndex, tc := range cases {
		err := tc.phrases.ValidNewPhrases(tc.newTotal, tc.totalApplied, tc.newPhrases)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
	}
}
