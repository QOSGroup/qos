package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/types"
	"github.com/tendermint/tendermint/config"
)

const (
	MapperName = "guardian"
)

var (
	guardianKey = []byte{0x00} // 系统账户存储前缀
	haltKey     = []byte{0x01} // 停网标志存储键值
)

// 系统账户存储键
func KeyGuardian(address btypes.AccAddress) []byte {
	return append(guardianKey, address...)
}

// 系统账户存储前缀
func KeyGuardiansSubspace() []byte {
	return guardianKey
}

// 系统账户模块数据库操作
type Mapper struct {
	*mapper.BaseMapper

	Metrics *Metrics
}

func (mapper *Mapper) Copy() mapper.IMapper {
	govMapper := &Mapper{}
	govMapper.BaseMapper = mapper.BaseMapper.Copy()
	govMapper.Metrics = mapper.Metrics
	return govMapper
}

var _ mapper.IMapper = (*Mapper)(nil)

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(MapperName).(*Mapper)
}

func NewMapper() *Mapper {
	var guardianMapper = Mapper{}
	guardianMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &guardianMapper
}

// 设置prometheus监控项
func (mapper *Mapper) SetUpMetrics(cfg *config.InstrumentationConfig) {
	mapper.Metrics = PrometheusMetrics(cfg)
}

// 添加系统账户
func (mapper Mapper) AddGuardian(guardian types.Guardian) {
	mapper.Set(KeyGuardian(guardian.Address), guardian)
}

// 删除系统账户
func (mapper Mapper) DeleteGuardian(address btypes.AccAddress) {
	mapper.Del(KeyGuardian(address))
}

// 获取系统账户
func (mapper Mapper) GetGuardian(address btypes.AccAddress) (guardian types.Guardian, exists bool) {
	exists = mapper.Get(KeyGuardian(address), &guardian)
	return guardian, exists
}

// 系统账户迭代器
func (mapper Mapper) GuardiansIterator() store.Iterator {
	return btypes.KVStorePrefixIterator(mapper.GetStore(), KeyGuardiansSubspace())
}

// 设置停网标识
func (mapper Mapper) SetHalt(reason string) {
	mapper.Set(haltKey, reason)
}

// 查询停止网络标志
func (mapper Mapper) GetHalt() (string, bool) {
	var reason string
	exists := mapper.Get(haltKey, &reason)
	return reason, exists
}
