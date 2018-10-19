package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
)

const (
	ApproveMapperName = "approvemapper"
	approveStoreKey   = "approve"
	approveKey        = "from:[%s]/to:[%s]"
)

type ApproveMapper struct {
	*mapper.BaseMapper
}

func NewApproveMapper() *ApproveMapper {
	var approveMapper = ApproveMapper{}
	approveMapper.BaseMapper = mapper.NewBaseMapper(store.NewKVStoreKey(approveStoreKey))
	return &approveMapper
}

func (mapper *ApproveMapper) Copy() mapper.IMapper {
	approveMapper := &ApproveMapper{}
	approveMapper.BaseMapper = mapper.BaseMapper.Copy()
	return approveMapper
}

func (mapper *ApproveMapper) Name() string {
	return ApproveMapperName
}

// 获取授权
func (mapper *ApproveMapper) GetApprove(from btypes.Address, to btypes.Address) (approve types.Approve, exists bool) {
	key := fmt.Sprintf(approveKey, from.String(), to.String())
	exists = mapper.BaseMapper.Get([]byte(key), &approve)
	return approve, exists
}

// 保存授权
func (mapper *ApproveMapper) SaveApprove(approve *types.Approve) error {
	key := fmt.Sprintf(approveKey, approve.From.String(), approve.To.String())
	mapper.BaseMapper.Set([]byte(key), approve)
	return nil
}

// 删除授权
func (mapper *ApproveMapper) DeleteApprove(approve *types.ApproveCancel) error {
	key := fmt.Sprintf(approveKey, approve.From.String(), approve.To.String())
	mapper.BaseMapper.GetStore().Delete([]byte(key))
	return nil
}
