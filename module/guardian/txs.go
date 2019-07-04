package guardian

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/guardian/types"
)

const (
	MaxDescriptionLen = 1000
)

type TxAddGuardian struct {
	Description string         `json:"description"`
	Address     btypes.Address `json:"address"`
	Creator     btypes.Address `json:"creator"`
}

func NewTxAddGuardian(description string, address, creator btypes.Address) *TxAddGuardian {
	return &TxAddGuardian{
		Description: description,
		Address:     address,
		Creator:     creator,
	}
}

var _ txs.ITx = (*TxAddGuardian)(nil)

func (tx TxAddGuardian) ValidateData(ctx context.Context) error {
	if len(tx.Description) > MaxDescriptionLen {
		return ErrInvalidInput("Description is too long")
	}

	if len(tx.Address) == 0 {
		return ErrInvalidInput("Address is empty")
	}

	if len(tx.Creator) == 0 {
		return ErrInvalidInput("Creator is empty")
	}

	mapper := GetGuardianMapper(ctx)
	if _, exists := mapper.GetGuardian(tx.Address); exists {
		return ErrGuardianAlreadyExists("")
	}

	guardian, exists := mapper.GetGuardian(tx.Creator)
	if !exists || guardian.GuardianType != gtypes.Genesis {
		return ErrInvalidCreator("Creator not exists or not init from genesis")
	}

	return nil
}

func (tx TxAddGuardian) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	GetGuardianMapper(ctx).AddGuardian(*gtypes.NewGuardian(tx.Description, gtypes.Ordinary, tx.Address, tx.Creator))

	result.Events = btypes.Events{
		btypes.NewEvent(
			EventTypeAddGuardian,
			btypes.NewAttribute(AttributeKeyCreator, tx.Creator.String()),
			btypes.NewAttribute(AttributeKeyGuardian, tx.Address.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxAddGuardian) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Creator}
}

func (tx TxAddGuardian) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxAddGuardian) GetGasPayer() btypes.Address {
	return tx.Creator
}

func (tx TxAddGuardian) GetSignData() (ret []byte) {
	ret = append(ret, tx.Description...)
	ret = append(ret, tx.Address...)
	ret = append(ret, tx.Creator...)

	return
}

type TxDeleteGuardian struct {
	Address   btypes.Address `json:"address"`    // this guardian's address
	DeletedBy btypes.Address `json:"deleted_by"` // address that initiated the AddGuardian tx
}

func NewTxDeleteGuardian(address, deletedBy btypes.Address) *TxDeleteGuardian {
	return &TxDeleteGuardian{
		Address:   address,
		DeletedBy: deletedBy,
	}
}

var _ txs.ITx = (*TxDeleteGuardian)(nil)

func (tx TxDeleteGuardian) ValidateData(ctx context.Context) error {
	if len(tx.Address) == 0 {
		return ErrInvalidInput("Address is empty")
	}

	if len(tx.DeletedBy) == 0 {
		return ErrInvalidInput("DeletedBy is empty")
	}

	mapper := GetGuardianMapper(ctx)
	guardian, exists := mapper.GetGuardian(tx.Address)
	if !exists {
		return ErrUnKnownGuardian("")
	}

	if guardian.GuardianType == gtypes.Genesis {
		return ErrInvalidInput("can not delete genesis guardian")
	}

	deletedBy, exists := mapper.GetGuardian(tx.DeletedBy)
	if !exists || deletedBy.GuardianType != gtypes.Genesis {
		return ErrInvalidCreator("DeletedBy not exists or not init from genesis")
	}

	return nil
}

func (tx TxDeleteGuardian) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	GetGuardianMapper(ctx).DeleteGuardian(tx.Address)

	result.Events = btypes.Events{
		btypes.NewEvent(
			EventTypeDeleteGuardian,
			btypes.NewAttribute(AttributeKeyDeleteBy, tx.DeletedBy.String()),
			btypes.NewAttribute(AttributeKeyGuardian, tx.Address.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}
	
	return
}

func (tx TxDeleteGuardian) GetSigner() []btypes.Address {
	return []btypes.Address{tx.DeletedBy}
}

func (tx TxDeleteGuardian) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxDeleteGuardian) GetGasPayer() btypes.Address {
	return tx.DeletedBy
}

func (tx TxDeleteGuardian) GetSignData() (ret []byte) {
	ret = append(ret, tx.Address...)
	ret = append(ret, tx.DeletedBy...)

	return
}
