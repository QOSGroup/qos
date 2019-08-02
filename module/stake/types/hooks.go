package types

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
)

type Hooks interface {
	AfterValidatorCreated(ctx context.Context, val btypes.Address)

	BeforeValidatorRemoved(ctx context.Context, val btypes.Address)

	AfterDelegationCreated(ctx context.Context, val btypes.Address, del btypes.Address)
	BeforeDelegationModified(ctx context.Context, val btypes.Address, del btypes.Address, updateTokes uint64)

	AfterValidatorSlashed(ctx context.Context, slashedTokes uint64)
}
