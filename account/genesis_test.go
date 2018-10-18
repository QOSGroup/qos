package account

import (
	"encoding/hex"
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func genGenesisState() *GenesisState {
	caPub := ed25519.GenPrivKey().PubKey()
	accPub := ed25519.GenPrivKey().PubKey()
	genesisState := GenesisState{
		CAPubKey: caPub,
		Accounts: []*GenesisAccount{
			{
				PubKey: accPub,
				QOS:    btypes.NewInt(100000000),
				QSC: []*types.QSC{
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

func TestToAccount(t *testing.T) {
	qosAcc := genNewAccount()
	genAcc := NewGenesisAccount(&qosAcc)
	genQosAcc, err := genAcc.ToAppAccount()
	require.Nil(t, err)
	require.Equal(t, qosAcc, *genQosAcc)
}

func TestGenesisStateJSONMarshal(t *testing.T) {
	genesisState := genGenesisState()
	genesisStateJson, err := cdc.MarshalJSON(genesisState)
	require.Nil(t, err)
	fmt.Println(string(genesisStateJson))
	umState := GenesisState{}
	err = cdc.UnmarshalJSON(genesisStateJson, &umState)
	require.Nil(t, err)
	require.Equal(t, *genesisState, umState)
}

func TestDefaultCoinKey(t *testing.T) {
	addr, pubKey, priKey, err := GenerateCoinKey()
	require.Nil(t, err)
	require.Equal(t, addr, btypes.Address(pubKey.Address()))
	priHex, _ := hex.DecodeString(priKey[2:])
	var discoverPriKey ed25519.PrivKeyEd25519
	cdc.MustUnmarshalBinaryBare(priHex, &discoverPriKey)
	require.Equal(t, discoverPriKey.PubKey(), pubKey)
}
