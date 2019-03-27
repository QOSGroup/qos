package gov

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
	"time"
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
	proposalID, err := mapper.getNewProposalID(ctx)
	if err != nil {
		return
	}

	submitTime := ctx.BlockHeader().Time
	depositPeriod := mapper.GetDepositParams().MaxDepositPeriod

	proposal = gtypes.Proposal{
		ProposalContent: content,
		ProposalID:      proposalID,

		Status:           gtypes.StatusDepositPeriod,
		FinalTallyResult: gtypes.EmptyTallyResult(),
		TotalDeposit:     content.GetDeposit(),
		SubmitTime:       submitTime,
		DepositEndTime:   submitTime.Add(depositPeriod),
	}

	mapper.SetProposal(ctx, proposal)
	mapper.InsertInactiveProposalQueue(ctx, proposal.DepositEndTime, proposalID)
	return
}

// Get Proposal from store by ProposalID
func (mapper GovMapper) GetProposal(ctx context.Context, proposalID uint64) (proposal gtypes.Proposal, ok bool) {
	ok = mapper.Get(KeyProposal(proposalID), &proposal)
	return
}

// Update proposal
func (mapper GovMapper) SetProposal(ctx context.Context, proposal gtypes.Proposal) {
	mapper.Set(KeyProposal(proposal.ProposalID), proposal)
}

// Delete proposal
func (mapper GovMapper) DeleteProposal(ctx context.Context, proposalID uint64) {
	proposal, ok := mapper.GetProposal(ctx, proposalID)
	if !ok {
		panic("DeleteProposal cannot fail to GetProposal")
	}
	mapper.RemoveFromInactiveProposalQueue(ctx, proposal.DepositEndTime, proposalID)
	mapper.RemoveFromActiveProposalQueue(ctx, proposal.VotingEndTime, proposalID)
	mapper.Del(KeyProposal(proposalID))
}

