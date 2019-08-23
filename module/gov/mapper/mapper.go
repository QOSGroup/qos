package mapper

import (
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/stake"
	"time"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

const (
	MapperName = "governance"
)

type Mapper struct {
	*mapper.BaseMapper
}

func (mapper *Mapper) Copy() mapper.IMapper {
	govMapper := &Mapper{}
	govMapper.BaseMapper = mapper.BaseMapper.Copy()
	return govMapper
}

var _ mapper.IMapper = (*Mapper)(nil)

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(MapperName).(*Mapper)
}

func NewMapper() *Mapper {
	var govMapper = Mapper{}
	govMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &govMapper
}

// Submit proposal
func (mapper Mapper) SubmitProposal(ctx context.Context, content types.ProposalContent) (proposal types.Proposal, err btypes.Error) {
	proposalID, err := mapper.getNewProposalID()
	if err != nil {
		return
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := mapper.GetParams(ctx).MaxDepositPeriod

	proposal = types.Proposal{
		ProposalContent: content,
		ProposalID:      proposalID,

		Status:           types.StatusDepositPeriod,
		FinalTallyResult: types.EmptyTallyResult(),
		TotalDeposit:     0,
		SubmitTime:       submitTime,
		DepositEndTime:   submitTime.Add(depositPeriod),
	}

	mapper.SetProposal(proposal)
	mapper.InsertInactiveProposalQueue(proposal.DepositEndTime, proposalID)
	return
}

// Get Proposal from store by ProposalID
func (mapper Mapper) GetProposal(proposalID uint64) (proposal types.Proposal, ok bool) {
	ok = mapper.Get(types.KeyProposal(proposalID), &proposal)
	return
}

// Update proposal
func (mapper Mapper) SetProposal(proposal types.Proposal) {
	mapper.Set(types.KeyProposal(proposal.ProposalID), proposal)
}

// Delete proposal
func (mapper Mapper) DeleteProposal(proposalID uint64) {
	proposal, ok := mapper.GetProposal(proposalID)
	if !ok {
		panic("DeleteProposal cannot fail to GetProposal")
	}
	mapper.RemoveFromInactiveProposalQueue(proposal.DepositEndTime, proposalID)
	mapper.RemoveFromActiveProposalQueue(proposal.VotingEndTime, proposalID)
	mapper.Del(types.KeyProposal(proposalID))
}

// Get Proposal from store by ProposalID
// voterAddr will filter proposals by whether or not that address has voted on them
// depositorAddr will filter proposals by whether or not that address has deposited to them
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (mapper Mapper) GetProposalsFiltered(ctx context.Context, voterAddr btypes.Address, depositorAddr btypes.Address, status types.ProposalStatus, numLatest uint64) []types.Proposal {

	maxProposalID, err := mapper.PeekCurrentProposalID()
	if err != nil {
		return nil
	}

	matchingProposals := []types.Proposal{}

	if numLatest == 0 {
		numLatest = maxProposalID
	}

	for proposalID := maxProposalID - numLatest; proposalID < maxProposalID; proposalID++ {
		if voterAddr != nil && len(voterAddr) != 0 {
			_, found := mapper.GetVote(proposalID, voterAddr)
			if !found {
				continue
			}
		}

		if depositorAddr != nil && len(depositorAddr) != 0 {
			_, found := mapper.GetDeposit(proposalID, depositorAddr)
			if !found {
				continue
			}
		}

		proposal, ok := mapper.GetProposal(proposalID)
		if !ok {
			continue
		}

		if types.ValidProposalStatus(status) {
			if proposal.Status != status {
				continue
			}
		}

		matchingProposals = append(matchingProposals, proposal)
	}
	return matchingProposals
}

func (mapper Mapper) GetProposals() []types.Proposal {

	var proposals []types.Proposal
	iterator := btypes.KVStorePrefixIterator(mapper.GetStore(), types.KeyProposalSubspace())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposal := types.Proposal{}
		mapper.DecodeObject(iterator.Value(), &proposal)
		proposals = append(proposals, proposal)
	}
	return proposals
}

// Set the initial proposal ID
func (mapper Mapper) SetInitialProposalID(proposalID uint64) btypes.Error {
	exists := mapper.Get(types.KeyNextProposalID, &proposalID)
	if exists {
		return types.ErrInvalidGenesis("Initial ProposalID already set")
	}
	mapper.Set(types.KeyNextProposalID, proposalID)
	return nil
}

