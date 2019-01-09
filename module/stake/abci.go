package stake

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	stakemapper "github.com/QOSGroup/qos/module/stake/mapper"
	staketypes "github.com/QOSGroup/qos/module/stake/types"
	"github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//1. 统计validator投票信息, 将不活跃的validator转成Inactive状态
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	validatorMapper := stakemapper.GetValidatorMapper(ctx)

	votingWindowLen := uint64(validatorMapper.GetParams().ValidatorVotingStatusLen)
	minVotingCounter := uint64(validatorMapper.GetParams().ValidatorVotingStatusLeast)

	for _, signingValidator := range req.LastCommitInfo.Votes {
		valAddr := btypes.Address(signingValidator.Validator.Address)
		voted := signingValidator.SignedLastBlock
		handleValidatorValidatorVoteInfo(ctx, valAddr, voted, votingWindowLen, minVotingCounter)
	}
}

//1. 将所有Inactive到一定期限的validator删除
//2. 统计新的validator
func EndBlocker(ctx context.Context) (res abci.ResponseEndBlock) {

	validatorMapper := stakemapper.GetValidatorMapper(ctx)
	survivalSecs := validatorMapper.GetParams().ValidatorSurvivalSecs
	maxValidatorCount := uint64(validatorMapper.GetParams().MaxValidatorCnt)

	closeExpireInactiveValidator(ctx, survivalSecs)
	res.ValidatorUpdates = GetUpdatedValidators(ctx, maxValidatorCount)
	return
}

func closeExpireInactiveValidator(ctx context.Context, survivalSecs uint32) {
	log := ctx.Logger()
	validatorMapper := stakemapper.GetValidatorMapper(ctx)
	voteInfoMapper := stakemapper.GetVoteInfoMapper(ctx)
	accountMapper := baseabci.GetAccountMapper(ctx)

	blockTimeSec := uint64(ctx.BlockHeader().Time.UTC().Unix())
	lastCloseValidatorSec := blockTimeSec - uint64(survivalSecs)

	iterator := validatorMapper.IteratorInactiveValidator(uint64(0), lastCloseValidatorSec)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		valAddress := btypes.Address(key[9:])

		log.Info("close validator", "height", ctx.BlockHeight(), "validator", valAddress.String())
		if validator, ok := validatorMapper.KickValidator(valAddress); ok {

			voteInfoMapper.DelValidatorVoteInfo(valAddress)
			voteInfoMapper.ClearValidatorVoteInfoInWindow(valAddress)

			//关闭validator后,归还绑定的token
			owner := accountMapper.GetAccount(validator.Owner)
			if qosAcc, ok := owner.(*types.QOSAccount); ok {
				backQOS := btypes.NewInt(int64(validator.BondTokens))
				qosAcc.SetQOS(qosAcc.GetQOS().NilToZero().Add(backQOS))
				accountMapper.SetAccount(qosAcc)
			}
		}
	}
}

func GetUpdatedValidators(ctx context.Context, maxValidatorCount uint64) []abci.ValidatorUpdate {
	log := ctx.Logger()
	validatorMapper := ctx.Mapper(stakemapper.ValidatorMapperName).(*stakemapper.ValidatorMapper)

	//获取当前的validator集合
	var currentValidators []staketypes.Validator
	validatorMapper.Get(stakemapper.BuildCurrentValidatorAddressKey(), &currentValidators)

	currentValidatorMap := make(map[string]staketypes.Validator)
	for _, curValidator := range currentValidators {
		curValidatorAddrString := curValidator.GetValidatorAddress().String()
		currentValidatorMap[curValidatorAddrString] = curValidator
	}

	//返回更新的validator
	updateValidators := make([]abci.ValidatorUpdate, 0, len(currentValidatorMap))

	i := uint64(0)
	newValidatorsMap := make(map[string]staketypes.Validator)
	newValidators := make([]staketypes.Validator, 0, len(currentValidators))

	iterator := validatorMapper.IteratorValidatrorByVoterPower(false)
	defer iterator.Close()

	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])

		if i >= maxValidatorCount {
			//超出MaxValidatorCnt的validator修改为Inactive状态
			validatorMapper.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time, staketypes.MaxValidator)

		} else {
			if validator, exsits := validatorMapper.GetValidator(valAddr); exsits {
				i++
				//保存数据
				newValidatorAddressString := validator.GetValidatorAddress().String()
				newValidatorsMap[newValidatorAddressString] = validator
				newValidators = append(newValidators, validator)

				//新增或修改
				curValidator, exsits := currentValidatorMap[newValidatorAddressString]
				if !exsits || (validator.BondTokens != curValidator.BondTokens) {
					updateValidators = append(updateValidators, validator.ToABCIValidatorUpdate(false))
				}
			}
		}
	}

	//删除
	for curValidatorAddr, curValidator := range currentValidatorMap {
		if _, ok := newValidatorsMap[curValidatorAddr]; !ok {
			// curValidator.Power = 0
			updateValidators = append(updateValidators, curValidator.ToABCIValidatorUpdate(true))
		}
	}

	//存储新的validator
	validatorMapper.Set(stakemapper.BuildCurrentValidatorAddressKey(), newValidators)

	log.Info("update Validators", "len", len(updateValidators))

	return updateValidators
}

