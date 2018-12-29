package qcp

import (
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
)

// init QCP
type TxInitQCP struct {
	Creator btypes.Address    `json:"creator"` //创建账户
	QCPCA   *cert.Certificate `json:"ca_qcp"`  //CA信息
}

func (tx TxInitQCP) ValidateData(ctx context.Context) error {
	// creator账户存在
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	creator := accountMapper.GetAccount(tx.Creator)
	if nil == creator {
		return errors.New("creator account not exists")
	}
	_, ok := creator.(*account.QOSAccount)
	if !ok {
		return errors.New("creator account is not a QOSAccount")
	}

	// CA 校验
	if tx.QCPCA == nil {
		return errors.New("QCPCA is empty")
	}
	subj, ok := tx.QCPCA.CSR.Subj.(cert.QCPSubject)
	if !ok {
		return errors.New("invalid QCPSubject")
	}
	if subj.ChainId != ctx.ChainID() {
		return errors.New(fmt.Sprintf("invalid CA，chainId %s not matches %s ", subj.ChainId, ctx.ChainID()))
	}
	if subj.QCPChain == "" {
		return errors.New("invalid CA, QCPChain is empty")
	}
	baseMapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.MainMapper)
	rootCA := baseMapper.GetRootCA()
	if !cert.VerityCrt([]crypto.PubKey{rootCA}, *tx.QCPCA) {
		return errors.New("invalid CA")
	}

	// 不存在初始化过的QCP信息
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	if pubKey := qcpMapper.GetChainInTrustPubKey(subj.QCPChain); pubKey != nil {
		return errors.New(fmt.Sprintf("duplicate chain %s", subj.QCPChain))
	}

	return nil
}

func (tx TxInitQCP) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	subj := tx.QCPCA.CSR.Subj.(cert.QCPSubject)

	// 保存QCP配置
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	qcpMapper.SetChainInTrustPubKey(subj.QCPChain, tx.QCPCA.CSR.PublicKey)
	qcpMapper.SetMaxChainInSequence(subj.QCPChain, 0)
	qcpMapper.SetMaxChainOutSequence(subj.QCPChain, 0)

	return
}

func (tx TxInitQCP) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Creator}
}

func (tx TxInitQCP) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxInitQCP) GetGasPayer() btypes.Address {
	return tx.Creator
}

func (tx TxInitQCP) GetSignData() (ret []byte) {
	ret = append(ret, tx.Creator...)
	ret = append(ret, cdc.MustMarshalBinaryBare(tx.QCPCA)...)

	return
}
