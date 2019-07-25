package guardian

import (
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	bmapper "github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/guardian/mapper"
	"github.com/QOSGroup/qos/module/guardian/txs"
	gtypes "github.com/QOSGroup/qos/module/guardian/types"
	"github.com/QOSGroup/qos/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func defaultContext() context.Context {
	mapperMap := make(map[string]bmapper.IMapper)

	approveMapper := mapper.NewMapper()
	approveMapper.SetCodec(txs.Cdc)
	approveKey := approveMapper.GetStoreKey()
	mapperMap[mapper.MapperName] = approveMapper

	accountMapper := bacc.NewAccountMapper(nil, types.ProtoQOSAccount)
	accountMapper.SetCodec(txs.Cdc)
	acountKey := accountMapper.GetStoreKey()
	mapperMap[bacc.AccountMapperName] = accountMapper

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(approveKey, btypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(acountKey, btypes.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := context.NewContext(cms, abci.Header{}, false, log.NewNopLogger(), mapperMap)
	return ctx
}

func TestValidateGenesis(t *testing.T) {
	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input    GenesisState
		expected bool
	}{
		{GenesisState{[]gtypes.Guardian{{"g1", gtypes.Genesis, addr1, nil}}}, true},
		{GenesisState{[]gtypes.Guardian{{"g1", gtypes.Ordinary, addr1, nil}}}, false},
		{GenesisState{[]gtypes.Guardian{{"g1", gtypes.Genesis, addr1, addr2}}}, false},
		{GenesisState{[]gtypes.Guardian{{"g1", gtypes.Genesis, addr1, nil}, {"g2", gtypes.Genesis, addr1, nil}}}, false},
	}

	for tcIndex, tc := range cases {
		res := ValidateGenesis(tc.input) == nil
		require.Equal(t, tc.expected, res, "tc #%d", tcIndex)
	}
}

func TestExportGenesis(t *testing.T) {
	addr1 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())
	addr2 := btypes.Address(ed25519.GenPrivKey().PubKey().Address())

	cases := []struct {
		input GenesisState
	}{
		{GenesisState{[]gtypes.Guardian{{"g1", gtypes.Genesis, addr1, nil}, {"g2", gtypes.Genesis, addr2, nil}}}},
		{GenesisState{[]gtypes.Guardian{{"g1", gtypes.Genesis, addr1, nil}}}},
	}

	for tcIndex, tc := range cases {
		ctx := defaultContext()
		InitGenesis(ctx, tc.input)
		export := ExportGenesis(ctx)
		require.True(t, tc.input.Equals(export), "tc #%d", tcIndex)
	}
}
