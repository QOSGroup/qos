package types

import (
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQOSAccount_GetQOS(t *testing.T) {
	cases := []struct {
		input    *QOSAccount
		expected btypes.BigInt
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), nil), btypes.ZeroInt()},
		{NewQOSAccount(nil, btypes.NewInt(1), nil), btypes.NewInt(1)},
		{NewQOSAccount(nil, btypes.NewInt(-1), nil), btypes.NewInt(-1)},
	}

	for tcIndex, tc := range cases {
		res := tc.input.GetQOS()
		require.Equal(t, tc.expected, res, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_SetQOS(t *testing.T) {
	cases := []struct {
		input1   *QOSAccount
		input2   btypes.BigInt
		expected btypes.BigInt
		correct  bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), nil), btypes.ZeroInt(), btypes.ZeroInt(), true},
		{NewQOSAccount(nil, btypes.ZeroInt(), nil), btypes.NewInt(1), btypes.NewInt(1), true},
		{NewQOSAccount(nil, btypes.ZeroInt(), nil), btypes.NewInt(-1), btypes.ZeroInt(), false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.SetQOS(tc.input2)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
		require.Equal(t, tc.expected, tc.input1.QOS, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_EnoughOfQOS(t *testing.T) {
	zero := btypes.ZeroInt()
	one := btypes.NewInt(1)
	negOne := btypes.NewInt(-1)

	cases := []struct {
		input1 *QOSAccount
		input2 btypes.BigInt
		enough bool
	}{
		{NewQOSAccount(nil, zero, nil), zero, true},
		{NewQOSAccount(nil, zero, nil), one, false},
		{NewQOSAccount(nil, zero, nil), negOne, true},
		{NewQOSAccount(nil, one, nil), one, true},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.EnoughOfQOS(tc.input2)
		require.Equal(t, tc.enough, res, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_PlusQOS(t *testing.T) {
	zero := btypes.ZeroInt()
	one := btypes.NewInt(1)
	negOne := btypes.NewInt(-1)
	two := btypes.NewInt(2)

	cases := []struct {
		input1  *QOSAccount
		input2  btypes.BigInt
		expect  btypes.BigInt
		correct bool
	}{
		{NewQOSAccount(nil, zero, nil), one, one, true},
		{NewQOSAccount(nil, zero, nil), negOne, zero, false},
		{NewQOSAccount(nil, one, nil), one, two, true},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.PlusQOS(tc.input2)
		require.Equal(t, tc.expect.Int64(), tc.input1.QOS.Int64(), "tc #%d", tcIndex)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_MinusQOS(t *testing.T) {
	zero := btypes.ZeroInt()
	one := btypes.NewInt(1)
	negOne := btypes.NewInt(-1)

	cases := []struct {
		input1  *QOSAccount
		input2  btypes.BigInt
		expect  btypes.BigInt
		correct bool
	}{
		{NewQOSAccount(nil, zero, nil), one, zero, false},
		{NewQOSAccount(nil, zero, nil), negOne, zero, false},
		{NewQOSAccount(nil, one, nil), one, zero, true},
		{NewQOSAccount(nil, one, nil), zero, one, true},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.MinusQOS(tc.input2)
		require.Equal(t, tc.expect.Int64(), tc.input1.QOS.Int64(), "tc #%d", tcIndex)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_GetQSCs(t *testing.T) {
	cases := []struct {
		input1   *QOSAccount
		expected QSCs
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), nil), nil},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs{&QSC{"QSC", btypes.NewInt(1)}}), QSCs{&QSC{"QSC", btypes.NewInt(1)}}},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.GetQSCs()
		require.Equal(t, res, tc.expected, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_PlusQSC(t *testing.T) {
	zero := QSC{"QSC", btypes.ZeroInt()}
	one := QSC{"QSC", btypes.NewInt(1)}
	negOne := QSC{"QSC", btypes.NewInt(-1)}
	two := QSC{"QSC", btypes.NewInt(2)}

	var emptyQSCs QSCs
	zeroQSCs := QSCs{&zero}
	oneQSCs := QSCs{&one}
	twoQSCs := QSCs{&two}

	cases := []struct {
		input1   *QOSAccount
		input2   QSC
		expected QSCs
		correct  bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), one, oneQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), negOne, emptyQSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), one, oneQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), negOne, zeroQSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), one, twoQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), negOne, oneQSCs, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.PlusQSC(tc.input2)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
		require.Equal(t, tc.expected, tc.input1.QSCs, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_MinusQSC(t *testing.T) {
	zero := QSC{"QSC", btypes.ZeroInt()}
	one := QSC{"QSC", btypes.NewInt(1)}
	negOne := QSC{"QSC", btypes.NewInt(-1)}

	var emptyQSCs QSCs
	zeroQSCs := QSCs{&zero}
	oneQSCs := QSCs{&one}

	cases := []struct {
		input1   *QOSAccount
		input2   QSC
		expected QSCs
		correct  bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), one, emptyQSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), negOne, emptyQSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), one, zeroQSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), negOne, zeroQSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), one, zeroQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), negOne, oneQSCs, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.MinusQSC(tc.input2)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
		if tc.correct {
			require.Equal(t, len(tc.expected), len(tc.input1.QSCs))
			for i := range tc.input1.QSCs {
				require.Equal(t, tc.expected[i].Name, tc.input1.QSCs[i].Name, "tc #%d", tcIndex)
				require.Equal(t, tc.expected[i].Amount.Int64(), tc.input1.QSCs[i].Amount.Int64(), "tc #%d", tcIndex)
			}
		}
	}
}

