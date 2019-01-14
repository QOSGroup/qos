package approve

import (
	"github.com/QOSGroup/qbase/context"
	approvetypes "github.com/QOSGroup/qos/module/approve/types"
)

type GenesisState struct {
	Approves []approvetypes.Approve `json:"approves"`
}

func NewGenesisState(approves []approvetypes.Approve) GenesisState {
	return GenesisState{
		approves,
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	for _, approve := range data.Approves {
		approveMapper.SaveApprove(approve)
	}
}

func ExportGenesis(ctx context.Context) GenesisState {
	approveMapper := ctx.Mapper(ApproveMapperName).(*ApproveMapper)
	return NewGenesisState(approveMapper.GetApproves())
}
