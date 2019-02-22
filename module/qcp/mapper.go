package qcp

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	qcptypes "github.com/QOSGroup/qos/module/qcp/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	QCPRootCAKey = "rootca"

	QCPTxExportLimit = 100 // QCP TX 导出条数限制
)

func GetQCPMapper(ctx context.Context) *qcp.QcpMapper {
	return ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
}

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

// TODO prefix定义到qbase中
func GetQCPTxs(ctx context.Context, chainId string) []txs.TxQcp {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpTxs := make([]txs.TxQcp, 0)
	qcpMapper.Iterator([]byte("tx/out/"+chainId), func(bz []byte) (stop bool) {
		tx := txs.TxQcp{}
		qcpMapper.DecodeObject(bz, &tx)
		qcpTxs = append(qcpTxs, tx)
		return false
	})

	return qcpTxs
}

// TODO prefix定义到qbase中
func GetQCPTxsWithLimit(ctx context.Context, chainId string, limit uint64, asc bool) []txs.TxQcp {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpTxs := make([]txs.TxQcp, 0)
	var iterator store.Iterator
	if asc {
		iterator = store.KVStorePrefixIterator(qcpMapper.GetStore(), []byte("tx/out/"+chainId))
	} else {
		iterator = store.KVStoreReversePrefixIterator(qcpMapper.GetStore(), []byte("tx/out/"+chainId))
	}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		tx := txs.TxQcp{}
		qcpMapper.DecodeObject(iterator.Value(), &tx)
		qcpTxs = append(qcpTxs, tx)
		if uint64(len(qcpTxs)) == limit {
			break
		}
	}

	return qcpTxs
}

// TODO prefix定义到qbase中
func ExportQCPs(ctx context.Context) []qcptypes.QCPInfo {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpChains := make([]string, 0)

	prefix := []byte("sequence/in/")
	iter := qcpMapper.GetStore().Iterator(prefix, store.PrefixEndBytes(prefix))
	defer iter.Close()
	for {
		if !iter.Valid() {
			break
		}

		qcpChains = append(qcpChains, string(iter.Key()[len(prefix):]))
		iter.Next()
	}

	qcps := make([]qcptypes.QCPInfo, 0)
	for _, chaiId := range qcpChains {
		qcp := qcptypes.NewQCPInfo(chaiId,
			qcpMapper.GetMaxChainOutSequence(chaiId),
			qcpMapper.GetMaxChainInSequence(chaiId),
			qcpMapper.GetChainInTrustPubKey(chaiId),
			GetQCPTxsWithLimit(ctx, chaiId, QCPTxExportLimit, false))
		qcps = append(qcps, *qcp)
	}

	return qcps
}
