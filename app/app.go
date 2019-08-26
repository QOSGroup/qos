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

	// 注册mappers and hooks, 初始化参数配置
	app.mm.RegisterMapperAndHooks(app, params.ModuleName, &stake.Params{}, &distribution.Params{}, &gov.Params{})

	// 设置gas处理逻辑
	app.SetGasHandler(app.GasHandler)

	// 设置BeginBlocker
	app.mm.SetOrderBeginBlockers(guardian.ModuleName, gov.ModuleName, mint.ModuleName, distribution.ModuleName, stake.ModuleName)
	app.SetBeginBlocker(app.BeginBlocker)

	// 设置EndBlocker
	app.mm.SetOrderEndBlockers(gov.ModuleName, distribution.ModuleName, stake.ModuleName, bank.ModuleName)
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
		if querier, ok := app.queryRoutes[route[0]]; ok {
			return querier(ctx, route[1:], req)
		} else {
			return nil, nil
		}
	})

	// Mount stores and load the latest state.
	err := app.LoadLatestVersion()
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
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

	// 收到检查事件或固定时间间隔进行数据检查
	check := false
	for _, event := range res.Events {
		if event.Type == types.EventTypeInvariantCheck {
			check = true
			break
		}
	}
	if check || app.invCheckPeriod == 0 || ctx.BlockHeight()%int64(app.invCheckPeriod) == 0 {
		app.AssertInvariants(ctx)
	}

	return res
}

func (app *QOSApp) ExportAppStates(forZeroHeight bool) (appState json.RawMessage, err error) {

	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	ctx = ctx.WithEventManager(btypes.NewEventManager())

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
	am := bank.GetMapper(ctx)

	// close all active validators
	var delegations []stake.Delegation
	var validators = make(map[string]stake.Validator)
	iterator := sm.IteratorValidatorByVoterPower(false)
	defer iterator.Close()
	var key []byte
	for ; iterator.Valid(); iterator.Next() {
		key = iterator.Key()
		valAddr := btypes.Address(key[9:])
		if validator, exists := sm.GetValidator(valAddr); exists {
			validators[valAddr.String()] = validator
			delegations = append(delegations, sm.GetDelegationsByValidator(valAddr)...)
			sm.MakeValidatorInactive(valAddr, uint64(ctx.BlockHeight()), ctx.BlockHeader().Time.UTC(), stake.Revoke)
		}
	}

	// close all inactive validators
	stake.CloseInactiveValidator(ctx, -1)

	// return unbond tokens
	for h := ctx.BlockHeight(); h <= (int64(sm.GetParams(ctx).DelegatorUnbondReturnHeight) + ctx.BlockHeight()); h++ {
		prePrefix := stake.BuildUnbondingDelegationByHeightPrefix(uint64(h))

		iter := btypes.KVStorePrefixIterator(sm.GetStore(), prePrefix)
		for ; iter.Valid(); iter.Next() {
			k := iter.Key()
			sm.Del(k)
			var unbond stake.UnbondingDelegationInfo
			sm.BaseMapper.DecodeObject(iter.Value(), &unbond)
			_, delAddr, _ := stake.GetUnbondingDelegationHeightAddress(k)
			delegator := am.GetAccount(delAddr).(*types.QOSAccount)
			delegator.PlusQOS(btypes.NewInt(int64(unbond.Amount)))
			am.SetAccount(delegator)
		}
		iter.Close()
	}

	// return redelegation tokens
	for h := ctx.BlockHeight(); h <= (int64(sm.GetParams(ctx).DelegatorRedelegationHeight) + ctx.BlockHeight()); h++ {
		prePrefix := stake.BuildRedelegationByHeightPrefix(uint64(h))

		iter := btypes.KVStorePrefixIterator(sm.GetStore(), prePrefix)
		for ; iter.Valid(); iter.Next() {
			k := iter.Key()
			sm.Del(k)
			var reDelegation stake.ReDelegationInfo
			sm.BaseMapper.DecodeObject(iter.Value(), &reDelegation)
			_, delAddr, _ := stake.GetRedelegationHeightAddress(k)
			delegator := am.GetAccount(delAddr).(*types.QOSAccount)
			delegator.PlusQOS(btypes.NewInt(int64(reDelegation.Amount)))
			am.SetAccount(delegator)
		}
		iter.Close()
	}

	// reinitialize validators
	for _, validator := range validators {
		val := validator.GetValidatorAddress()
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
		sm.Delegate(ctx, stake.NewDelegationInfo(delegation.DelegatorAddr, validators[delegation.ValidatorAddr.String()].GetValidatorAddress(), delegation.Amount, delegation.IsCompound), false)
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

	logger.Info(fmt.Sprintf("invariant check invariant:\t%s", totalCoins.String()))
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

func (app *QOSApp) RegisterHooksMapper(mhs map[string]types.MapperWithHooks) {
	for _, mh := range mhs {
		// register mapper hooks
		if mh.Hooks != nil {
			mhs[mh.Hooks.HookMapper()].Mapper.(types.HooksMapper).SetHooks(mh.Hooks)
		}
		// register mapper
		app.BaseApp.RegisterMapper(mh.Mapper)
	}

}
