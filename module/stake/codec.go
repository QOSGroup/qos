package stake

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/go-amino"
)

var cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	//validator相关
	cdc.RegisterConcrete(&TxCreateValidator{}, "qos/txs/TxCreateValidator", nil)
	cdc.RegisterConcrete(&TxRevokeValidator{}, "qos/txs/TxRevokeValidator", nil)
	cdc.RegisterConcrete(&TxActiveValidator{}, "qos/txs/TxActiveValidator", nil)

	//delegation相关
	cdc.RegisterConcrete(&TxCreateDelegation{}, "qos/txs/TxCreateDelegation", nil)
	cdc.RegisterConcrete(&TxModifyCompound{}, "qos/txs/TxModifyCompound", nil)
	cdc.RegisterConcrete(&TxUnbondDelegation{}, "qos/txs/TxUnbondDelegation", nil)
	cdc.RegisterConcrete(&TxCreateReDelegation{}, "qos/txs/TxCreateReDelegation", nil)
}
