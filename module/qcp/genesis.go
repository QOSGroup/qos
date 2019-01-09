package qcp

import (
	"github.com/QOSGroup/qbase/context"
	qcptypes "github.com/QOSGroup/qos/module/qcp/types"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	RootPubKey crypto.PubKey      `json:"ca_root_pub_key"`
	QCPs       []qcptypes.QCPInfo `json:"qcps""`
}

func NewGenesisState(pubKey crypto.PubKey, qcps []qcptypes.QCPInfo) GenesisState {
	return GenesisState{
		RootPubKey: pubKey,
		QCPs:       qcps,
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	if data.RootPubKey != nil {
		SetQCPRootCA(ctx, data.RootPubKey)
	}

	qcpMapper := GetQCPMapper(ctx)
	for _, qcp := range data.QCPs {
		qcpMapper.SetMaxChainInSequence(qcp.ChainId, qcp.SequenceIn)
		qcpMapper.SetMaxChainOutSequence(qcp.ChainId, qcp.SequenceOut)
		qcpMapper.SetChainInTrustPubKey(qcp.ChainId, qcp.PubKey)
		for _, tx := range qcp.OutTxs {
			qcpMapper.SetChainOutTxs(qcp.ChainId, tx.Sequence, &tx)
		}
	}
}

func ExportGenesis(ctx context.Context) GenesisState {
	return NewGenesisState(GetQCPRootCA(ctx), ExportQCPs(ctx))
}
