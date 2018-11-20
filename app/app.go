package app

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/types"
	qosacc "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/txs/qsc"
	"github.com/QOSGroup/qos/txs/validator"
	"github.com/QOSGroup/qos/x/miner"
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
		miner.BeginBlocker(ctx, req)
		return abci.ResponseBeginBlock{}
	})

	//设置endblocker
	app.SetEndBlocker(func(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return validator.EndBlocker(ctx)
	})

	// 账户mapper
	app.RegisterAccountProto(qosacc.ProtoQOSAccount)

	// 基础信息操作mapper
	app.RegisterMapper(mapper.NewMainMapper())

	// QCP mapper
	// qbase 默认已注入

	// QSC mapper
	app.RegisterMapper(qsc.NewQSCMapper())

	// 预授权mapper
	app.RegisterMapper(approve.NewApproveMapper())

	// Validator mapper
	app.RegisterMapper(validator.NewValidatorMapper())

	// Mount stores and load the latest state.
	err := app.LoadLatestVersion()
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

// 初始配置
func (app *QOSApp) initChainer(ctx context.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
	// 上下文中获取mapper
	mainMapper := ctx.Mapper(mapper.GetMainStoreKey()).(*mapper.MainMapper)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)

	// 反序列化app_state
	stateJSON := req.AppStateBytes
	genesisState := &GenesisState{}
	err := accountMapper.GetCodec().UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	// 保存CA
	mainMapper.SetRootCA(genesisState.CAPubKey)

	// 保存初始账户
	for _, acc := range genesisState.Accounts {
		accountMapper.SetAccount(acc)
	}
	
	// 保存Validators以及对应账户信息: validators信息从genesisState.Validators中获取
	if len(genesisState.Validators) > 0 {
		validatorMapper := ctx.Mapper(validator.ValidatorMapperName).(*validator.ValidatorMapper)
		for _, v := range genesisState.Validators {
			validatorMapper.SaveValidator(v)

			addr := types.Address(v.Operator)
			acc := accountMapper.GetAccount(addr)
			if acc.GetPubicKey() == nil {
				acc = accountMapper.NewAccountWithAddress(addr)
				accountMapper.SetAccount(acc)
			}

			res.Validators = append(res.Validators, v.ToABCIValidator())
		}
		validatorMapper.SetValidatorUnChanged()
	}

	return
}
