package qsc

import (
	"bytes"
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	qsctypes "github.com/QOSGroup/qos/module/qsc/types"
	"github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

const (
	MaxDescriptionLen = 1000
	MaxQSCNameLen     = 8
)

// create QSC
type TxCreateQSC struct {
	Creator     btypes.Address      `json:"creator"`     //QSC创建账户
	Extrate     string              `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	QSCCA       *cert.Certificate   `json:"qsc_crt"`     //CA信息
	Description string              `json:"description"` //描述信息
	Accounts    []*types.QOSAccount `json:"accounts"`
}

func (tx TxCreateQSC) ValidateData(ctx context.Context) error {
	if len(tx.Creator) == 0 || len(tx.Description) > MaxDescriptionLen {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	if _, err := strconv.ParseFloat(tx.Extrate, 64); err != nil {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	// CA校验
	if tx.QSCCA == nil {
		return ErrInvalidQSCCA(DefaultCodeSpace, "")
	}
	subj, ok := tx.QSCCA.CSR.Subj.(cert.QSCSubject)
	if !ok {
		return ErrInvalidQSCCA(DefaultCodeSpace, "")
	}
	if subj.ChainId != ctx.ChainID() {
		return ErrInvalidQSCCA(DefaultCodeSpace, "")
	}

	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	rootCA := qscMapper.GetQSCRootCA()
	if !cert.VerityCrt([]crypto.PubKey{rootCA}, *tx.QSCCA) {
		return ErrWrongQSCCA(DefaultCodeSpace, "")
	}

	// accounts校验
	for _, account := range tx.Accounts {
		if account.QOS.NilToZero().GT(btypes.ZeroInt()) ||
			len(account.QSCs) != 1 || account.QSCs[0].Name != subj.Name ||
			!account.QSCs[0].Amount.NilToZero().GT(btypes.ZeroInt()) {
			return ErrInvalidInitAccounts(DefaultCodeSpace, "")
		}
	}

	// QSC不存在
	if qscMapper.Exists(subj.Name) {
		return ErrQSCExists(DefaultCodeSpace, "")
	}

	// creator账户存在
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	creator := accountMapper.GetAccount(tx.Creator)
	if nil == creator {
		return ErrCreatorNotExists(DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxCreateQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	qscInfo := qsctypes.NewQSCInfoWithQSCCA(tx.QSCCA)
	qscInfo.Extrate = tx.Extrate
	qscInfo.Description = tx.Description

	// 保存QSC
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	qscMapper.SaveQsc(&qscInfo)

	// 保存账户信息
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	if qscInfo.Banker != nil {
		banker := qscInfo.Banker
		if nil == accountMapper.GetAccount(banker) {
			accountMapper.SetAccount(accountMapper.NewAccountWithAddress(banker))
		}
	}
	for _, acc := range tx.Accounts {
		if a := accountMapper.GetAccount(acc.AccountAddress); a != nil {
			qosAccount := a.(*types.QOSAccount)
			qosAccount.MustPlusQSCs(acc.QSCs)
			accountMapper.SetAccount(qosAccount)
		} else {
			accountMapper.SetAccount(acc)
		}
	}

	return
}

func (tx TxCreateQSC) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Creator}
}

func (tx TxCreateQSC) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxCreateQSC) GetGasPayer() btypes.Address {
	return tx.Creator
}

func (tx TxCreateQSC) GetSignData() (ret []byte) {
	ret = append(ret, tx.Creator...)
	ret = append(ret, tx.Extrate...)
	ret = append(ret, cdc.MustMarshalBinaryBare(tx.QSCCA)...)
	ret = append(ret, tx.Description...)

	for _, account := range tx.Accounts {
		ret = append(ret, fmt.Sprint(account)...)
	}

	return
}

// issue QSC
type TxIssueQSC struct {
	QSCName string         `json:"qsc_name"` //币名
	Amount  btypes.BigInt  `json:"amount"`   //金额
	Banker  btypes.Address `json:"banker"`   //banker地址
}

func (tx TxIssueQSC) ValidateData(ctx context.Context) error {
	// QscName不能为空，且不能超过8个字符
	if len(tx.QSCName) == 0 || len(tx.QSCName) > MaxQSCNameLen {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	// Amount大于0
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	// QSC存在
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	qscInfo := qscMapper.GetQsc(tx.QSCName)
	if nil == qscInfo {
		return ErrQSCNotExists(DefaultCodeSpace, "")
	}

	// qscInfo banker存在
	if qscInfo.Banker == nil {
		return ErrBankerNotExists(DefaultCodeSpace, "")
	}

	// QSC名称一致
	if tx.QSCName != qscInfo.Name {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	// banker 地址一致
	if !bytes.Equal(tx.Banker, qscInfo.Banker) {
		return ErrInvalidInput(DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxIssueQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)

	banker := accountMapper.GetAccount(tx.Banker).(*types.QOSAccount)
	banker.MustPlusQSCs(types.QSCs{btypes.NewBaseCoin(tx.QSCName, tx.Amount)})
	accountMapper.SetAccount(banker)

	return
}

func (tx TxIssueQSC) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Banker}
}

func (tx TxIssueQSC) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxIssueQSC) GetGasPayer() btypes.Address {
	return tx.Banker
}

func (tx TxIssueQSC) GetSignData() (ret []byte) {
	ret = append(ret, tx.QSCName...)
	ret = append(ret, btypes.Int2Byte(tx.Amount.Int64())...)
	ret = append(ret, tx.Banker...)

	return
}
