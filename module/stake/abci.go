package stake

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank"
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
			v, _ := sm.GetValidatorByConsensusAddr(btypes.ConsAddress(evidence.Validator.Address))
			handleDoubleSign(ctx, v.GetValidatorAddress(), evidence.Height-1, evidence.Time, evidence.Validator.Power)
		default:
			ctx.Logger().Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}

	// 统计validator投票信息, 将不活跃的validator转成Inactive状态
	params := sm.GetParams(ctx)
	for _, signingValidator := range req.LastCommitInfo.Votes {
		v, _ := sm.GetValidatorByConsensusAddr(btypes.ConsAddress(signingValidator.Validator.Address))
		handleValidatorValidatorVoteInfo(ctx, v.GetValidatorAddress(), signingValidator.SignedLastBlock, params)
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
	CloseInactiveValidator(ctx, survivalSecs)

	// return updated validators
	maxValidatorCount := sm.GetParams(ctx).MaxValidatorCnt
	return GetUpdatedValidators(ctx, maxValidatorCount)
}

func returnUnBondTokens(ctx context.Context) {
	sm := mapper.GetMapper(ctx)
	am := baseabci.GetAccountMapper(ctx)
	prePrefix := types.BuildUnbondingDelegationByHeightPrefix(ctx.BlockHeight())
	iter := btypes.KVStorePrefixIterator(sm.GetStore(), prePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		sm.Del(k)

		var unbonding types.UnbondingDelegationInfo
		sm.BaseMapper.DecodeObject(iter.Value(), &unbonding)

		height, delAddr, valAddr := types.GetUnbondingDelegationHeightDelegatorValidator(k)
		delegator := am.GetAccount(delAddr).(*qtypes.QOSAccount)
		delegator.PlusQOS(unbonding.Amount)
		am.SetAccount(delegator)
		sm.RemoveUnbondingDelegation(height, delAddr, valAddr)
	}
}

func handlerReDelegations(ctx context.Context) {
	sm := mapper.GetMapper(ctx)
	prePrefix := types.BuildRedelegationByHeightPrefix(ctx.BlockHeight())
	iter := btypes.KVStorePrefixIterator(sm.GetStore(), prePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		k := iter.Key()
		sm.Del(k)

		var reDelegation types.RedelegationInfo
		sm.BaseMapper.DecodeObject(iter.Value(), &reDelegation)

		height, delAddr, valAddr := types.GetRedelegationHeightDelegatorFromValidator(k)
		validator, exists := sm.GetValidator(reDelegation.ToValidator)
		if exists {
			sm.Delegate(ctx, NewDelegationInfo(reDelegation.DelegatorAddr, reDelegation.ToValidator, reDelegation.Amount, reDelegation.IsCompound), true)
			sm.ChangeValidatorBondTokens(validator, validator.GetBondTokens().Add(reDelegation.Amount))
		} else {
			// to validator 不存在时，返还待质押 tokens
			delegator := bank.GetAccount(ctx, reDelegation.DelegatorAddr)
			delegator.MustPlusQOS(reDelegation.Amount)
			bank.GetMapper(ctx).SetAccount(delegator)
		}
		sm.RemoveRedelegation(height, delAddr, valAddr)
	}
}

