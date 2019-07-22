package types

var (
	EventTypeCreateValidator   = "create-validator"
	EventTypeModifyValidator   = "modify-validator"
	EventTypeRevokeValidator   = "revoke-validator"
	EventTypeActiveValidator   = "active-validator"
	EventTypeInactiveValidator = "inactive-validator"
	EventTypeCloseValidator    = "close-validator"

	EventTypeCreateDelegation   = "create-delegation"
	EventTypeModifyCompound     = "modify-compound"
	EventTypeUnbondDelegation   = "unbond-delegation"
	EventTypeCreateReDelegation = "create-redelegation"

	AttributeKeyModule       = "stake"
	AttributeKeyHeight       = "height"
	AttributeKeyValidator    = "validator"
	AttributeKeyNewValidator = "new-validator"
	AttributeKeyOwner        = "owner"
	AttributeKeyDelegator    = "delegator"
)
