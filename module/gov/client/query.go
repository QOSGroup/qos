package client

import (
	"errors"
	"fmt"

	"strconv"
	"strings"

	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/spf13/viper"

	qcliacc "github.com/QOSGroup/qbase/client/account"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"

	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/spf13/cobra"
	go_amino "github.com/tendermint/go-amino"
)

func queryProposalCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a signal proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			result, err := getProposal(cliCtx, pID)
			if err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}
}

func getProposal(cliContext context.CLIContext, pid int64) (result types.Proposal, err error) {
	path := mapper.BuildQueryProposalPath(pid)
	res, err := cliContext.Query(path, []byte{})

	if err != nil {
		return types.Proposal{}, err
	}

	if len(res) == 0 {
		return types.Proposal{}, context.RecordsNotFoundError
	}

	err = cliContext.Codec.UnmarshalJSON(res, &result)
	return
}

func queryProposalsCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(`
Query for a all proposals. You can filter the returns with the following flags:

$ qos query gov proposals --depositor address1kfp9ueerfdjctqzfs4nlpqe6u0gefnz534xeef
$ qos query gov proposals --voter address1kfp9ueerfdjctqzfs4nlpqe6u0gefnz534xeef
$ qos query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
`),
		RunE: func(cmd *cobra.Command, args []string) error {

			limit := viper.GetInt64(flagNumLimit)
			if limit <= 0 {
				limit = int64(10)
			}

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var depositorAddr btypes.AccAddress
			var voterAddr btypes.AccAddress
			var status types.ProposalStatus

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagDepositor); err == nil {
				depositorAddr = d
			}

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagVoter); err == nil {
				voterAddr = d
			}

			status = toProposalStatus(viper.GetString(flagStatus))

			queryParam := mapper.QueryProposalsParam{
				Depositor: depositorAddr,
				Voter:     voterAddr,
				Status:    status,
				Limit:     limit,
			}

			result, err := queryProposalsByParams(cliCtx, queryParam)
			if err != nil {
				return err
			}

			if len(result) == 0 {
				return fmt.Errorf("no matching proposals found")
			}

			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagNumLimit, "", "(optional) limit to latest [number] proposals. Defaults to all proposals")
	cmd.Flags().String(flagDepositor, "", "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")

	return cmd
}

func queryProposalsByParams(ctx context.CLIContext, param mapper.QueryProposalsParam) ([]types.Proposal, error) {
	data, err := ctx.Codec.MarshalJSON(param)
	if err != nil {
		return nil, err
	}

	path := mapper.BuildQueryProposalsPath()
	res, err := ctx.Query(path, data)

	if len(res) == 0 {
		return nil, errors.New("no result found")
	}

	var result []types.Proposal
	err = ctx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func toProposalStatus(statusStr string) types.ProposalStatus {
	switch statusStr {
	case "DepositPeriod", "deposit_period":
		return types.StatusDepositPeriod
	case "VotingPeriod", "voting_period":
		return types.StatusVotingPeriod
	case "Passed", "passed":
		return types.StatusPassed
	case "Rejected", "rejected":
		return types.StatusRejected
	default:
		return types.StatusNil
	}
}

func queryVoteCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [voter]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			addr, err := qcliacc.GetAddrFromValue(cliCtx, args[1])
			if err != nil {
				return fmt.Errorf("voter %s is not a valid address value", args[1])
			}

			vote, err := getProposalVote(cliCtx, pID, addr)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(vote)
		},
	}
}

func getProposalVote(cliContext context.CLIContext, pid int64, voter btypes.AccAddress) (vote types.Vote, err error) {
	path := mapper.BuildQueryVotePath(pid, voter.String())
	res, err := cliContext.Query(path, []byte{})
	if err != nil {
		return vote, err
	}

	if len(res) == 0 {
		return vote, context.RecordsNotFoundError
	}

	if err = cliContext.Codec.UnmarshalJSON(res, &vote); err != nil {
		return vote, err
	}

	return
}

func queryVotesCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			votes, err := queryProposalVotes(cliCtx, pID)
			if err != nil {
				return err
			}

			if len(votes) == 0 {
				return errors.New("no votes found")
			}

			return cliCtx.PrintResult(votes)
		},
	}
}

