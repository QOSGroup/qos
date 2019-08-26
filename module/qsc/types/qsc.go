package types

import (
	"github.com/QOSGroup/kepler/cert"
	btypes "github.com/QOSGroup/qbase/types"
)

type Info struct {
	Name        string         `json:"name"`         //币名
	ChainId     string         `json:"chain_id"`     //证书可用链
	Extrate     string         `json:"extrate"`      //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	Description string         `json:"description"`  //描述信息
	Banker      btypes.Address `json:"banker"`       //Banker PubKey
	TotalAmount btypes.BigInt  `json:"total_amount"` //发行总量
}

func NewInfoWithQSCCA(cer *cert.Certificate) Info {
	subj := cer.CSR.Subj.(cert.QSCSubject)
	var banker btypes.Address
	if subj.Banker != nil {
		banker = btypes.Address(subj.Banker.Address())
	}
	return Info{
		Name:        subj.Name,
		ChainId:     subj.ChainId,
		Banker:      banker,
		TotalAmount: btypes.ZeroInt(),
	}
}
