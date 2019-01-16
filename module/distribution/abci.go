package distribution

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco/mapper"
	"github.com/QOSGroup/qos/module/eco/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//beginblocker根据Vote信息进行QOS分配
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	totalPower, signedTotalPower := int64(0), int64(0)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		totalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			signedTotalPower += voteInfo.Validator.Power
		}
	}

	distributionMapper := mapper.GetDistributionMapper(ctx)

	if ctx.BlockHeight() > 1 {
		var previousProposer btypes.Address
		distributionMapper.Get(types.BuildLastProposerKey(), &previousProposer)
		allocateQOS(ctx, uint64(signedTotalPower), uint64(totalPower), previousProposer, req.LastCommitInfo.GetVotes())
	}

	consAddr := btypes.Address(req.Header.ProposerAddress)
	distributionMapper.Set(types.BuildLastProposerKey(), consAddr)

}

//endblocker对delegator的收益进行发放,并决定是否有下一次收益
func EndBlocker(ctx context.Context, req abci.RequestEndBlock) {

}

//
// 2.  每块挖出的QOS数量:  `x%`proposer + `y%`validators + `z%`community
//        * `x%`proposer: 验证人获得的奖励
//        * `y%`validators: 根据每个validator的power占比平均分配
// 3.  validator奖励数 =  validator佣金 +  平分金额Fee
//        * validator佣金奖励: 佣金 = validator奖励数 * `commission rate`
//        * 平分金额Fee由validator,delegator根据各自绑定的stake平均分配
// 4.  validator的proposer奖励,佣金奖励 均按周期发放
//
func allocateQOS(ctx context.Context, signedTotalPower, totalPower uint64, proposerAddr btypes.Address, votes []abci.VoteInfo) {
	// distributionMapper := mapper.GetDistributionMapper(ctx)

	// params := distributionMapper.GetParams()

	// //获取待分配的QOS总量
	// totalAmount := distributionMapper.GetPreDistributionQOS()
	// distributionMapper.ClearPreDistributionQOS()

	// //社区奖励
	// communityFeePool := distributionMapper.GetCommunityFeePool()
	//TODO
}
