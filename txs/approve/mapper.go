package approve

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
)

const (
	approveStoreKey = "approve"
	approveKey      = "from:[%s]/to:[%s]"
)

type ApproveMapper struct {
	*mapper.BaseMapper
}

func NewApproveMapper() *ApproveMapper {
	var approveMapper = ApproveMapper{}
	approveMapper.BaseMapper = mapper.NewBaseMapper(nil, approveStoreKey)
	return &approveMapper
}

func GetApproveMapperStoreKey() string {
	return approveStoreKey
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
func (mapper *ApproveMapper) GetApprove(from btypes.Address, to btypes.Address) (Approve, bool) {
	approve := NewApprove(from, to, btypes.ZeroInt(), nil)
	key := BuildApproveKey(from.String(), to.String())
	exists := mapper.BaseMapper.Get(key, &approve)
	return approve, exists
}

// 保存授权
func (mapper *ApproveMapper) SaveApprove(approve Approve) error {
	key := BuildApproveKey(approve.From.String(), approve.To.String())
	mapper.BaseMapper.Set(key, approve)
	return nil
}

// 删除授权
func (mapper *ApproveMapper) DeleteApprove(from btypes.Address, to btypes.Address) error {
	key := BuildApproveKey(from.String(), to.String())
	mapper.BaseMapper.Del(key)
	return nil
}
