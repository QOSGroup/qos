package types

const (
	MintMapperName = "mint"
	MintParamsKey  = "mintparams"
)

var (
	firstBlockTimeKey = []byte("first_block_time")
)

func BuildFirstBlockTimeKey() []byte {
	return firstBlockTimeKey
}

func BuildMintParamsKey() []byte {
	return []byte(MintParamsKey)
}
