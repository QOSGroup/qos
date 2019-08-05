package params

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ types.AppModuleBasic   = AppModuleBasic{}
	_ types.AppModuleGenesis = AppModule{}
)

// app module basics object
type AppModuleBasic struct{}

func (amb AppModuleBasic) Name() string {
	return ModuleName
}

func (amb AppModuleBasic) RegisterCodec(cdc *amino.Codec) {
	RegisterCodec(cdc)
}

func (amb AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

func (amb AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

func (amb AppModuleBasic) GetTxCmds(cdc *amino.Codec) []*cobra.Command {
	return []*cobra.Command{}
}

func (amb AppModuleBasic) GetQueryCmds(cdc *amino.Codec) []*cobra.Command {
	return []*cobra.Command{}
}

func (amb AppModuleBasic) GetMapperAndHooks() types.MapperWithHooks {
	return types.NewMapperWithHooks(NewMapper(), nil)
}

// app module
type AppModule struct {
	AppModuleBasic
}

func NewAppModule() types.AppModule {
	return types.NewGenesisOnlyAppModule(AppModule{})
}

func (am AppModule) InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx context.Context) json.RawMessage {
	return nil
}
