package stake

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

//1. 双签惩罚
//2. 统计validator投票信息, 将不活跃的validator转成Inactive状态
func BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) {

	sm := mapper.GetMapper(ctx)

	// 双签惩罚
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height-1, evidence.Time, evidence.Validator.Power)
		default:
			ctx.Logger().Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}

	// 统计validator投票信息, 将不活跃的validator转成Inactive状态
	params := sm.GetParams(ctx)
	for _, signingValidator := range req.LastCommitInfo.Votes {
		handleValidatorValidatorVoteInfo(ctx, btypes.Address(signingValidator.Validator.Address), signingValidator.SignedLastBlock, params)
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

		var unbonding types.UnbondingDelegationInfo
		sm.BaseMapper.DecodeObject(iter.Value(), &unbonding)

		height, delAddr, valAddr := types.GetUnbondingDelegationHeightDelegatorValidator(k)
		delegator := am.GetAccount(delAddr).(*qtypes.QOSAccount)
		delegator.PlusQOS(btypes.NewInt(int64(unbonding.Amount)))
		am.SetAccount(delegator)
		sm.RemoveUnbondingDelegation(height, delAddr, valAddr)
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

		var reDelegation types.RedelegationInfo
		sm.BaseMapper.DecodeObject(iter.Value(), &reDelegation)

		height, delAddr, valAddr := types.GetRedelegationHeightDelegatorFromValidator(k)
		validator, _ := sm.GetValidator(reDelegation.ToValidator)
		sm.Delegate(ctx, NewDelegationInfo(reDelegation.DelegatorAddr, reDelegation.ToValidator, reDelegation.Amount, reDelegation.IsCompound), true)
		sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens()+reDelegation.Amount)

		sm.RemoveRedelegation(height, delAddr, valAddr)
	}
}

func CloseInactiveValidator(ctx context.Context, survivalSecs int32) {
	sm := mapper.GetMapper(ctx)

	blockTimeSec := uint64(ctx.BlockHeader().Time.UTC().Unix())
	var lastCloseValidatorSec uint64
	if survivalSecs >= 0 {
		lastCloseValidatorSec = blockTimeSec - uint64(survivalSecs)
	} else { // close all
		lastCloseValidatorSec = blockTimeSec + uint64(-survivalSecs)
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
		removeValidator(ctx, valAddress)
	}
	iterator.Close()
}

//删除和validator相关数据
//CONTRACT:
//delegator当前收益和收益发放信息数据不删除, 只是将bondTokens重置为0
//发放收益时,若delegator非validator的委托人, 或validator 不存在 则可以将delegator的收益相关数据删除
//发放收益时,validator的汇总数据可能会不存在
func removeValidator(ctx context.Context, valAddr btypes.Address) error {

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

func handleValidatorValidatorVoteInfo(ctx context.Context, valAddr btypes.Address, isVote bool, params types.Params) {

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

	index := voteInfo.IndexOffset % uint64(params.ValidatorVotingStatusLen)
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
	maxMissedCounter := params.ValidatorVotingStatusLen - params.ValidatorVotingStatusLeast

	// if height > minHeight && voteInfo.MissedBlocksCounter > maxMissedCounter
	if voteInfo.MissedBlocksCounter > uint64(maxMissedCounter) {

		// slash delegations
		delegationSlashTokens := slashDelegations(ctx, validator, params.SlashFractionDowntime, types.AttributeValueDoubleSign)
		updatedValidatorTokens := validator.BondTokens - delegationSlashTokens
		log.Debug("slash validator bond tokens", "validator", validator.GetValidatorAddress().String(), "preTokens", validator.BondTokens, "slashTokens", delegationSlashTokens, "afterTokens", updatedValidatorTokens)

		log.Info("validator gets inactive", "height", height, "validator", valAddr.String(), "missed counter", voteInfo.MissedBlocksCounter)
		sm.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.MissVoteBlock)
		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeInactiveValidator,
				btypes.NewAttribute(types.AttributeKeyHeight, string(height)),
				btypes.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
			),
		)

		// slash放入社区费池
		sm.AfterValidatorSlashed(ctx, delegationSlashTokens)
	}

	sm.SetValidatorVoteInfo(valAddr, voteInfo)
}

