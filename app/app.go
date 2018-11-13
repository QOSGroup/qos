package app

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	qosacc "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/test"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/txs/validator"
	"github.com/QOSGroup/qos/types"
	"github.com/QOSGroup/qos/x/miner"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"strconv"
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

	// 账户mapper
	app.RegisterAccountProto(qosacc.ProtoQOSAccount)

	// 基础信息操作mapper
	app.RegisterMapper(mapper.NewMainMapper())

	// QCP mapper
	// qbase 默认已注入

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
func (app *QOSApp) initChainer(ctx context.Context, req abci.RequestInitChain) abci.ResponseInitChain {
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

	// 设置初始账户(test only)
	// todo: remove later
	accret := test.InitKeys(app.GetCdc())
	for _, ac := range accret {
		accountMapper.SetAccount(&ac.Acc)
	}

	// 保存Validators以及对应账户信息
	if len(req.Validators) > 0 {
		validatorMapper := ctx.Mapper(validator.ValidatorMapperName).(*validator.ValidatorMapper)
		for i, v := range req.Validators {
			var pubKey ed25519.PubKeyEd25519
			copy(pubKey[:], v.PubKey.Data[:ed25519.PubKeyEd25519Size])
			validator := types.NewValidator("v-"+strconv.Itoa(i), pubKey, v.Power, 1)
			validatorMapper.SaveValidator(validator)
			accountMapper.SetAccount(accountMapper.NewAccountWithAddress(pubKey.Address().Bytes()))
		}
		validatorMapper.SetUpdated(false)
	}

	return abci.ResponseInitChain{}
}
