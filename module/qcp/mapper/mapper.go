package mapper

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/store"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	qcptypes "github.com/QOSGroup/qos/module/qcp/types"
	"github.com/tendermint/tendermint/crypto"
)

const (
	RootCAKey = "rootca"

	TxExportLimit = 100 // QCP TX 导出条数限制
)

func GetMapper(ctx context.Context) *qcp.QcpMapper {
	return ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
}

// 保存CA
func SetQCPRootCA(ctx context.Context, pubKey crypto.PubKey) {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpMapper.Set([]byte(RootCAKey), pubKey)
}

// 获取CA
func GetQCPRootCA(ctx context.Context) crypto.PubKey {
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	var pubKey crypto.PubKey
	qcpMapper.Get([]byte(RootCAKey), &pubKey)
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
		iterator = btypes.KVStorePrefixIterator(qcpMapper.GetStore(), []byte("tx/out/"+chainId))
	} else {
		iterator = btypes.KVStoreReversePrefixIterator(qcpMapper.GetStore(), []byte("tx/out/"+chainId))
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
	iter := qcpMapper.GetStore().Iterator(prefix, btypes.PrefixEndBytes(prefix))
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
			GetQCPTxsWithLimit(ctx, chaiId, TxExportLimit, false))
		qcps = append(qcps, *qcp)
	}

	return qcps
}
