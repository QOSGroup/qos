package types

// Governance tags
var (
	EventTypeSubmitProposal   = "submit-proposal"
	EventTypeDepositProposal  = "deposit-proposal"
	EventTypeVoteProposal     = "vote-proposal"
	EventTypeInactiveProposal = "inactive-proposal"
	EventTypeActiveProposal   = "active-proposal"

	AttributeKeyModule             = "gov"
	AttributeKeyProposer           = "proposer"
	AttributeKeyProposalID         = "proposal-id"
	AttributeKeyDepositor          = "depositor"
	AttributeKeyVoter              = "voter"
	AttributeKeyProposalResult     = "proposal-result"
	AttributeKeyProposalType       = "proposal-type"
	AttributeKeyDropped            = "proposal-dropped"
	AttributeKeyResultPassed       = "proposal-passed"
	AttributeKeyResultRejected     = "proposal-rejected"
	AttributeKeyResultVetoRejected = "proposal-veto-rejected"
)
