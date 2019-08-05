package types

import "github.com/QOSGroup/qos/types"

var (
	EventTypeSend           = "send"
	EventTypeReceive        = "receive"
	EventTypeInvariantCheck = types.EventTypeInvariantCheck

	AttributeKeyModule  = "transfer"
	AttributeKeyAddress = "address"
	AttributeKeyQOS     = "qos"
	AttributeKeyQSCs    = "qscs"
	AttributeKeySender  = "sender"
	AttributeKeyHeight  = "height"
)
