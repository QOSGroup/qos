package app

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/distribution"
	ecomapper "github.com/QOSGroup/qos/module/eco/mapper"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/qcp"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	"github.com/QOSGroup/qos/types"
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

	baseApp := baseabci.NewBaseApp(appName, logger, db, RegisterCodec)
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
		distribution.BeginBlocker(ctx, req)
		stake.BeginBlocker(ctx, req)
		mint.BeginBlocker(ctx, req)
		return abci.ResponseBeginBlock{}
	})

	//设置endblocker
	app.SetEndBlocker(func(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		distribution.EndBlocker(ctx, req)
		stake.EndBlockerByReturnUnbondTokens(ctx)
		return stake.EndBlocker(ctx)
	})

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

	res.Validators = InitGenesis(ctx, genesisState)

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
		stake.ExportGenesis(ctx, forZeroHeight),
		qcp.ExportGenesis(ctx),
		qsc.ExportGenesis(ctx),
		approve.ExportGenesis(ctx),
		distribution.ExportGenesis(ctx, forZeroHeight),
	)
	appState, err = app.GetCdc().MarshalJSONIndent(genState, "", " ")
	if err != nil {
		return nil, err
	}

	return appState, nil
}

// prepare for fresh start at zero height
func (app *QOSApp) prepForZeroHeightGenesis(ctx context.Context) {

	// close inactive validators
	stake.CloseExpireInactiveValidator(ctx, 0)

	// return unbond tokens
	stake.ReturnAllUnbondTokens(ctx)

	ecomapper.GetMintMapper(ctx).SetFirstBlockTime(0)
}

// gas
func (app *QOSApp) gasHandler(ctx context.Context, payer btypes.Address) btypes.Error {
	distributionMapper := ecomapper.GetDistributionMapper(ctx)
	gasFeeUsed := btypes.NewInt(int64(ctx.GasMeter().GasConsumed() / distributionMapper.GetParams().GasPerUnitCost))

	if gasFeeUsed.GT(btypes.ZeroInt()) {
		accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)
		account := accountMapper.GetAccount(payer).(*types.QOSAccount)

		if !account.EnoughOfQOS(gasFeeUsed) {
			log := fmt.Sprintf("%s no enough coins to pay the gas after this tx done", payer)
			return btypes.ErrInternal(log)
		}

		account.MustMinusQOS(gasFeeUsed)
		app.Logger.Info(fmt.Sprintf("cost %d QOS from %s for gas", gasFeeUsed.Int64(), payer))
		accountMapper.SetAccount(account)

		distributionMapper.AddPreDistributionQOS(gasFeeUsed)
	}

	return nil
}
