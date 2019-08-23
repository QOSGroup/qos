package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/mapper"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/params"
	qtypes "github.com/QOSGroup/qos/types"
)

const (
	MaxTitleLen       = 200
	MaxDescriptionLen = 1000
)

type TxProposal struct {
	Title          string             `json:"title"`           //  Title of the proposal
	Description    string             `json:"description"`     //  Description of the proposal
	ProposalType   types.ProposalType `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       btypes.Address     `json:"proposer"`        //  Address of the proposer
	InitialDeposit uint64             `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
}

func NewTxProposal(title, description string, proposer btypes.Address, deposit uint64) *TxProposal {
	return &TxProposal{
		Title:          title,
		Description:    description,
		ProposalType:   types.ProposalTypeText,
		Proposer:       proposer,
		InitialDeposit: deposit,
	}
}

var _ txs.ITx = (*TxProposal)(nil)

func (tx TxProposal) ValidateData(ctx context.Context) error {
	if len(tx.Title) == 0 || len(tx.Title) > MaxTitleLen {
		return types.ErrInvalidInput("invalid title")
	}
	if len(tx.Description) == 0 || len(tx.Description) > MaxDescriptionLen {
		return types.ErrInvalidInput("invalid description")
	}
	if !types.ValidProposalType(tx.ProposalType) {
		return types.ErrInvalidInput("unknown proposal type")
	}

	govMapper := mapper.GetMapper(ctx)
	params := govMapper.GetParams(ctx)
	if qtypes.NewDec(int64(tx.InitialDeposit)).LT(qtypes.NewDec(int64(params.MinDeposit)).Mul(params.MinProposerDepositRate)) {
		return types.ErrInvalidInput("initial deposit is too small")
	}

	accountMapper := baseabci.GetAccountMapper(ctx)
	account := accountMapper.GetAccount(tx.Proposer)
	if account == nil {
		return types.ErrInvalidInput("proposer not exists")
	}

	if !account.(*qtypes.QOSAccount).EnoughOfQOS(btypes.NewInt(int64(tx.InitialDeposit))) {
		return types.ErrInvalidInput("proposer has no enough qos")
	}

	return nil
}

func (tx TxProposal) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)

	textContent := types.NewTextProposal(tx.Title, tx.Description, tx.InitialDeposit)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeSubmitProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyProposer, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyProposalType, tx.ProposalType.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

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
	ret = append(ret, qtypes.Uint64ToBigEndian(tx.InitialDeposit)...)

	return
}

type TxTaxUsage struct {
	TxProposal
	DestAddress btypes.Address `json:"dest_address"`
	Percent     qtypes.Dec     `json:"percent"`
}

func NewTxTaxUsage(title, description string, proposer btypes.Address, deposit uint64, destAddress btypes.Address, percent qtypes.Dec) *TxTaxUsage {
	return &TxTaxUsage{
		TxProposal: TxProposal{
			Title:          title,
			Description:    description,
			ProposalType:   types.ProposalTypeTaxUsage,
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
		return types.ErrInvalidInput("DestAddress is empty")
	}

	if tx.Percent.LTE(qtypes.ZeroDec()) {
		return types.ErrInvalidInput("Percent lte zero")
	}

	// 接受账户必须是guardian
	if _, exists := guardian.GetMapper(ctx).GetGuardian(tx.DestAddress); !exists {
		return types.ErrInvalidInput("DestAddress must be guardian")
	}

	return nil
}

func (tx TxTaxUsage) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)

	textContent := types.NewTaxUsageProposal(tx.Title, tx.Description, tx.InitialDeposit, tx.DestAddress, tx.Percent)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeSubmitProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyProposer, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyProposalType, tx.ProposalType.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

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
	Params []types.Param `json:"params"`
}

