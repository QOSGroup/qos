package app

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/distribution"
	"github.com/QOSGroup/qos/module/gov"
	"github.com/QOSGroup/qos/module/guardian"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"io"
)

const (
	appName = "QOS"
)

type QOSApp struct {
	*baseabci.BaseApp
}

func NewApp(logger log.Logger, db dbm.DB, traceStore io.Writer) *QOSApp {

	baseApp := baseabci.NewBaseApp(appName, logger, db, RegisterCodec,
		baseabci.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
	baseApp.SetCommitMultiStoreTracer(traceStore)

	app := &QOSApp{
		BaseApp: baseApp,
	}

	// 设置 InitChainer
	app.SetInitChainer(app.initChainer)

	// 设置gas处理逻辑
	app.SetGasHandler(app.gasHandler)

	// 设置BeginBlocker
	app.SetBeginBlocker(func(ctx context.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		ctx = ctx.WithEventManager(btypes.NewEventManager())
		mint.BeginBlocker(ctx, req)
		distribution.BeginBlocker(ctx, req)
		stake.BeginBlocker(ctx, req)

		return abci.ResponseBeginBlock{
			Events: ctx.EventManager().ABCIEvents(),
		}
	})

	// 设置EndBlocker
	app.SetEndBlocker(func(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		ctx = ctx.WithEventManager(btypes.NewEventManager())
		gov.EndBlocker(ctx)
		distribution.EndBlocker(ctx, req)
		validators := stake.EndBlocker(ctx)
		confirmDataEveryHeight(ctx)
		return abci.ResponseEndBlock{
			ValidatorUpdates: validators,
			Events:           ctx.EventManager().ABCIEvents(),
		}
	})

	// 注册mappers
	app.RegisterMappers()

	// 注册自定义查询处理
	app.RegisterCustomQueryHandler(func(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err btypes.Error) {

		if len(route) == 0 {
			return nil, btypes.ErrInternal("miss custom subquery path")
		}

		if route[0] == stake.ModuleName {
			return stake.Query(ctx, route[1:], req)
		}

		if route[0] == distribution.ModuleName {
			return distribution.Query(ctx, route[1:], req)
		}

		if route[0] == gov.ModuleName {
			return gov.Query(ctx, route[1:], req)
		}

		return nil, nil
	})

	// Mount stores and load the latest state.
	err := app.LoadLatestVersion()
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

// 注册mappers
func (app *QOSApp) RegisterMappers() {
	//parameter mapper
	paramsMapper := params.NewMapper()
	//config params
	paramsMapper.RegisterParamSet(&stake.Params{}, &distribution.Params{}, &gov.Params{})
	app.RegisterMapper(paramsMapper)

	// 账户mapper
	app.RegisterAccountProto(types.ProtoQOSAccount)

	// QCP mapper
	// qbase 默认已注入

	// QSC mapper
	app.RegisterMapper(qsc.NewMapper())

	// 预授权mapper
	app.RegisterMapper(approve.NewMapper)

	// Staking mapper
	stakeMapper := stake.NewMapper()
	stakeMapper.SetHooks(distribution.NewStakingHooks())
	app.RegisterMapper(stakeMapper)

	// Mint mapper
	app.RegisterMapper(mint.NewMapper())

	//distribution mapper
	app.RegisterMapper(distribution.NewMapper())

	//gov mapper
	app.RegisterMapper(gov.NewMapper())

	//guardian mapper
	app.RegisterMapper(guardian.NewMapper())

}

func (app *QOSApp) initChainer(ctx context.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {

	stateJSON := req.AppStateBytes
	genesisState := GenesisState{}
	err := app.GetCdc().UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	if err = ValidGenesis(genesisState); err != nil {
		panic(err)
	}

	initAccounts(ctx, genesisState.Accounts)
	gov.InitGenesis(ctx, genesisState.GovData)
	guardian.InitGenesis(ctx, genesisState.GuardianData)
	mint.InitGenesis(ctx, genesisState.MintData)
	stake.InitGenesis(ctx, genesisState.StakeData)
	qcp.InitGenesis(ctx, genesisState.QCPData)
	qsc.InitGenesis(ctx, genesisState.QSCData)
	approve.InitGenesis(ctx, genesisState.ApproveData)
	distribution.InitGenesis(ctx, genesisState.DistributionData)
	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			bz := app.GetCdc().MustMarshalBinaryBare(genTx)
			reqDeliverTx := abci.RequestDeliverTx{Tx: bz}
			res := app.BaseApp.DeliverTx(reqDeliverTx)
			if !res.IsOK() {
				panic(res.Log)
			}
		}
	}

	res.Validators = stake.GetUpdatedValidators(ctx, uint64(genesisState.StakeData.Params.MaxValidatorCnt))

	return
}

func (app *QOSApp) ExportAppStates(forZeroHeight bool) (appState json.RawMessage, err error) {

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	if forZeroHeight {
		app.prepForZeroHeightGenesis(ctx)
	}

	accounts := []*types.QOSAccount{}
	appendAccount := func(acc account.Account) (stop bool) {
		accounts = append(accounts, acc.(*types.QOSAccount))
		return false
	}
	ctx.Mapper(account.AccountMapperName).(*account.AccountMapper).IterateAccounts(appendAccount)

	genState := NewGenesisState(
		accounts,
		mint.ExportGenesis(ctx),
		stake.ExportGenesis(ctx),
		qcp.ExportGenesis(ctx),
		qsc.ExportGenesis(ctx),
		approve.ExportGenesis(ctx),
		distribution.ExportGenesis(ctx),
		gov.ExportGenesis(ctx),
		guardian.ExportGenesis(ctx),
	)

	stateDataConsistencyCheck(ctx, genState)

	appState, err = app.GetCdc().MarshalJSONIndent(genState, "", " ")
	if err != nil {
		return nil, err
	}

	return appState, nil
}

// prepare for fresh start at zero height
func (app *QOSApp) prepForZeroHeightGenesis(ctx context.Context) {
	/*  reset staking && distribution */
	sm := stake.GetMapper(ctx)
	dm := distribution.GetMapper(ctx)
	am := baseabci.GetAccountMapper(ctx)
	// close all active validators
	var validators []stake.Validator
	var delegations []stake.Delegation
	var vals = make(map[string]stake.Validator)
	sm.IterateValidators(func(validator stake.Validator) {
		val := validator.GetValidatorAddress()
		vals[validator.GetValidatorAddress().String()] = validator
		sm.IterateDelegationsValDeleAddr(val, func(val btypes.Address, del btypes.Address) {
			delegation, exists := sm.GetDelegationInfo(del, val)
			if !exists {
				panic(fmt.Sprintf("delegation from %s to %s not exists", del, val))
			}
			delegations = append(delegations, delegation)
		})
		if validator.Status == stake.Active {
			sm.MakeValidatorInactive(val, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), stake.Revoke)
		}
	})
	// close all inactive validators
	stake.CloseInactiveValidator(ctx, -1)
	for _, delegation := range delegations {
		var info distribution.DelegatorEarningsStartInfo
		dm.Get(distribution.BuildDelegatorEarningStartInfoKey(delegation.ValidatorAddr, delegation.DelegatorAddr), &info)
		delegator := am.GetAccount(delegation.DelegatorAddr).(*types.QOSAccount)
		delegator.PlusQOS(info.HistoricalRewardFees.NilToZero())
		am.SetAccount(delegator)
		dm.MinusValidatorEcoFeePool(delegation.ValidatorAddr, info.HistoricalRewardFees.NilToZero())
	}
	// return unbond tokens
	for h := ctx.BlockHeight(); h <= (int64(sm.GetParams(ctx).DelegatorUnbondReturnHeight) + ctx.BlockHeight()); h++ {
		prePrefix := stake.BuildUnbondingDelegationByHeightPrefix(uint64(h))

		iter := btypes.KVStorePrefixIterator(dm.GetStore(), prePrefix)
		for ; iter.Valid(); iter.Next() {
			k := iter.Key()
			dm.Del(k)
			var unbonds []stake.UnbondingDelegationInfo
			dm.BaseMapper.DecodeObject(iter.Value(), &unbonds)
			for _, unbond := range unbonds {
				_, delAddr := stake.GetUnbondingDelegationHeightAddress(k)
				delegator := am.GetAccount(delAddr).(*types.QOSAccount)
				delegator.PlusQOS(btypes.NewInt(int64(unbond.Amount)))
				am.SetAccount(delegator)
			}
		}
		iter.Close()
	}
	// return redelegation tokens
	for h := ctx.BlockHeight(); h <= (int64(sm.GetParams(ctx).DelegatorRedelegationHeight) + ctx.BlockHeight()); h++ {
		prePrefix := stake.BuildRedelegationByHeightPrefix(uint64(h))

		iter := btypes.KVStorePrefixIterator(dm.GetStore(), prePrefix)
		for ; iter.Valid(); iter.Next() {
			k := iter.Key()
			dm.Del(k)
			var reDelegations []stake.ReDelegationInfo
			dm.BaseMapper.DecodeObject(iter.Value(), &reDelegations)
			for _, reDelegation := range reDelegations {
				_, delAddr := stake.GetRedelegationHeightAddress(k)
				delegator := am.GetAccount(delAddr).(*types.QOSAccount)
				delegator.PlusQOS(btypes.NewInt(int64(reDelegation.Amount)))
				am.SetAccount(delegator)
			}
		}
		iter.Close()
	}
	dm.DeleteDelegatorsIncomeHeight()
	// reinitialize validators
	for _, validator := range validators {
		val := validator.GetValidatorAddress()
		sm.DelValidatorVoteInfo(val)
		sm.ClearValidatorVoteInfoInWindow(val)
		dm.DeleteValidatorPeriodSummaryInfo(val)
		dm.InitValidatorPeriodSummaryInfo(val)
		sm.CreateValidator(validator)
	}
	// reset block height
	ctx = ctx.WithBlockHeight(0)
	// recreate delegations
	for _, delegation := range delegations {
		dm.DelDelegatorEarningStartInfo(delegation.ValidatorAddr, delegation.DelegatorAddr)
		sm.DelDelegationInfo(delegation.DelegatorAddr, delegation.ValidatorAddr)
		validator, _ := sm.GetValidator(delegation.ValidatorAddr)
		sm.ChangeValidatorBondTokens(validator, validator.BondTokens+delegation.Amount)
		sm.Delegate(ctx, stake.NewDelegationInfo(delegation.DelegatorAddr, vals[delegation.ValidatorAddr.String()].GetValidatorAddress(), delegation.Amount, delegation.IsCompound), false)
	}

	/* reset mint */
	mint.GetMapper(ctx).SetFirstBlockTime(0)

	/* reset gov */
	mapper := gov.GetMapper(ctx)
	proposals := mapper.GetProposalsFiltered(ctx, nil, nil, gov.StatusDepositPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.ProposalID
		mapper.RefundDeposits(ctx, proposalID, false)
		mapper.DeleteProposal(proposalID)
	}
	proposals = mapper.GetProposalsFiltered(ctx, nil, nil, gov.StatusVotingPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.ProposalID
		mapper.RefundDeposits(ctx, proposalID, false)
		mapper.DeleteVotes(proposalID)
		mapper.DeleteProposal(proposalID)
	}
}

