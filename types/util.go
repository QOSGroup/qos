package types

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.qoscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.qosd")
)

// Get default chain-id.
func GetDefaultChainId() (string, error) {
	chainId := viper.GetString("chain-id")
	if len(chainId) == 0 {
		chainId, err := btypes.GetChainID(DefaultNodeHome)
		if err != nil {
			return "", errors.New("use --chain-id flag to specify your chain")
		}
		return chainId, nil
	}
	return chainId, nil
}

// Parse QOS and QSCs from string
// str example : 100qos,100qstar
func ParseCoins(str string) (btypes.BigInt, QSCs, error) {
	if len(str) == 0 {
		return btypes.ZeroInt(), QSCs{}, nil
	}
	reDnm := `[[:alpha:]][[:alnum:]]{2,15}`
	reAmt := `[[:digit:]]+`
	reSpc := `[[:space:]]*`
	reCoin := regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))

	arr := strings.Split(str, ",")
	qos := btypes.ZeroInt()
	qscs := QSCs{}
	for _, q := range arr {
		coin := reCoin.FindStringSubmatch(q)
		if len(coin) != 3 {
			return btypes.ZeroInt(), nil, fmt.Errorf("coins str: %s parse faild", q)
		}
		coin[2] = strings.TrimSpace(coin[2])
		amount, ok := btypes.NewIntFromString(strings.TrimSpace(coin[1]))
		if !ok {
			return btypes.ZeroInt(), nil, fmt.Errorf("coins str: %s parse faild", q)
		}
		if strings.ToUpper(coin[2]) == QOSCoinName {
			qos = amount
		} else {
			qscs = append(qscs, &QSC{
				coin[2],
				amount,
			})
		}

	}

	return qos, qscs, nil
}

// Uint64ToBigEndian - marshals uint64 to a bigendian byte slice so it can be sorted
func Uint64ToBigEndian(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

func FileMD5(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// 从flag中解析Int值
func GetIntFromFlag(flag string, allowZero bool) (amount btypes.BigInt, err error) {
	intstr := viper.GetString(flag)
	if len(intstr) == 0 {
		if !allowZero {
			return btypes.ZeroInt(), fmt.Errorf("%s is empty", flag)
		} else {
			return btypes.ZeroInt(), nil
		}
	}
	amount, ok := btypes.NewIntFromString(strings.TrimSpace(intstr))
	if !ok {
		return btypes.ZeroInt(), fmt.Errorf("invalid value of %s", flag)
	}

	if !amount.GT(btypes.ZeroInt()) && !allowZero {
		return btypes.ZeroInt(), fmt.Errorf("%s must be positive", flag)
	}

	return amount, nil
}
