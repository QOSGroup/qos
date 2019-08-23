package txs

import (
	"github.com/QOSGroup/kepler/cert"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/qcp/mapper"
	"github.com/QOSGroup/qos/module/qcp/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/tendermint/tendermint/crypto"
)

const GasForCreateQCP = uint64(1.8*qtypes.QOSUnit) * qtypes.GasPerUnitCost // 1.8 QOS

// init QCP
type TxInitQCP struct {
	Creator btypes.Address    `json:"creator"` //创建账户
	QCPCA   *cert.Certificate `json:"ca_qcp"`  //CA信息
}

func (tx TxInitQCP) ValidateData(ctx context.Context) error {
	if len(tx.Creator) == 0 {
		return types.ErrInvalidInput(types.DefaultCodeSpace, "")
	}

	// creator账户存在
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	creator := accountMapper.GetAccount(tx.Creator)
	if nil == creator {
		return types.ErrCreatorNotExists(types.DefaultCodeSpace, "")
	}

	// CA 校验
	if tx.QCPCA == nil {
		return types.ErrInvalidQCPCA(types.DefaultCodeSpace, "")
	}
	subj, ok := tx.QCPCA.CSR.Subj.(cert.QCPSubject)
	if !ok {
		return types.ErrInvalidQCPCA(types.DefaultCodeSpace, "")
	}
	if subj.ChainId != ctx.ChainID() {
		return types.ErrInvalidQCPCA(types.DefaultCodeSpace, "")
	}
	if subj.QCPChain == "" {
		return types.ErrInvalidQCPCA(types.DefaultCodeSpace, "")
	}
	rootCA := mapper.GetQCPRootCA(ctx)
	if !cert.VerityCrt([]crypto.PubKey{rootCA}, *tx.QCPCA) {
		return types.ErrWrongQCPCA(types.DefaultCodeSpace, "")
	}

	// 不存在初始化过的QCP信息
	qcpMapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	if pubKey := qcpMapper.GetChainInTrustPubKey(subj.QCPChain); pubKey != nil {
		return types.ErrQCPExists(types.DefaultCodeSpace, "")
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

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeInitQcp,
			btypes.NewAttribute(types.AttributeKeyQcp, subj.QCPChain),
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

func (tx TxInitQCP) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Creator}
}

func (tx TxInitQCP) CalcGas() btypes.BigInt {
	return btypes.NewInt(int64(GasForCreateQCP))
}

func (tx TxInitQCP) GetGasPayer() btypes.Address {
	return tx.Creator
}

func (tx TxInitQCP) GetSignData() (ret []byte) {
	ret = append(ret, tx.Creator...)
	ret = append(ret, Cdc.MustMarshalBinaryBare(tx.QCPCA)...)

	return
}
