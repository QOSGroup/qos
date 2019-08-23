package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve/types"
)

type Mapper struct {
	*mapper.BaseMapper
}

func NewApproveMapper() *Mapper {
	var approveMapper = Mapper{}
	approveMapper.BaseMapper = mapper.NewBaseMapper(nil, types.MapperName)
	return &approveMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(types.MapperName).(*Mapper)
}

func (mapper *Mapper) Copy() mapper.IMapper {
	approveMapper := &Mapper{}
	approveMapper.BaseMapper = mapper.BaseMapper.Copy()
	return approveMapper
}

// 获取授权
func (mapper *Mapper) GetApprove(from btypes.Address, to btypes.Address) (types.Approve, bool) {
	approve := types.Approve{}
	key := types.BuildApproveKey(from, to)
	exists := mapper.BaseMapper.Get(key, &approve)
	return approve, exists
}

// 保存授权
func (mapper *Mapper) SaveApprove(approve types.Approve) {
	key := types.BuildApproveKey(approve.From, approve.To)
	mapper.BaseMapper.Set(key, approve)
}

// 删除授权
func (mapper *Mapper) DeleteApprove(from btypes.Address, to btypes.Address) {
	key := types.BuildApproveKey(from, to)
	mapper.BaseMapper.Del(key)
}

// 所有预授权
func (mapper *Mapper) GetApproves() []types.Approve {
	approves := make([]types.Approve, 0)
	mapper.Iterator(types.GetApprovePrefixKey(), func(bz []byte) (stop bool) {
		approve := types.Approve{}
		mapper.DecodeObject(bz, &approve)
		approves = append(approves, approve)
		return false
	})

	return approves
}
