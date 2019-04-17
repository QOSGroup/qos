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
	cdc.RegisterConcrete(&TxCreateValidator{}, "stake/txs/TxCreateValidator", nil)
	cdc.RegisterConcrete(&TxRevokeValidator{}, "stake/txs/TxRevokeValidator", nil)
	cdc.RegisterConcrete(&TxActiveValidator{}, "stake/txs/TxActiveValidator", nil)

	//delegation相关
	cdc.RegisterConcrete(&TxCreateDelegation{}, "stake/txs/TxCreateDelegation", nil)
	cdc.RegisterConcrete(&TxModifyCompound{}, "stake/txs/TxModifyCompound", nil)
	cdc.RegisterConcrete(&TxUnbondDelegation{}, "stake/txs/TxUnbondDelegation", nil)
	cdc.RegisterConcrete(&TxCreateReDelegation{}, "stake/txs/TxCreateReDelegation", nil)
}
