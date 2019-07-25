package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/types"
)

const (
	MapperName = "approve"
	approveKey = "from:[%s]/to:[%s]"
)

type Mapper struct {
	*mapper.BaseMapper
}

func NewApproveMapper() *Mapper {
	var approveMapper = Mapper{}
	approveMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &approveMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(MapperName).(*Mapper)
}

func BuildApproveKey(from string, to string) []byte {
	key := fmt.Sprintf(approveKey, from, to)
	return []byte(key)
}

func (mapper *Mapper) Copy() mapper.IMapper {
	approveMapper := &Mapper{}
	approveMapper.BaseMapper = mapper.BaseMapper.Copy()
	return approveMapper
}

// 获取授权
func (mapper *Mapper) GetApprove(from btypes.Address, to btypes.Address) (types.Approve, bool) {
	approve := types.Approve{}
	key := BuildApproveKey(from.String(), to.String())
	exists := mapper.BaseMapper.Get(key, &approve)
	return approve, exists
}

// 保存授权
func (mapper *Mapper) SaveApprove(approve types.Approve) {
	key := BuildApproveKey(approve.From.String(), approve.To.String())
	mapper.BaseMapper.Set(key, approve)
}

// 删除授权
func (mapper *Mapper) DeleteApprove(from btypes.Address, to btypes.Address) {
	key := BuildApproveKey(from.String(), to.String())
	mapper.BaseMapper.Del(key)
}

// 所有预授权
func (mapper *Mapper) GetApproves() []types.Approve {
	approves := make([]types.Approve, 0)
	mapper.Iterator([]byte("from:"), func(bz []byte) (stop bool) {
		approve := types.Approve{}
		mapper.DecodeObject(bz, &approve)
		if approve.IsPositive() {
			approves = append(approves, approve)
		}
		return false
	})

	return approves
}
