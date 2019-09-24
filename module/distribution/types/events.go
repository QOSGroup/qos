package types

// Distribution module event types
var (
	// 事件类型
	EventTypeProposerReward   = "proposer_reward"   // 提议奖励
	EventTypeDelegatorRewards = "delegator_rewards" // 所有委托奖励
	EventTypeCommunity        = "community"         // 社区费池
	EventTypeCommission       = "commission"        // 佣金
	EventTypeDelegatorReward  = "delegator_reward"  // 委托奖励
	EventTypeDelegate         = "delegate"          // 委托

	// 事件参数
	AttributeKeyTokens    = "tokens"    // tokens数量
	AttributeKeyValidator = "validator" // 验证节点
	AttributeKeyDelegator = "delegator" // 委托账户
)
