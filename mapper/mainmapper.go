package mapper

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/tendermint/tendermint/crypto"
)

const (
	BaseMapperName = "base"
	RootCAKey      = "rootca"
)

type MainMapper struct {
	*mapper.BaseMapper
}

func NewMainMapper() *MainMapper {
	var baseMapper = MainMapper{}
	baseMapper.BaseMapper = mapper.NewBaseMapper(nil, BaseMapperName)
	return &baseMapper
}

func GetMainStoreKey() string {
	return BaseMapperName
}

func (mapper *MainMapper) MapperName() string {
	return BaseMapperName
}

func (mapper *MainMapper) Copy() mapper.IMapper {
	cpyMapper := &MainMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}

// 保存CA
func (mapper *MainMapper) SetRootCA(pubKey crypto.PubKey) error {
	mapper.BaseMapper.Set([]byte(RootCAKey), pubKey)
	return nil
}

// 获取CA
func (mapper *MainMapper) GetRootCA() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte(RootCAKey), &pubKey)
	return pubKey
}
