package types

const (
	MapperName          = "mint"
	InflationPhrasesKey = "phrases"
)

var (
	firstBlockTimeKey  = []byte("first_block_time")
	allTotalMintQOSKey = []byte("total_mint_qos")
	totalQOSKey        = []byte("total_qos")
)

func BuildAllTotalMintQOSKey() []byte {
	return allTotalMintQOSKey
}

func BuildFirstBlockTimeKey() []byte {
	return firstBlockTimeKey
}

func BuildInflationPhrasesKey() []byte {
	return []byte(InflationPhrasesKey)
}

func BuildTotalQOSKey() []byte {
	return totalQOSKey
}
