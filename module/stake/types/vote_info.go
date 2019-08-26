package types

type ValidatorVoteInfo struct {
	StartHeight         uint64 `json:"start_height"`
	IndexOffset         uint64 `json:"index_offset"` //统计截止高度=StartHeight+IndexOffset-1
	MissedBlocksCounter uint64 `json:"missed_blocks_counter"`
}

func NewValidatorVoteInfo(startHeight, indexOffset, missedBlocksCounter uint64) ValidatorVoteInfo {
	return ValidatorVoteInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		MissedBlocksCounter: missedBlocksCounter,
	}
}
