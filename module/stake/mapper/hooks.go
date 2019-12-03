package mapper

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/types"
)

var _ types.Hooks = (*Mapper)(nil)

func (mapper *Mapper) AfterValidatorCreated(ctx context.Context, val btypes.ValAddress) {
	if mapper.hooks != nil {
		mapper.hooks.AfterValidatorCreated(ctx, val)
	}
}

func (mapper *Mapper) BeforeValidatorRemoved(ctx context.Context, val btypes.ValAddress) {
	if mapper.hooks != nil {
		mapper.hooks.BeforeValidatorRemoved(ctx, val)
	}
}

func (mapper *Mapper) AfterDelegationCreated(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress) {
	if mapper.hooks != nil {
		mapper.hooks.AfterDelegationCreated(ctx, val, del)
	}
}

func (mapper *Mapper) BeforeDelegationModified(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress, updateTokes btypes.BigInt) {
	if mapper.hooks != nil {
		mapper.hooks.BeforeDelegationModified(ctx, val, del, updateTokes)
	}
}

func (mapper *Mapper) AfterValidatorSlashed(ctx context.Context, slashedTokes btypes.BigInt) {
	if mapper.hooks != nil {
		mapper.hooks.AfterValidatorSlashed(ctx, slashedTokes)
	}
}
