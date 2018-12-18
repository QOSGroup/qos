package mapper

import (
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	BaseMapperName    = "base"
	RootCAKey         = "rootca"
	MinedQOSAmountKey = "minedqos"
	SPOConfigKey      = "spo"
	StakeConfigKey    = "stake"
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
func (mapper *MainMapper) SetRootCA(pubKey crypto.PubKey) {
	mapper.BaseMapper.Set([]byte(RootCAKey), pubKey)
}

// 获取CA
func (mapper *MainMapper) GetRootCA() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte(RootCAKey), &pubKey)
	return pubKey
}

// 设置SPOConfig
func (mapper *MainMapper) SetSPOConfig(config types.SPOConfig) {
	mapper.BaseMapper.Set([]byte(SPOConfigKey), config)
}

// 获取SPOConfig
func (mapper *MainMapper) GetSPOConfig() types.SPOConfig {
	config := types.SPOConfig{}
	mapper.BaseMapper.Get([]byte(SPOConfigKey), &config)
	return config
}

// 设置StakeConfig
func (mapper *MainMapper) SetStakeConfig(config types.StakeConfig) {
	mapper.BaseMapper.Set([]byte(StakeConfigKey), config)
}

// 获取StakeConfig
func (mapper *MainMapper) GetStakeConfig() types.StakeConfig {
	config := types.StakeConfig{}
	mapper.BaseMapper.Get([]byte(StakeConfigKey), &config)
	return config
}

// 获取 mined QOS amount
func (mapper *MainMapper) GetMinedQOS() (v uint64) {
	exists := mapper.Get([]byte(MinedQOSAmountKey), &v)
	if !exists {
		return 0
	}

	return v
}

// 设置 mined QOS amount
func (mapper *MainMapper) SetMinedQOS(amount uint64) {
	mapper.Set([]byte(MinedQOSAmountKey), amount)
}

// 增加 mined QOS amount
func (mapper *MainMapper) AddMinedQOS(amount uint64) {
	mined := mapper.GetMinedQOS()
	mined += amount
	mapper.SetMinedQOS(mined)
}
