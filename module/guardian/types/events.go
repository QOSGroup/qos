package types

var (
	// 事件类型
	EventTypeAddGuardian    = "add-guardian"    // 添加系统账户
	EventTypeDeleteGuardian = "delete-guardian" // 删除系统账户
	EventTypeHaltNetwork    = "halt-network"    // 停止网络

	// 事件参数
	AttributeKeyModule   = "guardian"  // 模块名称
	AttributeKeyGuardian = "guardian"  // 账户地址
	AttributeKeyCreator  = "creator"   // 创建账户地址
	AttributeKeyDeleteBy = "delete-by" // 删除账户地址
	AttributeKeyReason   = "reason"    // 操作原因
)
