package stake

import (
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/eco"
	"github.com/QOSGroup/qos/module/eco/mapper"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
)

type TxCreateDelegation struct {
	Delegator      btypes.Address
	ValidatorOwner btypes.Address
	Amount         uint64
	IsCompound     bool
}

var _ txs.ITx = (*TxCreateDelegation)(nil)

func (tx *TxCreateDelegation) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 || len(tx.ValidatorOwner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	// TODO: 是否应该在tx里做这种检查
	if tx.Amount == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Delegation amount must be a positive integer.")
	}

	_, err = validateValidator(ctx, tx.ValidatorOwner, true, staketypes.Active, true)
	if nil != err {
		return err
	}

	err = validateQOSAccount(ctx, tx.Delegator, tx.Amount)
	if nil != err {
		return err
	}

	err = validateQOSAccount(ctx, tx.ValidatorOwner, 0)
	if nil != err {
		return err
	}

	return nil
}

//创建或新增委托
//0. delegator账户扣减QOS
//1. validator增加周期 , 计算周期段内delegator收益,并更新收益信息
//2. 更新delegation信息
//3. 更新validator的bondTokens
func (tx *TxCreateDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	currentHeight := uint64(ctx.BlockHeight())
	ecoMapper := mapper.GetEcoMapper(ctx)

	validator, exists := ecoMapper.ValidatorMapper.GetValidatorByOwner(tx.ValidatorOwner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: "validator not exsits"}, nil
	}

	valAddr := btypes.Address(validator.ValidatorPubKey.Address())
	delegatorAddr := tx.Delegator

	//0. delegator账户扣减QOS, amount:qos = 1:1
	decrQOS := btypes.NewInt(int64(tx.Amount))
	err := eco.DecrAccountQOS(ctx, delegatorAddr, decrQOS)
	if err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	//获取delegation信息,若delegation信息不存在,则初始化degelator收益信息
	delegatedAmount := uint64(0)
	info, exsits := ecoMapper.DelegationMapper.GetDelegationInfo(delegatorAddr, valAddr)
	if !exsits {
		ecoMapper.DistributionMapper.InitDelegatorIncomeInfo(valAddr, delegatorAddr, uint64(0), currentHeight)
		info = staketypes.NewDelegationInfo(delegatorAddr, valAddr, uint64(0), false)
	} else {
		delegatedAmount = info.Amount
	}

	updatedAmount := delegatedAmount + tx.Amount
	//1. validator增加周期 , 计算周期段内delegator收益,并更新收益信息
	err = ecoMapper.DistributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, updatedAmount, currentHeight)
	if err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	//2. 更新delegation信息
	info.Amount = updatedAmount
	info.IsCompound = tx.IsCompound
	ecoMapper.DelegationMapper.SetDelegationInfo(info)

	//3. 更新validator的bondTokens, amount:token = 1:1
	validatorAddToken := tx.Amount
	updatedValidatorTokens := validator.BondTokens + validatorAddToken
	ecoMapper.ValidatorMapper.ChangeValidatorBondTokens(validator, updatedValidatorTokens)

	return btypes.Result{Code: btypes.CodeOK}, nil
}

func (tx *TxCreateDelegation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxCreateDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateDelegation) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Delegator)
}

func (tx *TxCreateDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorOwner...)
	ret = append(ret, btypes.Int2Byte(int64(tx.Amount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}

type TxModifyCompound struct {
	Delegator      btypes.Address
	ValidatorOwner btypes.Address
	IsCompound     bool
}

var _ txs.ITx = (*TxModifyCompound)(nil)

func (tx *TxModifyCompound) ValidateData(ctx context.Context) (err error) {

	if len(tx.Delegator) == 0 || len(tx.ValidatorOwner) == 0 {
		return ErrInvalidInput(DefaultCodeSpace, "Validator and Delegator must be specified.")
	}

	// TODO:是否允许validator为inactive/jailed时修改
	validator, err := validateValidator(ctx, tx.ValidatorOwner, true, staketypes.Active, true)
	if nil != err {
		return err
	}

	valAddr := btypes.Address(validator.ValidatorPubKey.Address())

	delegationMapper := mapper.GetDelegationMapper(ctx)
	info, exsits := delegationMapper.GetDelegationInfo(tx.Delegator, valAddr)
	if !exsits {
		return ErrInvalidInput(DefaultCodeSpace, "delegator not delegate the owner's validator")
	}

	if info.IsCompound == tx.IsCompound {
		return ErrInvalidInput(DefaultCodeSpace, "delegator's compound not changed")
	}

	return nil
}

//修改收益单复利
func (tx *TxModifyCompound) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	currentHeight := uint64(ctx.BlockHeight())
	ecoMapper := mapper.GetEcoMapper(ctx)

	validator, exists := ecoMapper.ValidatorMapper.GetValidatorByOwner(tx.ValidatorOwner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: "validator not exsits"}, nil
	}

	valAddr := btypes.Address(validator.ValidatorPubKey.Address())
	delegatorAddr := tx.Delegator
	info, exsits := ecoMapper.DelegationMapper.GetDelegationInfo(delegatorAddr, valAddr)
	if !exsits {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: "delegator not delegate the owner's validator"}, nil
	}

	//1. 计算delegator收益
	err := ecoMapper.DistributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, info.Amount, currentHeight)
	if err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	//2. 修改delegation信息
	info.IsCompound = tx.IsCompound
	ecoMapper.DelegationMapper.SetDelegationInfo(info)

	return btypes.Result{Code: btypes.CodeOK}, nil
}

