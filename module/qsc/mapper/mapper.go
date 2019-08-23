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
	Key        = "qsc/[%s]"
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

func (mapper *Mapper) SaveQsc(qscInfo *types.Info) {
	mapper.Set(BuildQSCKey(qscInfo.Name), qscInfo)
}

func (mapper *Mapper) Exists(qscName string) bool {
	return nil != mapper.GetQsc(qscName)
}

func (mapper *Mapper) GetQsc(qscName string) (qscinfo *types.Info) {
	var info types.Info
	exist := mapper.Get(BuildQSCKey(qscName), &info)
	if !exist {
		return nil
	}

	return &info
}

// 保存CA
func (mapper *Mapper) SetQSCRootCA(pubKey crypto.PubKey) {
	mapper.BaseMapper.Set([]byte(RootCAKey), pubKey)
}

// 获取CA
func (mapper *Mapper) GetQSCRootCA() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte(RootCAKey), &pubKey)
	return pubKey
}

func (mapper *Mapper) GetQSCs() []types.Info {
	qscs := make([]types.Info, 0)
	mapper.Iterator(BuildQSCKeyPrefix(), func(bz []byte) (stop bool) {
		qsc := types.Info{}
		mapper.DecodeObject(bz, &qsc)
		qscs = append(qscs, qsc)
		return false
	})

	return qscs
}
