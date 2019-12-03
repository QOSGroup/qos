package types

var (
	// 事件类型
	EventTypeCreateValidator   = "create-validator"   // 创建验证节点
	EventTypeModifyValidator   = "modify-validator"   // 修改验证节点
	EventTypeRevokeValidator   = "revoke-validator"   // 撤销验证节点
	EventTypeActiveValidator   = "active-validator"   // 激活验证节点
	EventTypeInactiveValidator = "inactive-validator" // 失活验证节点
	EventTypeCloseValidator    = "close-validator"    // 关闭验证节点

	EventTypeCreateDelegation   = "create-delegation"   // 创建委托
	EventTypeModifyCompound     = "modify-compound"     // 修改委托
	EventTypeUnbondDelegation   = "unbond-delegation"   // 取消委托
	EventTypeCreateReDelegation = "create-redelegation" // 转委托

	EventTypeMissingVote = "missing-vote" // 未参与投票

	EventTypeSlash = "slash" // 惩罚

	// 事件参数
	AttributeKeyModule       = "stake"         // 模块名称
	AttributeKeyHeight       = "height"        // 高度
	AttributeKeyValidator    = "validator"     // 验证节点
	AttributeKeyNewValidator = "new-validator" // 转委托目标验证节点
	AttributeKeyOwner        = "owner"         // 验证节点持有账户
	AttributeKeyDelegator    = "delegator"     // 委托账户
	AttributeKeyTokens       = "tokens"        // 质押量
	AttributeKeyReason       = "reason"        // 惩罚原因
	AttributeKeyMissedBlocks = "missed-blocks" // 漏块数

	// 常量
	AttributeValueDoubleSign = "double_sign" // 双签
	AttributeValueDownTime   = "down_time"   // 漏块
)
