package types

import (
	"github.com/QOSGroup/qbase/txs"
	"github.com/tendermint/tendermint/crypto"
)

type QCPInfo struct {
	ChainId     string        `json:"chain_id""`
	SequenceOut int64         `json:"sequence_out""`
	SequenceIn  int64         `json:"sequence_in""`
	PubKey      crypto.PubKey `json:"pub_key""`
	OutTxs      []txs.TxQcp   `json:"txs""`
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
