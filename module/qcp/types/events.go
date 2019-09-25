package types

import "github.com/QOSGroup/qbase/qcp"

var (
	// 事件类型
	EventTypeInitQcp = "init-qcp" // 初始化联盟链

	// 事件参数
	AttributeKeyModule  = qcp.EventModule // 模块名称
	AttributeKeyQcp     = "chain-id"      // 联盟链chain-id
	AttributeKeyCreator = "creator"       // 创建账户地址
)
