package types

import (
	"github.com/QOSGroup/kepler/cert"
	btypes "github.com/QOSGroup/qbase/types"
)

// 代币信息
type QSCInfo struct {
	Name         string            `json:"name"`          //币名
	ChainId      string            `json:"chain_id"`      //证书可用链
	ExchangeRate string            `json:"exchange_rate"` //qcs:qos汇率
	Description  string            `json:"description"`   //描述信息
	Banker       btypes.AccAddress `json:"banker"`        //Banker PubKey
	TotalAmount  btypes.BigInt     `json:"total_amount"`  //发行总量
}

func NewInfoWithQSCCA(cer *cert.Certificate) QSCInfo {
	subj := cer.CSR.Subj.(cert.QSCSubject)
	var banker btypes.AccAddress
	if subj.Banker != nil {
		banker = btypes.AccAddress(subj.Banker.Address())
	}
	return QSCInfo{
		Name:        subj.Name,
		ChainId:     subj.ChainId,
		Banker:      banker,
		TotalAmount: btypes.ZeroInt(),
	}
}
