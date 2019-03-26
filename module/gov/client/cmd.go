package gov

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands()
}

func TxCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.PostCommands(
		ProposalCmd(cdc),
		DepositCmd(cdc),
		VoteCmd(cdc),
	)
}

const (
	flagTitle        = "title"
	flagDescription  = "description"
	flagProposalType = "proposal-type"
	flagProposer     = "proposer"
	flagDeposit      = "deposit"

	flagProposalID = "proposal-id"
	flagDepositor  = "depositor"
	flagAmount     = "amount"

	flagVoter      = "voter"
	flagVoteOption = "option"
)
