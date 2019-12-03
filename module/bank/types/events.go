package types

import "github.com/QOSGroup/qos/types"

var (
	// 事件类型
	EventTypeTransfer       = "transfer"                    // 转账
	EventTypeSend           = "send"                        // 发送
	EventTypeReceive        = "receive"                     // 接收
	EventTypeInvariantCheck = types.EventTypeInvariantCheck // 数据检查
	EventTypeRelease        = "release"                     // 锁定释放

	// 事件参数
	AttributeKeyModule  = "bank"    // 模块名
	AttributeKeyAddress = "address" // 账户地址
	AttributeKeyQOS     = "qos"     // QOS
	AttributeKeyQSCs    = "qscs"    // QSC代币
	AttributeKeySender  = "sender"  // 数据检查发送账户地址
	AttributeKeyHeight  = "height"  // 高度
)