func handleDoubleSign(ctx context.Context, addr btypes.Address, infractionHeight int64, timestamp time.Time, power int64) {
	logger := ctx.Logger()
	sm := GetMapper(ctx)

	// validator should exists
	validator, exists := sm.GetValidator(addr)
	if !exists {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, the validator did not exist anymore", validator.GetValidatorAddress(), infractionHeight))
		return
	}

	// validator should not be inactive
	if validator.Status == types.Inactive {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, validator already inactive", validator.GetValidatorAddress(), infractionHeight))
		return
	}

	// calculate the age of the evidence
	age := ctx.BlockHeader().Time.Sub(timestamp)

	// Reject evidence if the double-sign is too old
	maxAge := sm.GetParams(ctx).MaxEvidenceAge
	if age > maxAge {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, age of %d past max age of %d",
			validator.GetValidatorAddress(), infractionHeight, age, maxAge))
		return
	}

	// double sign confirmed
	logger.Info(fmt.Sprintf("Confirmed double sign from %s at height %d, age of %d", validator.GetValidatorAddress(), infractionHeight, age))

	fraction := sm.GetParams(ctx).SlashFractionDoubleSign
	ctx.EventManager().EmitEvent(
		btypes.NewEvent(
			types.EventTypeSlash,
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
			btypes.NewAttribute(types.AttributeKeyOwner, validator.Owner.String()),
			btypes.NewAttribute(types.AttributeKeyReason, types.AttributeValueDoubleSign),
		),
	)

	if fraction.LT(qtypes.ZeroDec()) || fraction.GT(qtypes.OneDec()) {
		panic(fmt.Errorf("attempted to slash with a negative/gtone slash factor: %v", fraction))
	}

	slashAmount := fraction.MulInt(btypes.NewInt(power)).TruncateInt64()
	remainingSlashAmount := slashAmount
	switch {
	case infractionHeight > ctx.BlockHeight():

		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():

		// Special-case slash at current height for efficiency - we don't need to look through unbonding delegations or redelegations
		logger.Info(fmt.Sprintf(
			"slashing at current height %d, not scanning unbonding delegations & redelegations",
			infractionHeight))
	case infractionHeight < ctx.BlockHeight():

		// Iterate through unbonding delegations from slashed validator
		remainingSlashAmount = sm.SlashUnbondings(validator.GetValidatorAddress(), infractionHeight, fraction, remainingSlashAmount)

		// Iterate through redelegations from slashed source validator
		remainingSlashAmount = sm.SlashRedelegations(validator.GetValidatorAddress(), infractionHeight, fraction, remainingSlashAmount)
	}

	tokensToBurn := slashAmount - remainingSlashAmount
	if remainingSlashAmount == 0 {
		sm.AfterValidatorSlashed(ctx, uint64(tokensToBurn))
		return
	}

	// cannot decrease balance below zero
	if remainingSlashAmount > int64(validator.BondTokens) {
		remainingSlashAmount = int64(validator.BondTokens)
	}

	// calculate delegations slash fraction
	fraction = qtypes.NewDec(int64(remainingSlashAmount)).QuoInt(btypes.NewInt(int64(validator.BondTokens)))

	// slash delegations
	delegationSlashTokens := slashDelegations(ctx, validator, fraction, types.AttributeValueDoubleSign)

	// update validator
	updatedValidatorTokens := validator.BondTokens - delegationSlashTokens
	logger.Info("validator gets inactive", "height", ctx.BlockHeight(), "validator", validator.GetValidatorAddress().String())
	sm.MakeValidatorInactive(validator.GetValidatorAddress(), uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), types.DoubleSign)
	ctx.EventManager().EmitEvent(
		btypes.NewEvent(
			types.EventTypeInactiveValidator,
			btypes.NewAttribute(types.AttributeKeyHeight, string(ctx.BlockHeight())),
			btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
		),
	)

	logger.Debug("slash validator bond tokens", "validator", validator.GetValidatorAddress().String(), "preTokens", validator.BondTokens, "slashTokens", delegationSlashTokens, "afterTokens", updatedValidatorTokens)

	// slash放入社区费池
	sm.AfterValidatorSlashed(ctx, delegationSlashTokens)

}

func slashDelegations(ctx context.Context, validator types.Validator, fraction qtypes.Dec, reason string) uint64 {
	sm := mapper.GetMapper(ctx)
	logger := ctx.Logger()

	// get delegations
	var delegations []types.DelegationInfo
	sm.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr btypes.Address, delAddr btypes.Address) {
		if delegation, exists := sm.GetDelegationInfo(delAddr, valAddr); exists {
			delegations = append(delegations, delegation)
		}
	})

	delegationSlashTokens := int64(0)
	for _, delegation := range delegations {
		bondTokens := int64(delegation.Amount)

		// calculate slash amount
		amountSlash := qtypes.NewDec(bondTokens).Mul(fraction).TruncateInt64()
		if amountSlash > bondTokens {
			amountSlash = bondTokens
		}
		delegationSlashTokens += amountSlash

		// update delegation
		remainTokens := uint64(bondTokens - amountSlash)
		delegation.Amount = remainTokens
		sm.BeforeDelegationModified(ctx, validator.GetValidatorAddress(), delegation.DelegatorAddr, delegation.Amount)
		sm.SetDelegationInfo(delegation)

		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeSlash,
				btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
				btypes.NewAttribute(types.AttributeKeyDelegator, delegation.DelegatorAddr.String()),
				btypes.NewAttribute(types.AttributeKeyTokens, string(amountSlash)),
				btypes.NewAttribute(types.AttributeKeyReason, reason),
			),
		)
		logger.Debug("slash validator's delegators", "delegator", delegation.DelegatorAddr.String(), "preToken", bondTokens, "slashToken", amountSlash, "remainTokens", remainTokens)
	}

	// update validator
	sm.ChangeValidatorBondTokens(validator, validator.BondTokens-uint64(delegationSlashTokens))

	return uint64(delegationSlashTokens)
}
