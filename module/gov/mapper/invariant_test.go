package mapper

import (
	"github.com/QOSGroup/qbase/baseabci"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/params"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestDepositInvariant(t *testing.T) {
	ctx := defaultContext()
	addr := btypes.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	am := baseabci.GetAccountMapper(ctx)
	am.SetAccount(qtypes.NewQOSAccount(addr, btypes.NewInt(1000), nil))
	params.GetMapper(ctx).RegisterParamSet(&types.Params{})
	initGenesis(ctx, types.DefaultGenesisState())
	govMapper := GetMapper(ctx)

	textContent := types.NewTextProposal("p1", "p1", btypes.NewInt(10))
	proposal, _ := govMapper.SubmitProposal(ctx, textContent)
	govMapper.AddDeposit(ctx, proposal.ProposalID, addr, textContent.Deposit)
	_, coins, broken := DepositInvariant("gov")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(10))

	textContent = types.NewTextProposal("p2", "p2", btypes.NewInt(10))
	proposal, _ = govMapper.SubmitProposal(ctx, textContent)
	govMapper.AddDeposit(ctx, proposal.ProposalID, addr, textContent.Deposit)
	_, coins, broken = DepositInvariant("gov")(ctx)
	assert.False(t, broken)
	assert.Equal(t, coins.AmountOf(qtypes.QOSCoinName), btypes.NewInt(20))
}
