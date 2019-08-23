package qcp

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/qcp/mapper"
	"github.com/QOSGroup/qos/module/qcp/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	if data.RootPubKey != nil {
		mapper.SetQCPRootCA(ctx, data.RootPubKey)
	}

	qcpMapper := mapper.GetMapper(ctx)
	for _, qcp := range data.QCPs {
		qcpMapper.SetMaxChainInSequence(qcp.ChainId, qcp.SequenceIn)
		qcpMapper.SetMaxChainOutSequence(qcp.ChainId, qcp.SequenceOut)
		qcpMapper.SetChainInTrustPubKey(qcp.ChainId, qcp.PubKey)
		for _, tx := range qcp.OutTxs {
			qcpMapper.SetChainOutTxs(qcp.ChainId, tx.Sequence, &tx)
		}
	}
}

func ExportGenesis(ctx context.Context) types.GenesisState {
	return types.NewGenesisState(mapper.GetQCPRootCA(ctx), mapper.ExportQCPs(ctx))
}
