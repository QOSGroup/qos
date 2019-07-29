package stake

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//1. 统计validator投票信息, 将不活跃的validator转成Inactive状态
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	sm := mapper.GetMapper(ctx)

	votingWindowLen := uint64(sm.GetParams(ctx).ValidatorVotingStatusLen)
	minVotingCounter := uint64(sm.GetParams(ctx).ValidatorVotingStatusLeast)

	for _, signingValidator := range req.LastCommitInfo.Votes {
		valAddr := btypes.Address(signingValidator.Validator.Address)
		voted := signingValidator.SignedLastBlock
		handleValidatorValidatorVoteInfo(ctx, valAddr, voted, votingWindowLen, minVotingCounter)
	}
}

//1. 返还到期unbond tokens
//1. 处理到期redelegations
//2. 将所有Inactive到一定期限的validator删除
//3. 统计新的validator
func EndBlocker(ctx context.Context) []abci.ValidatorUpdate {

	// return unbond tokens
	returnUnBondTokens(ctx)

	// redelegations
	handlerReDelegations(ctx)

	// close inactive validators
	sm := mapper.GetMapper(ctx)
	survivalSecs := sm.GetParams(ctx).ValidatorSurvivalSecs
	CloseInactiveValidator(ctx, int32(survivalSecs))

	// return updated validators
	maxValidatorCount := uint64(sm.GetParams(ctx).MaxValidatorCnt)
	return GetUpdatedValidators(ctx, maxValidatorCount)
}

func returnUnBondTokens(ctx context.Context) {
	sm := mapper.GetMapper(ctx)
	am := baseabci.GetAccountMapper(ctx)
	prePrefix := types.BuildUnbondingDelegationByHeightPrefix(uint64(ctx.BlockHeight()))
	iter := btypes.KVStorePrefixIterator(sm.GetStore(), prePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		sm.Del(k)

		var unbondings []types.UnbondingDelegationInfo
		sm.BaseMapper.DecodeObject(iter.Value(), &unbondings)

		height, delAddr := types.GetUnbondingDelegationHeightAddress(k)
		for _, unbonding := range unbondings {
			delegator := am.GetAccount(delAddr).(*qtypes.QOSAccount)
			delegator.PlusQOS(btypes.NewInt(int64(unbonding.Amount)))
			am.SetAccount(delegator)
		}

		sm.RemoveUnbondingDelegations(height, delAddr)
	}
}

func handlerReDelegations(ctx context.Context) {
	sm := mapper.GetMapper(ctx)
	prePrefix := types.BuildRedelegationByHeightPrefix(uint64(ctx.BlockHeight()))
	iter := btypes.KVStorePrefixIterator(sm.GetStore(), prePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		sm.Del(k)

		var reDelegations []types.RedelegationInfo
		sm.BaseMapper.DecodeObject(iter.Value(), &reDelegations)

		height, delAddr := types.GetRedelegationHeightAddress(k)
		for _, reDelegation := range reDelegations {
			validator, _ := sm.GetValidator(reDelegation.ToValidator)
			sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens()+reDelegation.Amount)
			sm.Delegate(ctx, NewDelegationInfo(reDelegation.DelegatorAddr, reDelegation.ToValidator, reDelegation.Amount, reDelegation.IsCompound), true)
		}

		sm.RemoveRedelegations(height, delAddr)
	}
}

func CloseInactiveValidator(ctx context.Context, survivalSecs int32) {
	sm := mapper.GetMapper(ctx)

	blockTimeSec := uint64(ctx.BlockHeader().Time.UTC().Unix())
	var lastCloseValidatorSec uint64
	if survivalSecs >= 0 {
		lastCloseValidatorSec = blockTimeSec - uint64(survivalSecs)
	} else { // close all
		lastCloseValidatorSec = blockTimeSec + uint64(survivalSecs)
	}

	iterator := sm.IteratorInactiveValidator(uint64(0), lastCloseValidatorSec)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		valAddress := btypes.Address(key[9:])
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeCloseValidator,
				btypes.NewAttribute(types.AttributeKeyHeight, string(ctx.BlockHeight())),
				btypes.NewAttribute(types.AttributeKeyValidator, valAddress.String()),
			),
		)
		RemoveValidator(ctx, valAddress)
	}
	iterator.Close()
}

