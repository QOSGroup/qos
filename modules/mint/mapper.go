package mint

import (
	"github.com/QOSGroup/qbase/mapper"
	minttypes "github.com/QOSGroup/qos/modules/mint/types"
)

const (
	MintMapperName = "mint"
	MintParamsKey  = "params"
)

type MintMapper struct {
	*mapper.BaseMapper
}

func BuildMintParamsKey() []byte {
	return []byte(MintParamsKey)
}

func NewMintMapper() *MintMapper {
	var qscMapper = MintMapper{}
	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, MintMapperName)
	return &qscMapper
}

func (mapper *MintMapper) Copy() mapper.IMapper {
	qscMapper := &MintMapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

func (mapper *MintMapper) SetParams(params minttypes.Params) {
	mapper.Set(BuildMintParamsKey(), params)
}

func (mapper *MintMapper) GetParams() minttypes.Params {
	params := minttypes.Params{}
	mapper.Get(BuildMintParamsKey(), &params)
	return params
}
