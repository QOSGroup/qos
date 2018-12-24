package types

type ValidatorVoteInfo struct {
	StartHeight         uint64 `json:"startHeight"`
	IndexOffset         uint64 `json:"indexOffset"`
	MissedBlocksCounter uint64 `json:"missedBlocksCounter"`
}

func NewValidatorVoteInfo(startHeight, indexOffset, missedBlocksCounter uint64) ValidatorVoteInfo {
	return ValidatorVoteInfo{
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		MissedBlocksCounter: missedBlocksCounter,
	}
}
