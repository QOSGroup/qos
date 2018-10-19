package app

import (
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	qosacc "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs"
	"github.com/tendermint/go-amino"
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

	baseApp := baseabci.NewBaseApp(appName, logger, db, registerCdc)
	baseApp.SetCommitMultiStoreTracer(traceStore)

	app := &QOSApp{
		BaseApp: baseApp,
	}

	// 设置 InitChainer
	app.SetInitChainer(app.initChainer)

	// 账户mapper
	app.RegisterAccountProto(qosacc.ProtoQOSAccount)

	// 基础信息操作mapper
	app.RegisterMapper(mapper.NewBaseMapper())

	// QCP mapper
	// qbase 默认已注入

	// 预授权mapper
	app.RegisterMapper(mapper.NewApproveMapper())

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
	baseMapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.BaseMapper)
	accountMapper := ctx.Mapper(account.AccountMapperName).(*account.AccountMapper)

	// 反序列化app_state
	stateJSON := req.AppStateBytes
	genesisState := &qosacc.GenesisState{}
	err := accountMapper.GetCodec().UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	// 保存CA
	baseMapper.SetCA(genesisState.CAPubKey)

	// 保存初始账户
	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			panic(err)
		}
		accountMapper.SetAccount(acc)
	}

	return abci.ResponseInitChain{}
}

// 序列化反序列化相关注册
func MakeCodec() *amino.Codec {
	cdc := baseabci.MakeQBaseCodec()
	registerCdc(cdc)
	return cdc
}

func registerCdc(cdc *amino.Codec) {
	cdc.RegisterConcrete(&qosacc.QOSAccount{}, "qos/account/QOSAccount", nil)
	cdc.RegisterConcrete(&txs.TxCreateQSC{}, "qos/txs/TxCreateQSC", nil)
	cdc.RegisterConcrete(&txs.TxIssueQsc{}, "qos/txs/TxIssueQsc", nil)
	cdc.RegisterConcrete(&txs.TxTransform{}, "qos/txs/TxTransform", nil)
	cdc.RegisterConcrete(&txs.TxApproveCreate{}, "qos/txs/TxApproveCreate", nil)
	cdc.RegisterConcrete(&txs.TxApproveIncrease{}, "qos/txs/TxApproveIncrease", nil)
	cdc.RegisterConcrete(&txs.TxApproveDecrease{}, "qos/txs/TxApproveDecrease", nil)
	cdc.RegisterConcrete(&txs.TxApproveUse{}, "qos/txs/TxApproveUse", nil)
	cdc.RegisterConcrete(&txs.TxApproveCancel{}, "qos/txs/TxApproveCancel", nil)
}
