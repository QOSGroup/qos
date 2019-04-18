package gov

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/types"
)

const (
	MaxTitleLen       = 200
	MaxDescriptionLen = 1000
)

type TxProposal struct {
	Title          string              `json:"title"`           //  Title of the proposal
	Description    string              `json:"description"`     //  Description of the proposal
	ProposalType   gtypes.ProposalType `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       btypes.Address      `json:"proposer"`        //  Address of the proposer
	InitialDeposit uint64              `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
}

func NewTxProposal(title, description string, proposer btypes.Address, deposit uint64) *TxProposal {
	return &TxProposal{
		Title:          title,
		Description:    description,
		ProposalType:   gtypes.ProposalTypeText,
		Proposer:       proposer,
		InitialDeposit: deposit,
	}
}

var _ txs.ITx = (*TxProposal)(nil)

func (tx TxProposal) ValidateData(ctx context.Context) error {
	if len(tx.Title) == 0 || len(tx.Title) > MaxTitleLen {
		return ErrInvalidInput("invalid title")
	}
	if len(tx.Description) == 0 || len(tx.Description) > MaxDescriptionLen {
		return ErrInvalidInput("invalid description")
	}
	if !gtypes.ValidProposalType(tx.ProposalType) {
		return ErrInvalidInput("unknown proposal type")
	}

	govMapper := GetGovMapper(ctx)
	if types.NewDec(int64(tx.InitialDeposit)).LT(types.NewDec(int64(govMapper.GetParams(ctx).MinDeposit)).Mul(MinDepositRate)) {
		return ErrInvalidInput("initial deposit is too small")
	}

	accountMapper := baseabci.GetAccountMapper(ctx)
	account := accountMapper.GetAccount(tx.Proposer)
	if account == nil {
		return ErrInvalidInput("proposer not exists")
	}

	if !account.(*types.QOSAccount).EnoughOfQOS(btypes.NewInt(int64(tx.InitialDeposit))) {
		return ErrInvalidInput("proposer has no enough qos")
	}

	return nil
}

func (tx TxProposal) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := GetGovMapper(ctx)

	textContent := gtypes.NewTextProposal(tx.Title, tx.Description, tx.InitialDeposit)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	return
}

func (tx TxProposal) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Proposer}
}

func (tx TxProposal) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxProposal) GetGasPayer() btypes.Address {
	return tx.Proposer
}

func (tx TxProposal) GetSignData() (ret []byte) {
	ret = append(ret, tx.Title...)
	ret = append(ret, tx.Description...)
	ret = append(ret, byte(tx.ProposalType))
	ret = append(ret, tx.Proposer...)
	ret = append(ret, types.Uint64ToBigEndian(tx.InitialDeposit)...)

	return
}

type TxTaxUsage struct {
	TxProposal
	DestAddress btypes.Address `json:"dest_address"`
	Percent     types.Dec      `json:"percent"`
}

func NewTxTaxUsage(title, description string, proposer btypes.Address, deposit uint64, destAddress btypes.Address, percent types.Dec) *TxTaxUsage {
	return &TxTaxUsage{
		TxProposal: TxProposal{
			Title:          title,
			Description:    description,
			ProposalType:   gtypes.ProposalTypeTaxUsage,
			Proposer:       proposer,
			InitialDeposit: deposit,
		},
		DestAddress: destAddress,
		Percent:     percent,
	}
}

var _ txs.ITx = (*TxProposal)(nil)

func (tx TxTaxUsage) ValidateData(ctx context.Context) error {
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}

	if len(tx.DestAddress) == 0 {
		return ErrInvalidInput("DestAddress is empty")
	}

	if tx.Percent.LTE(types.ZeroDec()) {
		return ErrInvalidInput("Percent lte zero")
	}

	// 接受账户必须是guardian
	if _, exists := guardian.GetGuardianMapper(ctx).GetGuardian(tx.DestAddress); !exists {
		return ErrInvalidInput("DestAddress must be guardian")
	}

	return nil
}

func (tx TxTaxUsage) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := GetGovMapper(ctx)

	textContent := gtypes.NewTaxUsageProposal(tx.Title, tx.Description, tx.InitialDeposit, tx.DestAddress, tx.Percent)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	return
}

func (tx TxTaxUsage) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Proposer}
}

func (tx TxTaxUsage) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxTaxUsage) GetGasPayer() btypes.Address {
	return tx.Proposer
}

func (tx TxTaxUsage) GetSignData() (ret []byte) {
	ret = append(ret, tx.TxProposal.GetSignData()...)
	ret = append(ret, tx.DestAddress...)
	ret = append(ret, tx.Percent.String()...)

	return
}

type TxParameterChange struct {
	TxProposal
	Params []gtypes.Param `json:"params"`
}

func NewTxParameterChange(title, description string, proposer btypes.Address, deposit uint64, params []gtypes.Param) *TxParameterChange {
	return &TxParameterChange{
		TxProposal: TxProposal{
			Title:          title,
			Description:    description,
			ProposalType:   gtypes.ProposalTypeParameterChange,
			Proposer:       proposer,
			InitialDeposit: deposit,
		},
		Params: params,
	}
}

var _ txs.ITx = (*TxProposal)(nil)

func (tx TxParameterChange) ValidateData(ctx context.Context) error {
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}

	if len(tx.Params) == 0 {
		return ErrInvalidInput("Params is empty")
	}

	paramMapper := params.GetMapper(ctx)
	for _, param := range tx.Params {
		if err = paramMapper.Validate(param.Module, param.Key, param.Value); err != nil {
			return err
		}
	}

	return nil
}

func (tx TxParameterChange) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := GetGovMapper(ctx)

	textContent := gtypes.NewParameterProposal(tx.Title, tx.Description, tx.InitialDeposit, tx.Params)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	return
}

func (tx TxParameterChange) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Proposer}
}

func (tx TxParameterChange) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxParameterChange) GetGasPayer() btypes.Address {
	return tx.Proposer
}

func (tx TxParameterChange) GetSignData() (ret []byte) {
	ret = append(ret, tx.TxProposal.GetSignData()...)
	for _, param := range tx.Params {
		ret = append(ret, param.String()...)
	}
	return
}
