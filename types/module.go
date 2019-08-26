package types

import (
	"encoding/json"
	"github.com/QOSGroup/qbase/baseabci"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModuleBasic interface {
	Name() string
	RegisterCodec(*amino.Codec)

	// genesis
	DefaultGenesis() json.RawMessage
	ValidateGenesis(json.RawMessage) error

	// mapper and hooks
	GetMapperAndHooks() MapperWithHooks

	// client functionality
	GetTxCmds(*amino.Codec) []*cobra.Command
	GetQueryCmds(*amino.Codec) []*cobra.Command
}

// collections of AppModuleBasic
type BasicManager map[string]AppModuleBasic

func NewBasicManager(modules ...AppModuleBasic) BasicManager {
	moduleMap := make(map[string]AppModuleBasic)
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}
	return moduleMap
}

// RegisterCodecs registers all module codecs
func (bm BasicManager) RegisterCodec(cdc *amino.Codec) {
	for _, b := range bm {
		b.RegisterCodec(cdc)
	}
}

// Provided default genesis information for all modules
func (bm BasicManager) DefaultGenesis() map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range bm {
		genesis[b.Name()] = b.DefaultGenesis()
	}
	return genesis
}

// Provided default genesis information for all modules
func (bm BasicManager) ValidateGenesis(genesis map[string]json.RawMessage) error {
	for _, b := range bm {
		if err := b.ValidateGenesis(genesis[b.Name()]); err != nil {
			return err
		}
	}
	return nil
}

// add all tx commands to the rootTxCmd
func (bm BasicManager) AddTxCommands(rootTxCmd *cobra.Command, cdc *amino.Codec) {
	for _, b := range bm {
		if cmd := b.GetTxCmds(cdc); cmd != nil {
			rootTxCmd.AddCommand(cmd...)
			rootTxCmd.AddCommand(bctypes.LineBreak)
		}
	}
}

// add all query commands to the rootQueryCmd
func (bm BasicManager) AddQueryCommands(rootQueryCmd *cobra.Command, cdc *amino.Codec) {
	for _, b := range bm {
		if cmd := b.GetQueryCmds(cdc); cmd != nil {
			rootQueryCmd.AddCommand(cmd...)
			rootQueryCmd.AddCommand(bctypes.LineBreak)
		}
	}
}

//_________________________________________________________
// AppModuleGenesis is the standard form for an application module genesis functions
type AppModuleGenesis interface {
	AppModuleBasic
	InitGenesis(context.Context, *baseabci.BaseApp, json.RawMessage) []abci.ValidatorUpdate
	ExportGenesis(context.Context) json.RawMessage
}

// AppModule is the standard form for an application module
type AppModule interface {
	AppModuleGenesis

	// registers
	RegisterInvariants(InvariantRegistry)
	RegisterQuerier(QueryRegistry)

	BeginBlock(context.Context, abci.RequestBeginBlock)
	EndBlock(context.Context, abci.RequestEndBlock) []abci.ValidatorUpdate
}

//___________________________
// app module
type GenesisOnlyAppModule struct {
	AppModuleGenesis
}

// NewGenesisOnlyAppModule creates a new GenesisOnlyAppModule object
func NewGenesisOnlyAppModule(amg AppModuleGenesis) AppModule {
	return GenesisOnlyAppModule{
		AppModuleGenesis: amg,
	}
}

// register invariants
func (GenesisOnlyAppModule) RegisterInvariants(_ InvariantRegistry) {}

// register querier
func (GenesisOnlyAppModule) RegisterQuerier(_ QueryRegistry) {}

// module begin-block
func (gam GenesisOnlyAppModule) BeginBlock(ctx context.Context, req abci.RequestBeginBlock) {}

