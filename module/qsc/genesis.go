package qsc

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/qsc/types"
	"github.com/tendermint/tendermint/crypto"
)

type GenesisState struct {
	RootPubKey crypto.PubKey   `json:"ca_root_pub_key"`
	QSCs       []types.QSCInfo `json:"qscs"`
}

func NewGenesisState(pubKey crypto.PubKey, qscs []types.QSCInfo) GenesisState {
	return GenesisState{
		RootPubKey: pubKey,
		QSCs:       qscs,
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	if data.RootPubKey != nil {
		qscMapper.SetQSCRootCA(data.RootPubKey)
	}

	for _, qsc := range data.QSCs {
		qscMapper.SaveQsc(&qsc)
	}
}

func ExportGenesis(ctx context.Context) GenesisState {
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)

	return NewGenesisState(qscMapper.GetQSCRootCA(), qscMapper.GetQSCs())
}
