package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	MapperName = "qsc"
	Key        = "qsc/%s"
	RootCAKey  = "rootca"
)

type Mapper struct {
	*mapper.BaseMapper
}

func NewMapper() *Mapper {
	var qscMapper = Mapper{}
	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, MapperName)
	return &qscMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(MapperName).(*Mapper)
}

func BuildQSCKey(qscName string) []byte {
	return []byte(fmt.Sprintf(Key, qscName))
}

func BuildQSCKeyPrefix() []byte {
	return []byte("qsc/")
}

func (mapper *Mapper) Copy() mapper.IMapper {
	qscMapper := &Mapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

// 保存代币信息
func (mapper *Mapper) SaveQsc(qscInfo *types.QSCInfo) {
	mapper.Set(BuildQSCKey(qscInfo.Name), qscInfo)
}

// 是否存在代币信息
func (mapper *Mapper) Exists(qscName string) bool {
	_, exists := mapper.GetQsc(qscName)
	return exists
}

// 根据代币名称获取代币信息
func (mapper *Mapper) GetQsc(qscName string) (info types.QSCInfo, exists bool) {
	exists = mapper.Get(BuildQSCKey(qscName), &info)
	return
}

// 保存 kepler qsc 根证书公钥
func (mapper *Mapper) SetRootCAPubkey(pubKey crypto.PubKey) {
	mapper.BaseMapper.Set([]byte(RootCAKey), pubKey)
}

// 获取 kepler qsc 根证书公钥
func (mapper *Mapper) GetRootCAPubkey() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte(RootCAKey), &pubKey)
	return pubKey
}

// 获取所有代币
func (mapper *Mapper) GetQSCs() []types.QSCInfo {
	qscs := make([]types.QSCInfo, 0)
	mapper.Iterator(BuildQSCKeyPrefix(), func(bz []byte) (stop bool) {
		qsc := types.QSCInfo{}
		mapper.DecodeObject(bz, &qsc)
		qscs = append(qscs, qsc)
		return false
	})

	return qscs
}
