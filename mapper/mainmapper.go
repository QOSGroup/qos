package mapper

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/tendermint/tendermint/crypto"
)

const (
	storeKey = "base"
)

type MainMapper struct {
	*mapper.BaseMapper
}

func NewMainMapper() *MainMapper {
	var baseMapper = MainMapper{}
	baseMapper.BaseMapper = mapper.NewBaseMapper(nil, storeKey)
	return &baseMapper
}

func GetMainStoreKey() string {
	return storeKey
}

func (mapper *MainMapper) Copy() mapper.IMapper {
	cpyMapper := &MainMapper{}
	cpyMapper.BaseMapper = mapper.BaseMapper.Copy()
	return cpyMapper
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
