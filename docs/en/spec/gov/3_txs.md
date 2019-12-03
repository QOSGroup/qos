# Transactions

Governance module contains proposal, deposit and vote transactions.

## Proposal

proposal types in QOS:

- ProposalTypeText           // simple text
- ProposalTypeParameterChange // parameter change
- ProposalTypeTaxUsage       // community tax use
- ProposalTypeModifyInflation // inflation modify
- ProposalTypeSoftwareUpgrade // software upgrade

[Submiting proposal](../../command/qoscli.md#submit-proposal) to submit a proposal in QOS network.

### Text

Ordinary text proposal, the proposal can be QOS network constructive opinions, new or improved functions.

#### Struct

```go
type TxProposal struct {
	Title          string             `json:"title"`           //  title
	Description    string             `json:"description"`     //  description
	ProposalType   types.ProposalType `json:"proposal_type"`   //  proposal type
	Proposer       btypes.AccAddress  `json:"proposer"`        //  proposer
	InitialDeposit btypes.BigInt      `json:"initial_deposit"` //  initial deposit by proposer
}
```

#### Validations

Make sure to pass the following check:
- title can not be empty and lte `MaxTitleLen`(default 200)
- description can not be empty and lte `MaxDescriptionLen`(default 1000)
- proposal_type must be `ProposalTypeText`

#### Signer

`proposer`

#### Tx Gas

0

### ParameterChange

Parameters can queried by [Query params](../../command/qoscli.md#query-params) can be changed by this proposal.

#### Struct

```go
type TxParameterChange struct {
	TxProposal                           // `ProposalType` is `ProposalTypeParameterChange`
	Params []types.Param `json:"params"` // parameters to be changed
}

type Param struct {
	Module string `json:"module"`   // module name
	Key    string `json:"key"` // parameter name
	Value  string `json:"value"` // parameter value
}
```

#### Validations

Make sure to pass the following check:
- title can not be empty and lte `MaxTitleLen`(default 200)
- description can not be empty and lte `MaxDescriptionLen`(default 1000)
- proposal_type must be `ProposalTypeParameterChange`
- no unfinished `ProposalTypeParameterChange` proposal in QOS network
- params can not be empty
- parameters in params must have the right type and hold the right value


#### Signer

`proposer`

#### Tx Gas

0

### TaxUsage 

Only [Guardian](../guardian) can submit this proposal.

#### Struct

```go
type TxTaxUsage struct {
	TxProposal                                          
	DestAddress btypes.AccAddress `json:"dest_address"` // address to accept QOS from community fees
	Percent     qtypes.Dec        `json:"percent"`      // percent of community fees
}
```

#### Validations

Make sure to pass the following check:
- title can not be empty and lte `MaxTitleLen`(default 200)
- description can not be empty and lte `MaxDescriptionLen`(default 1000)
- proposal_type must be `ProposalTypeTaxUsage`
- dest_address can not be null and must be a guardian
- percent ranges (0, 1]

#### Signer

`proposer`

#### Tx Gas

0

### ModifyInflation

QOS network supports modification of inflation rules that have not yet begun.

#### Struct

```go
type TxModifyInflation struct {
	TxProposal              
	TotalAmount      btypes.BigInt         `json:"total_amount"`      // total amount of QOS
	InflationPhrases mint.InflationPhrases `json:"inflation_phrases"` // inflation phrases
}
```

#### Validations

Make sure to pass the following check:
- title can not be empty and lte `MaxTitleLen`(default 200)
- description can not be empty and lte `MaxDescriptionLen`(default 1000)
- proposal_type must be `ProposalTypeModifyInflation`
- current and finished phrases can not be changed
- total_amount must be right

#### Signer

`proposer`

#### Tx Gas

0

### SoftwareUpgrade

Submiting softwareupgrade to enhance existing functional or add new features.

#### Struct

```go
type TxSoftwareUpgrade struct {
	TxProposal
	Version       string `json:"version"`         // QOS version
	DataHeight    int64  `json:"data_height"`     // data height
	GenesisFile   string `json:"genesis_file"`    // `genesis.json` file url
	GenesisMD5    string `json:"genesis_md5"`     // `genesis.json` file md5
	ForZeroHeight bool   `json:"for_zero_height"` // whether start from zero height
}
```

#### Validations

Make sure to pass the following check:
- title can not be empty and lte `MaxTitleLen`(default 200)
- description can not be empty and lte `MaxDescriptionLen`(default 1000)
- proposal_type must be `ProposalTypeSoftwareUpgrade`
- version can not be empty
- if for_zero_height is true, make sure data_height must gt 0 ,genesis_file and genesis_md5 are not empty.

#### Signer

`proposer`

#### Tx Gas

0

## Deposit

Depositing for proposals in `StatusDepositPeriod` status.

### Struct

```go
type TxDeposit struct {
	ProposalID int64             `json:"proposal_id"` // proposal ID
	Depositor  btypes.AccAddress `json:"depositor"`   // depositor address
	Amount     btypes.BigInt     `json:"amount"`      // amount of QOS to deposit
}
```

### Validations

Make sure to pass the following check:
- proposal must exist and status is `StatusDepositPeriod`
- amount must be positive
- depositor has enough QOS to deposit

### Signer

`depositor`

### Tx Gas

0

## Vote

Voting for proposals in `StatusVotingPeriod` status.

### Struct

```go
type TxVote struct {
	ProposalID int64             `json:"proposal_id"` // proposal ID
	Voter      btypes.AccAddress `json:"voter"`       // voter address
	Option     types.VoteOption  `json:"option"`      // options: Yes/Abstain/No/Nowithveto
}
```

### Validations

Make sure to pass the following check:
- voter must exist
- vote option must be one of `Yes/Abstain/No/Nowithveto`
- proposal must exist and status is `StatusVotingPeriod`

### Signer

`voter`

### Tx Gas

0