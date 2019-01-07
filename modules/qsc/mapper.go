package qsc

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qos/types"
)

const (
	QSCMapperName = "qsc"
	QSCKey        = "qsc/[%s]"
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
