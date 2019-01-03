package init

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qos/app"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
)

func loadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc tmtypes.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

// ExportGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
func ExportGenesisFile(
	genFile, chainID string, appState json.RawMessage,
) error {

	genDoc := tmtypes.GenesisDoc{
		ChainID:  chainID,
		AppState: appState,
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genFile)
}

func initializeEmptyGenesis(
	cdc *amino.Codec, genFile, chainID string, overwrite bool,
) (appState json.RawMessage, err error) {

	if !overwrite && common.FileExists(genFile) {
		return nil, fmt.Errorf("genesis.json file already exists: %v", genFile)
	}

	return cdc.MarshalJSONIndent(app.NewDefaultGenesisState(), "", " ")
}
