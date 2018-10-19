package account

import (
	"github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

// QOSAccount定义基本账户之上的QOS和QSC
type QOSAccount struct {
	BaseAccount account.BaseAccount `json:"base_account"` // inherits BaseAccount
	Qos         btypes.BigInt       `json:"qos"`          // coins in public chain
	QscList     []*types.QSC        `json:"qsc"`          // varied QSCs
}

func ProtoQOSAccount() account.Account {
	return &QOSAccount{}
}

func (acc *QOSAccount) GetProto() account.Account {
	return ProtoQOSAccount()
}

func (account *QOSAccount) GetAddress() btypes.Address {
	return account.BaseAccount.GetAddress()
}

func (account *QOSAccount) SetAddress(addr btypes.Address) error {
	return account.BaseAccount.SetAddress(addr)
}

func (account *QOSAccount) GetPubicKey() crypto.PubKey {
	return account.BaseAccount.GetPubicKey()
}

func (account *QOSAccount) SetPublicKey(pubKey crypto.PubKey) error {
	return account.BaseAccount.SetPublicKey(pubKey)
}

func (account *QOSAccount) GetNonce() uint64 {
	return account.BaseAccount.GetNonce()
}

func (account *QOSAccount) SetNonce(nonce uint64) error {
	return account.BaseAccount.SetNonce(nonce)
}

// 获得账户QOS的数量
func (accnt *QOSAccount) GetQOS() btypes.BigInt {
	return accnt.Qos
}

// 设置账户QOS的数量
func (accnt *QOSAccount) SetQOS(amount btypes.BigInt) error {
	accnt.Qos = amount
	return nil
}

// 获取账户的名为QSCName的币的数量
func (accnt *QOSAccount) GetQSC(QSCName string) *types.QSC {
	for _, qsc := range accnt.QscList {
		if qsc.GetName() == QSCName {
			return qsc
		}
	}
	return nil
}

// 设置账户的名为QSCName的币
func (accnt *QOSAccount) SetQSC(newQSC *types.QSC) error {
	for _, qsc := range accnt.QscList {
		if qsc.GetName() == newQSC.GetName() {
			qsc.SetAmount(newQSC.GetAmount())
			return nil
		}
	}
	accnt.QscList = append(accnt.QscList, newQSC)
	return nil
}

// 删除账户中名为QSCName的币
func (accnt *QOSAccount) RemoveQSCByName(QSCName string) error {
	for i, qsc := range accnt.QscList {
		if qsc.GetName() == QSCName {
			if i == len(accnt.QscList)-1 {
				accnt.QscList = accnt.QscList[:i]
				return nil
			}
			accnt.QscList = append(accnt.QscList[:i], accnt.QscList[i+1:]...)
			return nil
		}
	}
	return btypes.ErrInvalidCoins(btypes.CodeToDefaultMsg(btypes.CodeInvalidCoins))
}

// Coins包含qos
func (accnt *QOSAccount) Coins() types.QSCS {
	l := []types.QSC{}
	for _, qsc := range accnt.QscList {
		l = append(l, *qsc)
	}
	l = append(l, types.QSC{Name: "qos", Amount: accnt.Qos})
	return types.QSCS(l).Sort()
}
