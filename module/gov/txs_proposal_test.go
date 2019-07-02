package gov

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/mapper"
	vtypes "github.com/QOSGroup/qos/module/eco/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/guardian"
	guardiantypes "github.com/QOSGroup/qos/module/guardian/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	govMapper := NewGovMapper()
	govMapper.SetCodec(cdc)
	approveKey := govMapper.GetStoreKey()
	mapperMap[GovMapperName] = govMapper

	accountMapper := bacc.NewAccountMapper(nil, types.ProtoQOSAccount)
	accountMapper.SetCodec(cdc)
	accountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	paramMapper := params.NewMapper()
	paramMapper.SetCodec(cdc)
	paramsKey := paramMapper.GetStoreKey()
	mapperMap[params.MapperName] = paramMapper

	delegationMapper := mapper.NewValidatorMapper()
	delegationMapper.SetCodec(cdc)
	validatorKey := delegationMapper.GetStoreKey()
	mapperMap[vtypes.ValidatorMapperName] = delegationMapper

	delegatorMapper := mapper.NewDelegationMapper()
	delegatorMapper.SetCodec(cdc)
	delegatorKey := delegatorMapper.GetStoreKey()
	mapperMap[vtypes.DelegationMapperName] = delegatorMapper

	guardianMapper := guardian.NewGuardianMapper()
	guardianMapper.SetCodec(cdc)
	guardianKey := guardianMapper.GetStoreKey()
	mapperMap[guardian.MapperName] = guardianMapper

	distributiongMapper := mapper.NewDistributionMapper()
	distributiongMapper.SetCodec(cdc)
	distributionKey := distributiongMapper.GetStoreKey()
	mapperMap[vtypes.DistributionMapperName] = distributiongMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(accountKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(validatorKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(delegatorKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(guardianKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(distributionKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestTxProposal_ValidateData(t *testing.T) {
	ctx := defaultContext()
	InitGenesis(ctx, DefaultGenesisState())
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input *TxProposal
		valid bool
	}{
		{NewTxProposal("p1", "p1", addr, 10), false},
		{NewTxProposal("p1", "p1", addr, 10), true},
		{NewTxProposal("", "p1", addr, 10), false},
		{NewTxProposal("p1", "", addr, 10), false},
		{NewTxProposal("p1", "p1", nil, 10), false},
		{NewTxProposal("p1", "p1", addr, 1), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
		if tcIndex == 0 {
			accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
		}
	}
}

func TestTxProposal_Exec(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	proposal := NewTxProposal("p1", "p1", addr, 10)
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposalMapper := GetGovMapper(ctx)
	p, exists := proposalMapper.GetProposal(1)
	require.True(t, exists)
	require.NotNil(t, p)
	require.Equal(t, gtypes.ProposalTypeText, p.GetProposalType())
}

func TestTxParameterChange_ValidateData(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())

	proposal := NewTxParameterChange("p1", "p1", addr, 10, nil)

	cases := []struct {
		input []gtypes.Param
		valid bool
	}{
		{[]gtypes.Param{{"gov", "min_deposit", "10"}}, true},
		{[]gtypes.Param{}, false},
		{[]gtypes.Param{{"gov", "min_deposit1", "10"}}, false},
		{[]gtypes.Param{{"m", "min_deposit1", "10"}}, false},
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
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())

	proposal := NewTxParameterChange("p1", "p1", addr, 10, []gtypes.Param{{"gov", "min_deposit", "10"}})
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposalMapper := GetGovMapper(ctx)
	p, exists := proposalMapper.GetProposal(1)
	require.True(t, exists)
	require.NotNil(t, p)
	require.Equal(t, gtypes.ProposalTypeParameterChange, p.GetProposalType())
}

func TestTxTaxUsage_ValidateData(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	dest := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	guardian.InitGenesis(ctx, guardian.GenesisState{[]guardiantypes.Guardian{{"g1", guardiantypes.Genesis, dest, nil}}})

	cases := []struct {
		input *TxTaxUsage
		valid bool
	}{
		{NewTxTaxUsage("p1", "p1", addr, 10, dest, types.MustNewDecFromStr("0.5")), true},
		{NewTxTaxUsage("p1", "p1", addr, 10, addr, types.MustNewDecFromStr("0.5")), false},
		{NewTxTaxUsage("p1", "p1", addr, 10, dest, types.MustNewDecFromStr("0")), false},
	}

	for tcIndex, tc := range cases {
		err := tc.input.ValidateData(ctx)
		require.Equal(t, tc.valid, err == nil, "tc #%d", tcIndex)
	}
}

func TestTxTaxUsage_Exec(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	dest := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(types.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())

	proposal := NewTxTaxUsage("p1", "p1", addr, 10, dest, types.MustNewDecFromStr("0.5"))
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposalMapper := GetGovMapper(ctx)
	p, exists := proposalMapper.GetProposal(1)
	require.True(t, exists)
	require.NotNil(t, p)
	require.Equal(t, gtypes.ProposalTypeTaxUsage, p.GetProposalType())
}
