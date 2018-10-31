package account

import (
	"github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
)

// QOSAccount定义基本账户之上的QOS和QSC
type QOSAccount struct {
	account.BaseAccount `json:"base_account"`       // inherits BaseAccount
	QOS                 btypes.BigInt `json:"qos"`  // coins in public chain
	QSCs                types.QSCs    `json:"qscs"` // varied QSCs
}

var _ account.Account = (*QOSAccount)(nil)

func ProtoQOSAccount() account.Account {
	return NewQOSAccount()
}

func NewQOSAccount() *QOSAccount {
	return &QOSAccount{QOS: btypes.ZeroInt()}
}

func (acc *QOSAccount) GetProto() account.Account {
	return ProtoQOSAccount()
}

// 获得账户QOS的数量
func (accnt *QOSAccount) GetQOS() btypes.BigInt {
	return accnt.QOS
}

// 设置账户QOS的数量
func (accnt *QOSAccount) SetQOS(amount btypes.BigInt) error {
	accnt.QOS = amount
	return nil
}

// 获取账户的名为QSCName的币的数量
func (accnt *QOSAccount) GetQSC(QSCName string) *types.QSC {
	for _, qsc := range accnt.QSCs {
		if qsc.GetName() == QSCName {
			return qsc
		}
	}
	return nil
}

// 设置账户的名为QSCName的币
func (accnt *QOSAccount) SetQSC(newQSC *types.QSC) error {
	for _, qsc := range accnt.QSCs {
		if qsc.GetName() == newQSC.GetName() {
			qsc.SetAmount(newQSC.GetAmount())
			return nil
		}
	}
	accnt.QSCs = append(accnt.QSCs, newQSC)
	return nil
}

// 删除账户中名为QSCName的币
func (accnt *QOSAccount) RemoveQSCByName(QSCName string) error {
	for i, qsc := range accnt.QSCs {
		if qsc.GetName() == QSCName {
			if i == len(accnt.QSCs)-1 {
				accnt.QSCs = accnt.QSCs[:i]
				return nil
			}
			accnt.QSCs = append(accnt.QSCs[:i], accnt.QSCs[i+1:]...)
			return nil
		}
	}
	return btypes.ErrInvalidCoins(btypes.CodeToDefaultMsg(btypes.CodeInvalidCoins))
}
