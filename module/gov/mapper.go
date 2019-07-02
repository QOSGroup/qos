package gov

import (
	"time"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	pmapper "github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
)

const (
	GovMapperName = "governance"
)

type GovMapper struct {
	*mapper.BaseMapper
}

func (mapper *GovMapper) Copy() mapper.IMapper {
	govMapper := &GovMapper{}
	govMapper.BaseMapper = mapper.BaseMapper.Copy()
	return govMapper
}

var _ mapper.IMapper = (*GovMapper)(nil)

func GetGovMapper(ctx context.Context) *GovMapper {
	return ctx.Mapper(GovMapperName).(*GovMapper)
}

func NewGovMapper() *GovMapper {
	var govMapper = GovMapper{}
	govMapper.BaseMapper = mapper.NewBaseMapper(nil, GovMapperName)
	return &govMapper
}

// Submit proposal
func (mapper GovMapper) SubmitProposal(ctx context.Context, content gtypes.ProposalContent) (proposal gtypes.Proposal, err btypes.Error) {
	proposalID, err := mapper.getNewProposalID()
	if err != nil {
		return
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := mapper.GetParams(ctx).MaxDepositPeriod

	proposal = gtypes.Proposal{
		ProposalContent: content,
		ProposalID:      proposalID,

		Status:           gtypes.StatusDepositPeriod,
		FinalTallyResult: gtypes.EmptyTallyResult(),
		TotalDeposit:     0,
		SubmitTime:       submitTime,
		DepositEndTime:   submitTime.Add(depositPeriod),
	}

	mapper.SetProposal(proposal)
	mapper.InsertInactiveProposalQueue(proposal.DepositEndTime, proposalID)
	return
}

// Get Proposal from store by ProposalID
func (mapper GovMapper) GetProposal(proposalID uint64) (proposal gtypes.Proposal, ok bool) {
	ok = mapper.Get(KeyProposal(proposalID), &proposal)
	return
}

// Update proposal
func (mapper GovMapper) SetProposal(proposal gtypes.Proposal) {
	mapper.Set(KeyProposal(proposal.ProposalID), proposal)
}

// Delete proposal
func (mapper GovMapper) DeleteProposal(proposalID uint64) {
	proposal, ok := mapper.GetProposal(proposalID)
	if !ok {
		panic("DeleteProposal cannot fail to GetProposal")
	}
	mapper.RemoveFromInactiveProposalQueue(proposal.DepositEndTime, proposalID)
	mapper.RemoveFromActiveProposalQueue(proposal.VotingEndTime, proposalID)
	mapper.Del(KeyProposal(proposalID))
}

// Get Proposal from store by ProposalID
// voterAddr will filter proposals by whether or not that address has voted on them
// depositorAddr will filter proposals by whether or not that address has deposited to them
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (mapper GovMapper) GetProposalsFiltered(ctx context.Context, voterAddr btypes.Address, depositorAddr btypes.Address, status gtypes.ProposalStatus, numLatest uint64) []gtypes.Proposal {

	maxProposalID, err := mapper.peekCurrentProposalID()
	if err != nil {
		return nil
	}

	matchingProposals := []gtypes.Proposal{}

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

		if gtypes.ValidProposalStatus(status) {
			if proposal.Status != status {
				continue
			}
		}

		matchingProposals = append(matchingProposals, proposal)
	}
	return matchingProposals
}

func (mapper GovMapper) GetProposals() []gtypes.Proposal {

	var proposals []gtypes.Proposal
	iterator := btypes.KVStorePrefixIterator(mapper.GetStore(), KeyProposalSubspace())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposal := gtypes.Proposal{}
		mapper.DecodeObject(iterator.Value(), &proposal)
		proposals = append(proposals, proposal)
	}
	return proposals
}

// Set the initial proposal ID
func (mapper GovMapper) setInitialProposalID(proposalID uint64) btypes.Error {
	exists := mapper.Get(KeyNextProposalID, &proposalID)
	if exists {
		return ErrInvalidGenesis("Initial ProposalID already set")
	}
	mapper.Set(KeyNextProposalID, proposalID)
	return nil
}

// Get the last used proposal ID
func (mapper GovMapper) GetLastProposalID() (proposalID uint64) {
	proposalID, err := mapper.peekCurrentProposalID()
	if err != nil {
		return 0
	}
	proposalID--
	return
}

