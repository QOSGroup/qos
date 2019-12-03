package types

import (
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
)

// 为解决gov和stake循环依赖问题，引入Hooks接口，stake mapper初始化时注入hooks distribution 实现
type Hooks interface {
	// 创建完验证节点后执行操作
	AfterValidatorCreated(ctx context.Context, val btypes.ValAddress)

	// 验证节点删除之前操作
	BeforeValidatorRemoved(ctx context.Context, val btypes.ValAddress)

	// 委托创建之后操作
	AfterDelegationCreated(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress)

	// 修改委托之前操作
	BeforeDelegationModified(ctx context.Context, val btypes.ValAddress, del btypes.AccAddress, updateTokes btypes.BigInt)

	// 验证节点惩罚之后操作
	AfterValidatorSlashed(ctx context.Context, slashedTokes btypes.BigInt)
}
