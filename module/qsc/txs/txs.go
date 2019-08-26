package txs

import (
	"bytes"
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

const (
	MaxDescriptionLen = 1000
	MaxQSCNameLen     = 8

	GasForCreateQSC = uint64(1.8*qtypes.QOSUnit) * qtypes.GasPerUnitCost  // 1.8 QOS
	GasForIssueQSC  = uint64(0.18*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 0.18 QOS
)

// create QSC
type TxCreateQSC struct {
	Creator     btypes.Address       `json:"creator"`     //QSC创建账户
	Extrate     string               `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	QSCCA       *cert.Certificate    `json:"qsc_crt"`     //CA信息
	Description string               `json:"description"` //描述信息
	Accounts    []*qtypes.QOSAccount `json:"accounts"`
}

func (tx TxCreateQSC) ValidateData(ctx context.Context) error {
	if len(tx.Creator) == 0 || len(tx.Description) > MaxDescriptionLen {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	if _, err := strconv.ParseFloat(tx.Extrate, 64); err != nil {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	// CA校验
	if tx.QSCCA == nil {
		return types.ErrInvalidQSCCA(types.DefaultCodeSpace, "")
	}
	subj, ok := tx.QSCCA.CSR.Subj.(cert.QSCSubject)
	if !ok {
		return types.ErrInvalidQSCCA(types.DefaultCodeSpace, "")
	}
	if subj.ChainId != ctx.ChainID() {
		return types.ErrInvalidQSCCA(types.DefaultCodeSpace, "")
	}

	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
	rootCA := qscMapper.GetQSCRootCA()
	if !cert.VerityCrt([]crypto.PubKey{rootCA}, *tx.QSCCA) {
		return types.ErrWrongQSCCA(types.DefaultCodeSpace, "")
	}

	// accounts校验
	for _, account := range tx.Accounts {
		if account.QOS.NilToZero().GT(btypes.ZeroInt()) ||
			len(account.QSCs) != 1 || account.QSCs[0].Name != subj.Name ||
			!account.QSCs[0].Amount.NilToZero().GT(btypes.ZeroInt()) {
			return types.ErrInvalidInitAccounts(types.DefaultCodeSpace, "")
		}
	}

	// QSC不存在
	if qscMapper.Exists(subj.Name) {
		return types.ErrQSCExists(types.DefaultCodeSpace, "")
	}

	// creator账户存在
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	creator := accountMapper.GetAccount(tx.Creator)
	if nil == creator {
		return types.ErrCreatorNotExists(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxCreateQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	qscInfo := types.NewInfoWithQSCCA(tx.QSCCA)
	qscInfo.Extrate = tx.Extrate
	qscInfo.Description = tx.Description

	// 保存QSC
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
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
			qosAccount := a.(*qtypes.QOSAccount)
			qosAccount.MustPlusQSCs(acc.QSCs)
			accountMapper.SetAccount(qosAccount)
		} else {
			accountMapper.SetAccount(acc)
		}
	}

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateQsc,
			btypes.NewAttribute(types.AttributeKeyQsc, qscInfo.Name),
			btypes.NewAttribute(types.AttributeKeyCreator, tx.Creator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxCreateQSC) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Creator}
}

func (tx TxCreateQSC) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForCreateQSC))
}

func (tx TxCreateQSC) GetGasPayer() btypes.Address {
	return tx.Creator
}

func (tx TxCreateQSC) GetSignData() (ret []byte) {
	ret = append(ret, tx.Creator...)
	ret = append(ret, tx.Extrate...)
	ret = append(ret, Cdc.MustMarshalBinaryBare(tx.QSCCA)...)
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
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	// Amount大于0
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	// QSC存在
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
	qscInfo := qscMapper.GetQsc(tx.QSCName)
	if nil == qscInfo {
		return types.ErrQSCNotExists(types.DefaultCodeSpace, "")
	}

	// qscInfo banker存在
	if qscInfo.Banker == nil {
		return types.ErrBankerNotExists(types.DefaultCodeSpace, "")
	}

	// QSC名称一致
	if tx.QSCName != qscInfo.Name {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	// banker 地址一致
	if !bytes.Equal(tx.Banker, qscInfo.Banker) {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	return nil
}

func (tx TxIssueQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
	qscInfo := qscMapper.GetQsc(tx.QSCName)
	qscInfo.TotalAmount = qscInfo.TotalAmount.Add(tx.Amount)
	qscMapper.SaveQsc(qscInfo)

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	banker := accountMapper.GetAccount(tx.Banker).(*qtypes.QOSAccount)
	banker.MustPlusQSCs(qtypes.QSCs{btypes.NewBaseCoin(tx.QSCName, tx.Amount)})
	accountMapper.SetAccount(banker)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeIssueQsc,
			btypes.NewAttribute(types.AttributeKeyQsc, tx.QSCName),
			btypes.NewAttribute(types.AttributeKeyBanker, tx.Banker.String()),
			btypes.NewAttribute(types.AttributeKeyTokens, tx.Amount.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxIssueQSC) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Banker}
}

func (tx TxIssueQSC) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForIssueQSC))
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
