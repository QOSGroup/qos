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
	return ctx.Mapper(qcp.MapperName).(*qcp.QcpMapper)
}

// 保存 kepler qcp 根证书公钥
func SetRootCaPubkey(ctx context.Context, pubKey crypto.PubKey) {
	qcpMapper := GetMapper(ctx)
	qcpMapper.Set([]byte(RootCAKey), pubKey)
}

// 获取 kepler qcp 根证书公钥
func GetRootCaPubkey(ctx context.Context) crypto.PubKey {
	qcpMapper := GetMapper(ctx)
	var pubKey crypto.PubKey
	qcpMapper.Get([]byte(RootCAKey), &pubKey)
	return pubKey
}

// 获取跨链结果交易集
func GetQCPTxs(ctx context.Context, chainId string) []txs.TxQcp {
	qcpMapper := GetMapper(ctx)
	qcpTxs := make([]txs.TxQcp, 0)
	qcpMapper.Iterator(append(qcp.BuildOutSequenceTxPrefixKey(), chainId...), func(bz []byte) (stop bool) {
		tx := txs.TxQcp{}
		qcpMapper.DecodeObject(bz, &tx)
		qcpTxs = append(qcpTxs, tx)
		return false
	})

	return qcpTxs
}

// 获取跨链结果交易集，带条数限制
func GetQCPTxsWithLimit(ctx context.Context, chainId string, limit uint64, asc bool) []txs.TxQcp {
	qcpMapper := GetMapper(ctx)
	qcpTxs := make([]txs.TxQcp, 0)
	var iterator store.Iterator
	if asc {
		iterator = btypes.KVStorePrefixIterator(qcpMapper.GetStore(), append(qcp.BuildOutSequenceTxPrefixKey(), chainId...))
	} else {
		iterator = btypes.KVStoreReversePrefixIterator(qcpMapper.GetStore(), append(qcp.BuildOutSequenceTxPrefixKey(), chainId...))
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

// 获取联盟链集合，结果交易集中最多保留 TxExportLimit 条数据
func GetQCPs(ctx context.Context) []qcptypes.QCPInfo {
	qcpMapper := GetMapper(ctx)
	qcpChains := make([]string, 0)

	prefix := qcp.BuildInSequencePrefixKey()
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