// gas
func (app *QOSApp) gasHandler(ctx context.Context, payer btypes.Address) (gasUsed uint64, err btypes.Error) {
	gasUsed = ctx.GasMeter().GasConsumed()
	// gas free for txs in the first block
	if ctx.BlockHeight() == 0 {
		return
	}

	// tax free for tx send by guardian
	if _, exists := guardian.GetMapper(ctx).GetGuardian(payer); exists {
		app.Logger.Info("tx send by guardian: %s", payer.String())
		return
	}

	dm := distribution.GetMapper(ctx)
	uint := dm.GetParams(ctx).GasPerUnitCost
	gasFeeUsed := btypes.NewInt(int64(gasUsed / uint))
	gasUsed = gasUsed / uint * uint

	if gasFeeUsed.GT(btypes.ZeroInt()) {
		accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
		account := accountMapper.GetAccount(payer).(*types.QOSAccount)

		if !account.EnoughOfQOS(gasFeeUsed) {
			log := fmt.Sprintf("%s no enough coins to pay the gas after this tx done", payer)
			err = btypes.ErrInternal(log)
			return
		}

		account.MustMinusQOS(gasFeeUsed)
		app.Logger.Info(fmt.Sprintf("cost %d QOS from %s for gas", gasFeeUsed.Int64(), payer))
		accountMapper.SetAccount(account)

		dm.AddPreDistributionQOS(gasFeeUsed)
	}

	return
}

