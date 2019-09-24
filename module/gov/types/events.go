package types

var (
	// 事件类型
	EventTypeSubmitProposal   = "submit-proposal"   // 提交提议
	EventTypeDepositProposal  = "deposit-proposal"  // 质押提议
	EventTypeVoteProposal     = "vote-proposal"     // 提议投票
	EventTypeInactiveProposal = "inactive-proposal" // 提议失效
	EventTypeActiveProposal   = "active-proposal"   // 提议结束

	// 事件参数
	AttributeKeyModule             = "gov"                    // 模块名
	AttributeKeyProposer           = "proposer"               // 提议账户
	AttributeKeyProposalID         = "proposal-id"            // 提议ID
	AttributeKeyDepositor          = "depositor"              // 质押账户
	AttributeKeyVoter              = "voter"                  // 投票账户
	AttributeKeyProposalResult     = "proposal-result"        // 提议结果
	AttributeKeyProposalType       = "proposal-type"          // 提议类型
	AttributeKeyDropped            = "proposal-dropped"       // 提议删除
	AttributeKeyResultPassed       = "proposal-passed"        // 提议通过
	AttributeKeyResultRejected     = "proposal-rejected"      // 提议拒绝
	AttributeKeyResultVetoRejected = "proposal-veto-rejected" // 提议拒绝（强烈反对）
)
