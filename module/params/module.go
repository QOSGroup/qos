package params

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/baseabci"
	cliContext "github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ types.AppModuleBasic   = AppModuleBasic{}
	_ types.AppModuleGenesis = AppModule{}
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
	return nil
}

// 校验初始状态数据
func (amb AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

func (amb AppModuleBasic) RegisterRestRoutes(ctx cliContext.CLIContext, routes *mux.Router) {
}

// 返回交易命令集合
func (amb AppModuleBasic) GetTxCmds(cdc *amino.Codec) []*cobra.Command {
	return []*cobra.Command{}
}

// 返回查询命令集合
func (amb AppModuleBasic) GetQueryCmds(cdc *amino.Codec) []*cobra.Command {
	return []*cobra.Command{}
}

// 返回数据库操作 Mapper
func (amb AppModuleBasic) GetMapperAndHooks() types.MapperWithHooks {
	return types.NewMapperWithHooks(NewMapper(), nil)
}

// 模块结构
type AppModule struct {
	AppModuleBasic
}

func NewAppModule() types.AppModule {
	return types.NewGenesisOnlyAppModule(AppModule{})
}

// 初始化本模块
func (am AppModule) InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// 导出状态数据
func (am AppModule) ExportGenesis(ctx context.Context) json.RawMessage {
	return nil
}
