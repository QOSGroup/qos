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

type MainMapper struct {
	*mapper.BaseMapper
}

func NewMainMapper() *MainMapper {
	var baseMapper = MainMapper{}
	baseMapper.BaseMapper = mapper.NewBaseMapper(store.NewKVStoreKey(storeKey))
	return &baseMapper
}

func GetStoreKey() string {
	return storeKey
}

func (mapper *MainMapper) Copy() mapper.IMapper {
	cpyMapper := &MainMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
}

func (mapper *MainMapper) Name() string {
	return BaseMapperName
}

// 保存CA
func (mapper *MainMapper) SetRootCA(pubKey crypto.PubKey) error {
	mapper.BaseMapper.Set([]byte("rootca"), pubKey)
	return nil
}

// 获取CA
func (mapper *MainMapper) GetRoot() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte("rootca"), &pubKey)
	return pubKey
}