// Get the last used proposal ID
func (mapper Mapper) GetLastProposalID() (proposalID uint64) {
	proposalID, err := mapper.PeekCurrentProposalID()
	if err != nil {
		return 0
	}
	proposalID--
	return
}

// Gets the next available ProposalID and increments it
func (mapper Mapper) getNewProposalID() (proposalID uint64, err btypes.Error) {
	exists := mapper.Get(types.KeyNextProposalID, &proposalID)
	if !exists {
		return 0, types.ErrInvalidGenesis("InitialProposalID never set")
	}
	mapper.Set(types.KeyNextProposalID, proposalID+1)
	return proposalID, nil
}

// Peeks the next available ProposalID without incrementing it
func (mapper Mapper) PeekCurrentProposalID() (proposalID uint64, err btypes.Error) {
	exists := mapper.Get(types.KeyNextProposalID, &proposalID)
	if !exists {
		return 0, types.ErrInvalidGenesis("InitialProposalID never set")
	}
	return proposalID, nil
}

func (mapper Mapper) activateVotingPeriod(ctx context.Context, proposal types.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	proposal.VotingStartHeight = uint64(ctx.BlockHeight())
	votingPeriod := mapper.GetParams(ctx).VotingPeriod
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	proposal.Status = types.StatusVotingPeriod
	mapper.SetProposal(proposal)

	mapper.RemoveFromInactiveProposalQueue(proposal.DepositEndTime, proposal.ProposalID)
	mapper.InsertActiveProposalQueue(proposal.VotingEndTime, proposal.ProposalID)

	mapper.saveValidatorSet(ctx, proposal.ProposalID)
}

// Save validator set when proposal entering voting period.
func (mapper Mapper) saveValidatorSet(ctx context.Context, proposalID uint64) {
	validators := stake.GetMapper(ctx).GetActiveValidatorSet(false)
	if validators != nil {
		mapper.Set(types.KeyVotingPeriodValidators(proposalID), validators)
	}
}

func (mapper Mapper) GetValidatorSet(proposalID uint64) (validators []btypes.Address, exists bool) {
	exists = mapper.Get(types.KeyVotingPeriodValidators(proposalID), &validators)
	return
}

func (mapper Mapper) DeleteValidatorSet(proposalID uint64) {
	mapper.Del(types.KeyVotingPeriodValidators(proposalID))
}

// Params

func (mapper Mapper) GetParams(ctx context.Context) types.Params {
	var p types.Params
	params.GetMapper(ctx).GetParamSet(&p)
	return p
}

func (mapper Mapper) SetParams(ctx context.Context, p types.Params) {
	params.GetMapper(ctx).SetParamSet(&p)
}

// Votes

// Adds a vote on a specific proposal
func (mapper Mapper) AddVote(proposalID uint64, voterAddr btypes.Address, option types.VoteOption) btypes.Error {
	vote := types.Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}
	mapper.SetVote(proposalID, voterAddr, vote)

	return nil
}

// Gets the vote of a specific voter on a specific proposal
func (mapper Mapper) GetVote(proposalID uint64, voterAddr btypes.Address) (vote types.Vote, exists bool) {
	exists = mapper.Get(types.KeyVote(proposalID, voterAddr), &vote)
	if !exists {
		return types.Vote{}, false
	}
	return
}

func (mapper Mapper) SetVote(proposalID uint64, voterAddr btypes.Address, vote types.Vote) {
	mapper.Set(types.KeyVote(proposalID, voterAddr), vote)
}

// Gets all the votes on a specific proposal
func (mapper Mapper) GetVotes(proposalID uint64) store.Iterator {
	return btypes.KVStorePrefixIterator(mapper.GetStore(), types.KeyVotesSubspace(proposalID))
}

func (mapper Mapper) deleteVote(proposalID uint64, voterAddr btypes.Address) {
	mapper.Del(types.KeyVote(proposalID, voterAddr))
}

// Delete votes
func (mapper Mapper) DeleteVotes(proposalID uint64) {
	iterator := mapper.GetVotes(proposalID)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		mapper.Del(iterator.Key())
	}
}

// Deposits

// Gets the deposit of a specific depositor on a specific proposal
func (mapper Mapper) GetDeposit(proposalID uint64, depositorAddr btypes.Address) (deposit types.Deposit, exists bool) {
	exists = mapper.Get(types.KeyDeposit(proposalID, depositorAddr), &deposit)
	if !exists {
		return types.Deposit{}, false
	}

	return deposit, true
}

func (mapper Mapper) SetDeposit(proposalID uint64, depositorAddr btypes.Address, deposit types.Deposit) {
	mapper.Set(types.KeyDeposit(proposalID, depositorAddr), deposit)
}

// Adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (mapper Mapper) AddDeposit(ctx context.Context, proposalID uint64, depositorAddr btypes.Address, depositAmount uint64) (btypes.Error, bool) {
	proposal, ok := mapper.GetProposal(proposalID)
	if !ok {
		return types.ErrUnknownProposal(proposalID), false
	}

	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	account := accountMapper.GetAccount(depositorAddr).(*qtypes.QOSAccount)
	account.MustMinusQOS(btypes.NewInt(int64(depositAmount)))
	accountMapper.SetAccount(account)

	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit + depositAmount
	mapper.SetProposal(proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == types.StatusDepositPeriod && proposal.TotalDeposit >= mapper.GetParams(ctx).MinDeposit {
		mapper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := mapper.GetDeposit(proposalID, depositorAddr)
	if !found {
		newDeposit := types.Deposit{depositorAddr, proposalID, depositAmount}
		mapper.SetDeposit(proposalID, depositorAddr, newDeposit)
	} else {
		currDeposit.Amount = currDeposit.Amount + depositAmount
		mapper.SetDeposit(proposalID, depositorAddr, currDeposit)
	}

	return nil, activatedVotingPeriod
}

// Gets all the deposits on a specific proposal as an sdk.Iterator
func (mapper Mapper) GetDeposits(proposalID uint64) store.Iterator {
	return btypes.KVStorePrefixIterator(mapper.GetStore(), types.KeyDepositsSubspace(proposalID))
}

// Refunds and deletes all the deposits on a specific proposal
func (mapper Mapper) RefundDeposits(ctx context.Context, proposalID uint64, deductDeposit bool) {

	log := ctx.Logger()
	params := mapper.GetParams(ctx)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	depositsIterator := mapper.GetDeposits(proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &types.Deposit{}
		mapper.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), deposit)

		depositAmount := int64(deposit.Amount)

		//需要扣除部分押金时
		burnAmount := int64(0)
		if deductDeposit {
			burnAmount = params.BurnRate.Mul(qtypes.NewDec(depositAmount)).TruncateInt64()
		}

		refundAmount := depositAmount - burnAmount

		// refund deposit
		depositor := accountMapper.GetAccount(deposit.Depositor).(*qtypes.QOSAccount)
		depositor.PlusQOS(btypes.NewInt(refundAmount))
		accountMapper.SetAccount(depositor)

		// burn deposit
		if deductDeposit {
			distribution.GetMapper(ctx).AddToCommunityFeePool(btypes.NewInt(burnAmount))
		}

		log.Debug("RefundDeposits", "depositAmount", depositAmount, "refundAmount", refundAmount, "burnAmount", burnAmount)

		mapper.Del(depositsIterator.Key())
	}
}

// Deletes all the deposits on a specific proposal without refunding them
func (mapper Mapper) DeleteDeposits(ctx context.Context, proposalID uint64) {
	depositsIterator := mapper.GetDeposits(proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &types.Deposit{}
		mapper.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), deposit)

		// burn deposit
		distribution.GetMapper(ctx).AddToCommunityFeePool(btypes.NewInt(int64(deposit.Amount)))

		mapper.Del(depositsIterator.Key())
	}
}

// ProposalQueues

// Returns an iterator for all the proposals in the Active Queue that expire by endTime
func (mapper Mapper) ActiveProposalQueueIterator(endTime time.Time) store.Iterator {
	return mapper.GetStore().Iterator(types.PrefixActiveProposalQueue, btypes.PrefixEndBytes(types.PrefixActiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (mapper Mapper) InsertActiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Set(types.KeyActiveProposalQueueProposal(endTime, proposalID), proposalID)
}

// removes a proposalID from the Active Proposal Queue
func (mapper Mapper) RemoveFromActiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Del(types.KeyActiveProposalQueueProposal(endTime, proposalID))
}

// Returns an iterator for all the proposals in the Inactive Queue that expire by endTime
func (mapper Mapper) InactiveProposalQueueIterator(endTime time.Time) store.Iterator {
	return mapper.GetStore().Iterator(types.PrefixInactiveProposalQueue, btypes.PrefixEndBytes(types.PrefixInactiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the inactive proposal queue at endTime
func (mapper Mapper) InsertInactiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Set(types.KeyInactiveProposalQueueProposal(endTime, proposalID), proposalID)
}

// removes a proposalID from the Inactive Proposal Queue
func (mapper Mapper) RemoveFromInactiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Del(types.KeyInactiveProposalQueueProposal(endTime, proposalID))
}
