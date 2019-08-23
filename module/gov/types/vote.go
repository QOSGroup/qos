package types

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/types"
	"strings"
)

// Vote
type Vote struct {
	Voter      types.Address `json:"voter"`       //  address of the voter
	ProposalID uint64        `json:"proposal_id"` //  proposalID of the proposal
	Option     VoteOption    `json:"option"`      //  option from OptionSet chosen by the voter
}

func (v Vote) String() string {
	return fmt.Sprintf("Voter %s voted with option %s on proposal %d", v.Voter, v.Option, v.ProposalID)
}

// Votes is a collection of Vote
type Votes []Vote

func (v Votes) String() string {
	out := fmt.Sprintf("Votes for Proposal %d:", v[0].ProposalID)
	for _, vot := range v {
		out += fmt.Sprintf("\n  %s: %s", vot.Voter, vot.Option)
	}
	return out
}

func (v Vote) Equals(comp Vote) bool {
	return v.Voter.EqualsTo(comp.Voter) && v.ProposalID == comp.ProposalID && v.Option == comp.Option
}

// Type that represents VoteOption as a byte
type VoteOption byte

//nolint
const (
	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionAbstain    VoteOption = 0x02
	OptionNo         VoteOption = 0x03
	OptionNoWithVeto VoteOption = 0x04
)

// Type that deduct deposits
type DeductOption byte

//nolint
const (
	DepositDeductNone DeductOption = 0x00
	DepositDeductPart DeductOption = 0x01
	DepositDeductAll  DeductOption = 0x02
)

// String to proposalType byte.  Returns ff if invalid.
func VoteOptionFromString(str string) (VoteOption, error) {
	switch strings.ToLower(str) {
	case "yes":
		return OptionYes, nil
	case "abstain":
		return OptionAbstain, nil
	case "no":
		return OptionNo, nil
	case "nowithveto":
		return OptionNoWithVeto, nil
	default:
		return VoteOption(0xff), fmt.Errorf("'%s' is not a valid vote option", str)
	}
}

// Is defined VoteOption
func ValidVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionAbstain ||
		option == OptionNo ||
		option == OptionNoWithVeto {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (vo VoteOption) Marshal() ([]byte, error) {
	return []byte{byte(vo)}, nil
}

// Unmarshal needed for protobuf compatibility
func (vo *VoteOption) Unmarshal(data []byte) error {
	*vo = VoteOption(data[0])
	return nil
}

// Marshals to JSON using string
func (vo VoteOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (vo *VoteOption) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := VoteOptionFromString(s)
	if err != nil {
		return err
	}
	*vo = bz2
	return nil
}

// Turns VoteOption byte to String
func (vo VoteOption) String() string {
	switch vo {
	case OptionYes:
		return "Yes"
	case OptionAbstain:
		return "Abstain"
	case OptionNo:
		return "No"
	case OptionNoWithVeto:
		return "NoWithVeto"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (vo VoteOption) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(vo.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(vo))))
	}
}
