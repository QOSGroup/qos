package init

import (
	"io/ioutil"

	amino "github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
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
