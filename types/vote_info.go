package types

type ValidatorVoteInfo struct {
	StartHeight         uint64 `json:"startHeight"`
	IndexOffset         uint64 `json:"indexOffset"` //统计截止高度=StartHeight+IndexOffset-1
	MissedBlocksCounter uint64 `json:"missedBlocksCounter"`
}

func NewValidatorVoteInfo(startHeight, indexOffset, missedBlocksCounter uint64) ValidatorVoteInfo {
	return ValidatorVoteInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		MissedBlocksCounter: missedBlocksCounter,
	}
}
