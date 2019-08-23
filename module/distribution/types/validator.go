package types

import (
	btypes "github.com/QOSGroup/qbase/types"
)

//ValidatorCurrentPeriodSummary validator当前周期收益信息
type ValidatorCurrentPeriodSummary struct {
	Fees   btypes.BigInt `json:"fees"`
	Period uint64        `json:"period"`
}

//ValidatorEcoFeePool validator收益信息
type ValidatorEcoFeePool struct {
	ProposerTotalRewardFee      btypes.BigInt `json:"proposerTotalRewardFee"`      //validator 通过proposer获取的总收益
	CommissionTotalRewardFee    btypes.BigInt `json:"commissionTotalRewardFee"`    //validator 通过投票获取的佣金总收益
	PreDistributeTotalRewardFee btypes.BigInt `json:"preDistributeTotalRewardFee"` //validator 通过投票获取的待分配金额总收益
	PreDistributeRemainTotalFee btypes.BigInt `json:"preDistributeRemainTotalFee"` //validator 待分配金额中剩余的收益
}

func NewValidatorEcoFeePool() ValidatorEcoFeePool {
	return ValidatorEcoFeePool{
		ProposerTotalRewardFee:      btypes.ZeroInt(),
		CommissionTotalRewardFee:    btypes.ZeroInt(),
		PreDistributeTotalRewardFee: btypes.ZeroInt(),
		PreDistributeRemainTotalFee: btypes.ZeroInt(),
	}
}