func TestQOSAccount_EnoughOfQSC(t *testing.T) {
	zero := QSC{"QSC", btypes.ZeroInt()}
	one := QSC{"QSC", btypes.NewInt(1)}
	negOne := QSC{"QSC", btypes.NewInt(-1)}
	two := QSC{"QSC", btypes.NewInt(2)}

	var emptyQSCs QSCs
	zeroQSCs := QSCs{&zero}
	oneQSCs := QSCs{&one}

	cases := []struct {
		input1 *QOSAccount
		input2 QSC
		enough bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), one, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), negOne, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), zero, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), zero, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), one, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSCs), negOne, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), one, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), negOne, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), two, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.EnoughOfQSC(tc.input2)
		require.Equal(t, tc.enough, res, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_EnoughOfQSCs(t *testing.T) {
	zeroQSC1 := QSC{"QSC1", btypes.ZeroInt()}
	oneQSC1 := QSC{"QSC1", btypes.NewInt(1)}
	twoQSC1 := QSC{"QSC1", btypes.NewInt(2)}
	zeroQSC2 := QSC{"QSC2", btypes.ZeroInt()}
	oneQSC2 := QSC{"QSC2", btypes.NewInt(1)}
	twoQSC2 := QSC{"QSC2", btypes.NewInt(2)}

	var emptyQSCs QSCs
	zeroQSC1QSCs := QSCs{&zeroQSC1}
	oneQSC1QSCs := QSCs{&oneQSC1}
	twoQSC1QSCs := QSCs{&twoQSC1}
	zeroQSC2QSCs := QSCs{&zeroQSC2}
	oneQSC2QSCs := QSCs{&oneQSC2}
	twoQSC2QSCs := QSCs{&twoQSC2}
	qscs := QSCs{&oneQSC1, &oneQSC2}

	cases := []struct {
		input1 *QOSAccount
		input2 QSCs
		enough bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), oneQSC1QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), zeroQSC1QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), emptyQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSC1QSCs), oneQSC1QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSC1QSCs), zeroQSC1QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSC1QSCs), emptyQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSC1QSCs), emptyQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSC1QSCs), zeroQSC1QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSC1QSCs), oneQSC1QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSC1QSCs), twoQSC1QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), zeroQSC2QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), oneQSC2QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSC1QSCs), zeroQSC2QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), zeroQSC1QSCs), oneQSC2QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSC1QSCs), zeroQSC2QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSC1QSCs), oneQSC2QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), emptyQSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), oneQSC1QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), zeroQSC1QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), twoQSC1QSCs, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), oneQSC2QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), zeroQSC2QSCs, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), qscs), twoQSC2QSCs, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.EnoughOfQSCs(tc.input2)
		require.Equal(t, tc.enough, res, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_PlusQSCs(t *testing.T) {
	zeroQSC1 := QSC{"QSC1", btypes.ZeroInt()}
	oneQSC1 := QSC{"QSC1", btypes.NewInt(1)}
	twoQSC1 := QSC{"QSC1", btypes.NewInt(2)}
	negOneQSC1 := QSC{"QSC1", btypes.NewInt(-1)}
	negTwoQSC1 := QSC{"QSC1", btypes.NewInt(-2)}
	oneQSC2 := QSC{"QSC2", btypes.NewInt(1)}
	twoQSC2 := QSC{"QSC2", btypes.NewInt(2)}

	QSCs1 := QSCs{&zeroQSC1}
	QSCs2 := QSCs{&oneQSC1}
	QSCs3 := QSCs{&oneQSC1, &oneQSC2}
	QSCs4 := QSCs{&twoQSC1}
	QSCs5 := QSCs{&twoQSC1, &oneQSC2}
	QSCs6 := QSCs{&twoQSC1, &twoQSC2}
	QSCs7 := QSCs{&negOneQSC1, &oneQSC2}
	QSCs8 := QSCs{&negTwoQSC1, &oneQSC2}

	cases := []struct {
		input1  *QOSAccount
		input2  QSCs
		expect  QSCs
		correct bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs1), QSCs1, QSCs1, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs1), QSCs2, QSCs2, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs1), QSCs3, QSCs3, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs2), QSCs2, QSCs4, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs2), QSCs3, QSCs5, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs3), QSCs3, QSCs6, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs1), QSCs7, QSCs1, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs2), QSCs8, QSCs2, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.PlusQSCs(tc.input2)
		require.Equal(t, tc.input1.QSCs, tc.expect, "tc #%d", tcIndex)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_MinusQSCs(t *testing.T) {
	zeroQSC1 := QSC{"QSC1", btypes.ZeroInt()}
	oneQSC1 := QSC{"QSC1", btypes.NewInt(1)}
	twoQSC1 := QSC{"QSC1", btypes.NewInt(2)}
	oneQSC2 := QSC{"QSC2", btypes.NewInt(1)}
	twoQSC2 := QSC{"QSC2", btypes.NewInt(2)}

	QSCs1 := QSCs{&zeroQSC1}
	QSCs2 := QSCs{&oneQSC1}
	QSCs3 := QSCs{&oneQSC1, &oneQSC2}
	QSCs4 := QSCs{&twoQSC1}
	QSCs5 := QSCs{&twoQSC1, &oneQSC2}
	QSCs6 := QSCs{&twoQSC1, &twoQSC2}

	cases := []struct {
		input1  *QOSAccount
		input2  QSCs
		expect  QSCs
		correct bool
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs1), QSCs2, QSCs1, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs1), QSCs1, QSCs1, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs2), QSCs1, QSCs2, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs3), QSCs1, QSCs3, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs3), QSCs5, QSCs3, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs4), QSCs2, QSCs2, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs4), QSCs5, QSCs4, false},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs5), QSCs2, QSCs3, true},
		{NewQOSAccount(nil, btypes.ZeroInt(), QSCs6), QSCs3, QSCs3, true},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.MinusQSCs(tc.input2)
		require.Equal(t, tc.input1.QSCs, tc.expect, "tc #%d", tcIndex)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_Plus(t *testing.T) {
	one := btypes.NewInt(1)
	two := btypes.NewInt(2)
	negOne := btypes.NewInt(-1)
	oneQSCs := QSCs{&QSC{"QSC", one}}
	twoQSCs := QSCs{&QSC{"QSC", two}}
	negOneQSCs := QSCs{&QSC{"QSC", negOne}}

	type Item struct {
		QOS  btypes.BigInt
		QSCs QSCs
	}

	item1 := Item{QOS: one}
	item2 := Item{QSCs: oneQSCs}
	item3 := Item{QOS: one, QSCs: oneQSCs}
	item4 := Item{QOS: two}
	item5 := Item{QOS: two, QSCs: oneQSCs}
	item6 := Item{QOS: one, QSCs: twoQSCs}
	item7 := Item{QOS: two, QSCs: twoQSCs}
	item8 := Item{QOS: negOne}
	item9 := Item{QOS: one, QSCs: negOneQSCs}

	cases := []struct {
		input1  *QOSAccount
		input2  Item
		expect  Item
		correct bool
	}{
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item1, item4, true},
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item2, item3, true},
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item3, item5, true},
		{NewQOSAccount(nil, item2.QOS, item2.QSCs), item3, item6, true},
		{NewQOSAccount(nil, item3.QOS, item3.QSCs), item3, item7, true},
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item8, item1, false},
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item9, item1, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.Plus(tc.input2.QOS, tc.input2.QSCs)
		require.Equal(t, Item{tc.input1.QOS, tc.input1.QSCs}, tc.expect, "tc #%d", tcIndex)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
	}
}

