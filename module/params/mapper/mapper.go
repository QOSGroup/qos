package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/params/types"
	qtypes "github.com/QOSGroup/qos/types"
	"reflect"
)

const MapperName = "params"

// 参数存储操作
type Mapper struct {
	*mapper.BaseMapper

	// 参数表：模块名-参数
	paramSets map[string]qtypes.ParamSet
}

func (mapper *Mapper) Copy() mapper.IMapper {
	paramMapper := &Mapper{}
	paramMapper.BaseMapper = mapper.BaseMapper.Copy()
	paramMapper.paramSets = mapper.paramSets
	return paramMapper
}

func BuildParamKey(paramSpace string, key []byte) []byte {
	return append([]byte(paramSpace), key...)
}

var _ mapper.IMapper = (*Mapper)(nil)

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(MapperName).(*Mapper)
}

func NewMapper() *Mapper {
	var paramsMapper = Mapper{}
	paramsMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	paramsMapper.paramSets = make(map[string]qtypes.ParamSet)
	return &paramsMapper
}

// 参数校验
func (mapper Mapper) Validate(paramSpace string, key string, value string) btypes.Error {
	module, ok := mapper.paramSets[paramSpace]
	if !ok {
		return types.ErrInvalidParam("unknown module")
	}
	_, err := module.ValidateKeyValue(key, value)
	return err
}

// 参数注册
func (mapper Mapper) RegisterParamSet(ps ...qtypes.ParamSet) {
	for _, ps := range ps {
		if ps != nil {
			// 禁止重复注册
			if _, ok := mapper.paramSets[ps.GetParamSpace()]; ok {
				panic(fmt.Sprintf("<%s> already registered ", ps.GetParamSpace()))
			}
			mapper.paramSets[ps.GetParamSpace()] = ps
		}
	}
}

// 保存参数集
func (mapper Mapper) SetParamSet(params qtypes.ParamSet) {
	for _, pair := range params.KeyValuePairs() {
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		mapper.Set(BuildParamKey(params.GetParamSpace(), pair.Key), v)
	}
}

// 获取参数集
func (mapper Mapper) GetParamSet(params qtypes.ParamSet) {
	for _, pair := range params.KeyValuePairs() {
		mapper.Get(BuildParamKey(params.GetParamSpace(), pair.Key), pair.Value)
	}
}

// 设置单个参数
func (mapper Mapper) SetParam(paramSpace string, key string, value interface{}) {
	mapper.Set(BuildParamKey(paramSpace, []byte(key)), value)
}

// 获取单个参数
func (mapper Mapper) GetParam(paramSpace string, key string) (value interface{}, exists bool) {
	for _, pair := range mapper.paramSets[paramSpace].KeyValuePairs() {
		if key == string(pair.Key) {
			mapper.Get(BuildParamKey(paramSpace, pair.Key), pair.Value)
			return pair.Value, true
		}
	}

	return
}

// 获取模块参数集
func (mapper Mapper) GetModuleParams(module string) (set qtypes.ParamSet, exists bool) {
	set, ok := mapper.paramSets[module]
	if !ok {
		return nil, false
	}
	mapper.GetParamSet(set)
	return set, true
}

// 获取模块参数结构，并非保存在数据库中参数数据
func (mapper Mapper) GetModuleParamSet(module string) (set qtypes.ParamSet, exists bool) {
	set, ok := mapper.paramSets[module]
	if !ok {
		return nil, false
	}
	return set, true
}

// 获取所有参数
func (mapper Mapper) GetParams() (params []qtypes.ParamSet) {
	for _, set := range mapper.paramSets {
		mapper.GetParamSet(set)
		params = append(params, set)
	}

	return params
}
