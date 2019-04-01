package types

import (
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"strings"
	"time"
)

type Proposal struct {
	ProposalContent `json:"proposal_content"` // Proposal content interface

	ProposalID uint64 `json:"proposal_id"` //  ID of the proposal

	Status           ProposalStatus `json:"proposal_status"`    //  Status of the Proposal {Pending, Active, Passed, Rejected}
	FinalTallyResult TallyResult    `json:"final_tally_result"` //  Result of Tallys

	SubmitTime     time.Time `json:"submit_time"`      //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime time.Time `json:"deposit_end_time"` // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit   uint64    `json:"total_deposit"`    //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime   time.Time `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingStartHeight uint64    `json:"voting_start_height"`
	VotingEndTime     time.Time `json:"voting_end_time"` // Time that the VotingPeriod for this proposal will end and votes will be tallied
}

type ProposalContent interface {
	GetTitle() string
	GetDescription() string
	GetDeposit() uint64
	GetProposalType() ProposalType
}

type ProposalResult string

const (
	PASS       ProposalResult = "pass"
	REJECT     ProposalResult = "reject"
	REJECTVETO ProposalResult = "reject-veto"
)

// Type that represents Proposal Status as a byte
type ProposalStatus byte

//nolint
const (
	StatusNil           ProposalStatus = 0x00
	StatusDepositPeriod ProposalStatus = 0x01
	StatusVotingPeriod  ProposalStatus = 0x02
	StatusPassed        ProposalStatus = 0x03
	StatusRejected      ProposalStatus = 0x04
)

func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusDepositPeriod ||
		status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected {
		return true
	}
	return false
}

// Tally Results
type TallyResult struct {
	Yes        int64 `json:"yes"`
	Abstain    int64 `json:"abstain"`
	No         int64 `json:"no"`
	NoWithVeto int64 `json:"no_with_veto"`
}

func NewTallyResult(yes, abstain, no, noWithVeto int64) TallyResult {
	return TallyResult{
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
	}
}

// checks if two proposals are equal
func EmptyTallyResult() TallyResult {
	return TallyResult{
		Yes:        0,
		Abstain:    0,
		No:         0,
		NoWithVeto: 0,
	}
}

// checks if two proposals are equal
func (tr TallyResult) Equals(comp TallyResult) bool {
	return tr.Yes == comp.Yes &&
		tr.Abstain == comp.Abstain &&
		tr.No == comp.No &&
		tr.NoWithVeto == comp.NoWithVeto
}

func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:        %d
  Abstain:    %d
  No:         %d
  NoWithVeto: %d`, tr.Yes, tr.Abstain, tr.No, tr.NoWithVeto)
}

// Type that represents Proposal Type as a byte
type ProposalType byte

//nolint
const (
	ProposalTypeNil             ProposalType = 0x00
	ProposalTypeText            ProposalType = 0x01
	ProposalTypeParameterChange ProposalType = 0x02
	ProposalTypeTaxUsage        ProposalType = 0x03
)

// String to proposalType byte. Returns 0xff if invalid.
func ProposalTypeFromString(str string) (ProposalType, error) {
	switch strings.ToLower(str) {
	case "text":
		return ProposalTypeText, nil
	case "parameterchange":
		return ProposalTypeParameterChange, nil
	case "taxusage":
		return ProposalTypeTaxUsage, nil
	default:
		return ProposalType(0xff), fmt.Errorf("'%s' is not a valid proposal type", str)
	}
}

// Turns VoteOption byte to String
func (pt ProposalType) String() string {
	switch pt {
	case ProposalTypeText:
		return "Text"
	case ProposalTypeParameterChange:
		return "Parameter"
	case ProposalTypeTaxUsage:
		return "TaxUsage"
	default:
		return ""
	}
}

// is defined GetProposalType?
func ValidProposalType(pt ProposalType) bool {
	if pt == ProposalTypeText ||
		pt == ProposalTypeParameterChange ||
		pt == ProposalTypeTaxUsage {
		return true
	}
	return false
}

// Text Proposals
type TextProposal struct {
	Title       string `json:"title"`       //  Title of the proposal
	Description string `json:"description"` //  Description of the proposal
	Deposit     uint64 `json:"deposit"`     //	Deposit of the proposal
}

func NewTextProposal(title, description string, deposit uint64) TextProposal {
	return TextProposal{
		Title:       title,
		Description: description,
		Deposit:     deposit,
	}
}

// Implements Proposal Interface
var _ ProposalContent = TextProposal{}

// nolint
func (tp TextProposal) GetTitle() string              { return tp.Title }
func (tp TextProposal) GetDescription() string        { return tp.Description }
func (tp TextProposal) GetDeposit() uint64            { return tp.Deposit }
func (tp TextProposal) GetProposalType() ProposalType { return ProposalTypeText }

// TaxUsage Proposals
type TaxUsageProposal struct {
	TextProposal
	DestAddress btypes.Address `json:"dest_address"`
	Percent     types.Dec      `json:"percent"`
}

func NewTaxUsageProposal(title, description string, deposit uint64, destAddress btypes.Address, percent types.Dec) TaxUsageProposal {
	return TaxUsageProposal{
		TextProposal: TextProposal{
			Title:       title,
			Description: description,
			Deposit:     deposit,
		},
		DestAddress: destAddress,
		Percent:     percent,
	}
}

// Implements Proposal Interface
var _ ProposalContent = TaxUsageProposal{}

// nolint
func (tp TaxUsageProposal) GetTitle() string              { return tp.Title }
func (tp TaxUsageProposal) GetDescription() string        { return tp.Description }
func (tp TaxUsageProposal) GetDeposit() uint64            { return tp.Deposit }
func (tp TaxUsageProposal) GetProposalType() ProposalType { return ProposalTypeTaxUsage }

// Parameters change Proposals
type ParameterProposal struct {
	TextProposal
	Params []Param `json:"params"`
}

func NewParameterProposal(title, description string, deposit uint64, params []Param) ParameterProposal {
	return ParameterProposal{
		TextProposal: TextProposal{
			Title:       title,
			Description: description,
			Deposit:     deposit,
		},
		Params: params,
	}
}

// Implements Proposal Interface
var _ ProposalContent = ParameterProposal{}

// nolint
func (tp ParameterProposal) GetTitle() string              { return tp.Title }
func (tp ParameterProposal) GetDescription() string        { return tp.Description }
func (tp ParameterProposal) GetDeposit() uint64            { return tp.Deposit }
func (tp ParameterProposal) GetProposalType() ProposalType { return ProposalTypeParameterChange }

type Param struct {
	Module string `json:"module"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func (param Param) String() string {
	return fmt.Sprintf(`
  Module:     %s
  Key:    	  %s
  Value:      %s`, param.Module, param.Key, param.Value)
}
