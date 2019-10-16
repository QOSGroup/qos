# Transactions

## Validator

### TxCreateValidator

[Sending TxCreateValidator](../../command/qoscli.md#create-validator) to make a full node a validator node.

#### Struct

```go
type TxCreateValidator struct {
	Owner       btypes.AccAddress      `json:"owner"`        // owner, self delegator
	ConsPubKey  crypto.PubKey          `json:"cons_pub_key"` // validator public key
	BondTokens  btypes.BigInt          `json:"bond_tokens"`  // tokens
	IsCompound  bool                   `json:"is_compound"`  // whether the income is calculated as compound interest
	Description types.Description      `json:"description"`  // description
	Commission  types.CommissionRates  `json:"commission"`   // commission
	Delegations []types.DelegationInfo `json:"delegations"`  // initial delegations, only valid in iniChainer
}
```

#### Validations

This tx needs to pass the following validations:
- `moniker` can not be empty and `length(moniker) < 300`
- `owner` can not be empty and `owenr` has no other validator node
- `cons_pub_key` can not be empty and no other validator that with the same `cons_pub_key`
- length of `logo` and `website` in `description` must less than 255
- `bond_tokens` must be positive, and `owner` has enough qos to bond
- `commission` must be valid

#### Signer

`owner`

#### Tx Gas

`1.8QOS`

### TxModifyValidator

[Sending TxModifyValidator](../../command/qoscli.md#modify-validator) to modify the information of a validator.

#### Struct

```go
type TxModifyValidator struct {
	Owner          btypes.AccAddress `json:"owner"`           // owner address
	ValidatorAddr  btypes.ValAddress `json:"validator_addr"`  // validator address
	Description    types.Description `json:"description"`     // description
	CommissionRate *qtypes.Dec       `json:"commission_rate"` // commission
}
```

#### Validations

This tx needs to pass the following validations:
- `owner` can not be empty
- `validator_addr` can not be empty, and there has the validator that with this `owner` and `validator_addr`
- length of `logo` and `website` in `description` must less than 255
- `commission` must be valid

#### Signer

`owner`

#### Tx Gas

`0.18QOS`

### TxRevokeValidator

[Sending TxRevokeValidator](../../command/qoscli.md#revoke-validator) to inactive a validator.

#### Struct

```go
type TxRevokeValidator struct {
	Owner         btypes.AccAddress `json:"owner"`          // owner address
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // validator address
}
```

#### Validations

This tx needs to pass the following validations:
- `owner` can not be empty
- `validator_addr` can not be empty, and there has the validator that with this `owner` and `validator_addr`

#### Signer

`owner`

#### Tx Gas

`18QOS`

### TxActiveValidator

[Sending TxActiveValidator](../../command/qoscli.md#active-validator) to active a inactive validator.

#### Struct

```go
type TxActiveValidator struct {
	Owner         btypes.AccAddress `json:"owner"`          // owner
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // validator address
	BondTokens    btypes.BigInt     `json:"bond_tokens"`    // increase tokens to bond
}
```

#### Validations

This tx needs to pass the following validations:
- `owner` can not be empty
- `validator_addr` can not be empty, and there has the validator that with this `owner` and `validator_addr`
- `bond_tokens` can not be negative, and `owner` has enough qos to bond

#### Signer

`owner`

#### Tx Gas

`0`

## Delegation

### TxCreateDelegation

[Sending TxCreateDelegation](../../command/qoscli.md#delegate) to create delegation relationship with a validator.

#### Struct

```go
type TxCreateDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // delegator address
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // validator address
	Amount        btypes.BigInt     `json:"amount"`         // tokens
	IsCompound    bool              `json:"is_compound"`    // whether the income is calculated as compound interest
}
```
#### Validations

This tx needs to pass the following validations:
- `owner` can not be empty
- `validator_addr` can not be empty, and there has the validator that with this `owner` and `validator_addr`
- `amount` must be positive, and `owner` has enough qos to bond

#### Signer

`Delegator`

#### Tx Gas

`0` 

### TxModifyCompound

[Sending TxModifyCompound](../../command/qoscli.md#modify-compound) to modify the compound of a delegation.

#### Struct

```go
type TxModifyCompound struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // delegator address
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // validator address
	IsCompound    bool              `json:"is_compound"`    // whether the income is calculated as compound interest
}
```

#### Validations

This tx needs to pass the following validations:
- `delegator` can not be empty
- `validator_addr` can not be empty, and there has the delegation that with this `delegator` and `validator_addr`
- `is_compound` must not the same as the old value

#### Signer

`delegator`

#### Tx Gas

`0` 

### TxUnbondDelegation

[Sending TxUnbondDelegation](../../command/qoscli.md#unbond) to unbond tokens from delegations.

#### Struct

```go
type TxUnbondDelegation struct {
	Delegator     btypes.AccAddress `json:"delegator"`      // delegator address
	ValidatorAddr btypes.ValAddress `json:"validator_addr"` // validator address
	UnbondAmount  btypes.BigInt     `json:"unbond_amount"`  // unbond tokens
	UnbondAll     bool              `json:"unbond_all"`     // whether unbond all
}
```

#### Validations

This tx needs to pass the following validations:
- if `unbond_all` not true, `unbond_amount` must be positive and less than tokens in the delegation
- there must exist delegation that with this `delegator` and `validator_addr`

#### SIgner

`delegator`

#### Tx Gas

`0.18QOS`

### TxCreateReDelegation

[Sending TxCreateReDelegation](../../command/qoscli.md#redelegate) to create a redelegation information.

#### Struct

```go
type TxCreateReDelegation struct {
	Delegator         btypes.AccAddress `json:"delegator"`           // delegator address
	FromValidatorAddr btypes.ValAddress `json:"from_validator_addr"` // source validator address
	ToValidatorAddr   btypes.ValAddress `json:"to_validator_addr"`   // target validator address
	Amount            btypes.BigInt     `json:"amount"`              // tokens
	RedelegateAll     bool              `json:"redelegate_all"`      // whether redelegate all
	Compound          bool              `json:"compound"`            // whether the income is calculated as compound interest
}
```

#### Validations

This tx needs to pass the following validations:
- if `redelegate_all` not true, `Amount` must be positive
- there has a delegation that with this `delegator` and `from_validator_addr`
- `to_validator_addr` can not be empty, and there has validator that with this `to_validator_addr`

#### Signer

`delegator`

#### Tx Gas

`0`
