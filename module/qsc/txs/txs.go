package txs

import (
	"bytes"
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/qsc/mapper"
	"github.com/QOSGroup/qos/module/qsc/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

var (
	MaxDescriptionLen = 1000

	GasForCreateQSC = int64(1.8*qtypes.UnitQOS) * qtypes.UnitQOSGas  // 1.8 QOS
	GasForIssueQSC  = int64(0.18*qtypes.UnitQOS) * qtypes.UnitQOSGas // 0.18 QOS
)

// 初始化代币Tx
type TxCreateQSC struct {
	Creator      btypes.AccAddress    `json:"creator"`       // QSC创建账户
	ExchangeRate string               `json:"exchange_rate"` // qcs:qos汇率
	QSCCA        *cert.Certificate    `json:"qsc_crt"`       // CA信息
	Description  string               `json:"description"`   // 描述信息
	Accounts     []*qtypes.QOSAccount `json:"accounts"`      // 初始账户
}

func (tx TxCreateQSC) ValidateInputs() error {
	// 创建账户校验
	if len(tx.Creator) == 0 {
		return types.ErrEmptyCreator()
	}
	// 描述信息最大1000字节
	if len(tx.Description) > MaxDescriptionLen {
		return types.ErrDescriptionTooLong()
	}
	// 汇率值校验
	if _, err := strconv.ParseFloat(tx.ExchangeRate, 64); err != nil {
		return types.ErrInvalidExchangeRate()
	}

	// 证书校验
	if tx.QSCCA == nil {
		return types.ErrInvalidQSCCA()
	}
	subj, ok := tx.QSCCA.CSR.Subj.(cert.QSCSubject)
	if !ok {
		return types.ErrInvalidQSCCA()
	}

	// 初始账户校验，只能包含即将初始化的代币
	for _, account := range tx.Accounts {
		if account.QOS.NilToZero().GT(btypes.ZeroInt()) ||
			len(account.QSCs) != 1 || account.QSCs[0].Name != subj.Name ||
			!account.QSCs[0].Amount.NilToZero().GT(btypes.ZeroInt()) {
			return types.ErrInvalidInitAccounts()
		}
	}

	return nil
}

func (tx TxCreateQSC) ValidateData(ctx context.Context) error {
	// 校验基础数据
	err := tx.ValidateInputs()
	if err != nil {
		return err
	}

	subj, _ := tx.QSCCA.CSR.Subj.(cert.QSCSubject)
	if subj.ChainId != ctx.ChainID() {
		return types.ErrInvalidQSCCA()
	}
	qscMapper := mapper.GetMapper(ctx)
	rootCA := qscMapper.GetRootCAPubkey()
	if rootCA == nil || len(rootCA.Bytes()) == 0 {
		return types.ErrRootCANotConfigure()
	}

	if !cert.VerityCrt([]crypto.PubKey{rootCA}, *tx.QSCCA) {
		return types.ErrWrongQSCCA()
	}

	// 校验已存在代币
	if qscMapper.Exists(subj.Name) {
		return types.ErrQSCExists()
	}

	// 创建账户必须存在
	creator := bank.GetMapper(ctx).GetAccount(tx.Creator)
	if nil == creator {
		return types.ErrCreatorNotExists()
	}

	return nil
}

func (tx TxCreateQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	qscInfo := types.NewInfoWithQSCCA(tx.QSCCA)
	qscInfo.ExchangeRate = tx.ExchangeRate
	qscInfo.Description = tx.Description

	// 保存初始账户信息
	bankMapper := bank.GetMapper(ctx)
	if qscInfo.Banker != nil {
		banker := qscInfo.Banker
		if nil == bankMapper.GetAccount(banker) {
			bankMapper.SetAccount(bankMapper.NewAccountWithAddress(banker))
		}
	}
	qscInfo.TotalAmount = btypes.ZeroInt()
	for _, acc := range tx.Accounts {
		if a := bankMapper.GetAccount(acc.AccountAddress); a != nil {
			qosAccount := a.(*qtypes.QOSAccount)
			qosAccount.MustPlusQSCs(acc.QSCs)
			bankMapper.SetAccount(qosAccount)
		} else {
			bankMapper.SetAccount(acc)
		}
		qscInfo.TotalAmount = qscInfo.TotalAmount.Add(acc.QSCs[0].Amount)
	}

	// 保存代币信息
	qscMapper := mapper.GetMapper(ctx)
	qscMapper.SaveQsc(&qscInfo)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeCreateQsc,
			btypes.NewAttribute(types.AttributeKeyQsc, qscInfo.Name),
			btypes.NewAttribute(types.AttributeKeyCreator, tx.Creator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeCreateQsc),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxCreateQSC) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Creator}
}

func (tx TxCreateQSC) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForCreateQSC)
}

func (tx TxCreateQSC) GetGasPayer() btypes.AccAddress {
	return tx.Creator
}

func (tx TxCreateQSC) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 发行代币Tx，多次执行表现为代币累加
type TxIssueQSC struct {
	QSCName string            `json:"qsc_name"` //币名
	Amount  btypes.BigInt     `json:"amount"`   //币量
	Banker  btypes.AccAddress `json:"banker"`   //banker地址
}

func (tx TxIssueQSC) ValidateData(ctx context.Context) error {
	// 币量校验
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return types.ErrAmountLTZero()
	}

	// 代币校验
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
	qscInfo, exists := qscMapper.GetQsc(tx.QSCName)
	if !exists {
		return types.ErrQSCNotExists()
	}

	// banker地址校验
	if len(qscInfo.Banker) == 0 {
		return types.ErrBankerNotExists()
	}
	if !bytes.Equal(tx.Banker, qscInfo.Banker) {
		return types.ErrInvalidBanker()
	}

	return nil
}

func (tx TxIssueQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}
	qscMapper := ctx.Mapper(mapper.MapperName).(*mapper.Mapper)
	qscInfo, _ := qscMapper.GetQsc(tx.QSCName)
	qscInfo.TotalAmount = qscInfo.TotalAmount.Add(tx.Amount)
	qscMapper.SaveQsc(&qscInfo)

	bankMapper := bank.GetMapper(ctx)
	// 发行代币
	banker := bankMapper.GetAccount(tx.Banker).(*qtypes.QOSAccount)
	banker.MustPlusQSCs(qtypes.QSCs{btypes.NewBaseCoin(tx.QSCName, tx.Amount)})
	bankMapper.SetAccount(banker)

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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeIssueQsc),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxIssueQSC) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Banker}
}

func (tx TxIssueQSC) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForIssueQSC)
}

func (tx TxIssueQSC) GetGasPayer() btypes.AccAddress {
	return tx.Banker
}

func (tx TxIssueQSC) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}
