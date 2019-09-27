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

// 文本提议
type TxProposal struct {
	Title          string             `json:"title"`           //  Title of the proposal
	Description    string             `json:"description"`     //  Description of the proposal
	ProposalType   types.ProposalType `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       btypes.AccAddress  `json:"proposer"`        //  Address of the proposer
	InitialDeposit btypes.BigInt      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
}

func NewTxProposal(title, description string, proposer btypes.AccAddress, deposit btypes.BigInt) *TxProposal {
	return &TxProposal{
		Title:          title,
		Description:    description,
		ProposalType:   types.ProposalTypeText,
		Proposer:       proposer,
		InitialDeposit: deposit,
	}
}

var _ txs.ITx = (*TxProposal)(nil)

// 基础数据校验
func (tx TxProposal) ValidateInputs() error {
	// 标题不能为空且不能超过最大长度
	if len(tx.Title) == 0 || len(tx.Title) > MaxTitleLen {
		return types.ErrInvalidInput("invalid title")
	}
	// 描述信息不能为空且不能超过最大长度
	if len(tx.Description) == 0 || len(tx.Description) > MaxDescriptionLen {
		return types.ErrInvalidInput("invalid description")
	}
	// 提议类型校验
	if !types.ValidProposalType(tx.ProposalType) {
		return types.ErrInvalidInput("unknown proposal type")
	}

	return nil
}

// 数据校验
func (tx TxProposal) ValidateData(ctx context.Context) error {
	// 输入校验
	if err := tx.ValidateInputs(); err != nil {
		return err
	}

	// 初始质押不能小于`MinProposerDepositRate`参数这顶值
	govMapper := mapper.GetMapper(ctx)
	params := govMapper.GetLevelParams(ctx, tx.ProposalType.Level())
	if qtypes.NewDecFromInt(tx.InitialDeposit).LT(qtypes.NewDecFromInt(params.MinDeposit).Mul(params.MinProposerDepositRate)) {
		return types.ErrInvalidInput("initial deposit is too small")
	}
	// 提议账户存在且有足够QOS可以质押
	accountMapper := baseabci.GetAccountMapper(ctx)
	account := accountMapper.GetAccount(tx.Proposer)
	if account == nil {
		return types.ErrInvalidInput("proposer not exists")
	}
	if !account.(*qtypes.QOSAccount).EnoughOfQOS(tx.InitialDeposit) {
		return types.ErrInvalidInput("proposer has no enough qos")
	}

	return nil
}

// 交易执行
func (tx TxProposal) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)
	// 保存提议
	textContent := types.NewTextProposal(tx.Title, tx.Description, tx.InitialDeposit)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}
	// 初始化质押
	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	// 发送事件
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeSubmitProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户, Proposer
func (tx TxProposal) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Proposer}
}

// Tx gas, 0
func (tx TxProposal) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, Proposal
func (tx TxProposal) GetGasPayer() btypes.AccAddress {
	return tx.Proposer
}

// 签名字节
func (tx TxProposal) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 提取社区费池，从社区非池中提取QOS到指定账户，仅guardian账户可提交此提议
type TxTaxUsage struct {
	TxProposal                                          // 基础提议信息
	DestAddress btypes.AccAddress `json:"dest_address"` // 接收账户
	Percent     qtypes.Dec        `json:"percent"`      // 提取比例
}

func NewTxTaxUsage(title, description string, proposer btypes.AccAddress, deposit btypes.BigInt, destAddress btypes.AccAddress, percent qtypes.Dec) *TxTaxUsage {
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

// 数据校验
func (tx TxTaxUsage) ValidateInputs() error {
	// 接收地址不能为空
	if len(tx.DestAddress) == 0 {
		return types.ErrInvalidInput("dest_address is empty")
	}
	// 提取
	if tx.Percent.LTE(qtypes.ZeroDec()) {
		return types.ErrInvalidInput("percent lte zero")
	}
	if tx.Percent.GT(qtypes.OneDec()) {
		return types.ErrInvalidInput("percent gte 100%")
	}

	return nil
}

// 数据校验
func (tx TxTaxUsage) ValidateData(ctx context.Context) error {
	// 基础信息校验
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}
	err = tx.ValidateInputs()
	if err != nil {
		return err
	}

	// 接受账户必须是guardian
	if _, exists := guardian.GetMapper(ctx).GetGuardian(tx.DestAddress); !exists {
		return types.ErrInvalidInput("dest_address must be guardian")
	}

	return nil
}

// 交易执行
func (tx TxTaxUsage) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)
	guardianMapper := guardian.GetMapper(ctx)
	// 保存提议
	textContent := types.NewTaxUsageProposal(tx.Title, tx.Description, tx.InitialDeposit, tx.DestAddress, tx.Percent)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)

	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}
	// 初始化质押
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeSubmitProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	// metrics
	guardianMapper.Metrics.Guardian.With(guardian.AddressLabel, tx.Proposer.String(), guardian.OperationLabel, "TxTaxUsage").Set(1)

	return
}

// 签名账户, Proposer
func (tx TxTaxUsage) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Proposer}
}

// Tx gas, 0
func (tx TxTaxUsage) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer
func (tx TxTaxUsage) GetGasPayer() btypes.AccAddress {
	return tx.Proposer
}

// 签名字节
func (tx TxTaxUsage) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 参数修改提议
type TxParameterChange struct {
	TxProposal                           // 基础数据
	Params []types.Param `json:"params"` // 参数变更
}

func NewTxParameterChange(title, description string, proposer btypes.AccAddress, deposit btypes.BigInt, params []types.Param) *TxParameterChange {
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

// 数据校验
func (tx TxParameterChange) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}
	// 参数变更不能为空
	if len(tx.Params) == 0 {
		return types.ErrInvalidInput("params is empty")
	}

	// 不存在质押或投票期参数修改提议
	existsUnfinished := mapper.GetMapper(ctx).ExistsUnfinishedProposals(ctx, types.ProposalTypeParameterChange)
	if existsUnfinished {
		return types.ErrInvalidInput("there are unfinished parameter change proposals")
	}

	paramMapper := params.GetMapper(ctx)
	paramSets := paramMapper.GetParams()
	for _, param := range tx.Params {
		exists := false
		for _, paramSet := range paramSets {
			if param.Module == paramSet.GetParamSpace() {
				exists = true
				// 参数值类型校验
				value, err := paramSet.ValidateKeyValue(param.Key, param.Value)
				if err != nil {
					return err
				}
				// 设置新值
				err = paramSet.SetKeyValue(param.Key, value)
				if err != nil {
					return err
				}
				break
			}
		}
		if !exists {
			return types.ErrInvalidInput(fmt.Sprintf("no params in module:%s", param.Module))
		}
	}

	for _, paramSet := range paramSets {
		// 模块参数整体校验
		if paramSet.Validate() != nil {
			return err
		}
	}

	return nil
}

// 交易执行
func (tx TxParameterChange) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)

	// 保存提议
	textContent := types.NewParameterProposal(tx.Title, tx.Description, tx.InitialDeposit, tx.Params)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	// 初始化质押
	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	// 发送事件
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeSubmitProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户, Proposer
func (tx TxParameterChange) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Proposer}
}

// Tx gas, 0
func (tx TxParameterChange) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer
func (tx TxParameterChange) GetGasPayer() btypes.AccAddress {
	return tx.Proposer
}

// 签名字节
func (tx TxParameterChange) GetSignData() (ret []byte) {
	ret = Cdc.MustMarshalBinaryBare(tx)

	return
}

// 修改通胀提议，已结束和当前通胀阶段不可更改
type TxModifyInflation struct {
	TxProposal
	TotalAmount      btypes.BigInt         `json:"total_amount"`      // 总发行量
	InflationPhrases mint.InflationPhrases `json:"inflation_phrases"` // 通胀阶段
}

func NewTxModifyInflation(title, description string, proposer btypes.AccAddress, deposit btypes.BigInt, totalAmount btypes.BigInt, phrases []mint.InflationPhrase) *TxModifyInflation {
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

// 基础数据检验
func (tx TxModifyInflation) ValidateInputs() error {
	// 校验QOS发行总量
	if !tx.TotalAmount.GT(btypes.ZeroInt()) {
		return types.ErrInvalidInput("total_amount must be positive")
	}

	// 校验通胀规则
	err := tx.InflationPhrases.Valid()
	if err != nil {
		return types.ErrInvalidInput(err.Error())
	}

	return nil
}

// 数据检验
func (tx TxModifyInflation) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}
	err = tx.TxProposal.ValidateInputs()
	if err != nil {
		return err
	}

	applied := mint.GetMapper(ctx).GetAllTotalMintQOSAmount()
	phrases := mint.GetMapper(ctx).MustGetInflationPhrases()
	// 校验当前通胀时间， 当前通胀结束时间 > 当前时间+质押期+投票期 或 当前无通胀
	currentPhrase, exists := phrases.GetPhrase(ctx.BlockHeader().Time.UTC())
	params := mapper.GetMapper(ctx).GetLevelParams(ctx, tx.ProposalType.Level())
	if exists && currentPhrase.EndTime.UTC().Before(ctx.BlockHeader().Time.UTC().Add(params.MaxDepositPeriod).Add(params.VotingPeriod)) {
		return types.ErrInvalidInput("cannot submit proposal at current time")
	}
	// 根据当前通胀规则校验新规则
	err = phrases.ValidNewPhrases(tx.TotalAmount, applied, tx.InflationPhrases)
	if err != nil {
		return types.ErrInvalidInput(err.Error())
	}

	return nil
}

// 交易执行
func (tx TxModifyInflation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)
	// 保存提议
	textContent := types.NewAddInflationPhrase(tx.Title, tx.Description, tx.InitialDeposit, tx.TotalAmount, tx.InflationPhrases)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	// 初始化质押
	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	// 发送事件
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeSubmitProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户, Proposer
func (tx TxModifyInflation) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Proposer}
}

// Tx gas, 0
func (tx TxModifyInflation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer
func (tx TxModifyInflation) GetGasPayer() btypes.AccAddress {
	return tx.Proposer
}

// 签名字节
func (tx TxModifyInflation) GetSignData() (ret []byte) {
	ret, _ = Cdc.MarshalBinaryBare(tx)

	return
}

// 软件升级提议
type TxSoftwareUpgrade struct {
	TxProposal
	Version       string `json:"version"`         // qosd version
	DataHeight    int64  `json:"data_height"`     // data version
	GenesisFile   string `json:"genesis_file"`    // url of genesis file
	GenesisMD5    string `json:"genesis_md5"`     // signature of genesis.json
	ForZeroHeight bool   `json:"for_zero_height"` // restart from zero height
}

func NewTxSoftwareUpgrade(title, description string, proposer btypes.AccAddress, deposit btypes.BigInt,
	version string, dataHeight int64, genesisFile string, genesisMd5 string, forZeroHeight bool) *TxSoftwareUpgrade {
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

// 基础数据校验
func (tx TxSoftwareUpgrade) ValidateInputs() error {
	// 版本信息不能为空
	if len(tx.Version) == 0 {
		return types.ErrInvalidInput("Version is empty")
	}

	// 从0高度升级时，DataHeight大于0，GenesisFile和GenesisMD5不能为空
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

// 数据校验
func (tx TxSoftwareUpgrade) ValidateData(ctx context.Context) error {
	// 基础数据校验
	err := tx.TxProposal.ValidateData(ctx)
	if err != nil {
		return err
	}
	err = tx.ValidateInputs()
	if err != nil {
		return err
	}

	return nil
}

// 交易执行
func (tx TxSoftwareUpgrade) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.CodeOK,
	}

	govMapper := mapper.GetMapper(ctx)

	// 保存提议
	textContent := types.NewSoftwareUpgradeProposal(tx.Title, tx.Description, tx.InitialDeposit,
		tx.Version, tx.DataHeight, tx.GenesisFile, tx.GenesisMD5, tx.ForZeroHeight)
	proposal, err := govMapper.SubmitProposal(ctx, textContent)
	if err != nil {
		result = btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}
	}

	// 初始化质押
	govMapper.AddDeposit(ctx, proposal.ProposalID, tx.Proposer, tx.InitialDeposit)

	// 发送事件
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
			btypes.NewAttribute(btypes.AttributeKeyAction, types.EventTypeSubmitProposal),
			btypes.NewAttribute(btypes.AttributeKeyGasPayer, tx.GetGasPayer().String()),
		),
	}

	return
}

// 签名账户, Proposer
func (tx TxSoftwareUpgrade) GetSigner() []btypes.AccAddress {
	return []btypes.AccAddress{tx.Proposer}
}

// Tx gas, 0
func (tx TxSoftwareUpgrade) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

// Gas payer, Proposer
func (tx TxSoftwareUpgrade) GetGasPayer() btypes.AccAddress {
	return tx.Proposer
}

// 签名字节
func (tx TxSoftwareUpgrade) GetSignData() (ret []byte) {
	Cdc.MustMarshalBinaryBare(tx)

	return
}
