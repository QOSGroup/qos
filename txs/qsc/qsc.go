package qsc

import (
	"bytes"
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strings"
	"time"
)

var (
	CreatorQOSLimit = btypes.NewInt(0) //TODO
)

// types and funcs for CA
// copy from QOSGroup/kepler
// TODO remove
//-----------------------------------------------------------------
type Subject struct {
	CN string `json:"cn"`
}

type Issuer struct {
	Subj      Subject               `json:"subj"`
	PublicKey ed25519.PubKeyEd25519 `json:"public_key"`
}

type CertificateSigningRequest struct {
	Subj      Subject               `json:"subj"`
	IsCa      bool                  `json:"is_ca"`
	IsBanker  bool                  `json:"is_banker"`
	NotBefore time.Time             `json:"not_before"`
	NotAfter  time.Time             `json:"not_after"`
	PublicKey ed25519.PubKeyEd25519 `json:"public_key"`
}

type Certificate struct {
	CSR       CertificateSigningRequest `json:"csr"`
	CA        Issuer                    `json:"ca"`
	Signature []byte                    `json:"signature"`
}

func VerityCrt(pubKey crypto.PubKey, crt *Certificate) bool {
	// Check issue
	if pubKey.Equals(crt.CA.PublicKey) {
		sign, err := cdc.MarshalBinaryBare(crt.CSR)
		if err != nil {
			return false
		}
		if !crt.CA.PublicKey.VerifyBytes(sign, crt.Signature) {
			return false
		}
	} else {
		return false
	}

	// Check timestamp
	now := time.Now().Unix()
	if now <= crt.CSR.NotBefore.Unix() || now >= crt.CSR.NotAfter.Unix() {
		return false
	}

	return true
}

//-----------------------------------------------------------------

type QSCInfo = TxCreateQSC

// create QSC
type TxCreateQSC struct {
	ChainID     string               `json:"chain_id"`    //chain-id
	Creator     btypes.Address       `json:"creator"`     //QSC创建账户
	Extrate     string               `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	QSCCA       *Certificate         `json:"ca_qsc"`      //CA信息
	BankerCA    *Certificate         `json:"ca_banker"`   //CA信息
	Description string               `json:"description"` //描述信息
	Accounts    []account.QOSAccount `json:"accounts"`    //初始化时接受qsc的账户
}

func (tx TxCreateQSC) ValidateData(ctx context.Context) error {
	// chainId 不能为空
	if len(strings.TrimSpace(tx.ChainID)) == 0 {
		return errors.New("chain-id is empty")
	}

	// QSC不存在
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	if qscMapper.Exists(tx.QSCCA.CSR.Subj.CN) {
		return errors.New("qsc already exists")
	}

	// creator账户存在，且有足够的QOS
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	creator := accountMapper.GetAccount(tx.Creator)
	if nil == creator {
		return errors.New("creator account not exists")
	}

	qosAcc, ok := creator.(*account.QOSAccount)
	if !ok {
		return errors.New("creator account is not a QOSAccount")
	}

	if qosAcc.QOS.LT(CreatorQOSLimit) {
		return errors.New("creator account does not have enough qos")
	}

	// CA校验
	if tx.QSCCA.CSR.IsBanker {
		return errors.New("invalid QSC CA, is_banker can not be true")
	}
	baseMapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.MainMapper)
	rootCA := baseMapper.GetRootCA()
	if !VerityCrt(rootCA, tx.QSCCA) {
		return errors.New("invalid QSC CA")
	}
	if !(tx.BankerCA == nil) {
		if !VerityCrt(rootCA, tx.BankerCA) {
			return errors.New("invalid banker CA")
		}
	}

	// accounts校验
	for _, account := range tx.Accounts {
		if account.QOS.NilToZero().GT(btypes.ZeroInt()) {
			return errors.New(fmt.Sprintf("invalid Accounts, %s QOS must be zero", account.AccountAddress))
		}
		if len(account.QSCs) != 1 || account.QSCs[0].Name != tx.QSCCA.CSR.Subj.CN {
			return errors.New(fmt.Sprintf("invalid Accounts, %s len(QSCs) must be 1 and QSCs[0].Name must be %s", account.AccountAddress, tx.QSCCA.CSR.Subj.CN))
		}
		if !account.QSCs[0].Amount.NilToZero().GT(btypes.ZeroInt()) {
			return errors.New(fmt.Sprintf("invalid Accounts, %s QSCs[0].Amount must gt zero", account.AccountAddress))
		}
	}

	return nil
}

func (tx TxCreateQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	// 保存QSC
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	qscMapper.SaveQsc(&tx)

	// 保存QCP配置
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpMapper.SetChainInTrustPubKey(tx.ChainID, tx.QSCCA.CSR.PublicKey)
	qcpMapper.SetMaxChainInSequence(tx.ChainID, 0)
	qcpMapper.SetMaxChainOutSequence(tx.ChainID, 0)

	// 保存账户信息
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	banker := btypes.Address(tx.BankerCA.CSR.PublicKey.Address())
	if nil == accountMapper.GetAccount(banker) {
		accountMapper.SetAccount(accountMapper.NewAccountWithAddress(banker))
	}
	for _, account := range tx.Accounts {
		accountMapper.SetAccount(accountMapper.NewAccountWithAddress(account.AccountAddress))
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
	ret = append(ret, []byte(tx.ChainID)...)
	ret = append(ret, tx.Creator...)
	ret = append(ret, tx.Extrate...)
	ret = append(ret, fmt.Sprint(tx.QSCCA)...)
	ret = append(ret, fmt.Sprint(tx.BankerCA)...)
	ret = append(ret, tx.Description...)

	for _, account := range tx.Accounts {
		ret = append(ret, fmt.Sprint(account)...)
	}

	return
}

// issue QSC
type TxIssueQSC struct {
	QscName string         `json:"qsc_name"` //币名
	Amount  btypes.BigInt  `json:"amount"`   //金额
	Banker  btypes.Address `json:"banker"`   //banker地址
}

func (tx TxIssueQSC) ValidateData(ctx context.Context) error {
	// QscName不能为空
	if len(strings.TrimSpace(tx.QscName)) < 0 {
		return errors.New("qsc-name is empty")
	}

	// Amount大于0
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return errors.New("amount is lte zero")
	}

	// QSC存在
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	qsc := qscMapper.GetQsc(tx.QscName)
	if nil == qsc {
		return errors.New("qsc does not exists")
	}

	// QSC名称一致
	if tx.QscName != qsc.QSCCA.CSR.Subj.CN {
		return errors.New("wrong qsc name")
	}

	// banker 存在
	if nil == qsc.BankerCA {
		return errors.New("banker not exists")
	}

	// banker 地址一致
	if !bytes.Equal(tx.Banker, btypes.Address(qsc.BankerCA.CSR.PublicKey.Address())) {
		return errors.New("wrong banker address")
	}

	return nil
}

func (tx TxIssueQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)

	banker := accountMapper.GetAccount(tx.Banker).(*account.QOSAccount)
	if nil != banker.QSCs {
		banker.QSCs = banker.QSCs.Plus(types.QSCs{btypes.NewBaseCoin(tx.QscName, tx.Amount)})
	} else {
		banker.QSCs = types.QSCs{btypes.NewBaseCoin(tx.QscName, tx.Amount)}
	}
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
	ret = append(ret, []byte(tx.QscName)...)
	ret = append(ret, btypes.Int2Byte(tx.Amount.Int64())...)
	ret = append(ret, []byte(tx.Banker)...)

	return
}
