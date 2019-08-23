package stake

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/baseabci"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ types.AppModuleBasic = AppModuleBasic{}
	_ types.AppModule      = AppModule{}
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
	return Cdc.MustMarshalJSON(DefaultGenesis())
}

func (amb AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := Cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (amb AppModuleBasic) GetTxCmds(cdc *amino.Codec) []*cobra.Command {
	return TxCommands(cdc)
}

func (amb AppModuleBasic) GetQueryCmds(cdc *amino.Codec) []*cobra.Command {
	return QueryCommands(cdc)
}

func (amb AppModuleBasic) GetMapperAndHooks() types.MapperWithHooks {
	return types.NewMapperWithHooks(NewMapper(), nil)
}

// app module
type AppModule struct {
	AppModuleBasic
}

func NewAppModule() types.AppModule {
	return AppModule{
		AppModuleBasic{},
	}
}

func (am AppModule) InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	Cdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, bapp, genesisState)
}

func (am AppModule) ExportGenesis(ctx context.Context) json.RawMessage {
	gs := ExportGenesis(ctx)
	return Cdc.MustMarshalJSON(gs)
}

func (am AppModule) RegisterInvariants(ir types.InvariantRegistry) {
	ir.RegisterInvarRoute(ModuleName, "delegation", mapper.DelegationInvariant(ModuleName))
	ir.RegisterInvarRoute(ModuleName, "unbonding", mapper.UnbondingInvariant(ModuleName))
	ir.RegisterInvarRoute(ModuleName, "redelegation", mapper.RedelegationInvariant(ModuleName))
}

func (am AppModule) BeginBlock(ctx context.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req)
}

func (am AppModule) EndBlock(ctx context.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx)
}

func (am AppModule) RegisterQuerier(qr types.QueryRegistry) {
	qr.RegisterQueryRoute(ModuleName, mapper.Query)
}
