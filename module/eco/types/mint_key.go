package types

const (
	MintMapperName = "mint"
	MintParamsKey  = "mintparams"
)

var (
	firstBlockTimeKey  = []byte("first_block_time")
	allTotalMintQOSKey = []byte("total_mint_qos")
)

func BuildAllTotalMintQOSKey() []byte {
	return allTotalMintQOSKey
}

func BuildFirstBlockTimeKey() []byte {
	return firstBlockTimeKey
}

func BuildMintParamsKey() []byte {
	return []byte(MintParamsKey)
}
