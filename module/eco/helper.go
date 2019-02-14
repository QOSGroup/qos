package eco

import (
	"fmt"

	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qtypes "github.com/QOSGroup/qos/types"
)

func IncrAccountQOS(ctx context.Context, addr btypes.Address, amount btypes.BigInt) error {
	accountMapper := baseabci.GetAccountMapper(ctx)

	acc := accountMapper.GetAccount(addr)
	if qosAcc, ok := acc.(*qtypes.QOSAccount); ok {
		err := qosAcc.SetQOS(qosAcc.GetQOS().NilToZero().Add(amount))
		if err != nil {
			return err
		}
		accountMapper.SetAccount(acc)
		return nil
	}

	return fmt.Errorf("addr: %s not a QOSAccount", addr)
}

func DecrAccountQOS(ctx context.Context, addr btypes.Address, amount btypes.BigInt) error {
	accountMapper := baseabci.GetAccountMapper(ctx)

	acc := accountMapper.GetAccount(addr)
	if qosAcc, ok := acc.(*qtypes.QOSAccount); ok {
		current := qosAcc.GetQOS().NilToZero()
		if current.LT(amount) {
			return fmt.Errorf("addr: %s has not much OQS to decrease. expect: %d , actual: %d", addr, amount, current)
		}

		err := qosAcc.SetQOS(current.Sub(amount))
		if err != nil {
			return err
		}
		accountMapper.SetAccount(acc)
		return nil
	}

	return fmt.Errorf("addr: %s not a QOSAccount", addr)
}
