package client

import (
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func QueryCommands(cdc *amino.Codec) []*cobra.Command {
	return bctypes.GetCommands(
		queryProposalCommand(cdc),
		queryProposalsCommand(cdc),
		queryVoteCommand(cdc),
		queryVotesCommand(cdc),
		queryDepositCommand(cdc),
		queryDepositsCommand(cdc),
		queryTallyCommand(cdc),
		queryParamsCommand(cdc),
	)
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

	flagDestAddress = "dest-address"
	flagPercent     = "percent"

	flagParams = "params"

	flagInflationPhrases = "inflation-phrases"
	flagTotalAmount      = "total-amount"

	flagStatus   = "status"
	flagNumLimit = "limit"

	flagModule   = "module"
	flagParamKey = "key"

	flagVersion       = "version"
	flagDataHeight    = "data-height"
	flagGenesisFile   = "genesis-file"
	flagGenesisMD5    = "genesis-md5"
	flagForZeroHeight = "for-zero-height"
)