func (tx *TxModifyCompound) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxModifyCompound) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxModifyCompound) GetGasPayer() btypes.Address {
	return btypes.Address(tx.Delegator)
}

func (tx *TxModifyCompound) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorOwner...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}

type TxUnbondDelegation struct {
	Delegator      btypes.Address
	ValidatorOwner btypes.Address
	UnbondAmount   uint64
}

var _ txs.ITx = (*TxUnbondDelegation)(nil)

func (tx *TxUnbondDelegation) ValidateData(ctx context.Context) error {

	validator, err := validateValidator(ctx, tx.ValidatorOwner, false, staketypes.Active, true)
	if nil != err {
		return err
	}

	valAddr := btypes.Address(validator.ValidatorPubKey.Address())
	delegationMapper := mapper.GetDelegationMapper(ctx)

	delegationInfo, exsits := delegationMapper.GetDelegationInfo(tx.Delegator, valAddr)
	if !exsits {
		return ErrInvalidInput(DefaultCodeSpace, "delegator not delegate the owner's validator")
	}

	if delegationInfo.Amount < tx.UnbondAmount {
		return ErrInvalidInput(DefaultCodeSpace, "delegator does't have enough amount to unbond")
	}

	if validator.BondTokens < tx.UnbondAmount {
		return ErrInvalidInput(DefaultCodeSpace, "validator does't have enough tokens")
	}

	return nil
}

//unbond delegator tokens
//1. 计算当前delegator收益
//2. 更新delegation信息,若全部unbond 则删除对于的delegation信息
//3. 增加unbond信息
//4. 更新validator bondTokens信息
func (tx *TxUnbondDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {

	currentHeight := uint64(ctx.BlockHeight())
	ecoMapper := mapper.GetEcoMapper(ctx)

	validator, exists := ecoMapper.ValidatorMapper.GetValidatorByOwner(tx.ValidatorOwner)
	if !exists {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: "validator not exsits"}, nil
	}

	valAddr := btypes.Address(validator.ValidatorPubKey.Address())
	delegatorAddr := tx.Delegator

	info, _ := ecoMapper.DelegationMapper.GetDelegationInfo(delegatorAddr, valAddr)

	//1. 计算当前delegator收益
	updatedTokens := info.Amount - tx.UnbondAmount
	err := ecoMapper.DistributionMapper.ModifyDelegatorTokens(validator, delegatorAddr, updatedTokens, currentHeight)
	if err != nil {
		return btypes.Result{Code: btypes.CodeInternal, Codespace: btypes.CodespaceType(err.Error())}, nil
	}

	//2. 更新delegation信息
	info.Amount = info.Amount - tx.UnbondAmount
	ecoMapper.DelegationMapper.SetDelegationInfo(info)
	if info.Amount == uint64(0) { //全部unbond,删除delegation info
		ecoMapper.DelegationMapper.DelDelegationInfo(delegatorAddr, valAddr)
	}

	//3. 增加unbond信息
	stakeParams := ecoMapper.ValidatorMapper.GetParams()
	unbondHeight := uint64(stakeParams.DelegatorUnbondDistributeHeight) + currentHeight
	ecoMapper.DelegationMapper.AddDelegatorUnbondingQOSatHeight(unbondHeight, delegatorAddr, tx.UnbondAmount)

	//4. 更新validator的bondTokens, amount:token = 1:1
	validatorMinusToken := tx.UnbondAmount
	updatedValidatorTokens := validator.BondTokens - validatorMinusToken
	ecoMapper.ValidatorMapper.ChangeValidatorBondTokens(validator, updatedValidatorTokens)

	return btypes.Result{Code: btypes.CodeOK}, nil
}

func (tx *TxUnbondDelegation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxUnbondDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxUnbondDelegation) GetGasPayer() btypes.Address {
	return tx.Delegator
}

func (tx *TxUnbondDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.ValidatorOwner...)
	ret = append(ret, btypes.Int2Byte(int64(tx.UnbondAmount))...)
	return
}

type TxCreateReDelegation struct {
	Delegator          btypes.Address
	FromValidatorOwner btypes.Address
	ToValidatorOwner   btypes.Address
	Amount             uint64
	IsCompound         bool
}

var _ txs.ITx = (*TxCreateReDelegation)(nil)

func (tx *TxCreateReDelegation) ValidateData(ctx context.Context) error {
	//1. 校验fromValidator是否存在
	//2. 校验toValidator是否存在 且 状态为active
	//3. 校验当前用户是否委托了fromValidator
	//4. 校验用户委托的amount是否大于amount

	return nil
}

//delegate from one to another
func (tx *TxCreateReDelegation) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	panic("not implemented")
}

func (tx *TxCreateReDelegation) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Delegator}
}

func (tx *TxCreateReDelegation) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx *TxCreateReDelegation) GetGasPayer() btypes.Address {
	return tx.Delegator
}

func (tx *TxCreateReDelegation) GetSignData() (ret []byte) {
	ret = append(ret, tx.Delegator...)
	ret = append(ret, tx.FromValidatorOwner...)
	ret = append(ret, tx.ToValidatorOwner...)
	ret = append(ret, btypes.Int2Byte(int64(tx.Amount))...)
	ret = append(ret, btypes.Bool2Byte(tx.IsCompound)...)
	return
}
