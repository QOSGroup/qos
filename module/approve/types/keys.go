package types

const (
	MapperName = "approve"
)

var (
	ApproveKey = []byte{0x01} // 预授权存储前缀
)

func BuildApproveByFromKey(from []byte) []byte {
	return append(ApproveKey, from...)
}

func BuildApproveKey(from []byte, to []byte) []byte {
	key := append(ApproveKey, from...)
	return append(key, to...)
}

func GetApprovePrefixKey() []byte {
	return ApproveKey
}
