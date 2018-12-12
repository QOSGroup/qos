package qsc

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
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

func (mapper *QSCMapper) SaveQsc(qscInfo *QSCInfo) {
	mapper.Set(BuildQSCKey(qscInfo.QSCCA.CSR.Subj.CN), qscInfo)
}

func (mapper *QSCMapper) Exists(qscName string) bool {
	return nil != mapper.GetQsc(qscName)
}

func (mapper *QSCMapper) GetQsc(qscName string) (qscinfo *QSCInfo) {
	var info QSCInfo
	exist := mapper.Get(BuildQSCKey(qscName), &info)
	if !exist {
		return nil
	}

	return &info
}
