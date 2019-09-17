package types

const (
	MapperName = "mint"
)

var (
	firstBlockTimeKey   = []byte("first_block_time")
	allTotalMintQOSKey  = []byte("total_mint_qos")
	totalQOSKey         = []byte("total_qos")
	InflationPhrasesKey = []byte("phrases")
)

func BuildAllTotalMintQOSKey() []byte {
	return allTotalMintQOSKey
}

func BuildFirstBlockTimeKey() []byte {
	return firstBlockTimeKey
}

func BuildInflationPhrasesKey() []byte {
	return InflationPhrasesKey
}

func BuildTotalQOSKey() []byte {
	return totalQOSKey
}
