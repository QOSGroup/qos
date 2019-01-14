package mint

import (
	"github.com/QOSGroup/qbase/mapper"
	minttypes "github.com/QOSGroup/qos/module/mint/types"
)

const (
	MintMapperName      = "mint"
	MintParamsKey       = "params"
	AppliedQOSAmountKey = "appliedqos"
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

// 获取已分配QOS总数
func (mapper *MintMapper) GetAppliedQOSAmount() (v uint64) {
	exists := mapper.Get([]byte(AppliedQOSAmountKey), &v)
	if !exists {
		return 0
	}

	return v
}

// 设置 已分配 QOS amount
func (mapper *MintMapper) SetAppliedQOSAmount(amount uint64) {
	mapper.Set([]byte(AppliedQOSAmountKey), amount)
}

// 增加 已分配 QOS amount
func (mapper *MintMapper) AddAppliedQOSAmount(amount uint64) {
	mined := mapper.GetAppliedQOSAmount()
	mined += amount
	mapper.SetAppliedQOSAmount(mined)
}
