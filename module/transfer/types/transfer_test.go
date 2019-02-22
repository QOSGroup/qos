package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestTransItem_IsValid(t *testing.T) {
	empty := TransItem{}
	validQSCs := TransItem{QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}
	invalidQSCs := TransItem{QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(-1)}}}
	validQOS := TransItem{QOS: btypes.NewInt(1)}
	invalidQOS := TransItem{QOS: btypes.NewInt(-1)}
	valid := TransItem{QOS: btypes.NewInt(1), QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}

	cases := []struct {
		input TransItem
		valid bool
	}{
		{empty, false},
		{validQSCs, true},
		{invalidQSCs, false},
		{validQOS, true},
		{invalidQOS, false},
		{valid, true},
	}

	for tcIndex, tc := range cases {
		valid, err := tc.input.IsValid()
		require.Equal(t, valid, err == nil, "tc #%d", tcIndex)
		require.Equal(t, valid, tc.valid)
	}
}

func TestTransItems_IsEmpty(t *testing.T) {
	empty := TransItems{}
	valid := TransItems{TransItem{QOS: btypes.NewInt(1), QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}}

	cases := []struct {
		input TransItems
		valid bool
	}{
		{empty, true},
		{valid, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.valid, tc.input.IsEmpty(), "tc #%d", tcIndex)
	}
}

func TestTransItems_IsValid(t *testing.T) {
	addr := ed25519.GenPrivKey().PubKey().Address()
	item := TransItem{Address: btypes.Address(addr), QOS: btypes.NewInt(1), QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}

	empty := TransItems{}
	repeat := TransItems{item, item}
	valid := TransItems{item}

	cases := []struct {
		input TransItems
		valid bool
	}{
		{empty, false},
		{repeat, false},
		{valid, true},
	}

	for tcIndex, tc := range cases {
		valid, err := tc.input.IsValid()
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
		require.Equal(t, tc.valid, valid)
	}
}

func TestTransItems_Match(t *testing.T) {
	addr1 := ed25519.GenPrivKey().PubKey().Address()
	addr2 := ed25519.GenPrivKey().PubKey().Address()
	item1 := TransItems{TransItem{Address: btypes.Address(addr1), QOS: btypes.NewInt(1), QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}}
	item2 := TransItems{TransItem{Address: btypes.Address(addr2), QOS: btypes.NewInt(2), QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}}
	item3 := TransItems{TransItem{Address: btypes.Address(addr2), QOS: btypes.NewInt(1), QSCs: types.QSCs{&types.QSC{"QSTARS", btypes.NewInt(1)}}}}
	item4 := TransItems{TransItem{Address: btypes.Address(addr2), QOS: btypes.NewInt(1), QSCs: types.QSCs{&types.QSC{"QSTARS1", btypes.NewInt(1)}}}}

	cases := []struct {
		input1 TransItems
		input2 TransItems
		match  bool
	}{
		{item1, item2, false},
		{item1, item3, true},
		{item1, item4, false},
	}

	for tcIndex, tc := range cases {
		match, err := tc.input1.Match(tc.input2)
		require.Equal(t, tc.match, err == nil, "tc #%d", tcIndex)
		require.Equal(t, tc.match, match)
	}
}