type item struct {
	QOS  btypes.BigInt
	QSCs QSCs
}

func (i item) Equals(iB item) bool {
	if !i.QOS.NilToZero().Equal(iB.QOS.NilToZero()) {
		return false
	}

	if i.QSCs.IsEqual(iB.QSCs) {
		return true
	}

	return false
}

func TestQOSAccount_Minus(t *testing.T) {
	one := btypes.NewInt(1)
	two := btypes.NewInt(2)
	negOne := btypes.NewInt(-1)
	oneQSCs := QSCs{&QSC{"QSC", one}}
	twoQSCs := QSCs{&QSC{"QSC", two}}
	negOneQSCs := QSCs{&QSC{"QSC", negOne}}

	item1 := item{QOS: one}
	item2 := item{QSCs: oneQSCs}
	item3 := item{QOS: one, QSCs: oneQSCs}
	item4 := item{QOS: two}
	item5 := item{QOS: two, QSCs: oneQSCs}
	item6 := item{QOS: one, QSCs: twoQSCs}
	item7 := item{QOS: two, QSCs: twoQSCs}
	item8 := item{QOS: negOne}
	item9 := item{QOS: one, QSCs: negOneQSCs}

	cases := []struct {
		input1  *QOSAccount
		input2  item
		expect  item
		correct bool
	}{
		{NewQOSAccount(nil, item4.QOS, item4.QSCs), item1, item1, true},
		{NewQOSAccount(nil, item3.QOS, item3.QSCs), item1, item2, true},
		{NewQOSAccount(nil, item5.QOS, item5.QSCs), item1, item3, true},
		{NewQOSAccount(nil, item6.QOS, item6.QSCs), item2, item3, true},
		{NewQOSAccount(nil, item7.QOS, item7.QSCs), item3, item3, true},
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item8, item1, false},
		{NewQOSAccount(nil, item1.QOS, item1.QSCs), item9, item1, false},
	}

	for tcIndex, tc := range cases {
		res := tc.input1.Minus(tc.input2.QOS, tc.input2.QSCs)
		require.True(t, item{tc.input1.QOS, tc.input1.QSCs}.Equals(tc.expect), "tc #%d", tcIndex)
		require.Equal(t, tc.correct, res == nil, "tc #%d", tcIndex)
	}
}

