package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestValidVoteOption(t *testing.T) {
	cases := []struct {
		input    VoteOption
		expected bool
	}{
		{OptionYes, true},
		{OptionAbstain, true},
		{OptionNo, true},
		{OptionNoWithVeto, true},
		{0x05, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.expected, ValidVoteOption(tc.input), "tc #%d", tcIndex)
	}
}

func TestVote_Equals(t *testing.T) {
	addr1 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input1   Vote
		input2   Vote
		expected bool
	}{
		{Vote{addr1, 1, OptionYes}, Vote{addr1, 1, OptionYes}, true},
		{Vote{addr1, 1, OptionYes}, Vote{addr2, 1, OptionYes}, false},
		{Vote{addr1, 1, OptionYes}, Vote{addr1, 2, OptionYes}, false},
		{Vote{addr1, 1, OptionYes}, Vote{addr1, 1, OptionNo}, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.expected, tc.input1.Equals(tc.input2), "tc #%d", tcIndex)
	}
}

func TestVoteOptionFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected VoteOption
	}{
		{"yes", OptionYes},
		{"YES", OptionYes},
		{"Yes", OptionYes},
		{"no", OptionNo},
		{"NO", OptionNo},
		{"No", OptionNo},
		{"nowithveto", OptionNoWithVeto},
		{"NOWITHVETO", OptionNoWithVeto},
		{"NoWithVeto", OptionNoWithVeto},
		{"abstain", OptionAbstain},
		{"ABSTAIN", OptionAbstain},
		{"Abstain", OptionAbstain},
		{"qos", VoteOption(0xff)},
	}

	for tcIndex, tc := range cases {
		voteOption, _ := VoteOptionFromString(tc.input)
		require.Equal(t, tc.expected, voteOption, "tc #%d", tcIndex)
	}
}

func TestVoteOption_MarshalUnMarshal(t *testing.T) {
	cdc := amino.NewCodec()

	// binary
	o1 := OptionYes
	o1b := cdc.MustMarshalBinaryBare(o1)
	var o2 VoteOption
	cdc.UnmarshalBinaryBare(o1b, &o2)
	require.Equal(t, o1, o2)

	o1 = OptionNo
	o1b = cdc.MustMarshalBinaryBare(o1)
	cdc.UnmarshalBinaryBare(o1b, &o2)
	require.Equal(t, o1, o2)

	o1 = OptionAbstain
	o1b = cdc.MustMarshalBinaryBare(o1)
	cdc.UnmarshalBinaryBare(o1b, &o2)
	require.Equal(t, o1, o2)

	o1 = OptionNoWithVeto
	o1b = cdc.MustMarshalBinaryBare(o1)
	cdc.UnmarshalBinaryBare(o1b, &o2)
	require.Equal(t, o1, o2)

	o1b = cdc.MustMarshalBinaryBare(o1)
	cdc.UnmarshalBinaryBare(o1b, &o2)
	require.Equal(t, o1, o2)

	// json
	o1b = cdc.MustMarshalJSON(o1)
	cdc.MustUnmarshalJSON(o1b, &o2)
	require.Equal(t, o1, o2)

	o1 = OptionNo
	o1b = cdc.MustMarshalJSON(o1)
	cdc.MustUnmarshalJSON(o1b, &o2)
	require.Equal(t, o1, o2)

	o1 = OptionAbstain
	o1b = cdc.MustMarshalJSON(o1)
	cdc.MustUnmarshalJSON(o1b, &o2)
	require.Equal(t, o1, o2)

	o1 = OptionYes
	o1b = cdc.MustMarshalJSON(o1)
	cdc.MustUnmarshalJSON(o1b, &o2)
	require.Equal(t, o1, o2)
}