func stateDataConsistencyCheck(ctx context.Context, state GenesisState) bool {

	qosInAccounts := btypes.ZeroInt()
	for _, account := range state.Accounts {
		qosInAccounts = qosInAccounts.Add(account.QOS)
	}
	qosInDelegation := btypes.ZeroInt()
	for _, delegation := range state.StakeData.DelegatorsInfo {
		qosInDelegation = qosInDelegation.Add(btypes.NewInt(int64(delegation.Amount)))
	}
	preDistributionRemainTotal := btypes.ZeroInt()
	for _, data := range state.DistributionData.ValidatorEcoFeePools {
		preDistributionRemainTotal = preDistributionRemainTotal.Add(data.EcoFeePool.PreDistributeRemainTotalFee)
	}
	qosUnbond := btypes.ZeroInt()
	for _, unbond := range state.StakeData.DelegatorsUnbondInfo {
		qosUnbond = qosUnbond.Add(btypes.NewInt(int64(unbond.Amount)))
	}
	redelegations := btypes.ZeroInt()
	for _, reDelegation := range state.StakeData.ReDelegationsInfo {
		redelegations = redelegations.Add(btypes.NewInt(int64(reDelegation.Amount)))
	}
	govDeposit := btypes.ZeroInt()
	for _, proposal := range state.GovData.Proposals {
		if proposal.Proposal.Status != gov.StatusPassed && proposal.Proposal.Status != gov.StatusRejected {
			govDeposit = govDeposit.Add(btypes.NewInt(int64(proposal.Proposal.TotalDeposit)))
		}
	}

	qosFeePool := state.DistributionData.CommunityFeePool
	qosPreQOS := state.DistributionData.PreDistributionQOSAmount

	qosTotal := qosInAccounts.Add(qosInDelegation).Add(qosUnbond).Add(redelegations).Add(qosFeePool).Add(qosPreQOS).Add(preDistributionRemainTotal).Add(govDeposit)
	qosApplied := state.MintData.AppliedQOSAmount
	diff := qosTotal.Sub(btypes.NewInt(int64(qosApplied)))

	ctx.Logger().Info("DATA CONFIRM",
		"height", ctx.BlockHeight(),
		"accounts", qosInAccounts,
		"delegations", qosInDelegation,
		"unbond", qosUnbond,
		"redelegation", redelegations,
		"feepool", qosFeePool,
		"pre", qosPreQOS,
		"valshared", preDistributionRemainTotal,
		"total", qosTotal,
		"applied", qosApplied,
		"diff", diff)

	return diff.Equal(btypes.ZeroInt())
}

func confirmDataEveryHeight(ctx context.Context) {
	accounts := []*types.QOSAccount{}
	ctx.Mapper(account.AccountMapperName).(*account.AccountMapper).IterateAccounts(func(acc account.Account) (stop bool) {
		accounts = append(accounts, acc.(*types.QOSAccount))
		return false
	})
	genState := NewGenesisState(
		accounts,
		mint.ExportGenesis(ctx),
		stake.ExportGenesis(ctx),
		qcp.ExportGenesis(ctx),
		qsc.ExportGenesis(ctx),
		approve.ExportGenesis(ctx),
		distribution.ExportGenesis(ctx),
		gov.ExportGenesis(ctx),
		guardian.ExportGenesis(ctx),
	)

	isSame := stateDataConsistencyCheck(ctx, genState)
	if !isSame {
		panic("DATA NOT CONSISTENCY")
	}
}