func TestQOSAccount_RemoveQSC(t *testing.T) {
	qstars := "QSC"
	qstars1 := "QSTARS1"

	emptyQSCs := QSCs{}
	oneQSCs := QSCs{&QSC{"QSC", btypes.NewInt(1)}}

	cases := []struct {
		input1   *QOSAccount
		input2   string
		expected QSCs
	}{
		{NewQOSAccount(nil, btypes.ZeroInt(), emptyQSCs), qstars, emptyQSCs},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), qstars, emptyQSCs},
		{NewQOSAccount(nil, btypes.ZeroInt(), oneQSCs), qstars1, oneQSCs},
	}

	for tcIndex, tc := range cases {
		tc.input1.RemoveQSC(tc.input2)
		require.Equal(t, tc.expected, tc.input1.QSCs, "tc #%d", tcIndex)
	}
}

func TestParseAccounts(t *testing.T) {
	bech1 := "address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd"
	bech2 := "address1czkqg0ekmdaj3xpazkzr5kmsatg3fx27qg609m"
	str1 := bech1 + ",1QOS,2QSTARS"
	str2 := bech1 + ",1QOS,2QSTARS;" + bech2 + ",1QOS"

	addr1, _ := btypes.GetAddrFromBech32(bech1)
	addr2, _ := btypes.GetAddrFromBech32(bech2)
	accs1 := []*QOSAccount{NewQOSAccount(addr1, btypes.NewInt(1), QSCs{&QSC{"QSTARS", btypes.NewInt(2)}})}
	accs2 := []*QOSAccount{
		NewQOSAccount(addr1, btypes.NewInt(1), QSCs{&QSC{"QSTARS", btypes.NewInt(2)}}),
		NewQOSAccount(addr2, btypes.NewInt(1), QSCs{}),
	}

	cases := []struct {
		input1   string
		expected []*QOSAccount
		correct  bool
	}{
		{str1, accs1, true},
		{str2, accs2, true},
	}

	for tcIndex, tc := range cases {
		res, err := ParseAccounts(tc.input1)
		require.Equal(t, tc.correct, err == nil, "tc #%d", tcIndex)
		require.Equal(t, tc.expected, res, "tc #%d", tcIndex)
	}
}