// module end-block
func (GenesisOnlyAppModule) EndBlock(_ context.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

//____________________________________________________________________________
// module manager provides the high level utility for managing and executing
// operations for a group of modules
type Manager struct {
	Modules            map[string]AppModule
	OrderInitGenesis   []string
	OrderExportGenesis []string
	OrderBeginBlockers []string
	OrderEndBlockers   []string
}

// NewModuleManager creates a new Manager object
func NewManager(modules ...AppModule) *Manager {

	moduleMap := make(map[string]AppModule)
	var modulesStr []string
	for _, module := range modules {
		moduleMap[module.Name()] = module
		modulesStr = append(modulesStr, module.Name())
	}

	return &Manager{
		Modules:            moduleMap,
		OrderInitGenesis:   modulesStr,
		OrderExportGenesis: modulesStr,
		OrderBeginBlockers: modulesStr,
		OrderEndBlockers:   modulesStr,
	}
}

// set the order of init genesis calls
func (m *Manager) SetOrderInitGenesis(moduleNames ...string) {
	m.OrderInitGenesis = moduleNames
}

// set the order of export genesis calls
func (m *Manager) SetOrderExportGenesis(moduleNames ...string) {
	m.OrderExportGenesis = moduleNames
}

// set the order of set begin-blocker calls
func (m *Manager) SetOrderBeginBlockers(moduleNames ...string) {
	m.OrderBeginBlockers = moduleNames
}

// set the order of set end-blocker calls
func (m *Manager) SetOrderEndBlockers(moduleNames ...string) {
	m.OrderEndBlockers = moduleNames
}

// register all module routes and module invariant routes
func (m *Manager) RegisterInvariants(ir InvariantRegistry) {
	for _, module := range m.Modules {
		module.RegisterInvariants(ir)
	}
}

// register all module querier routes
func (m *Manager) RegisterQueriers(qr QueryRegistry) {
	for _, module := range m.Modules {
		module.RegisterQuerier(qr)
	}
}

// register all module mapper and hooks routes
func (m *Manager) RegisterMapperAndHooks(hmr HooksMapperRegistry, paramsInitializerModule string, ps ...ParamSet) {
	mhs := make(map[string]MapperWithHooks)
	for _, module := range m.Modules {
		mh := module.GetMapperAndHooks()
		if !mh.IsNil() {
			if module.Name() == paramsInitializerModule {
				mh.Mapper.(ParamsInitializer).RegisterParamSet(ps...)
			}
			mhs[module.Name()] = mh
		}
	}
	hmr.RegisterHooksMapper(mhs)
}

// perform init genesis functionality for modules
func (m *Manager) InitGenesis(ctx context.Context, bapp *baseabci.BaseApp, genesisData map[string]json.RawMessage) abci.ResponseInitChain {
	var validatorUpdates []abci.ValidatorUpdate
	for _, moduleName := range m.OrderInitGenesis {
		if genesisData[moduleName] == nil {
			continue
		}
		moduleValUpdates := m.Modules[moduleName].InitGenesis(ctx, bapp, genesisData[moduleName])

		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator InitGenesis updates already set by a previous module")
			}
			validatorUpdates = moduleValUpdates
		}
	}
	return abci.ResponseInitChain{
		Validators: validatorUpdates,
	}
}

// perform export genesis functionality for modules
func (m *Manager) ExportGenesis(ctx context.Context) map[string]json.RawMessage {
	genesisData := make(map[string]json.RawMessage)
	for _, moduleName := range m.OrderExportGenesis {
		genesisData[moduleName] = m.Modules[moduleName].ExportGenesis(ctx)
	}
	return genesisData
}

// BeginBlock performs begin block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) BeginBlock(ctx context.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	ctx = ctx.WithEventManager(btypes.NewEventManager())

	for _, moduleName := range m.OrderBeginBlockers {
		m.Modules[moduleName].BeginBlock(ctx, req)
	}

	return abci.ResponseBeginBlock{
		Events: ctx.EventManager().ABCIEvents(),
	}
}

// EndBlock performs end block functionality for all modules. It creates a
// child context with an event manager to aggregate events emitted from all
// modules.
func (m *Manager) EndBlock(ctx context.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	ctx = ctx.WithEventManager(btypes.NewEventManager())
	validatorUpdates := []abci.ValidatorUpdate{}

	for _, moduleName := range m.OrderEndBlockers {
		moduleValUpdates := m.Modules[moduleName].EndBlock(ctx, req)

		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator EndBlock updates already set by a previous module")
			}

			validatorUpdates = moduleValUpdates
		}
	}

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Events:           ctx.EventManager().ABCIEvents(),
	}
}
