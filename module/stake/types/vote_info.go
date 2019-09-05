package types

type ValidatorVoteInfo struct {
	StartHeight         int64 `json:"start_height"`
	IndexOffset         int64 `json:"index_offset"` //统计截止高度=StartHeight+IndexOffset-1
	MissedBlocksCounter int64 `json:"missed_blocks_counter"`
}

func NewValidatorVoteInfo(startHeight, indexOffset, missedBlocksCounter int64) ValidatorVoteInfo {
	return ValidatorVoteInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		MissedBlocksCounter: missedBlocksCounter,
	}
}
