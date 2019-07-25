package approve

import (
	"github.com/QOSGroup/qos/module/approve/mapper"
	"github.com/QOSGroup/qos/module/approve/txs"
	"github.com/QOSGroup/qos/module/approve/types"
)

var (
	ModuleName    = "approve"
	RegisterCodec = txs.RegisterCodec

	NewMapper  = mapper.NewApproveMapper()
	MapperName = mapper.MapperName
)

type (
	Mapper = mapper.Mapper

	GenesisState = types.GenesisState

	TxCreateApprove = txs.TxCreateApprove
	TxIncreaseApprove = txs.TxIncreaseApprove
	TxDecreaseApprove = txs.TxDecreaseApprove
	TxUseApprove = txs.TxUseApprove
	TxCancelApprove = txs.TxCancelApprove
)
