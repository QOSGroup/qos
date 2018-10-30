package mapper

import (
	"fmt"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	BaseMapperName = "base"
	QSCName        = "qsc/[%s]"
)

type MainMapper struct {
	*mapper.BaseMapper
}

type QscInfo struct {
	Qscname     string         `json:"qscname"`
	BankAddr    btypes.Address `json:"pubkeybank"`
	CreateAddr  btypes.Address `json:"createaddr"`
	QscPubkey   crypto.PubKey  `json:"qscpubkey"`
	Extrate     string         `json:"extrate"`
	CAqsc       []byte         `json:"caqsc"`
	CAbanker    []byte         `json:"cabanker"`
	Description string         `json:"description"`
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

func (mapper *MainMapper) GetQsc(qscname string) (qscinfo *QscInfo) {
	key := fmt.Sprintf(QSCName, qscname)

	var qinfo QscInfo
	exist := mapper.Get([]byte(key), &qinfo)
	if !exist {
		return nil
	}

	return &qinfo
}

func (mapper *MainMapper) SetQsc(qscname string, qscinfo *QscInfo) bool {
	key := fmt.Sprintf(QSCName, qscname)
	mapper.Set([]byte(key), qscinfo)

	return true
}
