package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func DepositInvariant(module string) qtypes.Invariant {
	return func(ctx context.Context) (string, btypes.BaseCoins, bool) {
		gm := GetMapper(ctx)

		var depositTokens = btypes.ZeroInt()
		for _, proposal := range gm.GetProposals() {
			proposalID := proposal.ProposalID
			depositsIterator := gm.GetDeposits(proposalID)
			for ; depositsIterator.Valid(); depositsIterator.Next() {
				var deposit types.Deposit
				gm.GetCodec().MustUnmarshalBinaryBare(depositsIterator.Value(), &deposit)
				depositTokens = depositTokens.Add(deposit.Amount)
			}
			depositsIterator.Close()
		}

		return qtypes.FormatInvariant(module, "deposit",
			fmt.Sprintf("total deposit %d\n", depositTokens), btypes.BaseCoins{btypes.NewBaseCoin(qtypes.QOSCoinName, depositTokens)}, false)
	}
}
