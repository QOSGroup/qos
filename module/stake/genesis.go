package stake

import (
	"fmt"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	validatorMapper := mapper.GetMapper(ctx)

	if len(data.CurrentValidators) > 0 {
		validatorMapper.Set(types.BuildCurrentValidatorsAddressKey(), data.CurrentValidators)
	}

	initValidators(ctx, data.Validators)
	initParams(ctx, data.Params)
	initValidatorsVotesInfo(ctx, data.ValidatorsVoteInfo, data.ValidatorsVoteInWindow)
	initDelegatorsInfo(ctx, data.DelegatorsInfo, data.DelegatorsUnbondInfo)
}

func initValidators(ctx context.Context, validators []types.Validator) {
	validatorMapper := mapper.GetMapper(ctx)
	for _, v := range validators {

		if validatorMapper.Exists(v.ValidatorPubKey.Address().Bytes()) {
			panic(fmt.Errorf("validator %s already exists", v.ValidatorPubKey.Address()))
		}
		if validatorMapper.ExistsWithOwner(v.Owner) {
			panic(fmt.Errorf("owner %s already bind a validator", v.Owner))
		}
		validatorMapper.CreateValidator(v)
		if !v.IsActive() {
			validatorMapper.MakeValidatorInactive(v.GetValidatorAddress(), v.InactiveHeight, v.InactiveTime, v.InactiveCode)
		}
	}
}

func initValidatorsVotesInfo(ctx context.Context, voteInfos []types.ValidatorVoteInfoState, voteWindowInfos []types.ValidatorVoteInWindowInfoState) {
	sm := mapper.GetMapper(ctx)
	for _, voteInfo := range voteInfos {
		sm.SetValidatorVoteInfo(btypes.Address(voteInfo.ValidatorPubKey.Address()), voteInfo.VoteInfo)
	}

	for _, voteWindowInfo := range voteWindowInfos {
		sm.SetVoteInfoInWindow(btypes.Address(voteWindowInfo.ValidatorPubKey.Address()), voteWindowInfo.Index, voteWindowInfo.Vote)
	}
}

func initDelegatorsInfo(ctx context.Context, delegatorsInfo []types.DelegationInfoState, delegatorsUnbondInfo []types.UnbondingDelegationInfo) {
	sm := mapper.GetMapper(ctx)

	for _, info := range delegatorsInfo {
		sm.SetDelegationInfo(types.DelegationInfo{
			DelegatorAddr: info.DelegatorAddr,
			ValidatorAddr: btypes.Address(info.ValidatorPubKey.Address()),
			Amount:        info.Amount,
			IsCompound:    info.IsCompound,
		})
	}

	for _, info := range delegatorsUnbondInfo {
		sm.AddUnbondingDelegations(info.Height, []types.UnbondingDelegationInfo{info})
	}
}

func initParams(ctx context.Context, params types.Params) {
	mapper := ctx.Mapper(types.MapperName).(*mapper.Mapper)
	mapper.SetParams(ctx, params)
}

func ExportGenesis(ctx context.Context) types.GenesisState {

	validatorMapper := mapper.GetMapper(ctx)
	sm := mapper.GetMapper(ctx)

	var currentValidators []types.Validator
	validatorMapper.Get(types.BuildCurrentValidatorsAddressKey(), &currentValidators)

	params := validatorMapper.GetParams(ctx)

	var validators []types.Validator
	validatorMapper.IterateValidators(func(validator types.Validator) {
		validators = append(validators, validator)
	})

	var validatorsVoteInfo []types.ValidatorVoteInfoState
	sm.IterateVoteInfos(func(valAddr btypes.Address, info types.ValidatorVoteInfo) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			vvis := ValidatorVoteInfoState{
				ValidatorPubKey: validator.GetValidatorPubKey(),
				VoteInfo:        info,
			}
			validatorsVoteInfo = append(validatorsVoteInfo, vvis)
		}
	})

	var validatorsVoteInWindow []types.ValidatorVoteInWindowInfoState
	sm.IterateVoteInWindowsInfos(func(index uint64, valAddr btypes.Address, vote bool) {

		validator, exists := validatorMapper.GetValidator(valAddr)
		if exists {
			validatorsVoteInWindow = append(validatorsVoteInWindow, ValidatorVoteInWindowInfoState{
				ValidatorPubKey: validator.GetValidatorPubKey(),
				Index:           index,
				Vote:            vote,
			})
		}
	})

	var delegatorsInfo []types.DelegationInfoState
	sm.IterateDelegationsInfo(btypes.Address{}, func(info types.DelegationInfo) {

		validator, exists := validatorMapper.GetValidator(info.ValidatorAddr)
		if !exists {
			panic(fmt.Sprintf("validator:%s not exists", info.ValidatorAddr.String()))
		}

		delegatorsInfo = append(delegatorsInfo, DelegationInfoState{
			DelegatorAddr:   info.DelegatorAddr,
			ValidatorPubKey: validator.GetValidatorPubKey(),
			Amount:          info.Amount,
			IsCompound:      info.IsCompound,
		})
	})

	var delegatorsUnbondInfo []types.UnbondingDelegationInfo
	sm.IterateUnbondingDelegations(func(deleAddr btypes.Address, height uint64, unbondings []types.UnbondingDelegationInfo) {
		delegatorsUnbondInfo = append(delegatorsUnbondInfo, unbondings...)
	})

	var reDelegationsInfo []types.RedelegationInfo
	sm.IterateRedelegationsInfo(func(deleAddr btypes.Address, height uint64, reDelegations []types.RedelegationInfo) {
		reDelegationsInfo = append(reDelegationsInfo, reDelegations...)
	})

	return GenesisState{
		Params:                 params,
		Validators:             validators,
		ValidatorsVoteInfo:     validatorsVoteInfo,
		ValidatorsVoteInWindow: validatorsVoteInWindow,
		DelegatorsInfo:         delegatorsInfo,
		DelegatorsUnbondInfo:   delegatorsUnbondInfo,
		ReDelegationsInfo:      reDelegationsInfo,
		CurrentValidators:      currentValidators,
	}
}
