package types

import (
	"fmt"
	"github.com/QOSGroup/qbase/types"
)

// Deposit
type Deposit struct {
	Depositor  types.Address `json:"depositor"`   //  Address of the depositor
	ProposalID uint64        `json:"proposal_id"` //  proposalID of the proposal
	Amount     uint64        `json:"amount"`      //  Deposit amount
}

func (d Deposit) String() string {
	return fmt.Sprintf("Deposit by %s on Proposal %d is for the amount %d",
		d.Depositor, d.ProposalID, d.Amount)
}

type Deposits []Deposit

func (d Deposits) String() string {
	if len(d) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Deposits for Proposal %d:", d[0].ProposalID)
	for _, dep := range d {
		out += fmt.Sprintf("\n  %s: %d", dep.Depositor, dep.Amount)
	}
	return out
}

func (d Deposit) Equals(comp Deposit) bool {
	return d.Depositor.EqualsTo(comp.Depositor) && d.ProposalID == comp.ProposalID && d.Amount == comp.Amount
}
