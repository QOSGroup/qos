package mapper

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/stake"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	cdc := baseabci.MakeQBaseCodec()
	qtypes.RegisterCodec(cdc)
	cdc.RegisterInterface((*types.ProposalContent)(nil), nil)
	cdc.RegisterConcrete(&types.TextProposal{}, "gov/TextProposal", nil)
	cdc.RegisterConcrete(&types.TaxUsageProposal{}, "gov/TaxUsageProposal", nil)
	cdc.RegisterConcrete(&types.ParameterProposal{}, "gov/ParameterProposal", nil)

	govMapper := NewMapper()
	govMapper.SetCodec(cdc)
	approveKey := govMapper.GetStoreKey()
	mapperMap[MapperName] = govMapper

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	accountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	paramMapper := params.NewMapper()
	paramMapper.SetCodec(cdc)
	paramsKey := paramMapper.GetStoreKey()
	mapperMap[params.MapperName] = paramMapper

	validatorMapper := stake.NewMapper()
	validatorMapper.SetCodec(cdc)
	validatorKey := validatorMapper.GetStoreKey()
	mapperMap[stake.MapperName] = validatorMapper

	guardianMapper := guardian.NewMapper()
	guardianMapper.SetCodec(cdc)
	guardianKey := guardianMapper.GetStoreKey()
	mapperMap[guardian.MapperName] = guardianMapper

	distributionMapper := distribution.NewMapper()
	distributionMapper.SetCodec(cdc)
	distributionKey := distributionMapper.GetStoreKey()
	mapperMap[distribution.MapperName] = distributionMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(accountKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(validatorKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(guardianKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(distributionKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func initGenesis(ctx context.Context, data types.GenesisState) {
	err := types.ValidateGenesis(data)
	if err != nil {
		panic(err)
	}
	mapper := GetMapper(ctx)
	err = mapper.SetInitialProposalID(data.StartingProposalID)
	if err != nil {
		panic(err)
	}
	mapper.SetParams(ctx, data.Params)
	for _, proposal := range data.Proposals {
		switch proposal.Proposal.Status {
		case types.StatusDepositPeriod:
			mapper.InsertInactiveProposalQueue(proposal.Proposal.DepositEndTime, proposal.Proposal.ProposalID)
		case types.StatusVotingPeriod:
			mapper.InsertActiveProposalQueue(proposal.Proposal.VotingEndTime, proposal.Proposal.ProposalID)
		}
		for _, deposit := range proposal.Deposits {
			mapper.SetDeposit(deposit.ProposalID, deposit.Depositor, deposit)
		}
		for _, vote := range proposal.Votes {
			mapper.SetVote(vote.ProposalID, vote.Voter, vote)
		}
		mapper.SetProposal(proposal.Proposal)
	}
}

func TestGovMapper_SetGetProposal(t *testing.T) {
	ctx := defaultContext()
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())

	content := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	govMapper := GetMapper(ctx)
	p1, err := govMapper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	p2, exists := govMapper.GetProposal(p1.ProposalID)
	require.True(t, exists)
	require.Equal(t, p1.ProposalID, p2.ProposalID)
}

func TestGovMapper_DeleteProposal(t *testing.T) {
	ctx := defaultContext()
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())

	content := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	govMapper := GetMapper(ctx)
	p1, err := govMapper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	p2, exists := govMapper.GetProposal(p1.ProposalID)
	require.True(t, exists)
	require.Equal(t, p1.ProposalID, p2.ProposalID)

	govMapper.DeleteProposal(p1.ProposalID)
	_, exists = govMapper.GetProposal(p1.ProposalID)
	require.False(t, exists)
}

func TestGovMapper_GetProposalsFiltered(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	govMapper := GetMapper(ctx)
	addr1 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr1, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())

	textContent := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addr1, textContent.Deposit)
	require.Nil(t, err)

	proposals := govMapper.GetProposalsFiltered(ctx, nil, addr1, types.StatusVotingPeriod, 0)
	require.Equal(t, 1, len(proposals))

	proposals = govMapper.GetProposalsFiltered(ctx, nil, addr1, types.StatusDepositPeriod, 0)
	require.Equal(t, 0, len(proposals))

	proposals = govMapper.GetProposalsFiltered(ctx, nil, addr2, types.StatusVotingPeriod, 0)
	require.Equal(t, 0, len(proposals))
}

