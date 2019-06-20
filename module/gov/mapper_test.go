package gov

import (
	"github.com/QOSGroup/qbase/baseabci"
	btypes "github.com/QOSGroup/qbase/types"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	vtypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestGovMapper_SetGetProposal(t *testing.T) {
	ctx := defaultContext()
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())

	content := types.NewTextProposal("p1", "p1", 10)
	govMapper := GetGovMapper(ctx)
	p1, err := govMapper.SubmitProposal(ctx, content)
	require.Nil(t, err)

	p2, exists := govMapper.GetProposal(p1.ProposalID)
	require.True(t, exists)
	require.Equal(t, p1.ProposalID, p2.ProposalID)
}

func TestGovMapper_DeleteProposal(t *testing.T) {
	ctx := defaultContext()
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())

	content := types.NewTextProposal("p1", "p1", 10)
	govMapper := GetGovMapper(ctx)
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
	govMapper := GetGovMapper(ctx)
	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr1, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	proposal := NewTxProposal("p1", "p1", addr1, 10)
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

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
	govMapper := GetGovMapper(ctx)
	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr1, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	proposal := NewTxProposal("p1", "p1", addr1, 10)
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	proposals := govMapper.GetProposals()
	require.Equal(t, 1, len(proposals))
}

func TestGovMapper_ProposalID(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetGovMapper(ctx)
	_, err := govMapper.getNewProposalID()
	require.NotNil(t, err)

	err = govMapper.setInitialProposalID(1)
	require.Nil(t, err)

	id, err := govMapper.getNewProposalID()
	require.Nil(t, err)
	require.Equal(t, uint64(1), id)

	id = govMapper.GetLastProposalID()
	require.Equal(t, uint64(1), id)

	id, err = govMapper.peekCurrentProposalID()
	require.Nil(t, err)
	require.Equal(t, uint64(2), id)
}

func TestValidatorSet(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetGovMapper(ctx)
	validatorMapper := ecomapper.GetValidatorMapper(ctx)
	validatorMapper.CreateValidator(vtypes.Validator{Description: vtypes.Description{Moniker: "qos"}, Owner: btypes.Address(""), ValidatorPubKey: ed25519.PubKeyEd25519{}})

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
	govMapper := GetGovMapper(ctx)
	params.GetMapper(ctx).RegisterParamSet(&Params{})

	params := govMapper.GetParams(ctx)
	require.Zero(t, params.MinDeposit)

	params = DefaultParams()
	govMapper.SetParams(ctx, params)

	params = govMapper.GetParams(ctx)
	require.NotNil(t, params)
	require.NotZero(t, params.MinDeposit)
}

func TestGovMapper_GetSetVote(t *testing.T) {
	ctx := defaultContext()
	govMapper := GetGovMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())

	govMapper.setVote(1, addr, types.Vote{addr, 1, types.OptionYes})

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
	govMapper := GetGovMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	proposal := NewTxProposal("p1", "p1", addr, 10)
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	account := accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	require.Equal(t, btypes.NewInt(10), account.QOS)

	govMapper.RefundDeposits(ctx, 1, true)
	account = accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	burn := BurnRate.Mul(qtypes.MustNewDecFromStr("10")).TruncateInt()
	require.Equal(t, btypes.NewInt(10).Add(btypes.NewInt(10).Sub(burn)), account.QOS)
	pool := ecomapper.GetDistributionMapper(ctx).GetCommunityFeePool()
	require.Equal(t, burn, pool)
	_, exists := govMapper.GetDeposit(1, addr)
	require.False(t, exists)
}

func TestGovMapper_DeleteDeposits(t *testing.T) {
	ctx := defaultContext()
	accountMapper := baseabci.GetAccountMapper(ctx)
	govMapper := GetGovMapper(ctx)
	addr := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	accountMapper.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(20), nil))
	params.GetMapper(ctx).RegisterParamSet(&Params{})
	InitGenesis(ctx, DefaultGenesisState())
	proposal := NewTxProposal("p1", "p1", addr, 10)
	result, _ := proposal.Exec(ctx)
	require.Equal(t, result.Code, btypes.CodeOK)

	account := accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	require.Equal(t, btypes.NewInt(10), account.QOS)

	govMapper.DeleteDeposits(ctx, 1)
	account = accountMapper.GetAccount(addr).(*qtypes.QOSAccount)
	require.Equal(t, btypes.NewInt(10), account.QOS)
	pool := ecomapper.GetDistributionMapper(ctx).GetCommunityFeePool()
	require.Equal(t, btypes.NewInt(10), pool)
	_, exists := govMapper.GetDeposit(1, addr)
	require.False(t, exists)

}
