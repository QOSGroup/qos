package app

import (
	"fmt"
	"testing"

	bacc "github.com/QOSGroup/qbase/account"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func genGenesisState() *GenesisState {
	caPub := ed25519.GenPrivKey().PubKey()
	accPub := ed25519.GenPrivKey().PubKey()
	genesisState := GenesisState{
		CAPubKey: caPub,
		Accounts: []*account.QOSAccount{
			{
				BaseAccount: bacc.BaseAccount{
					AccountAddress: accPub.Address().Bytes(),
					Publickey:      accPub,
				},
				QOS: btypes.NewInt(100000000),
				QSCs: types.QSCs{
					{
						Name:   "QSC1",
						Amount: btypes.NewInt(100000000),
					},
				},
			},
		},
	}
	return &genesisState
}

func TestGenesisStateJSONMarshal(t *testing.T) {
	genesisState := genGenesisState()
	genesisStateJson, err := cdc.MarshalJSON(genesisState)
	fmt.Println(string(genesisStateJson))
	require.Nil(t, err)
	umState := GenesisState{}
	err = cdc.UnmarshalJSON(genesisStateJson, &umState)
	require.Nil(t, err)
	require.Equal(t, *genesisState, umState)
}
