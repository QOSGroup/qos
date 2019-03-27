package gov

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	qcliacc "github.com/QOSGroup/qbase/client/account"
	btypes "github.com/QOSGroup/qbase/types"

	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qos/module/gov"
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
			pID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			path := gov.BuildQueryProposalPath(pID)
			res, err := cliCtx.Query(path, []byte{})

			if err != nil {
				return nil
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var result types.Proposal
			cliCtx.Codec.UnmarshalJSON(res, &result)

			return cliCtx.PrintResult(result)
		},
	}
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

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var depositorAddr btypes.Address
			var voterAddr btypes.Address
			var status types.ProposalStatus

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagDepositor); err == nil {
				depositorAddr = d
			}

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagVoter); err == nil {
				voterAddr = d
			}

			status = toProposalStatus(viper.GetString(flagStatus))

			queryParam := gov.QueryProposalsParam{
				Depositor: depositorAddr,
				Voter:     voterAddr,
				Status:    status,
				Limit:     uint64(viper.GetInt64(flagNumLimit)),
			}

			data, err := cliCtx.Codec.MarshalJSON(queryParam)
			if err != nil {
				return err
			}

			path := gov.BuildQueryProposalsPath()
			res, err := cliCtx.Query(path, data)

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var result []types.Proposal
			err = cdc.UnmarshalJSON(res, &result)
			if err != nil {
				return err
			}

			if len(result) == 0 {
				return fmt.Errorf("No matching proposals found")
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
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			if _, err := btypes.GetAddrFromBech32(args[1]); err != nil {
				return fmt.Errorf("voter-addr %s is not a valid address value", args[1])
			}

			path := gov.BuildQueryVotePath(pID, args[1])
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var vote types.Vote
			if err := cliCtx.Codec.UnmarshalJSON(res, &vote); err != nil {
				return err
			}

			return cliCtx.PrintResult(vote)
		},
	}
}

func queryVotesCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			path := gov.BuildQueryVotesPath(pID)
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var votes []types.Vote
			if err := cliCtx.Codec.UnmarshalJSON(res, &votes); err != nil {
				return err
			}

			if len(votes) == 0 {
				return errors.New("no votes found")
			}

			return cliCtx.PrintResult(votes)
		},
	}
}

func queryDepositCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [proposal-id] [depositer-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			if _, err := btypes.GetAddrFromBech32(args[1]); err != nil {
				return fmt.Errorf("depositer-addr %s is not a valid address value", args[1])
			}

			path := gov.BuildQueryDepositPath(pID, args[1])
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var deposit types.Deposit
			if err := cliCtx.Codec.UnmarshalJSON(res, &deposit); err != nil {
				return nil
			}

			return cliCtx.PrintResult(deposit)
		},
	}
}

func queryDepositsCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposits [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			path := gov.BuildQueryDepositsPath(pID)
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var deposits []types.Deposit
			if err := cliCtx.Codec.UnmarshalJSON(res, &deposits); err != nil {
				return err
			}

			return cliCtx.PrintResult(deposits)
		},
	}
}

func queryTallyCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a proposal vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			pID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal id %s is not a valid uint value", args[0])
			}

			path := gov.BuildQueryTallyPath(pID)
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var result types.TallyResult
			if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}
}

func queryParamsCommand(cdc *go_amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			path := gov.BuildQueryParamsPath()
			res, err := cliCtx.Query(path, []byte{})
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("no result found")
			}

			var result gov.Params
			if err := cliCtx.Codec.UnmarshalJSON(res, &result); err != nil {
				return err
			}

			return cliCtx.PrintResult(result)
		},
	}
}
