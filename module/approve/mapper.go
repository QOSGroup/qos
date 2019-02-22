package approve

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/types"
)

const (
	ApproveMapperName = "approve"
	approveKey        = "from:[%s]/to:[%s]"
)

type ApproveMapper struct {
	*mapper.BaseMapper
}

func NewApproveMapper() *ApproveMapper {
	var approveMapper = ApproveMapper{}
	approveMapper.BaseMapper = mapper.NewBaseMapper(nil, ApproveMapperName)
	return &approveMapper
}

func BuildApproveKey(from string, to string) []byte {
	key := fmt.Sprintf(approveKey, from, to)
	return []byte(key)
}

func (mapper *ApproveMapper) Copy() mapper.IMapper {
	approveMapper := &ApproveMapper{}
	approveMapper.BaseMapper = mapper.BaseMapper.Copy()
	return approveMapper
}

// 获取授权
func (mapper *ApproveMapper) GetApprove(from btypes.Address, to btypes.Address) (types.Approve, bool) {
	approve := types.Approve{}
	key := BuildApproveKey(from.String(), to.String())
	exists := mapper.BaseMapper.Get(key, &approve)
	return approve, exists
}

// 保存授权
func (mapper *ApproveMapper) SaveApprove(approve types.Approve) {
	key := BuildApproveKey(approve.From.String(), approve.To.String())
	mapper.BaseMapper.Set(key, approve)
}

// 删除授权
func (mapper *ApproveMapper) DeleteApprove(from btypes.Address, to btypes.Address) {
	key := BuildApproveKey(from.String(), to.String())
	mapper.BaseMapper.Del(key)
}

// 所有预授权
func (mapper *ApproveMapper) GetApproves() []types.Approve {
	approves := make([]types.Approve, 0)
	mapper.Iterator([]byte("from:"), func(bz []byte) (stop bool) {
		approve := types.Approve{}
		mapper.DecodeObject(bz, &approve)
		approves = append(approves, approve)
		return false
	})

	return approves
}