func TestGovMapper_GetProposals(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	govMapper := GetMapper(ctx)
	addr1 := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr1, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	textContent := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addr1, textContent.Deposit)
	require.Nil(t, err)

	proposals := govMapper.GetProposals()
	require.Equal(t, 1, len(proposals))
}

func TestGovMapper_ProposalID(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetMapper(ctx)
	_, err := govMapper.getNewProposalID()
	require.NotNil(t, err)

	err = govMapper.SetInitialProposalID(1)
	require.Nil(t, err)

	id, err := govMapper.getNewProposalID()
	require.Nil(t, err)
	require.Equal(t, int64(1), id)

	id = govMapper.GetLastProposalID()
	require.Equal(t, int64(1), id)

	id, err = govMapper.PeekCurrentProposalID()
	require.Nil(t, err)
	require.Equal(t, int64(2), id)
}

func TestValidatorSet(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetMapper(ctx)
	sm := stake.GetMapper(ctx)
	sm.CreateValidator(stake.Validator{Description: stake.Description{Moniker: "qos"}, OperatorAddress: btypes.ValAddress(""), ConsPubKey: ed25519.PubKeyEd25519{}, BondTokens: btypes.ZeroInt()})

	govMapper.saveValidatorSet(ctx, 1)
	validators, exists := govMapper.GetValidatorSet(1)
	require.True(t, exists)
	require.Equal(t, 1, len(validators))
	govMapper.DeleteValidatorSet(1)
	validators, exists = govMapper.GetValidatorSet(1)
	require.False(t, exists)
}

func TestGovMapper_GetSetParams(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetMapper(ctx)
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})

	params := govMapper.GetParams(ctx)
	require.Zero(t, params.NormalMinDeposit)

	params = types.DefaultParams()
	govMapper.SetParams(ctx, params)

	params = govMapper.GetParams(ctx)
	require.NotNil(t, params)
	require.NotZero(t, params.NormalMinDeposit)
}

func TestGovMapper_GetSetVote(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	govMapper.SetVote(1, addr, types.Vote{addr, 1, types.OptionYes})

	vote, exists := govMapper.GetVote(1, addr)
	require.True(t, exists)
	require.Equal(t, types.OptionYes, vote.Option)

	govMapper.deleteVote(1, addr)
	vote, exists = govMapper.GetVote(1, addr)
	require.False(t, exists)
}

func TestGovMapper_RefundDeposits(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	govMapper := GetMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	textContent := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addr, textContent.Deposit)
	require.Nil(t, err)

	account := accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	require.Equal(t, btypes.NewInt(10), account.QOS)

	govMapper.RefundDeposits(ctx, 1, textContent.GetProposalLevel(), true)
	govParams := govMapper.GetLevelParams(ctx, textContent.GetProposalLevel())
	account = accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	burn := govParams.BurnRate.Mul(qtypes.MustNewDecFromStr("10")).TruncateInt()
	require.Equal(t, btypes.NewInt(10).Add(btypes.NewInt(10).Sub(burn)), account.QOS)
	pool := distribution.GetMapper(ctx).GetCommunityFeePool()
	require.Equal(t, burn, pool)
	_, exists := govMapper.GetDeposit(1, addr)
	require.False(t, exists)
}

func TestGovMapper_DeleteDeposits(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	govMapper := GetMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	textContent := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	require.Nil(t, err)
	err, _ = govMapper.AddDeposit(ctx, proposal.ProposalID, addr, textContent.Deposit)
	require.Nil(t, err)

	account := accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	require.Equal(t, btypes.NewInt(10), account.QOS)

	govMapper.DeleteDeposits(ctx, 1)
	account = accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	require.Equal(t, btypes.NewInt(10), account.QOS)
	pool := distribution.GetMapper(ctx).GetCommunityFeePool()
	require.Equal(t, btypes.NewInt(10), pool)
	_, exists := govMapper.GetDeposit(1, addr)
	require.False(t, exists)

}