// Gets the next available ProposalID and increments it
func (mapper GovMapper) getNewProposalID() (proposalID uint64, err btypes.Error) {
	exists := mapper.Get(KeyNextProposalID, &proposalID)
	if !exists {
		return 0, ErrInvalidGenesis("InitialProposalID never set")
	}
	mapper.Set(KeyNextProposalID, proposalID+1)
	return proposalID, nil
}

// Peeks the next available ProposalID without incrementing it
func (mapper GovMapper) peekCurrentProposalID() (proposalID uint64, err btypes.Error) {
	exists := mapper.Get(KeyNextProposalID, &proposalID)
	if !exists {
		return 0, ErrInvalidGenesis("InitialProposalID never set")
	}
	return proposalID, nil
}

func (mapper GovMapper) activateVotingPeriod(ctx context.Context, proposal gtypes.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	proposal.VotingStartHeight = uint64(ctx.BlockHeight())
	votingPeriod := mapper.GetParams(ctx).VotingPeriod
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	proposal.Status = gtypes.StatusVotingPeriod
	mapper.SetProposal(proposal)

	mapper.RemoveFromInactiveProposalQueue(proposal.DepositEndTime, proposal.ProposalID)
	mapper.InsertActiveProposalQueue(proposal.VotingEndTime, proposal.ProposalID)

	mapper.saveValidatorSet(ctx, proposal.ProposalID)
}

// Save validator set when proposal entering voting period.
func (mapper GovMapper) saveValidatorSet(ctx context.Context, proposalID uint64) {
	validators := ecomapper.GetValidatorMapper(ctx).GetActiveValidatorSet(false)
	if validators != nil {
		mapper.Set(KeyVotingPeriodValidators(proposalID), validators)
	}
}

func (mapper GovMapper) GetValidatorSet(proposalID uint64) (validators []btypes.Address, exists bool) {
	exists = mapper.Get(KeyVotingPeriodValidators(proposalID), &validators)
	return
}

func (mapper GovMapper) DeleteValidatorSet(proposalID uint64) {
	mapper.Del(KeyVotingPeriodValidators(proposalID))
}

// Params

func (mapper GovMapper) GetParams(ctx context.Context) Params {
	var params Params
	pmapper.GetMapper(ctx).GetParamSet(&params)
	return params
}

func (mapper GovMapper) SetParams(ctx context.Context, params Params) {
	pmapper.GetMapper(ctx).SetParamSet(&params)
}

// Votes

// Adds a vote on a specific proposal
func (mapper GovMapper) AddVote(proposalID uint64, voterAddr btypes.Address, option gtypes.VoteOption) btypes.Error {
	vote := gtypes.Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}
	mapper.setVote(proposalID, voterAddr, vote)

	return nil
}

// Gets the vote of a specific voter on a specific proposal
func (mapper GovMapper) GetVote(proposalID uint64, voterAddr btypes.Address) (vote gtypes.Vote, exists bool) {
	exists = mapper.Get(KeyVote(proposalID, voterAddr), &vote)
	if !exists {
		return gtypes.Vote{}, false
	}
	return
}

func (mapper GovMapper) setVote(proposalID uint64, voterAddr btypes.Address, vote gtypes.Vote) {
	mapper.Set(KeyVote(proposalID, voterAddr), vote)
}

// Gets all the votes on a specific proposal
func (mapper GovMapper) GetVotes(proposalID uint64) store.Iterator {
	return btypes.KVStorePrefixIterator(mapper.GetStore(), KeyVotesSubspace(proposalID))
}

func (mapper GovMapper) deleteVote(proposalID uint64, voterAddr btypes.Address) {
	mapper.Del(KeyVote(proposalID, voterAddr))
}

// Delete votes
func (mapper GovMapper) DeleteVotes(proposalID uint64) {
	iterator := mapper.GetVotes(proposalID)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		mapper.Del(iterator.Key())
	}
}

// Deposits

// Gets the deposit of a specific depositor on a specific proposal
func (mapper GovMapper) GetDeposit(proposalID uint64, depositorAddr btypes.Address) (deposit gtypes.Deposit, exists bool) {
	exists = mapper.Get(KeyDeposit(proposalID, depositorAddr), &deposit)
	if !exists {
		return gtypes.Deposit{}, false
	}

	return deposit, true
}

func (mapper GovMapper) setDeposit(proposalID uint64, depositorAddr btypes.Address, deposit gtypes.Deposit) {
	mapper.Set(KeyDeposit(proposalID, depositorAddr), deposit)
}

// Adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (mapper GovMapper) AddDeposit(ctx context.Context, proposalID uint64, depositorAddr btypes.Address, depositAmount uint64) (btypes.Error, bool) {
	proposal, ok := mapper.GetProposal(proposalID)
	if !ok {
		return ErrUnknownProposal(proposalID), false
	}

	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	account := accountMapper.GetAccount(depositorAddr).(*types.QOSAccount)
	account.MustMinusQOS(btypes.NewInt(int64(depositAmount)))
	accountMapper.SetAccount(account)

	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit + depositAmount
	mapper.SetProposal(proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == gtypes.StatusDepositPeriod && proposal.TotalDeposit >= mapper.GetParams(ctx).MinDeposit {
		mapper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := mapper.GetDeposit(proposalID, depositorAddr)
	if !found {
		newDeposit := gtypes.Deposit{depositorAddr, proposalID, depositAmount}
		mapper.setDeposit(proposalID, depositorAddr, newDeposit)
	} else {
		currDeposit.Amount = currDeposit.Amount + depositAmount
		mapper.setDeposit(proposalID, depositorAddr, currDeposit)
	}

	return nil, activatedVotingPeriod
}

// Gets all the deposits on a specific proposal as an sdk.Iterator
func (mapper GovMapper) GetDeposits(proposalID uint64) store.Iterator {
	return btypes.KVStorePrefixIterator(mapper.GetStore(), KeyDepositsSubspace(proposalID))
}

// Refunds and deletes all the deposits on a specific proposal
func (mapper GovMapper) RefundDeposits(ctx context.Context, proposalID uint64, burnDeposit bool) {

	log := ctx.Logger()
	params := mapper.GetParams(ctx)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	depositsIterator := mapper.GetDeposits(proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &gtypes.Deposit{}
		mapper.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), deposit)

		depositAmount := int64(deposit.Amount)

		//需要扣除部分押金时
		burnAmount := int64(0)
		if burnDeposit {
			burnAmount = params.BurnRate.Mul(types.NewDec(depositAmount)).TruncateInt64()
		}

		refundAmount := depositAmount - burnAmount

		// refund deposit
		depositor := accountMapper.GetAccount(deposit.Depositor).(*types.QOSAccount)
		depositor.PlusQOS(btypes.NewInt(refundAmount))
		accountMapper.SetAccount(depositor)

		// burn deposit
		if burnDeposit {
			ecomapper.GetDistributionMapper(ctx).AddToCommunityFeePool(btypes.NewInt(burnAmount))
		}

		log.Debug("RefundDeposits", "depositAmount", depositAmount, "refundAmount", refundAmount, "burnAmount", burnAmount)

		mapper.Del(depositsIterator.Key())
	}
}

// Deletes all the deposits on a specific proposal without refunding them
func (mapper GovMapper) DeleteDeposits(ctx context.Context, proposalID uint64) {
	depositsIterator := mapper.GetDeposits(proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &gtypes.Deposit{}
		mapper.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), deposit)

		// burn deposit
		ecomapper.GetDistributionMapper(ctx).AddToCommunityFeePool(btypes.NewInt(int64(deposit.Amount)))

		mapper.Del(depositsIterator.Key())
	}
}

// ProposalQueues

// Returns an iterator for all the proposals in the Active Queue that expire by endTime
func (mapper GovMapper) ActiveProposalQueueIterator(endTime time.Time) store.Iterator {
	return mapper.GetStore().Iterator(PrefixActiveProposalQueue, btypes.PrefixEndBytes(PrefixActiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (mapper GovMapper) InsertActiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Set(KeyActiveProposalQueueProposal(endTime, proposalID), proposalID)
}

// removes a proposalID from the Active Proposal Queue
func (mapper GovMapper) RemoveFromActiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Del(KeyActiveProposalQueueProposal(endTime, proposalID))
}

// Returns an iterator for all the proposals in the Inactive Queue that expire by endTime
func (mapper GovMapper) InactiveProposalQueueIterator(endTime time.Time) store.Iterator {
	return mapper.GetStore().Iterator(PrefixInactiveProposalQueue, btypes.PrefixEndBytes(PrefixInactiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the inactive proposal queue at endTime
func (mapper GovMapper) InsertInactiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Set(KeyInactiveProposalQueueProposal(endTime, proposalID), proposalID)
}

// removes a proposalID from the Inactive Proposal Queue
func (mapper GovMapper) RemoveFromInactiveProposalQueue(endTime time.Time, proposalID uint64) {
	mapper.Del(KeyInactiveProposalQueueProposal(endTime, proposalID))
}
