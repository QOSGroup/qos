package app

import (
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	qosacc "github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/txs/approve"
	"github.com/QOSGroup/qos/txs/qsc"
	"github.com/QOSGroup/qos/txs/staking"
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
		staking.BeginBlocker(ctx, req)

		return abci.ResponseBeginBlock{}
	})

	//设置endblocker
	app.SetEndBlocker(func(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		return staking.EndBlocker(ctx)
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

	// Staking Validator mapper
	app.RegisterMapper(staking.NewValidatorMapper())

	// Staking mapper
	app.RegisterMapper(staking.NewVoteInfoMapper())

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
	if genesisState.CAPubKey != nil {
		mainMapper.SetRootCA(genesisState.CAPubKey)
	}

	// 保存SPOConfig
	mainMapper.SetSPOConfig(genesisState.SPOConfig)

	// 保存StakeConfig
	mainMapper.SetStakeConfig(genesisState.StakeConfig)

	var appliedQOSAmount uint64

	// 保存初始账户
	for _, acc := range genesisState.Accounts {
		accountMapper.SetAccount(acc)
		appliedQOSAmount += uint64(acc.QOS.Int64())
	}

	// 保存 QOS amount
	mainMapper.SetAppliedQOSAmount(appliedQOSAmount)

	// 保存Validators以及对应账户信息: validators信息从genesisState.Validators中获取
	if len(genesisState.Validators) > 0 {
		validatorMapper := ctx.Mapper(staking.ValidatorMapperName).(*staking.ValidatorMapper)
		for _, v := range genesisState.Validators {

			if validatorMapper.Exists(v.ValidatorPubKey.Address().Bytes()) {
				panic(fmt.Errorf("validator %s already exists", v.ValidatorPubKey.Address()))
			}
			if validatorMapper.ExistsWithOwner(v.Owner) {
				panic(fmt.Errorf("owner %s already bind a validator", v.Owner))
			}

			validatorMapper.CreateValidator(v)

			acc := accountMapper.GetAccount(v.Owner)
			if acc == nil {
				panic(fmt.Errorf("owner of %s not exists", v.Name))
			}
			owner := acc.(*qosacc.QOSAccount)
			tokens := btypes.NewInt(int64(v.BondTokens))
			if !owner.EnoughOfQOS(tokens) {
				panic(fmt.Errorf("%s no enough QOS", v.Name))
			}
			owner.MustMinusQOS(tokens)
			accountMapper.SetAccount(acc)

			// res.Validators = append(res.Validators, v.ToABCIValidator())
		}
	}

	res.Validators = staking.GetUpdatedValidators(ctx, uint64(genesisState.StakeConfig.MaxValidatorCnt))
	return
}
