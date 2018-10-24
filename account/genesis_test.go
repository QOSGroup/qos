package account

import (
	"encoding/hex"
	"github.com/QOSGroup/qbase/account"
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
		Accounts: []*QOSAccount{
			{
				BaseAccount: account.BaseAccount{
					AccountAddress: accPub.Address().Bytes(),
					Publickey:      accPub,
				},
				Qos: btypes.NewInt(100000000),
				QscList: []*types.QSC{
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
	require.Nil(t, err)
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