func handleValidatorValidatorVoteInfo(ctx context.Context, valAddr btypes.Address, isVote bool, votingWindowLen, minVotingCounter uint64) {

	log := ctx.Logger()
	height := uint64(ctx.BlockHeight())
	validatorMapper := stakemapper.GetValidatorMapper(ctx)
	voteInfoMapper := stakemapper.GetVoteInfoMapper(ctx)

	validator, exsits := validatorMapper.GetValidator(valAddr)
	if !exsits {
		log.Info("validatorVoteInfo", valAddr.String(), "not exsits,may be closed")
		return
	}

	//非Active状态不处理
	if !validator.IsActive() {
		log.Info("validatorVoteInfo", valAddr.String(), "is Inactive")
		return
	}

	voteInfo, exsits := voteInfoMapper.GetValidatorVoteInfo(valAddr)
	if !exsits {
		voteInfo = staketypes.NewValidatorVoteInfo(height, 0, 0)
	}

	index := voteInfo.IndexOffset % votingWindowLen
	voteInfo.IndexOffset++

	previousVote := voteInfoMapper.GetVoteInfoInWindow(valAddr, index)

	switch {
	case previousVote && !isVote:
		voteInfoMapper.SetVoteInfoInWindow(valAddr, index, false)
		voteInfo.MissedBlocksCounter++
	case !previousVote && isVote:
		voteInfoMapper.SetVoteInfoInWindow(valAddr, index, true)
		voteInfo.MissedBlocksCounter--
	default:
		//nothing
	}

	if !isVote {
		log.Info("validatorVoteInfo", "height", height, valAddr.String(), "not vote")
	}

	// minHeight := voteInfo.StartHeight + votingWindowLen
	maxMissedCounter := votingWindowLen - minVotingCounter

	// if height > minHeight && voteInfo.MissedBlocksCounter > maxMissedCounter
	if voteInfo.MissedBlocksCounter > maxMissedCounter {
		log.Info("validator gets inactive", "height", height, "validator", valAddr.String(), "missed counter", voteInfo.MissedBlocksCounter)

		blockValidator(ctx, valAddr)

		// voteInfo.IndexOffset = 0
		// voteInfo.MissedBlocksCounter = 0
		// voteInfoMapper.ClearValidatorVoteInfoInWindow(valAddr)
	}

	voteInfoMapper.SetValidatorVoteInfo(valAddr, voteInfo)
}

//
func blockValidator(ctx context.Context, valAddr btypes.Address) {
	validatorMapper := stakemapper.GetValidatorMapper(ctx)
	validatorMapper.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time, staketypes.MissVoteBlock)
}

func closeActiveValidators(ctx context.Context) {
	validatorMapper := stakemapper.GetValidatorMapper(ctx)
	voteInfoMapper := stakemapper.GetVoteInfoMapper(ctx)
	accountMapper := baseabci.GetAccountMapper(ctx)

	iterator := validatorMapper.IteratorValidatrorByVoterPower(false)
	defer iterator.Close()

	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddress := btypes.Address(key[9:])
		if validator, ok := validatorMapper.KickValidator(valAddress); ok {

			voteInfoMapper.DelValidatorVoteInfo(valAddress)
			voteInfoMapper.ClearValidatorVoteInfoInWindow(valAddress)

			owner := accountMapper.GetAccount(validator.Owner)
			if qosAcc, ok := owner.(*types.QOSAccount); ok {
				qosAcc.MustPlusQOS(btypes.NewInt(int64(validator.BondTokens)))
				accountMapper.SetAccount(qosAcc)
			}

			validatorMapper.Del(stakemapper.BuildValidatorByVotePower(validator.BondTokens, valAddress))
		}
	}
}

func CloseAllValidators(ctx context.Context) {
	closeExpireInactiveValidator(ctx, 0)
	closeActiveValidators(ctx)
}