func queryProposalVotes(cliContext context.CLIContext, pid int64) ([]types.Vote, error) {
	path := mapper.BuildQueryVotesPath(pid)
	res, err := cliContext.Query(path, []byte{})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, context.RecordsNotFoundError
	}

	var votes []types.Vote
	if err := cliContext.Codec.UnmarshalJSON(res, &votes); err != nil {
		return nil, err
	}

	return votes, nil
}

func queryDepositCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [proposal-id] [depositor]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			addr, err := qcliacc.GetAddrFromValue(cliCtx, args[1])
			if err != nil {
				return fmt.Errorf("depositer %s is not a valid address value", args[1])
			}

			deposit, err := getProposalDeposit(cliCtx, pID, addr)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(deposit)
		},
	}
}

func getProposalDeposit(cliContext context.CLIContext, pid int64, addr btypes.AccAddress) (deposit types.Deposit, err error) {
	path := mapper.BuildQueryDepositPath(pid, addr.String())
	res, err := cliContext.Query(path, []byte{})
	if err != nil {
		return deposit, err
	}

	if len(res) == 0 {
		return deposit, context.RecordsNotFoundError
	}

	if err = cliContext.Codec.UnmarshalJSON(res, &deposit); err != nil {
		return deposit, err
	}

	return
}

func queryDepositsCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposits [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			deposits, err := queryProposalDeposits(cliCtx, pID)
			if err != nil {
				return nil
			}

			if len(deposits) == 0 {
				return errors.New("no result found")
			}

			return cliCtx.PrintResult(deposits)
		},
	}
}

func queryProposalDeposits(cliContext context.CLIContext, pid int64) ([]types.Deposit, error) {
	path := mapper.BuildQueryDepositsPath(pid)
	res, err := cliContext.Query(path, []byte{})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.New("no result found")
	}

	var deposits []types.Deposit
	if err := cliContext.Codec.UnmarshalJSON(res, &deposits); err != nil {
		return nil, err
	}

	return deposits, nil
}

func queryTallyCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a proposal vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			result, err := getProposalTally(cliCtx, pID)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(result)
		},
	}
}

func getProposalTally(cliContext context.CLIContext, pid int64) (result types.TallyResult, err error) {
	path := mapper.BuildQueryTallyPath(pid)
	res, err := cliContext.Query(path, []byte{})
	if err != nil {
		return result, err
	}

	if len(res) == 0 {
		return result, context.RecordsNotFoundError
	}

	if err = cliContext.Codec.UnmarshalJSON(res, &result); err != nil {
		return result, err
	}

	return
}

func queryParamsCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			module := viper.GetString(flagModule)
			key := viper.GetString(flagParamKey)

			result, err := queryModuleParams(cliCtx, module, key)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagModule, "", "(optional) module name.")
	cmd.Flags().String(flagParamKey, "", "(optional) parameter name")

	return cmd
}

func queryModuleParams(cliContext context.CLIContext, module, key string) (rest interface{}, err error) {
	if len(key) != 0 && len(module) == 0 {
		return rest, errors.New("module is empty")
	}

	mod := 0
	var path string
	if len(module) == 0 {
		path = mapper.BuildQueryParamsPath()
	} else if len(key) == 0 {
		mod = 1
		path = mapper.BuildQueryModuleParamsPath(module)
	} else {
		mod = 2
		path = mapper.BuildQueryParamPath(module, key)
	}
	res, err := cliContext.Query(path, []byte{})
	if err != nil {
		return rest, err
	}

	if len(res) == 0 {
		return rest, context.RecordsNotFoundError
	}

	if mod == 0 {
		var result []qtypes.ParamSet
		if err = cliContext.Codec.UnmarshalJSON(res, &result); err != nil {
			return result, err
		}

		if len(result) == 0 {
			return result, context.RecordsNotFoundError
		}

		return result, nil
	} else if mod == 1 {
		var result qtypes.ParamSet
		if err = cliContext.Codec.UnmarshalJSON(res, &result); err != nil {
			return result, err
		}

		return result, nil
	} else {
		return string(res), nil
	}
}
