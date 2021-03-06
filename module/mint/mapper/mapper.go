package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/mint/types"
	"github.com/tendermint/tendermint/config"
)

// 通胀相关存储操作Mapper
type Mapper struct {
	*mapper.BaseMapper

	Metrics *Metrics // prometheus监控项
}

func NewMapper() *Mapper {
	var qscMapper = Mapper{}

	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, types.MapperName)
	return &qscMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(types.MapperName).(*Mapper)
}

// 设置prometheus监控项
func (mapper *Mapper) SetUpMetrics(cfg *config.InstrumentationConfig) {
	mapper.Metrics = PrometheusMetrics(cfg)
}

func (mapper *Mapper) Copy() mapper.IMapper {
	qscMapper := &Mapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	qscMapper.Metrics = mapper.Metrics
	return qscMapper
}

// 保存通胀规则
func (mapper *Mapper) SetInflationPhrases(phrases types.InflationPhrases) {
	mapper.Set(types.BuildInflationPhrasesKey(), phrases)
}

// 获取通胀规则
func (mapper *Mapper) GetInflationPhrases() (phrases types.InflationPhrases, exists bool) {
	phrases = types.InflationPhrases{}
	exists = mapper.Get(types.BuildInflationPhrasesKey(), &phrases)
	return
}

// 获取通胀规则，不存在时panic
func (mapper *Mapper) MustGetInflationPhrases() types.InflationPhrases {
	phrases := types.InflationPhrases{}
	exists := mapper.Get(types.BuildInflationPhrasesKey(), &phrases)
	if !exists {
		panic("inflation phrases not exists")
	}
	return phrases
}

// 设置第一块时间
func (mapper *Mapper) SetFirstBlockTime(t int64) {
	mapper.Set(types.BuildFirstBlockTimeKey(), t)
}

// 获取第一块时间
func (mapper *Mapper) GetFirstBlockTime() (t int64) {
	mapper.Get(types.BuildFirstBlockTimeKey(), &t)
	return
}

//获取流通QOS总数
func (mapper *Mapper) GetAllTotalMintQOSAmount() (amount btypes.BigInt) {
	mapper.Get(types.BuildAllTotalMintQOSKey(), &amount)
	return
}

// 删除流通QOS总数
func (mapper *Mapper) DelAllTotalMintQOSAmount() {
	mapper.Del(types.BuildAllTotalMintQOSKey())
}

// 设置流通QOS总数
func (mapper *Mapper) SetAllTotalMintQOSAmount(amount btypes.BigInt) {
	mapper.Set(types.BuildAllTotalMintQOSKey(), amount)
}

// 增加流通QOS总数
func (mapper *Mapper) AddAllTotalMintQOSAmount(amount btypes.BigInt) {

	totalAmount := mapper.GetAllTotalMintQOSAmount()
	totalAmount = totalAmount.Add(amount)

	mapper.SetAllTotalMintQOSAmount(totalAmount)
}

// 设置QOS发行总量
func (mapper *Mapper) SetTotalQOSAmount(amount btypes.BigInt) {
	mapper.Set(types.BuildTotalQOSKey(), amount)
	return
}

// 获取QOS发行总量（已发行+待发行）
func (mapper *Mapper) GetTotalQOSAmount() (amount btypes.BigInt) {
	mapper.Get(types.BuildTotalQOSKey(), &amount)
	return
}
