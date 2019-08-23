package types

const (
	MapperName = "approve"
)

var (
	ApproveKey = []byte{0x01}
)

func BuildApproveKey(from []byte, to []byte) []byte {
	key := append(ApproveKey, from...)
	return append(key, to...)
}

func GetApprovePrefixKey() []byte {
	return ApproveKey
}
