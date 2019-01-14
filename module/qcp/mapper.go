package qcp

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/tendermint/tendermint/crypto"
)

const (
	QCPRootCAKey = "rootca"
)

// 保存CA
func SetQCPRootCA(ctx context.Context, pubKey crypto.PubKey) {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpMapper.Set([]byte(QCPRootCAKey), pubKey)
}

// 获取CA
func GetQCPRootCA(ctx context.Context) crypto.PubKey {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	var pubKey crypto.PubKey
	qcpMapper.Get([]byte(QCPRootCAKey), &pubKey)
	return pubKey
}
