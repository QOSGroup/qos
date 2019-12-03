package types

var (
	// 事件类型
	EventTypeCreateApprove   = "create-approve"   // 创建预授权
	EventTypeIncreaseApprove = "increase-approve" // 增加预授权
	EventTypeDecreaseApprove = "decrease-approve" // 减少预授权
	EventTypeUseApprove      = "use-approve"      // 使用预授权
	EventTypeCancelApprove   = "cancel-approve"   // 取消预授权

	// 事件参数
	AttributeKeyModule      = "approve"      // 模块
	AttributeKeyApproveFrom = "approve-from" // 授权账户
	AttributeKeyApproveTo   = "approve-to"   // 被授权账户
)
