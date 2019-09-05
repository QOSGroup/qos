package types

import (
	qtypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidProposalStatus(t *testing.T) {
	cases := []struct {
		input    ProposalStatus
		expected bool
	}{
		{StatusNil, false},
		{StatusDepositPeriod, true},
		{StatusVotingPeriod, true},
		{StatusPassed, true},
		{StatusRejected, true},
		{0x05, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.expected, ValidProposalStatus(tc.input), "tc #%d", tcIndex)
	}
}

func TestTallyResult_Equals(t *testing.T) {
	cases := []struct {
		input1   TallyResult
		input2   TallyResult
		expected bool
	}{
		{TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, true},
		{TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, TallyResult{qtypes.NewInt(1), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, false},
		{TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, TallyResult{qtypes.NewInt(0), qtypes.NewInt(1), qtypes.NewInt(0), qtypes.NewInt(0)}, false},
		{TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(1), qtypes.NewInt(0)}, false},
		{TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0)}, TallyResult{qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(0), qtypes.NewInt(1)}, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.expected, tc.input1.Equals(tc.input2), "tc #%d", tcIndex)
	}
}

func TestProposalTypeFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected ProposalType
	}{
		{"text", ProposalTypeText},
		{"TEXT", ProposalTypeText},
		{"Text", ProposalTypeText},
		{"parameterchange", ProposalTypeParameterChange},
		{"PARAMETERCHANGE", ProposalTypeParameterChange},
		{"ParameterChange", ProposalTypeParameterChange},
		{"taxusage", ProposalTypeTaxUsage},
		{"TAXUSAGE", ProposalTypeTaxUsage},
		{"TaxUsage", ProposalTypeTaxUsage},
		{"qos", ProposalType(0xff)},
	}

	for tcIndex, tc := range cases {
		proposalType, _ := ProposalTypeFromString(tc.input)
		require.Equal(t, tc.expected, proposalType, "tc #%d", tcIndex)
	}
}

func TestProposalStatusFromString(t *testing.T) {
	cases := []struct {
		input    string
		expected ProposalStatus
	}{
		{"deposit", StatusDepositPeriod},
		{"DEPOSIT", StatusDepositPeriod},
		{"Deposit", StatusDepositPeriod},
		{"voting", StatusVotingPeriod},
		{"VOTING", StatusVotingPeriod},
		{"Voting", StatusVotingPeriod},
		{"passed", StatusPassed},
		{"PASSED", StatusPassed},
		{"Passed", StatusPassed},
		{"rejected", StatusRejected},
		{"REJECTED", StatusRejected},
		{"Rejected", StatusRejected},
		{"qos", ProposalStatus(0xff)},
	}

	for tcIndex, tc := range cases {
		proposalStatus, _ := ProposalStatusFromString(tc.input)
		require.Equal(t, tc.expected, proposalStatus, "tc #%d", tcIndex)
	}
}

func TestProposalType_String(t *testing.T) {
	cases := []struct {
		input    ProposalType
		expected string
	}{
		{ProposalTypeText, "Text"},
		{ProposalTypeParameterChange, "Parameter"},
		{ProposalTypeTaxUsage, "TaxUsage"},
		{ProposalTypeModifyInflation, "ModifyInflation"},
		{ProposalTypeSoftwareUpgrade, "SoftwareUpgrade"},
		{0x06, ""},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.expected, tc.input.String(), "tc #%d", tcIndex)
	}
}

func TestGetProposalType(t *testing.T) {
	p1 := NewTextProposal("p1", "p1", qtypes.NewInt(1))
	require.Equal(t, ProposalTypeText, p1.GetProposalType())

	p2 := NewParameterProposal("p2", "p2", qtypes.NewInt(1), []Param{})
	require.Equal(t, ProposalTypeParameterChange, p2.GetProposalType())

	p3 := NewTaxUsageProposal("p3", "p3", qtypes.NewInt(1), nil, types.NewDec(1))
	require.Equal(t, ProposalTypeTaxUsage, p3.GetProposalType())
}
