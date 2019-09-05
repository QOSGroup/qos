package types

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
)

type Hooks interface {
	AfterValidatorCreated(ctx context.Context, val btypes.ValAddress)

	BeforeValidatorRemoved(ctx context.Context, val btypes.ValAddress)

	AfterDelegationCreated(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress)
	BeforeDelegationModified(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress, updateTokes btypes.BigInt)

	AfterValidatorSlashed(ctx context.Context, slashedTokes btypes.BigInt)
}