//删除和validator相关数据
//CONTRACT:
//delegator当前收益和收益发放信息数据不删除, 只是将bondTokens重置为0
//发放收益时,若delegator非validator的委托人, 或validator 不存在 则可以将delegator的收益相关数据删除
//发放收益时,validator的汇总数据可能会不存在
func RemoveValidator(ctx context.Context, valAddr btypes.Address) error {

	sm := mapper.GetMapper(ctx)

	// 处理分配逻辑
	sm.BeforeValidatorRemoved(ctx, valAddr)

	// 删除validator相关数据
	sm.KickValidator(valAddr)
	sm.DelValidatorVoteInfo(valAddr)
	sm.ClearValidatorVoteInfoInWindow(valAddr)

	return nil
}

func GetUpdatedValidators(ctx context.Context, maxValidatorCount uint64) []abci.ValidatorUpdate {
	sm := mapper.GetMapper(ctx)

	//获取当前的validator集合
	var currentValidators []types.Validator
	sm.Get(types.BuildCurrentValidatorsAddressKey(), &currentValidators)

	currentValidatorMap := make(map[string]types.Validator)
	for _, curValidator := range currentValidators {
		curValidatorAddrString := curValidator.GetValidatorAddress().String()
		currentValidatorMap[curValidatorAddrString] = curValidator
	}

	//返回更新的validator
	updateValidators := make([]abci.ValidatorUpdate, 0, len(currentValidatorMap))

	i := uint64(0)
	newValidatorsMap := make(map[string]types.Validator)
	newValidators := make([]types.Validator, 0, len(currentValidators))

	iterator := sm.IteratorValidatorByVoterPower(false)
	defer iterator.Close()

	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])

		if i >= maxValidatorCount {
			//超出MaxValidatorCnt的validator修改为Inactive状态
			if validator, exists := sm.GetValidator(valAddr); exists {
				sm.MakeValidatorInactive(validator.GetValidatorAddress(), uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.MaxValidator)
			}
		} else {
			if validator, exists := sm.GetValidator(valAddr); exists {
				if !validator.IsActive() {
					continue
				}
				i++
				//保存数据
				newValidatorAddressString := validator.GetValidatorAddress().String()
				newValidatorsMap[newValidatorAddressString] = validator
				newValidators = append(newValidators, validator)

				//新增或修改
				curValidator, exists := currentValidatorMap[newValidatorAddressString]
				if !exists || (validator.GetBondTokens() != curValidator.BondTokens) {
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

	if len(newValidators) == 0 {
		panic("consens error. no validator exists")
	}

	//存储新的validator
	sm.Set(types.BuildCurrentValidatorsAddressKey(), newValidators)

	return updateValidators
}

func handleValidatorValidatorVoteInfo(ctx context.Context, valAddr btypes.Address, isVote bool, votingWindowLen, minVotingCounter uint64) {

	log := ctx.Logger()
	height := uint64(ctx.BlockHeight())
	sm := mapper.GetMapper(ctx)

	validator, exists := sm.GetValidator(valAddr)
	if !exists {
		log.Info("validatorVoteInfo", valAddr.String(), "not exists,may be closed")
		return
	}

	//非Active状态不处理
	if !validator.IsActive() {
		log.Info("validatorVoteInfo", valAddr.String(), "is Inactive")
		return
	}

	voteInfo, exists := sm.GetValidatorVoteInfo(valAddr)
	if !exists {
		voteInfo = types.NewValidatorVoteInfo(height, 0, 0)
	}

	index := voteInfo.IndexOffset % votingWindowLen
	voteInfo.IndexOffset++

	previousVote := sm.GetVoteInfoInWindow(valAddr, index)

	switch {
	case previousVote && !isVote:
		sm.SetVoteInfoInWindow(valAddr, index, false)
		voteInfo.MissedBlocksCounter++
	case !previousVote && isVote:
		sm.SetVoteInfoInWindow(valAddr, index, true)
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
		sm.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.MissVoteBlock)
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeInactiveValidator,
				btypes.NewAttribute(types.AttributeKeyHeight, string(height)),
				btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			),
		)
	}

	sm.SetValidatorVoteInfo(valAddr, voteInfo)
}
