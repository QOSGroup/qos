# State

## Guardian

Struct:
```go
type Guardian struct {
	Description  string            `json:"description"`   // description
	GuardianType GuardianType      `json:"guardian_type"` // guardian type: Genesis, config in genesis.json; Ordinary created by guardian. 
	Address      btypes.AccAddress `json:"address"`       // address
	Creator      btypes.AccAddress `json:"creator"`       // address of creator
}
```
When `GuardianType` is `Genesis`, `Creator` is empty.

storage:
- guardian: `0x00 address -> amino(guardian)`

## Halt network

storage reason for network shutdown:
- halt network: `0x01 -> amino(reason)`