package types

type ValidatorSignInfo struct {
	StartHeight         uint64 `json:"startHeight"`
	IndexOffset         uint64 `json:"indexOffset"`
	MissedBlocksCounter uint64 `json:"missedBlocksCounter"`
}

func NewValidatorSignInfo(startHeight uint64, indexOffset uint64, missedBlocksCounter uint64) ValidatorSignInfo {
	return ValidatorSignInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		MissedBlocksCounter: missedBlocksCounter,
	}
}