// Get Proposal from store by ProposalID
// voterAddr will filter proposals by whether or not that address has voted on them
// depositorAddr will filter proposals by whether or not that address has deposited to them
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (mapper GovMapper) GetProposalsFiltered(ctx context.Context, voterAddr btypes.Address, depositorAddr btypes.Address, status gtypes.ProposalStatus, numLatest uint64) []gtypes.Proposal {

	maxProposalID, err := mapper.peekCurrentProposalID(ctx)
	if err != nil {
		return nil
	}

	matchingProposals := []gtypes.Proposal{}

	if numLatest == 0 {
		numLatest = maxProposalID
	}

	for proposalID := maxProposalID - numLatest; proposalID < maxProposalID; proposalID++ {
		if voterAddr != nil && len(voterAddr) != 0 {
			_, found := mapper.GetVote(ctx, proposalID, voterAddr)
			if !found {
				continue
			}
		}

		if depositorAddr != nil && len(depositorAddr) != 0 {
			_, found := mapper.GetDeposit(ctx, proposalID, depositorAddr)
			if !found {
				continue
			}
		}

		proposal, ok := mapper.GetProposal(ctx, proposalID)
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

// Set the initial proposal ID
func (mapper GovMapper) setInitialProposalID(ctx context.Context, proposalID uint64) btypes.Error {
	exists := mapper.Get(KeyNextProposalID, &proposalID)
	if exists {
		return ErrInvalidGenesis("Initial ProposalID already set")
	}
	mapper.Set(KeyNextProposalID, proposalID)
	return nil
}

// Get the last used proposal ID
func (mapper GovMapper) GetLastProposalID(ctx context.Context) (proposalID uint64) {
	proposalID, err := mapper.peekCurrentProposalID(ctx)
	if err != nil {
		return 0
	}
	proposalID--
	return
}

// Gets the next available ProposalID and increments it
func (mapper GovMapper) getNewProposalID(ctx context.Context) (proposalID uint64, err btypes.Error) {
	exists := mapper.Get(KeyNextProposalID, &proposalID)
	if !exists {
		return 0, ErrInvalidGenesis("InitialProposalID never set")
	}
	mapper.Set(KeyNextProposalID, proposalID+1)
	return proposalID, nil
}

// Peeks the next available ProposalID without incrementing it
func (mapper GovMapper) peekCurrentProposalID(ctx context.Context) (proposalID uint64, err btypes.Error) {
	exists := mapper.Get(KeyNextProposalID, &proposalID)
	if !exists {
		return 0, ErrInvalidGenesis("InitialProposalID never set")
	}
	return proposalID, nil
}

func (mapper GovMapper) activateVotingPeriod(ctx context.Context, proposal gtypes.Proposal) {
	proposal.VotingStartTime = ctx.BlockHeader().Time
	votingPeriod := mapper.GetVotingParams().VotingPeriod
	proposal.VotingEndTime = proposal.VotingStartTime.Add(votingPeriod)
	proposal.Status = gtypes.StatusVotingPeriod
	mapper.SetProposal(ctx, proposal)

	mapper.RemoveFromInactiveProposalQueue(ctx, proposal.DepositEndTime, proposal.ProposalID)
	mapper.InsertActiveProposalQueue(ctx, proposal.VotingEndTime, proposal.ProposalID)
}

// Params

// Returns the current DepositParams from the global param store
func (mapper GovMapper) GetDepositParams() DepositParams {
	var depositParams DepositParams
	mapper.Get(ParamStoreKeyDepositParams, &depositParams)
	return depositParams
}

// Returns the current VotingParams from the global param store
func (mapper GovMapper) GetVotingParams() VotingParams {
	var votingParams VotingParams
	mapper.Get(ParamStoreKeyVotingParams, &votingParams)
	return votingParams
}

// Returns the current TallyParam from the global param store
func (mapper GovMapper) GetTallyParams() TallyParams {
	var tallyParams TallyParams
	mapper.Get(ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

func (mapper GovMapper) setDepositParams(depositParams DepositParams) {
	mapper.Set(ParamStoreKeyDepositParams, &depositParams)
}

func (mapper GovMapper) setVotingParams(votingParams VotingParams) {
	mapper.Set(ParamStoreKeyVotingParams, &votingParams)
}

func (mapper GovMapper) setTallyParams(tallyParams TallyParams) {
	mapper.Set(ParamStoreKeyTallyParams, &tallyParams)
}

// Votes

// Adds a vote on a specific proposal
func (mapper GovMapper) AddVote(ctx context.Context, proposalID uint64, voterAddr btypes.Address, option gtypes.VoteOption) btypes.Error {
	vote := gtypes.Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}
	mapper.setVote(ctx, proposalID, voterAddr, vote)

	return nil
}

// Gets the vote of a specific voter on a specific proposal
func (mapper GovMapper) GetVote(ctx context.Context, proposalID uint64, voterAddr btypes.Address) (vote gtypes.Vote, exists bool) {
	exists = mapper.Get(KeyVote(proposalID, voterAddr), &vote)
	if !exists {
		return gtypes.Vote{}, false
	}
	return
}

func (mapper GovMapper) setVote(ctx context.Context, proposalID uint64, voterAddr btypes.Address, vote gtypes.Vote) {
	mapper.Set(KeyVote(proposalID, voterAddr), vote)
}

// Gets all the votes on a specific proposal
func (mapper GovMapper) GetVotes(ctx context.Context, proposalID uint64) store.Iterator {
	return store.KVStorePrefixIterator(mapper.GetStore(), KeyVotesSubspace(proposalID))
}

func (mapper GovMapper) deleteVote(ctx context.Context, proposalID uint64, voterAddr btypes.Address) {
	mapper.Del(KeyVote(proposalID, voterAddr))
}

// Deposits

// Gets the deposit of a specific depositor on a specific proposal
func (mapper GovMapper) GetDeposit(ctx context.Context, proposalID uint64, depositorAddr btypes.Address) (deposit gtypes.Deposit, exists bool) {
	exists = mapper.Get(KeyDeposit(proposalID, depositorAddr), &deposit)
	if !exists {
		return gtypes.Deposit{}, false
	}

	return deposit, true
}

func (mapper GovMapper) setDeposit(ctx context.Context, proposalID uint64, depositorAddr btypes.Address, deposit gtypes.Deposit) {
	mapper.Set(KeyDeposit(proposalID, depositorAddr), deposit)
}

// Adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (mapper GovMapper) AddDeposit(ctx context.Context, proposalID uint64, depositorAddr btypes.Address, depositAmount uint64) (btypes.Error, bool) {
	proposal, ok := mapper.GetProposal(ctx, proposalID)
	if !ok {
		return ErrUnknownProposal(proposalID), false
	}

	// add gov deposit
	mapper.addGovDeposit(depositAmount)

	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	account := accountMapper.GetAccount(depositorAddr).(*types.QOSAccount)
	account.MustMinusQOS(btypes.NewInt(int64(depositAmount)))
	accountMapper.SetAccount(account)

	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit + depositAmount
	mapper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == gtypes.StatusDepositPeriod && proposal.TotalDeposit >= mapper.GetDepositParams().MinDeposit {
		mapper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := mapper.GetDeposit(ctx, proposalID, depositorAddr)
	if !found {
		newDeposit := gtypes.Deposit{depositorAddr, proposalID, depositAmount}
		mapper.setDeposit(ctx, proposalID, depositorAddr, newDeposit)
	} else {
		currDeposit.Amount = currDeposit.Amount + depositAmount
		mapper.setDeposit(ctx, proposalID, depositorAddr, currDeposit)
	}

	return nil, activatedVotingPeriod
}

// Gets all the deposits on a specific proposal as an sdk.Iterator
func (mapper GovMapper) GetDeposits(ctx context.Context, proposalID uint64) store.Iterator {
	return store.KVStorePrefixIterator(mapper.GetStore(), KeyDepositsSubspace(proposalID))
}

// Refunds and deletes all the deposits on a specific proposal
func (mapper GovMapper) RefundDeposits(ctx context.Context, proposalID uint64) {
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	depositsIterator := mapper.GetDeposits(ctx, proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &gtypes.Deposit{}
		mapper.GetCodec().MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		// refund deposit
		depositor := accountMapper.GetAccount(deposit.Depositor).(*types.QOSAccount)
		depositor.PlusQOS(btypes.NewInt(int64(deposit.Amount)))
		mapper.minusGovDeposit(deposit.Amount)

		mapper.Del(depositsIterator.Key())
	}
}

// Deletes all the deposits on a specific proposal without refunding them
func (mapper GovMapper) DeleteDeposits(ctx context.Context, proposalID uint64) {
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
	depositsIterator := mapper.GetDeposits(ctx, proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &gtypes.Deposit{}
		mapper.GetCodec().MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		// refund deposit
		depositor := accountMapper.GetAccount(deposit.Depositor).(*types.QOSAccount)
		depositor.PlusQOS(btypes.NewInt(int64(deposit.Amount)))
		mapper.minusGovDeposit(deposit.Amount)

		mapper.Del(depositsIterator.Key())
	}
}

// ProposalQueues

// Returns an iterator for all the proposals in the Active Queue that expire by endTime
func (mapper GovMapper) ActiveProposalQueueIterator(ctx context.Context, endTime time.Time) store.Iterator {
	return mapper.GetStore().Iterator(PrefixActiveProposalQueue, store.PrefixEndBytes(PrefixActiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (mapper GovMapper) InsertActiveProposalQueue(ctx context.Context, endTime time.Time, proposalID uint64) {
	mapper.Set(KeyActiveProposalQueueProposal(endTime, proposalID), proposalID)
}

// removes a proposalID from the Active Proposal Queue
func (mapper GovMapper) RemoveFromActiveProposalQueue(ctx context.Context, endTime time.Time, proposalID uint64) {
	mapper.Del(KeyActiveProposalQueueProposal(endTime, proposalID))
}

// Returns an iterator for all the proposals in the Inactive Queue that expire by endTime
func (mapper GovMapper) InactiveProposalQueueIterator(ctx context.Context, endTime time.Time) store.Iterator {
	return mapper.GetStore().Iterator(PrefixInactiveProposalQueue, store.PrefixEndBytes(PrefixInactiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the inactive proposal queue at endTime
func (mapper GovMapper) InsertInactiveProposalQueue(ctx context.Context, endTime time.Time, proposalID uint64) {
	mapper.Set(KeyInactiveProposalQueueProposal(endTime, proposalID), proposalID)
}

// removes a proposalID from the Inactive Proposal Queue
func (mapper GovMapper) RemoveFromInactiveProposalQueue(ctx context.Context, endTime time.Time, proposalID uint64) {
	mapper.Del(KeyInactiveProposalQueueProposal(endTime, proposalID))
}

// get deposit
func (mapper *GovMapper) getGovDeposit() (deposit uint64) {
	exists := mapper.Get(GovDepositKey, &deposit)
	if !exists {
		return 0
	}

	return deposit
}

// add deposit
func (mapper *GovMapper) addGovDeposit(deposit uint64) {
	current := mapper.getGovDeposit()
	mapper.Set(GovDepositKey, current+deposit)
}

// minus deposit
func (mapper *GovMapper) minusGovDeposit(deposit uint64) {
	current := mapper.getGovDeposit()
	if current < deposit {
		panic("no enough deposit to minus")
	}
	mapper.Set(GovDepositKey, current-deposit)
}

// get burned deposit
func (mapper *GovMapper) getGovBurnedDeposit() (deposit uint64) {
	exists := mapper.Get(BurnedGovDepositKey, &deposit)
	if !exists {
		return 0
	}

	return deposit
}

// burn deposit
func (mapper *GovMapper) burnGovDeposit(deposit uint64) {
	current := mapper.getGovDeposit()
	if deposit > current {
		panic("no enough deposit to burn")
	}
	mapper.Set(GovDepositKey, current-deposit)
	burned := mapper.getGovBurnedDeposit()
	mapper.Set(BurnedGovDepositKey, burned+deposit)
}
