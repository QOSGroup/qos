package txs

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
)

type TransItem struct {
	Address btypes.Address `json:"addr"` // 账户地址
	QOS     btypes.BigInt  `json:"qos"`  // QOS
	QSCs    types.QSCs     `json:"qscs"` // QSCs
}

type TransferTx struct {
	Senders   []TransItem `json:"senders"`   // 发送集合
	Receivers []TransItem `json:"receivers"` // 接收集合
}

// 数据校验
// 1.Senders、Receivers不为空，地址不重复，币值大于0
// 2.Senders、Receivers 币值总和对应币种相等
// 3.Senders中账号对应币种币值足够
func (tx TransferTx) ValidateData(ctx context.Context) bool {
	if tx.Senders == nil || len(tx.Senders) == 0 {
		return false
	}

	if tx.Receivers == nil || len(tx.Receivers) == 0 {
		return false
	}

	accountMapper := ctx.Mapper(bacc.GetAccountKVStoreName()).(*bacc.AccountMapper)
	smap := map[string]bool{}
	sumsqos := btypes.ZeroInt()
	sumsqscs := types.QSCs{}
	for _, sender := range tx.Senders {
		if _, ok := smap[sender.Address.String()]; ok {
			return false
		}
		smap[sender.Address.String()] = true
		sender.QOS = types.ZeroNilBigInt(sender.QOS)
		if sender.QOS.IsZero() && sender.QSCs.IsZero() {
			return false
		}
		if btypes.ZeroInt().GT(sender.QOS) {
			return false
		}
		if !sender.QSCs.IsNotNegative() {
			return false
		}
		a := accountMapper.GetAccount(sender.Address)
		if a == nil {
			return false
		}
		acc := a.(*account.QOSAccount)
		if acc == nil || acc.QOS.LT(sender.QOS) || acc.QSCs.IsLT(sender.QSCs) {
			return false
		}
		sumsqos = sumsqos.Add(sender.QOS)
		if nil != sender.QSCs {
			sumsqscs = sumsqscs.Plus(sender.QSCs)
		}
	}
	rmap := map[string]bool{}
	sumrqos := btypes.ZeroInt()
	sumrqscs := types.QSCs{}
	for _, receiver := range tx.Receivers {
		if _, ok := rmap[receiver.Address.String()]; ok {
			return false
		}
		rmap[receiver.Address.String()] = true
		receiver.QOS = types.ZeroNilBigInt(receiver.QOS)
		if receiver.QOS.IsZero() && receiver.QSCs.IsZero() {
			return false
		}
		if btypes.ZeroInt().GT(receiver.QOS) {
			return false
		}
		if !receiver.QSCs.IsNotNegative() {
			return false
		}
		sumrqos = sumrqos.Add(receiver.QOS)
		if nil != receiver.QSCs {
			sumrqscs = sumrqscs.Plus(receiver.QSCs)
		}
	}

	// 转入转出相等
	if !sumsqos.Equal(sumrqos) || !sumsqscs.IsEqual(sumrqscs) {
		return false
	}

	return true
}

// 转账
func (tx TransferTx) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	accountMapper := ctx.Mapper(bacc.GetAccountKVStoreName()).(*bacc.AccountMapper)

	for _, sender := range tx.Senders {
		acc := accountMapper.GetAccount(sender.Address).(*account.QOSAccount)
		acc.QOS = types.ZeroNilBigInt(acc.QOS)
		sender.QOS = types.ZeroNilBigInt(sender.QOS)
		acc.QOS = acc.QOS.Add(sender.QOS.Neg())
		acc.QSCs = acc.QSCs.Minus(sender.QSCs)
		accountMapper.SetAccount(acc)
	}
	for _, receiver := range tx.Receivers {
		a := accountMapper.GetAccount(receiver.Address)
		var acc *account.QOSAccount
		if a != nil {
			acc = a.(*account.QOSAccount)
		} else {
			acc = accountMapper.NewAccountWithAddress(receiver.Address).(*account.QOSAccount)
			accountMapper.SetAccount(acc)
		}
		acc.QOS = types.ZeroNilBigInt(acc.QOS)
		receiver.QOS = types.ZeroNilBigInt(receiver.QOS)
		acc.QOS = acc.QOS.Add(receiver.QOS)
		acc.QSCs = acc.QSCs.Plus(receiver.QSCs)
		accountMapper.SetAccount(acc)
	}

	return btypes.Result{Code: btypes.ABCICodeOK}, nil
}

// 所有Senders
func (tx TransferTx) GetSigner() []btypes.Address {
	addrs := make([]btypes.Address, 0)
	for _, sender := range tx.Senders {
		addrs = append(addrs, sender.Address)
	}

	return addrs
}

// Gas TODO
func (tx TransferTx) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Senders[0]
func (tx TransferTx) GetGasPayer() btypes.Address {
	return tx.Senders[0].Address
}

// 签名字节
func (tx TransferTx) GetSignData() (ret []byte) {
	for _, sender := range tx.Senders {
		ret = append(ret, sender.Address...)
		ret = append(ret, sender.QOS.String()...)
		ret = append(ret, sender.QSCs.String()...)
	}
	for _, receiver := range tx.Receivers {
		ret = append(ret, receiver.Address...)
		ret = append(ret, receiver.QOS.String()...)
		ret = append(ret, receiver.QSCs.String()...)
	}

	return ret
}
