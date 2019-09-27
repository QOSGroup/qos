# State

`MapperName` is `governance`

## Proposal

Self-increase proposal ID:

- newProposalID: `newProposalID -> amino(newProposalID)`


proposal struct:
```go
type Proposal struct {
	ProposalContent `json:"proposal_content"` // proposal content, different types of proposal may be different

	ProposalID int64 `json:"proposal_id"` //  proposal ID

	Status           ProposalStatus `json:"proposal_status"`    //  proposal status
	FinalTallyResult TallyResult    `json:"final_tally_result"` //  proposal result

	SubmitTime     time.Time     `json:"submit_time"`      // proposal submit time
	DepositEndTime time.Time     `json:"deposit_end_time"` // deposit end time
	TotalDeposit   btypes.BigInt `json:"total_deposit"`    // amount of deposit

	VotingStartTime   time.Time `json:"voting_start_time"` // voting start time
	VotingStartHeight int64     `json:"voting_start_height"` // voting start height
	VotingEndTime     time.Time `json:"voting_end_time"` // voting end time
}
```

- proposal `proposals:{proposal_id} -> amino(Proposal)`
- proposal in depositing period `inactiveProposalQueue:{deposit_end_time}:{proposal_id} -> amino(proposal_id)`
- proposal in voting period  `activeProposalQueue:{voting_end_time}:{proposal_id} -> amino(proposal_id)`

## Deposit

struct:
```go
type Deposit struct {
	Depositor  types.AccAddress `json:"depositor"`   //  depositor address
	ProposalID int64            `json:"proposal_id"` //  proposal ID
	Amount     types.BigInt     `json:"amount"`      //  amount of QOS to be deposited
}
```

- deposit `deposits:{proposal_id}:{depositor} -> amino(Deposit)`

## Vote

struct:
```go
type Vote struct {
	Voter      types.AccAddress `json:"voter"`       //  voter address
	ProposalID int64            `json:"proposal_id"` //  proposal ID
	Option     VoteOption       `json:"option"`      //  vote option
}
```

- vote `votes:{proposal_id}:{voter} -> amino(Vote)`