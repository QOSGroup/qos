package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/mint/types"
)

type Mapper struct {
	*mapper.BaseMapper
}

func NewMapper() *Mapper {
	var qscMapper = Mapper{}

	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, types.MapperName)
	return &qscMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(types.MapperName).(*Mapper)
}

func (mapper *Mapper) Copy() mapper.IMapper {
	qscMapper := &Mapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

func (mapper *Mapper) SetInflationPhrases(phrases types.InflationPhrases) {
	mapper.Set(types.BuildInflationPhrasesKey(), phrases)
}

func (mapper *Mapper) GetInflationPhrases() (phrases types.InflationPhrases, exists bool) {
	phrases = types.InflationPhrases{}
	exists = mapper.Get(types.BuildInflationPhrasesKey(), &phrases)
	return
}

func (mapper *Mapper) MustGetInflationPhrases() types.InflationPhrases {
	phrases := types.InflationPhrases{}
	exists := mapper.Get(types.BuildInflationPhrasesKey(), &phrases)
	if !exists {
		panic("inflation phrases not exists")
	}
	return phrases
}

func (mapper *Mapper) SetFirstBlockTime(t int64) {
	mapper.Set(types.BuildFirstBlockTimeKey(), t)
}

func (mapper *Mapper) GetFirstBlockTime() (t int64) {
	mapper.Get(types.BuildFirstBlockTimeKey(), &t)
	return
}

//获取总分配的QOS总数
func (mapper *Mapper) GetAllTotalMintQOSAmount() (amount btypes.BigInt) {
	mapper.Get(types.BuildAllTotalMintQOSKey(), &amount)
	return
}

func (mapper *Mapper) DelAllTotalMintQOSAmount() {
	mapper.Del(types.BuildAllTotalMintQOSKey())
}

//设置总分配的QOS总数
func (mapper *Mapper) SetAllTotalMintQOSAmount(amount btypes.BigInt) {
	mapper.Set(types.BuildAllTotalMintQOSKey(), amount)
}

//增加总分配的QOS总数
func (mapper *Mapper) AddAllTotalMintQOSAmount(amount btypes.BigInt) {

	totalAmount := mapper.GetAllTotalMintQOSAmount()
	totalAmount = totalAmount.Add(amount)

	mapper.SetAllTotalMintQOSAmount(totalAmount)
}

//设置QOS发行总量
func (mapper *Mapper) SetTotalQOSAmount(amount btypes.BigInt) {
	mapper.Set(types.BuildTotalQOSKey(), amount)
	return
}

//获取QOS发行总量（已发行+待发行）
func (mapper *Mapper) GetTotalQOSAmount() (amount btypes.BigInt) {
	mapper.Get(types.BuildTotalQOSKey(), &amount)
	return
}
