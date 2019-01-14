package app

import (
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/approve"
	"github.com/QOSGroup/qos/module/mint"
	mintmapper "github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/module/qsc"
	"github.com/QOSGroup/qos/module/stake"
	stakemapper "github.com/QOSGroup/qos/module/stake/mapper"
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

	app.SetBeginBlocker(func(ctx context.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
		mint.BeginBlocker(ctx, req)
		stake.BeginBlocker(ctx, req)

		return abci.ResponseBeginBlock{}
	})

	//设置endblocker
	app.SetEndBlocker(func(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
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
	app.RegisterMapper(stakemapper.NewValidatorMapper())

	// Staking mapper
	app.RegisterMapper(stakemapper.NewVoteInfoMapper())

	// Mint mapper
	app.RegisterMapper(mintmapper.NewMintMapper())

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
