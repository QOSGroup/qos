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
	"github.com/QOSGroup/qos/module/bank"
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
	"time"
)

const (
	appName = "QOS"
)

var (
	ModuleBasics = types.NewBasicManager(
		approve.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.AppModuleBasic{},
		guardian.AppModuleBasic{},
		mint.AppModuleBasic{},
		params.AppModuleBasic{},
		qcp.AppModuleBasic{},
		qsc.AppModuleBasic{},
		stake.AppModuleBasic{},
		bank.AppModuleBasic{},
	)
)

type QOSApp struct {
	*baseabci.BaseApp

	// module manager
	mm *types.Manager

	// invariants
	invarRoutes    []types.InvarRoute
	invCheckPeriod uint

	// query router
	queryRoutes map[string]types.Querier
}

func NewApp(logger log.Logger, db dbm.DB, traceStore io.Writer, invCheckPeriod uint) *QOSApp {

	baseApp := baseabci.NewBaseApp(appName, logger, db, RegisterCodec,
		baseabci.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
	baseApp.SetCommitMultiStoreTracer(traceStore)

	app := &QOSApp{
		BaseApp: baseApp,
		mm: types.NewManager(
			bank.NewAppModule(),
			approve.NewAppModule(),
			distribution.NewAppModule(),
			gov.NewAppModule(),
			guardian.NewAppModule(),
			mint.NewAppModule(),
			params.NewAppModule(),
			qcp.NewAppModule(),
			qsc.NewAppModule(),
			stake.NewAppModule(),
		),
		invCheckPeriod: invCheckPeriod,
		queryRoutes:    make(map[string]types.Querier),
	}

	// 注册invariants
	app.mm.RegisterInvariants(app)

	// 注册mappers
	app.RegisterMappers()

	// 设置gas处理逻辑
	app.SetGasHandler(app.GasHandler)

	// 设置BeginBlocker
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distribution.ModuleName, stake.ModuleName)
	app.SetBeginBlocker(app.BeginBlocker)

	// 设置EndBlocker
	app.mm.SetOrderEndBlockers(gov.ModuleName, distribution.ModuleName, stake.ModuleName)
	app.SetEndBlocker(app.EndBlocker)

	// 设置 InitChainer
	// !!! accounts first, stake last
	app.mm.SetOrderInitGenesis(bank.ModuleName, gov.ModuleName, guardian.ModuleName, mint.ModuleName, qcp.ModuleName, qsc.ModuleName, approve.ModuleName, distribution.ModuleName, stake.ModuleName)
	app.SetInitChainer(app.InitChainer)

	// 注册自定义查询处理
	app.mm.RegisterQueriers(app)
	app.RegisterCustomQueryHandler(func(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err btypes.Error) {
		if len(route) == 0 {
			return nil, btypes.ErrInternal("miss custom subquery path")
		}
		return app.queryRoutes[route[0]](ctx, route[1:], req)
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
	app.RegisterMapper(approve.NewMapper())

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

func (app *QOSApp) InitChainer(ctx context.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {

	stateJSON := req.AppStateBytes
	genesisState := types.GenesisState{}
	err := app.GetCdc().UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	res = app.mm.InitGenesis(ctx, app.BaseApp, genesisState)

	return
}

func (app *QOSApp) BeginBlocker(ctx context.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *QOSApp) EndBlocker(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	res := app.mm.EndBlock(ctx, req)

	if app.invCheckPeriod == 0 || ctx.BlockHeight()%int64(app.invCheckPeriod) == 0 {
		app.AssertInvariants(ctx)
	}
	return res
}

func (app *QOSApp) ExportAppStates(forZeroHeight bool) (appState json.RawMessage, err error) {

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	if forZeroHeight {
		app.prepForZeroHeightGenesis(ctx)
	}

	genState := app.mm.ExportGenesis(ctx)

	//TODO imuge 数据校验

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
			validators = append(validators, validator)
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
				_, delAddr, _ := stake.GetUnbondingDelegationHeightAddress(k)
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
				_, delAddr, _ := stake.GetRedelegationHeightAddress(k)
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
func (app *QOSApp) GasHandler(ctx context.Context, payer btypes.Address) (gasUsed uint64, err btypes.Error) {
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

func (app *QOSApp) RegisterInvarRoute(module string, route string, invar types.Invariant) {
	invarRoute := types.NewInvarRoute(module, route, invar)
	app.invarRoutes = append(app.invarRoutes, invarRoute)
}

func (app *QOSApp) AssertInvariants(ctx context.Context) {
	logger := app.Logger

	start := time.Now()

	totalCoins := btypes.BaseCoins{}
	for _, invarRoute := range app.invarRoutes {
		msg, coins, stop := invarRoute.Invar(ctx)
		if stop {
			panic(msg)
		}
		totalCoins = totalCoins.Plus(coins)
		logger.Info(fmt.Sprintf("invariant check %s\t%s:\t%s", invarRoute.ModuleName, invarRoute.Route, coins.String()))
	}

	if !totalCoins.IsZero() {
		panic("invariant check not pass")
	}

	end := time.Now()
	diff := end.Sub(start)

	logger.Info("asserted all invariants", "duration", diff, "height", ctx.BlockHeight())
}

func (app *QOSApp) RegisterQueryRoute(module string, query types.Querier) {
	app.queryRoutes[module] = query
}
