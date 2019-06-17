package app

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/distribution"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
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

	// abci:
	// begin blocker:
	// 1. validator奖励分配(distribution)
	// 2. 不活跃validator置为inactive(stake)
	// 3. 计算本块挖出的QOS数量(mint)
	// end blocker:
	// 1. delegator收益发放: 计算下一发放周期(distribution)
	// 2. unbond QOS 返还 (stake)
	// 3. validator period 旧数据删除(distribution) //TODO
	// 4. close inactive  validator(stake),统计新的validator (stake)

	app.SetBeginBlocker(func(ctx context.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		mint.BeginBlocker(ctx, req)
		distribution.BeginBlocker(ctx, req)
		stake.BeginBlocker(ctx, req)

		return abci.ResponseBeginBlock{}
	})

	//设置endblocker
	app.SetEndBlocker(func(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		gov.EndBlocker(ctx)
		distribution.EndBlocker(ctx, req)
		stake.EndBlockerByReturnUnbondTokens(ctx)
		validators := stake.EndBlocker(ctx)
		confirmDataEveryHeight(ctx)
		return validators
	})

	//parameter mapper
	paramsMapper := params.NewMapper()
	//config params
	paramsMapper.RegisterParamSet(&ecotypes.StakeParams{}, &ecotypes.DistributionParams{}, &gov.Params{})
	app.RegisterMapper(paramsMapper)

	// 账户mapper
	app.RegisterAccountProto(types.ProtoQOSAccount)

	// QCP mapper
	// qbase 默认已注入

	// QSC mapper
	app.RegisterMapper(qsc.NewQSCMapper())

	// 预授权mapper
	app.RegisterMapper(approve.NewApproveMapper())

	// Staking Validator mapper
	app.RegisterMapper(ecomapper.NewValidatorMapper())

	// Staking mapper
	app.RegisterMapper(ecomapper.NewVoteInfoMapper())

	// Mint mapper
	app.RegisterMapper(ecomapper.NewMintMapper())

	//distributionMapper
	app.RegisterMapper(ecomapper.NewDistributionMapper())

	//delegationMapper
	app.RegisterMapper(ecomapper.NewDelegationMapper())

	//gov mapper
	app.RegisterMapper(gov.NewGovMapper())

	//guardian mapper
	app.RegisterMapper(guardian.NewGuardianMapper())

	app.RegisterCustomQueryHandler(func(ctx context.Context, route []string, req abci.RequestQuery) (res []byte, err btypes.Error) {

		if len(route) == 0 {
			return nil, btypes.ErrInternal("miss custom subquery path")
		}

		if route[0] == ecotypes.Stake {
			return stake.Query(ctx, route[1:], req)
		}

		if route[0] == ecotypes.Distribution {
			return distribution.Query(ctx, route[1:], req)
		}

		if route[0] == gov.GOV {
			return gov.Query(ctx, route[1:], req)
		}

		if route[0] == mint.Mint {
			return mint.Query(ctx, route[1:], req)
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
			res := app.BaseApp.DeliverTx(bz)
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

	stake.PrepForZeroHeightGenesis(ctx)

	mint.PrepForZeroHeightGenesis(ctx)

	gov.PrepForZeroHeightGenesis(ctx)
}

// gas
func (app *QOSApp) gasHandler(ctx context.Context, payer btypes.Address) (gasUsed uint64, err btypes.Error) {
	gasUsed = ctx.GasMeter().GasConsumed()
	// gas free for txs in the first block
	if ctx.BlockHeight() == 0 {
		return
	}

	// tax free for tx send by guardian
	if _, exists := guardian.GetGuardianMapper(ctx).GetGuardian(payer); exists {
		app.Logger.Info("tx send by guardian: %s", payer.String())
		return
	}

	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	uint := distributionMapper.GetParams(ctx).GasPerUnitCost
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

		distributionMapper.AddPreDistributionQOS(gasFeeUsed)
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
	qosFeePool := state.DistributionData.CommunityFeePool
	qosPreQOS := state.DistributionData.PreDistributionQOSAmount

	qosTotal := qosInAccounts.Add(qosInDelegation).Add(qosUnbond).Add(qosFeePool).Add(qosPreQOS).Add(preDistributionRemainTotal)
	qosApplied := state.MintData.AppliedQOSAmount
	diff := qosTotal.Sub(btypes.NewInt(int64(qosApplied)))

	ctx.Logger().Info("DATA CONFIRM",
		"height", ctx.BlockHeight(),
		"accounts", qosInAccounts,
		"delegations", qosInDelegation,
		"unbond", qosUnbond,
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
