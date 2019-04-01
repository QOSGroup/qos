package gov

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

/*
	custom query : /custom/gov/$queryPath

	$queryPath:
		1. /proposals
		2. /proposal/:pid
		3. /votes/:pid
		4. /vote/:pid/:addr
		5. /deposit/:pid/:addr
		6. /deposits/:pid
		7. /tally/:pid
		8. /params
*/

//nolint
const (
	GOV        = "gov"
	Proposals  = "proposals"
	Proposal   = "proposal"
	Votes      = "votes"
	Vote       = "vote"
	Deposits   = "deposits"
	Deposit    = "deposit"
	Tally      = "tally"
	ParamsPath = "params"
)

//nolint
func Query(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err btypes.Error) {

	defer func() {
		if r := recover(); r != nil {
			err = btypes.ErrInternal(string(debug.Stack()))
			return
		}
	}()

	var data []byte
	var e error

	switch route[0] {
	case Proposals:
		data, e = queryProposalsByParams(ctx, req.Data)
	case Proposal:
		data, e = queryProposal(ctx, route[1])
	case Votes:
		data, e = queryVotesByProposalID(ctx, route[1])
	case Vote:
		data, e = queryVoteByProposalIDAndVoter(ctx, route[1], route[2])
	case Deposits:
		data, e = queryDepositsByProposalID(ctx, route[1])
	case Deposit:
		data, e = queryDepositsByProposalIDAndDepositer(ctx, route[1], route[2])
	case Tally:
		data, e = queryTallyByProposalID(ctx, route[1])
	case ParamsPath:
		data, e = queryParams(ctx)
	default:
		data = nil
		e = errors.New("not found match path")
	}

	if e != nil {
		return nil, btypes.ErrInternal(e.Error())
	}

	return data, nil
}

func queryProposal(ctx context.Context, pid string) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	proposal, exsits := govMapper.GetProposal(ctx, proposalID)
	if !exsits {
		return nil, fmt.Errorf("proposal id %d not exsits", proposalID)
	}

	return govMapper.GetCodec().MarshalJSON(proposal)
}

func queryProposalsByParams(ctx context.Context, paramsData []byte) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	var params QueryProposalsParam
	if err := govMapper.GetCodec().UnmarshalJSON(paramsData, &params); err != nil {
		return nil, errors.New("params can not unmarshal")
	}

	result := govMapper.GetProposalsFiltered(ctx, params.Voter, params.Depositor, params.Status, params.Limit)

	return govMapper.GetCodec().MarshalJSON(result)
}

func queryVotesByProposalID(ctx context.Context, pid string) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exsits := govMapper.GetProposal(ctx, proposalID)
	if !exsits {
		return nil, fmt.Errorf("proposal id %d not exsits", proposalID)
	}

	iter := govMapper.GetVotes(ctx, proposalID)
	defer iter.Close()

	var votes []types.Vote
	for ; iter.Valid(); iter.Next() {
		vote := types.Vote{}
		govMapper.DecodeObject(iter.Value(), &vote)
		votes = append(votes, vote)
	}
	return govMapper.GetCodec().MarshalJSON(votes)
}

func queryVoteByProposalIDAndVoter(ctx context.Context, pid string, voter string) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exsits := govMapper.GetProposal(ctx, proposalID)
	if !exsits {
		return nil, fmt.Errorf("proposal id %d not exsits", proposalID)
	}

	voterAddress, err := btypes.GetAddrFromBech32(voter)
	if err != nil {
		return nil, fmt.Errorf("voter %s is not valid address", voter)
	}

	result, exsits := govMapper.GetVote(ctx, proposalID, voterAddress)
	if !exsits {
		return nil, fmt.Errorf("voter %s is not vote on proposal %s", voterAddress, pid)
	}

	return govMapper.GetCodec().MarshalJSON(result)
}

func queryDepositsByProposalIDAndDepositer(ctx context.Context, pid string, depositer string) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exsits := govMapper.GetProposal(ctx, proposalID)
	if !exsits {
		return nil, fmt.Errorf("proposal id %d not exsits", proposalID)
	}

	depositerAddress, err := btypes.GetAddrFromBech32(depositer)
	if err != nil {
		return nil, fmt.Errorf("depositer %s is not valid address", depositer)
	}

	result, exsits := govMapper.GetDeposit(ctx, proposalID, depositerAddress)
	if !exsits {
		return nil, fmt.Errorf("depositer %s is not deposit on proposal %s", depositer, pid)
	}
	return govMapper.GetCodec().MarshalJSON(result)
}

func queryDepositsByProposalID(ctx context.Context, pid string) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exsits := govMapper.GetProposal(ctx, proposalID)
	if !exsits {
		return nil, fmt.Errorf("proposal id %d not exsits", proposalID)
	}

	iter := govMapper.GetDeposits(ctx, proposalID)
	defer iter.Close()

	var deposits []types.Deposit
	for ; iter.Valid(); iter.Next() {
		deposit := types.Deposit{}
		govMapper.DecodeObject(iter.Value(), &deposit)
		deposits = append(deposits, deposit)
	}

	return govMapper.GetCodec().MarshalJSON(deposits)
}

func queryTallyByProposalID(ctx context.Context, pid string) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	proposal, exsits := govMapper.GetProposal(ctx, proposalID)
	if !exsits {
		return nil, fmt.Errorf("proposal id %d not exsits", proposalID)
	}

	_, result, _ := tally(ctx, govMapper, proposal)
	return govMapper.GetCodec().MarshalJSON(result)
}

func queryParams(ctx context.Context) ([]byte, error) {
	govMapper := GetGovMapper(ctx)

	var result = govMapper.GetParams()

	return govMapper.GetCodec().MarshalJSON(result)
}

//nolint
type QueryProposalsParam struct {
	Depositor btypes.Address       `json:"depositor"`
	Voter     btypes.Address       `json:"voter"`
	Status    types.ProposalStatus `json:"status"`
	Limit     uint64               `json:"limit"`
}

//nolint
func BuildQueryProposalPath(pid uint64) string {
	return fmt.Sprintf("custom/%s/%s/%d", GOV, Proposal, pid)
}

//nolint
func BuildQueryProposalsPath() string {
	return fmt.Sprintf("custom/%s/%s", GOV, Proposals)
}

//nolint
func BuildQueryVotePath(pid uint64, voterAddr string) string {
	return fmt.Sprintf("custom/%s/%s/%d/%s", GOV, Vote, pid, voterAddr)
}

//nolint
func BuildQueryVotesPath(pid uint64) string {
	return fmt.Sprintf("custom/%s/%s/%d", GOV, Votes, pid)
}

//nolint
func BuildQueryDepositPath(pid uint64, depositAddr string) string {
	return fmt.Sprintf("custom/%s/%s/%d/%s", GOV, Deposit, pid, depositAddr)
}

//nolint
func BuildQueryDepositsPath(pid uint64) string {
	return fmt.Sprintf("custom/%s/%s/%d", GOV, Deposits, pid)
}

//nolint
func BuildQueryTallyPath(pid uint64) string {
	return fmt.Sprintf("custom/%s/%s/%d", GOV, Tally, pid)
}

//nolint
func BuildQueryParamsPath() string {
	return fmt.Sprintf("custom/%s/%s", GOV, ParamsPath)
}
