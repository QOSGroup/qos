package mapper

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	"github.com/tendermint/tendermint/crypto"
)

const (
	BaseMapperName = "basemapper"
	storeKey       = "base"
)

type BaseMapper struct {
	*mapper.BaseMapper
}

func NewBaseMapper() *BaseMapper {
	var baseMapper = BaseMapper{}
	baseMapper.BaseMapper = mapper.NewBaseMapper(store.NewKVStoreKey(storeKey))
	return &baseMapper
}

func (mapper *BaseMapper) Copy() mapper.IMapper {
	cpyMapper := &BaseMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}

func (mapper *BaseMapper) Name() string {
	return BaseMapperName
}

// 保存CA
func (mapper *BaseMapper) SetCA(pubKey crypto.PubKey) error {
	mapper.BaseMapper.Set([]byte("rootca"), pubKey)
	return nil
}

// 获取CA
func (mapper *BaseMapper) GetCA() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte("rootca"), &pubKey)
	return pubKey
}
