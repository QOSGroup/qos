package txs

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank"
	"github.com/QOSGroup/qos/module/qcp/mapper"
	"github.com/QOSGroup/qos/module/qcp/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

const GasForCreateQCP = int64(1.8*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 1.8 QOS

// 初始化联盟链Tx
type TxInitQCP struct {
	Creator btypes.AccAddress `json:"creator"` //创建账户
	QCPCA   *cert.Certificate `json:"ca_qcp"`  //CA信息
}

// 校验基础数据
func (tx TxInitQCP) ValidateInputs() error {
	// 校验创建账户
	if len(tx.Creator) == 0 {
		return types.ErrEmptyCreator()
	}

	// 校验证书
	if tx.QCPCA == nil {
		return types.ErrInvalidQCPCA()
	}
	subj, ok := tx.QCPCA.CSR.Subj.(cert.QCPSubject)
	if !ok {
		return types.ErrInvalidQCPCA()
	}
	if len(subj.QCPChain) == 0 {
		return types.ErrInvalidQCPCA()
	}

	return nil
}

func (tx TxInitQCP) ValidateData(ctx context.Context) error {
	// 校验基础数据
	err := tx.ValidateInputs()
	if err != nil {
		return err
	}

	bankMapper := bank.GetMapper(ctx)
	creator := bankMapper.GetAccount(tx.Creator)
	if nil == creator {
		return types.ErrCreatorNotExists()
	}

	subj, _ := tx.QCPCA.CSR.Subj.(cert.QCPSubject)
	if subj.ChainId != ctx.ChainID() {
		return types.ErrInvalidQCPCA()
	}
	rootCA := mapper.GetRootCaPubkey(ctx)
	if rootCA == nil || len(rootCA.Bytes()) == 0 {
		return types.ErrRootCANotConfigure()
	}
	if !cert.VerityCrt([]crypto.PubKey{rootCA}, *tx.QCPCA) {
		return types.ErrWrongQCPCA()
	}

	// 校验已存在的联盟链信息
	qcpMapper := mapper.GetMapper(ctx)
	if pubKey := qcpMapper.GetChainInTrustPubKey(subj.QCPChain); pubKey != nil {
		return types.ErrQCPExists()
	}

	return nil
}

func (tx TxInitQCP) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	// 联盟链基础信息
	subj := tx.QCPCA.CSR.Subj.(cert.QCPSubject)

	// 保存/初始化联盟链信息
	qcpMapper := ctx.Mapper(qcp.MapperName).(*qcp.QcpMapper)
	qcpMapper.SetChainInTrustPubKey(subj.QCPChain, tx.QCPCA.CSR.PublicKey)
	qcpMapper.SetMaxChainInSequence(subj.QCPChain, 0)
	qcpMapper.SetMaxChainOutSequence(subj.QCPChain, 0)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeInitQcp,
			btypes.NewAttribute(types.AttributeKeyQcp, subj.QCPChain),
			btypes.NewAttribute(types.AttributeKeyCreator, tx.Creator.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeInitQcp),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxInitQCP) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Creator}
}

func (tx TxInitQCP) CalcGas() btypes.BigInt {
	return btypes.NewInt(GasForCreateQCP)
}

func (tx TxInitQCP) GetGasPayer() btypes.AccAddress {
	return tx.Creator
}

func (tx TxInitQCP) GetSignData() (ret []byte) {
	ret = append(ret, Cdc.MustMarshalBinaryBare(tx)...)

	return
}
