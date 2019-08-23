package qsc

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
)

func InitGenesis(ctx context.Context, data types.GenesisState) {
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
	if data.RootPubKey != nil {
		qscMapper.SetQSCRootCA(data.RootPubKey)
	}

	for _, qsc := range data.QSCs {
		qscMapper.SaveQsc(&qsc)
	}
}

func ExportGenesis(ctx context.Context) types.GenesisState {
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)

	return types.NewGenesisState(qscMapper.GetQSCRootCA(), qscMapper.GetQSCs())
}
