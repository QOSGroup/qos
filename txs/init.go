package txs

import (
	baccount "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

//功能：序列化设置
func RegisterCodec(cdc *go_amino.Codec) {
	cryptoAmino.RegisterAmino(cdc)
	cdc.RegisterConcrete(&CA{}, "qos/txs/ca", nil)
	cdc.RegisterConcrete(&TxCreateQSC{}, "qos/txs/createqsc", nil)
	cdc.RegisterConcrete(&TxIssueQsc{}, "qos/txs/issueqsc", nil)
	cdc.RegisterConcrete(&TxTransform{}, "qos/txs/transform", nil)
}

//todo: 依赖ctx中store操作，稍后更新(暂模拟)
func GetAccount(addr btypes.Address) (acc *account.QOSAccount) {
	accmapper := baccount.NewAccountMapper(account.ProtoQOSAccount)
	addrKey := baccount.AddressStoreKey(addr)
	accmapper.Get(addrKey, &acc)

	return
}

// todo: 暂模拟，稍后实现
func CreateAccount(addr btypes.Address) (acc *account.QOSAccount) {
	acc = GetAccount(addr)
	//if acc != nil {
	//	err = errors.New("Error; address()  account existed. ")
	//	return acc, err
	//}

	//accmapper := baccount.NewAccountMapper(account.ProtoQOSAccount)
	//addrKey := baccount.AddressStoreKey(addr)
	//accQos := accmapper.
	return acc
}

func GetBanker(qscname string) (ret *account.QOSAccount) {
	return nil
}

func FetchQscCA() (caQsc []byte) {
	return nil
}

func FetchBankerCA() (caBanker []byte) {
	return nil
}