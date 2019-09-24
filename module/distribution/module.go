package distribution

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/distribution/mapper"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ types.AppModuleBasic = AppModuleBasic{}
	_ types.AppModule      = AppModule{}
)

// 基础模块结构
type AppModuleBasic struct{}

// 模块名
func (amb AppModuleBasic) Name() string {
	return ModuleName
}

// amino 相关类/对象注册
func (amb AppModuleBasic) RegisterCodec(cdc *amino.Codec) {
	RegisterCodec(cdc)
}

// 默认初始状态数据
func (amb AppModuleBasic) DefaultGenesis() json.RawMessage {
	return Cdc.MustMarshalJSON(DefaultGenesis())
}

// 校验初始状态数据
func (amb AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := Cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// 返回交易命令集合
func (amb AppModuleBasic) GetTxCmds(cdc *amino.Codec) []*cobra.Command {
	return TxCommands(cdc)
}

// 返回查询命令集合
func (amb AppModuleBasic) GetQueryCmds(cdc *amino.Codec) []*cobra.Command {
	return QueryCommands(cdc)
}

// 返回数据库操作 Mapper
func (amb AppModuleBasic) GetMapperAndHooks() types.MapperWithHooks {
	return types.NewMapperWithHooks(NewMapper(), NewStakingHooks())
}

// 模块结构
type AppModule struct {
	AppModuleBasic
}

func NewAppModule() types.AppModule {
	return AppModule{
		AppModuleBasic{},
	}
}

// 初始化本模块
func (am AppModule) InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	Cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, genesisState)
	return []abci.ValidatorUpdate{}
}

// 导出状态数据
func (am AppModule) ExportGenesis(ctx context.Context) json.RawMessage {
	gs := ExportGenesis(ctx)
	return Cdc.MustMarshalJSON(gs)
}

// 注册数据验证
func (am AppModule) RegisterInvariants(ir types.InvariantRegistry) {
	ir.RegisterInvarRoute(ModuleName, "fee-pool", mapper.FeePoolInvariant(ModuleName))
	ir.RegisterInvarRoute(ModuleName, "pre-distribution", mapper.PreDistributionInvariant(ModuleName))
	ir.RegisterInvarRoute(ModuleName, "validator-fee-pool", mapper.ValidatorFeePoolInvariant(ModuleName))
}

// App BeginBlocker 中执行操作
func (am AppModule) BeginBlock(ctx context.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req)
}

// App EndBlocker 中执行操作
func (am AppModule) EndBlock(ctx context.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, req)
	return []abci.ValidatorUpdate{}
}

// 返回本模块自定义应用查询路由信息
func (am AppModule) RegisterQuerier(qr types.QueryRegistry) {
	qr.RegisterQueryRoute(ModuleName, mapper.Query)
}
