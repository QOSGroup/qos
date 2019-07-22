package txs

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qos/module/stake/types"
	"github.com/tendermint/go-amino"
)

var Cdc = baseabci.MakeQBaseCodec()

func init() {
	RegisterCodec(Cdc)
}

func RegisterCodec(cdc *amino.Codec) {
	// validator
	cdc.RegisterConcrete(&TxCreateValidator{}, "stake/txs/TxCreateValidator", nil)
	cdc.RegisterConcrete(&TxModifyValidator{}, "stake/txs/TxModifyValidator", nil)
	cdc.RegisterConcrete(&TxRevokeValidator{}, "stake/txs/TxRevokeValidator", nil)
	cdc.RegisterConcrete(&TxActiveValidator{}, "stake/txs/TxActiveValidator", nil)

	// delegation
	cdc.RegisterConcrete(&TxCreateDelegation{}, "stake/txs/TxCreateDelegation", nil)
	cdc.RegisterConcrete(&TxModifyCompound{}, "stake/txs/TxModifyCompound", nil)
	cdc.RegisterConcrete(&TxUnbondDelegation{}, "stake/txs/TxUnbondDelegation", nil)
	cdc.RegisterConcrete(&TxCreateReDelegation{}, "stake/txs/TxCreateReDelegation", nil)

	// params
	cdc.RegisterConcrete(&types.Params{}, "stake/params", nil)
}
