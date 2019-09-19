package types

var (
	// 事件类型
	EventTypeCreateQsc = "create-qsc" // 创建QSC
	EventTypeIssueQsc  = "issue-qsc"  // 发行QSC

	// 事件参数
	AttributeKeyModule  = "qsc"     // 模块名
	AttributeKeyQsc     = "name"    // QSC名称
	AttributeKeyCreator = "creator" // QSC创建账户
	AttributeKeyBanker  = "banker"  // QSC发行接收账户
	AttributeKeyTokens  = "tokens"  // QSC发行量
)
