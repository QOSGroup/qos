package types

import (
	"encoding/json"
	"github.com/tendermint/go-amino"
)

type GenesisState map[string]json.RawMessage

func (gs GenesisState) UnmarshalModuleState(cdc *amino.Codec, module string, state interface{}) interface{} {
	cdc.MustUnmarshalJSON(gs[module], state)

	return state
}
