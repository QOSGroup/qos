package qsc

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	QSCMapperName = "qsc"
	QSCKey        = "qsc/[%s]"
	QSCRootCAKey  = "rootca"
)

type QSCMapper struct {
	*mapper.BaseMapper
}

func NewQSCMapper() *QSCMapper {
	var qscMapper = QSCMapper{}
	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, QSCMapperName)
	return &qscMapper
}

func BuildQSCKey(qscName string) []byte {
	return []byte(fmt.Sprintf(QSCKey, qscName))
}

func BuildQSCKeyPrefix() []byte {
	return []byte("qsc/")
}

func (mapper *QSCMapper) Copy() mapper.IMapper {
	qscMapper := &QSCMapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

func (mapper *QSCMapper) SaveQsc(qscInfo *types.QSCInfo) {
	mapper.Set(BuildQSCKey(qscInfo.Name), qscInfo)
}

func (mapper *QSCMapper) Exists(qscName string) bool {
	return nil != mapper.GetQsc(qscName)
}

func (mapper *QSCMapper) GetQsc(qscName string) (qscinfo *types.QSCInfo) {
	var info types.QSCInfo
	exist := mapper.Get(BuildQSCKey(qscName), &info)
	if !exist {
		return nil
	}

	return &info
}

// 保存CA
func (mapper *QSCMapper) SetQSCRootCA(pubKey crypto.PubKey) {
	mapper.BaseMapper.Set([]byte(QSCRootCAKey), pubKey)
}

// 获取CA
func (mapper *QSCMapper) GetQSCRootCA() crypto.PubKey {
	var pubKey crypto.PubKey
	mapper.BaseMapper.Get([]byte(QSCRootCAKey), &pubKey)
	return pubKey
}

func (mapper *QSCMapper) GetQSCs() []types.QSCInfo {
	qscs := make([]types.QSCInfo, 0)
	mapper.Iterator(BuildQSCKeyPrefix(), func(bz []byte) (stop bool) {
		qsc := types.QSCInfo{}
		mapper.DecodeObject(bz, &qsc)
		qscs = append(qscs, qsc)
		return false
	})

	return qscs
}