func NewTxParameterChange(title, description string, proposer btypes.Address, deposit uint64, params []types.Param) *TxParameterChange {
	return &TxParameterChange{
		TxProposal: TxProposal{
			Title:          title,
			Description:    description,
			ProposalType:   types.ProposalTypeParameterChange,
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
		return types.ErrInvalidInput("Params is empty")
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

	govMapper := mapper.GetMapper(ctx)

	textContent := types.NewParameterProposal(tx.Title, tx.Description, tx.InitialDeposit, tx.Params)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeSubmitProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyProposer, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyProposalType, tx.ProposalType.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

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

type TxModifyInflation struct {
	TxProposal
	TotalAmount      uint64                `json:"total_amount"`      // 总发行量
	InflationPhrases mint.InflationPhrases `json:"inflation_phrases"` // 通胀阶段
}

func NewTxModifyInflation(title, description string, proposer btypes.Address, deposit uint64, totalAmount uint64, phrases []mint.InflationPhrase) *TxModifyInflation {
	return &TxModifyInflation{
		TxProposal: TxProposal{
			Title:          title,
			Description:    description,
			ProposalType:   types.ProposalTypeParameterChange,
			Proposer:       proposer,
			InitialDeposit: deposit,
		},
		TotalAmount:      totalAmount,
		InflationPhrases: phrases,
	}
}

func (tx TxModifyInflation) ValidateData(ctx context.Context) error {
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}

	// 校验QOS发行总量
	if tx.TotalAmount <= 0 {
		return types.ErrInvalidInput("TotalAmount must be positive")
	}

	// 校验通胀
	err = tx.InflationPhrases.Valid()
	if err != nil {
		return types.ErrInvalidInput(err.Error())
	}
	applied := mint.GetMapper(ctx).GetAllTotalMintQOSAmount()
	phrases := mint.GetMapper(ctx).MustGetInflationPhrases()
	// 校验当前通胀时间， 当前通胀结束时间 > 当前时间+质押期+投票期 或 当前无通胀
	currentPhrase, exists := phrases.GetPhrase(ctx.BlockHeader().Time.UTC())
	params := mapper.GetMapper(ctx).GetParams(ctx)
	if exists && currentPhrase.EndTime.UTC().Before(ctx.BlockHeader().Time.UTC().Add(params.MaxDepositPeriod).Add(params.VotingPeriod)) {
		return types.ErrInvalidInput("cannot submit proposal at current time")
	}
	err = phrases.ValidNewPhrases(tx.TotalAmount, applied, tx.InflationPhrases)
	if err != nil {
		return types.ErrInvalidInput(err.Error())
	}

	return nil
}

func (tx TxModifyInflation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)

	textContent := types.NewAddInflationPhrase(tx.Title, tx.Description, tx.InitialDeposit, tx.TotalAmount, tx.InflationPhrases)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeSubmitProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyProposer, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyProposalType, tx.ProposalType.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxModifyInflation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Proposer}
}

func (tx TxModifyInflation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxModifyInflation) GetGasPayer() btypes.Address {
	return tx.Proposer
}

func (tx TxModifyInflation) GetSignData() (ret []byte) {
	ret, _ = Cdc.MarshalBinaryBare(tx)
	return
}

type TxSoftwareUpgrade struct {
	TxProposal
	Version       string `json:"version"`         // qosd version
	DataHeight    uint64 `json:"data_height"`     // data version
	GenesisFile   string `json:"genesis_file"`    // url of genesis file
	GenesisMD5    string `json:"genesis_md5"`     // signature of genesis.json
	ForZeroHeight bool   `json:"for_zero_height"` // restart from zero height
}

func NewTxSoftwareUpgrade(title, description string, proposer btypes.Address, deposit uint64,
	version string, dataHeight uint64, genesisFile string, genesisMd5 string, forZeroHeight bool) *TxSoftwareUpgrade {
	return &TxSoftwareUpgrade{
		TxProposal: TxProposal{
			Title:          title,
			Description:    description,
			ProposalType:   types.ProposalTypeSoftwareUpgrade,
			Proposer:       proposer,
			InitialDeposit: deposit,
		},
		Version:       version,
		DataHeight:    dataHeight,
		GenesisFile:   genesisFile,
		GenesisMD5:    genesisMd5,
		ForZeroHeight: forZeroHeight,
	}
}

var _ txs.ITx = (*TxSoftwareUpgrade)(nil)

func (tx TxSoftwareUpgrade) ValidateData(ctx context.Context) error {
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}

	if len(tx.Version) == 0 {
		return types.ErrInvalidInput("Version is empty")
	}

	if tx.ForZeroHeight {
		if tx.DataHeight <= 0 {
			return types.ErrInvalidInput("DataHeight must be positive")
		}

		if len(tx.GenesisFile) == 0 {
			return types.ErrInvalidInput("GenesisFile is empty")
		}

		if len(tx.GenesisMD5) == 0 {
			return types.ErrInvalidInput("GenesisFileMD5 is empty")
		}
	}

	return nil
}

func (tx TxSoftwareUpgrade) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)

	textContent := types.NewSoftwareUpgradeProposal(tx.Title, tx.Description, tx.InitialDeposit,
		tx.Version, tx.DataHeight, tx.GenesisFile, tx.GenesisMD5, tx.ForZeroHeight)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	result.Events = btypes.Events{
		btypes.NewEvent(
			types.EventTypeSubmitProposal,
			btypes.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
			btypes.NewAttribute(types.AttributeKeyProposer, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyDepositor, tx.Proposer.String()),
			btypes.NewAttribute(types.AttributeKeyProposalType, tx.ProposalType.String()),
		),
		btypes.NewEvent(
			btypes.EventTypeMessage,
			btypes.NewAttribute(btypes.AttributeKeyModule, types.AttributeKeyModule),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetSigner()[0].String()),
		),
	}

	return
}

func (tx TxSoftwareUpgrade) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Proposer}
}

func (tx TxSoftwareUpgrade) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxSoftwareUpgrade) GetGasPayer() btypes.Address {
	return tx.Proposer
}

func (tx TxSoftwareUpgrade) GetSignData() (ret []byte) {
	Cdc.MustMarshalBinaryBare(tx)
	return
}
