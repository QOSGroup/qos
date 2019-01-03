package init

import (
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qos/app"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenQOSGenesisDoc(ctx *server.Context, cdc *amino.Codec, chainID string, nodeValidatorPubKey crypto.PubKey) (tmtypes.GenesisDoc, error) {

	appState, _ := cdc.MarshalJSONIndent(app.NewDefaultGenesisState(), "", " ")

	return tmtypes.GenesisDoc{
		ChainID:  chainID,
		AppState: appState,
	}, nil

}