func CloseInactiveValidator(ctx context.Context, survivalSecs int64) {
	sm := mapper.GetMapper(ctx)

	blockTimeSec := ctx.BlockHeader().Time.UTC().Unix()
	lastCloseValidatorSec := blockTimeSec - survivalSecs

	iterator := sm.IteratorInactiveValidator(0, lastCloseValidatorSec)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		valAddress := btypes.ValAddress(key[9:])
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
func removeValidator(ctx context.Context, valAddr btypes.ValAddress) error {

	sm := mapper.GetMapper(ctx)

	// 处理分配逻辑
	sm.BeforeValidatorRemoved(ctx, valAddr)

	// 删除validator相关数据
	sm.KickValidator(valAddr)
	sm.DelValidatorVoteInfo(valAddr)
	sm.ClearValidatorVoteInfoInWindow(valAddr)

	return nil
}

func GetUpdatedValidators(ctx context.Context, maxValidatorCount int64) []abci.ValidatorUpdate {
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

	i := int64(0)
	newValidatorsMap := make(map[string]types.Validator)
	newValidators := make([]types.Validator, 0, len(currentValidators))

	iterator := sm.IteratorValidatorByVoterPower(false)
	defer iterator.Close()

	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()

		power, valAddr, err := types.ParseValidatorVotePowerKey(key)
		tokens := types.PowerReduction.MulRaw(power)
		if err != nil {
			ctx.Logger().Error("parse validatorVotePowerKey error", "key", key)
			panic(err)
		}

		if i >= maxValidatorCount {
			//超出MaxValidatorCnt的validator修改为Inactive状态
			if validator, exists := sm.GetValidator(valAddr); exists {
				sm.MakeValidatorInactive(validator.GetValidatorAddress(), ctx.BlockHeight(), ctx.BlockHeader().Time.UTC(), types.MaxValidator)
			}
		} else {
			if validator, exists := sm.GetValidator(valAddr); exists {
				if !validator.IsActive() {
					continue
				}

				if !validator.BondTokens.Equal(tokens) {
					ctx.Logger().Error("validator votePower list may have dup record. if you forgot delete?",
						"validator", validator.OperatorAddress.String(),
						"tokens", validator.BondTokens,
						"recordTokens", tokens)
					continue
				}

				i++
				//保存数据
				newValidatorAddressString := validator.GetValidatorAddress().String()
				newValidatorsMap[newValidatorAddressString] = validator
				newValidators = append(newValidators, validator)

				//新增或修改
				curValidator, exists := currentValidatorMap[newValidatorAddressString]
				if !exists || !(validator.GetBondTokens().Equal(curValidator.BondTokens)) {
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

func handleValidatorValidatorVoteInfo(ctx context.Context, valAddr btypes.ValAddress, isVote bool, params types.Params) {

	log := ctx.Logger()
	height := ctx.BlockHeight()
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

	index := voteInfo.IndexOffset % params.ValidatorVotingStatusLen
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
	if voteInfo.MissedBlocksCounter > maxMissedCounter {

		// slash delegations
		delegationSlashTokens := slashDelegations(ctx, validator, params.SlashFractionDowntime, types.AttributeValueDoubleSign)
		updatedValidatorTokens := validator.BondTokens.Sub(delegationSlashTokens)
		log.Debug("slash validator bond tokens", "validator", validator.GetValidatorAddress().String(), "preTokens", validator.BondTokens, "slashTokens", delegationSlashTokens, "afterTokens", updatedValidatorTokens)

		log.Info("validator gets inactive", "height", height, "validator", valAddr.String(), "missed counter", voteInfo.MissedBlocksCounter)
		sm.MakeValidatorInactive(valAddr, ctx.BlockHeight(), ctx.BlockHeader().Time.UTC(), types.MissVoteBlock)
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

func handleDoubleSign(ctx context.Context, addr btypes.ValAddress, infractionHeight int64, timestamp time.Time, power int64) {
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

	slashAmount := fraction.MulInt(btypes.NewInt(power)).MulInt(types.PowerReduction).TruncateInt()
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

	tokensToBurn := slashAmount.Sub(remainingSlashAmount)
	if remainingSlashAmount.Equal(btypes.ZeroInt()) {
		sm.AfterValidatorSlashed(ctx, tokensToBurn)
		return
	}

	// cannot decrease balance below zero
	if remainingSlashAmount.GT(validator.BondTokens) {
		remainingSlashAmount = validator.BondTokens
	}

	// calculate delegations slash fraction
	fraction = qtypes.NewDecFromInt(remainingSlashAmount).QuoInt(validator.BondTokens)

	// slash delegations
	delegationSlashTokens := slashDelegations(ctx, validator, fraction, types.AttributeValueDoubleSign)

	// update validator
	updatedValidatorTokens := validator.BondTokens.Sub(delegationSlashTokens)
	logger.Info("validator gets inactive", "height", ctx.BlockHeight(), "validator", validator.GetValidatorAddress().String())
	sm.MakeValidatorInactive(validator.GetValidatorAddress(), ctx.BlockHeight(), ctx.BlockHeader().Time.UTC(), types.DoubleSign)
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

func slashDelegations(ctx context.Context, validator types.Validator, fraction qtypes.Dec, reason string) btypes.BigInt {
	sm := mapper.GetMapper(ctx)
	logger := ctx.Logger()

	// get delegations
	var delegations []types.DelegationInfo
	sm.IterateDelegationsValDeleAddr(validator.GetValidatorAddress(), func(valAddr btypes.ValAddress, delAddr btypes.AccAddress) {
		if delegation, exists := sm.GetDelegationInfo(delAddr, valAddr); exists {
			delegations = append(delegations, delegation)
		}
	})

	delegationSlashTokens := btypes.ZeroInt()
	for _, delegation := range delegations {
		bondTokens := delegation.Amount

		// calculate slash amount
		amountSlash := qtypes.NewDecFromInt(bondTokens).Mul(fraction).TruncateInt()
		if amountSlash.GT(bondTokens) {
			amountSlash = bondTokens
		}
		delegationSlashTokens = delegationSlashTokens.Add(amountSlash)

		// update delegation
		remainTokens := bondTokens.Sub(amountSlash)
		delegation.Amount = remainTokens
		sm.BeforeDelegationModified(ctx, validator.GetValidatorAddress(), delegation.DelegatorAddr, delegation.Amount)
		sm.SetDelegationInfo(delegation)

		ctx.EventManager().EmitEvent(
			btypes.NewEvent(
				types.EventTypeSlash,
				btypes.NewAttribute(types.AttributeKeyValidator, validator.GetValidatorAddress().String()),
				btypes.NewAttribute(types.AttributeKeyDelegator, delegation.DelegatorAddr.String()),
				btypes.NewAttribute(types.AttributeKeyTokens, amountSlash.String()),
				btypes.NewAttribute(types.AttributeKeyReason, reason),
			),
		)
		logger.Debug("slash validator's delegators", "delegator", delegation.DelegatorAddr.String(), "preToken", bondTokens, "slashToken", amountSlash, "remainTokens", remainTokens)
	}

	// update validator
	sm.ChangeValidatorBondTokens(validator, validator.BondTokens.Sub(delegationSlashTokens))

	return delegationSlashTokens
}
