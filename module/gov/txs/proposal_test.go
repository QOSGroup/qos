package txs

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov/mapper"
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

	govMapper := mapper.NewMapper()
	govMapper.SetCodec(Cdc)
	approveKey := govMapper.GetStoreKey()
	mapperMap[mapper.MapperName] = govMapper

	accountMapper := bacc.NewAccountMapper(nil, qtypes.ProtoQOSAccount)
	accountMapper.SetCodec(Cdc)
	accountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	paramMapper := params.NewMapper()
	paramMapper.SetCodec(Cdc)
	paramsKey := paramMapper.GetStoreKey()
	mapperMap[params.MapperName] = paramMapper

	stakingMapper := stake.NewMapper()
	stakingMapper.SetCodec(Cdc)
	validatorKey := stakingMapper.GetStoreKey()
	mapperMap[stake.MapperName] = stakingMapper

	guardianMapper := guardian.NewMapper()
	guardianMapper.SetCodec(Cdc)
	guardianKey := guardianMapper.GetStoreKey()
	mapperMap[guardian.MapperName] = guardianMapper

	distributionMapper := distribution.NewMapper()
	distributionMapper.SetCodec(Cdc)
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
	mapper := mapper.GetMapper(ctx)
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

func TestTxProposal_ValidateData(t *testing.T) {
	ctx := defaultContext()
	initGenesis(ctx, types.DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input *TxProposal
		valid bool
	}{
		{NewTxProposal("p1", "p1", addr, btypes.NewInt(10)), false},
		{NewTxProposal("p1", "p1", addr, btypes.NewInt(10)), true},
		{NewTxProposal("", "p1", addr, btypes.NewInt(10)), false},
		{NewTxProposal("p1", "", addr, btypes.NewInt(10)), false},
		{NewTxProposal("p1", "p1", nil, btypes.NewInt(10)), false},
		{NewTxProposal("p1", "p1", addr, btypes.NewInt(1)), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
		if tcIndex == 0 {
			accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
		}
	}
}

func TestTxProposal_Exec(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	proposal := NewTxProposal("p1", "p1", addr, btypes.NewInt(10))
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposalMapper := mapper.GetMapper(ctx)
	p, exists := proposalMapper.GetProposal(1)
	require.True(t, exists)
	require.NotNil(t, p)
	require.Equal(t, types.ProposalTypeText, p.GetProposalType())
}

func TestTxParameterChange_ValidateData(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())

	proposal := NewTxParameterChange("p1", "p1", addr, btypes.NewInt(10), nil)

	cases := []struct {
		input []types.Param
		valid bool
	}{
		{[]types.Param{{"gov", "min_deposit", "10"}}, true},
		{[]types.Param{}, false},
		{[]types.Param{{"gov", "min_deposit1", "10"}}, false},
		{[]types.Param{{"m", "min_deposit1", "10"}}, false},
	}

	for tcIndex, tc := range cases {
		proposal.Params = tc.input
		err := proposal.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
	}
}

func TestTxParameterChange_Exec(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())

	proposal := NewTxParameterChange("p1", "p1", addr, btypes.NewInt(10), []types.Param{{"gov", "min_deposit", "10"}})
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposalMapper := mapper.GetMapper(ctx)
	p, exists := proposalMapper.GetProposal(1)
	require.True(t, exists)
	require.NotNil(t, p)
	require.Equal(t, types.ProposalTypeParameterChange, p.GetProposalType())
}

func TestTxTaxUsage_ValidateData(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	dest := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	guardian.InitGenesis(ctx, guardian.GenesisState{[]guardian.Guardian{{"g1", guardian.Genesis, dest, nil}}})

	cases := []struct {
		input *TxTaxUsage
		valid bool
	}{
		{NewTxTaxUsage("p1", "p1", addr, btypes.NewInt(10), dest, qtypes.MustNewDecFromStr("0.5")), true},
		{NewTxTaxUsage("p1", "p1", addr, btypes.NewInt(10), addr, qtypes.MustNewDecFromStr("0.5")), false},
		{NewTxTaxUsage("p1", "p1", addr, btypes.NewInt(10), dest, qtypes.MustNewDecFromStr("0")), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
	}
}

func TestTxTaxUsage_Exec(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	dest := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())

	proposal := NewTxTaxUsage("p1", "p1", addr, btypes.NewInt(10), dest, qtypes.MustNewDecFromStr("0.5"))
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposalMapper := mapper.GetMapper(ctx)
	p, exists := proposalMapper.GetProposal(1)
	require.True(t, exists)
	require.NotNil(t, p)
	require.Equal(t, types.ProposalTypeTaxUsage, p.GetProposalType())
}
