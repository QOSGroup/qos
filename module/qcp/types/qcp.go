package types

import (
	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/tendermint/crypto"
)

// 联盟链信息
type QCPInfo struct {
	ChainId     string        `json:"chain_id"`     // 联盟链标识
	SequenceOut int64         `json:"sequence_out"` // 结果输出序号
	SequenceIn  int64         `json:"sequence_in"`  // 交易接收序号
	PubKey      crypto.PubKey `json:"pub_key"`      // 公钥
	OutTxs      []txs.TxQcp   `json:"txs"`          // 跨链交易结果集
}

func NewQCPInfo(chainId string, sequenceOut int64, sequenceIn int64, pubKey crypto.PubKey, txs []txs.TxQcp) *QCPInfo {
	return &QCPInfo{
		ChainId:     chainId,
		SequenceIn:  sequenceIn,
		SequenceOut: sequenceOut,
		PubKey:      pubKey,
		OutTxs:      txs,
	}
}
