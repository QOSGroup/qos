package params

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/params/types"
	"reflect"
)

const MapperName = "params"

type Mapper struct {
	*mapper.BaseMapper

	paramSets map[string]types.ParamSet
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
	paramsMapper.paramSets = make(map[string]types.ParamSet)
	return &paramsMapper
}

func (mapper Mapper) Validate(paramSpace string, key string, value string) btypes.Error {
	module, ok := mapper.paramSets[paramSpace]
	if !ok {
		return ErrInvalidParam("unknown module")
	}
	_, err := module.Validate(key, value)
	return err
}

func (mapper Mapper) RegisterParamSet(ps ...types.ParamSet) {
	for _, ps := range ps {
		if ps != nil {
			if _, ok := mapper.paramSets[ps.GetParamSpace()]; ok {
				panic(fmt.Sprintf("<%s> already registered ", ps.GetParamSpace()))
			}
			mapper.paramSets[ps.GetParamSpace()] = ps
		}
	}
}

func (mapper Mapper) SetParamSet(params types.ParamSet) {
	for _, pair := range params.KeyValuePairs() {
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		mapper.Set(BuildParamKey(params.GetParamSpace(), []byte(pair.Key)), v)
	}
}

func (mapper Mapper) GetParamSet(params types.ParamSet) {
	for _, pair := range params.KeyValuePairs() {
		mapper.Get(BuildParamKey(params.GetParamSpace(), pair.Key), pair.Value)
	}
}

func (mapper Mapper) SetParam(paramSpace string, key string, value interface{}) {
	mapper.Set(BuildParamKey(paramSpace, []byte(key)), value)
}

func (mapper Mapper) GetParam(paramSpace string, key string) (value interface{}, exists bool) {
	for _, pair := range mapper.paramSets[paramSpace].KeyValuePairs() {
		if key == string(pair.Key) {
			mapper.Get(BuildParamKey(paramSpace, pair.Key), pair.Value)
			return pair.Value, true
		}
	}

	return
}

func (mapper Mapper) GetModuleParams(module string) (set types.ParamSet, exists bool) {
	set, ok := mapper.paramSets[module]
	if !ok {
		return nil, false
	}
	mapper.GetParamSet(set)
	return set, true
}

func (mapper Mapper) GetModuleParamSet(module string) (set types.ParamSet, exists bool) {
	set, ok := mapper.paramSets[module]
	if !ok {
		return nil, false
	}
	return set, true
}

func (mapper Mapper) GetParams() (params []types.ParamSet) {
	for _, set := range mapper.paramSets {
		mapper.GetParamSet(set)
		params = append(params, set)
	}

	return params
}
