package mapper

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qos/module/params/mapper"
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
	TallyQuery = "tally"
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
	case TallyQuery:
		data, e = queryTallyByProposalID(ctx, route[1])
	case ParamsPath:
		if len(route) == 1 {
			data, e = queryParams(ctx)
		} else if len(route) == 2 {
			data, e = queryModuleParams(ctx, route[1])
		} else {
			data, e = queryParam(ctx, route[1], route[2])
		}
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
	govMapper := GetMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	proposal, exists := govMapper.GetProposal(proposalID)
	if !exists {
		return nil, fmt.Errorf("proposal id %d not exists", proposalID)
	}

	return govMapper.GetCodec().MarshalJSON(proposal)
}

func queryProposalsByParams(ctx context.Context, paramsData []byte) ([]byte, error) {
	govMapper := GetMapper(ctx)

	var params QueryProposalsParam
	if err := govMapper.GetCodec().UnmarshalJSON(paramsData, &params); err != nil {
		return nil, errors.New("params can not unmarshal")
	}

	result := govMapper.GetProposalsFiltered(ctx, params.Voter, params.Depositor, params.Status, params.Limit)

	return govMapper.GetCodec().MarshalJSON(result)
}

func queryVotesByProposalID(ctx context.Context, pid string) ([]byte, error) {
	govMapper := GetMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exists := govMapper.GetProposal(proposalID)
	if !exists {
		return nil, fmt.Errorf("proposal id %d not exists", proposalID)
	}

	iter := govMapper.GetVotes(proposalID)
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
	govMapper := GetMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exists := govMapper.GetProposal(proposalID)
	if !exists {
		return nil, fmt.Errorf("proposal id %d not exists", proposalID)
	}

	voterAddress, err := btypes.GetAddrFromBech32(voter)
	if err != nil {
		return nil, fmt.Errorf("voter %s is not valid address", voter)
	}

	result, exists := govMapper.GetVote(proposalID, voterAddress)
	if !exists {
		return nil, fmt.Errorf("voter %s is not vote on proposal %s", voterAddress, pid)
	}

	return govMapper.GetCodec().MarshalJSON(result)
}

func queryDepositsByProposalIDAndDepositer(ctx context.Context, pid string, depositer string) ([]byte, error) {
	govMapper := GetMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exists := govMapper.GetProposal(proposalID)
	if !exists {
		return nil, fmt.Errorf("proposal id %d not exists", proposalID)
	}

	depositerAddress, err := btypes.GetAddrFromBech32(depositer)
	if err != nil {
		return nil, fmt.Errorf("depositer %s is not valid address", depositer)
	}

	result, exists := govMapper.GetDeposit(proposalID, depositerAddress)
	if !exists {
		return nil, fmt.Errorf("depositer %s is not deposit on proposal %s", depositer, pid)
	}
	return govMapper.GetCodec().MarshalJSON(result)
}

func queryDepositsByProposalID(ctx context.Context, pid string) ([]byte, error) {
	govMapper := GetMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	_, exists := govMapper.GetProposal(proposalID)
	if !exists {
		return nil, fmt.Errorf("proposal id %d not exists", proposalID)
	}

	iter := govMapper.GetDeposits(proposalID)
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
	govMapper := GetMapper(ctx)

	proposalID, err := strconv.ParseUint(pid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("pid %s is not a valid uint value", pid)
	}

	proposal, exists := govMapper.GetProposal(proposalID)
	if !exists {
		return nil, fmt.Errorf("proposal id %d not exists", proposalID)
	}

	_, result, _, _ := Tally(ctx, govMapper, proposal)
	return govMapper.GetCodec().MarshalJSON(result)
}

func queryParams(ctx context.Context) ([]byte, error) {
	paramMapper := mapper.GetMapper(ctx)
	return paramMapper.GetCodec().MarshalJSON(paramMapper.GetParams())
}

func queryModuleParams(ctx context.Context, module string) ([]byte, error) {
	paramMapper := mapper.GetMapper(ctx)
	result, exists := paramMapper.GetModuleParams(module)
	if !exists {
		return nil, errors.New(fmt.Sprintf("no parameter in module %s", module))
	}
	return paramMapper.GetCodec().MarshalJSON(result)
}

func queryParam(ctx context.Context, module string, key string) ([]byte, error) {
	paramMapper := mapper.GetMapper(ctx)
	result, exists := paramMapper.GetParam(module, key)
	if !exists {
		return nil, errors.New(fmt.Sprintf("no %s in module %s", key, module))
	}

	return paramMapper.GetCodec().MarshalJSON(result)
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
	return fmt.Sprintf("custom/%s/%s/%d", GOV, TallyQuery, pid)
}

//nolint
func BuildQueryParamsPath() string {
	return fmt.Sprintf("custom/%s/%s", GOV, ParamsPath)
}

//nolint
func BuildQueryModuleParamsPath(module string) string {
	return fmt.Sprintf("custom/%s/%s/%s", GOV, ParamsPath, module)
}

//nolint
func BuildQueryParamPath(module string, key string) string {
	return fmt.Sprintf("custom/%s/%s/%s/%s", GOV, ParamsPath, module, key)
}
