package client

import (
	"fmt"
	"github.com/QOSGroup/qbase/client"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"regexp"
	"strconv"
	"strings"
)

// Get default chain-id.
func GetDefaultChainId() (string, error) {
	chainId := viper.GetString(client.FlagChainID)
	if len(chainId) == 0 {
		chainId, err := btypes.GetChainID(app.DefaultNodeHome)
		if err != nil {
			return "", errors.New("use --chain-id flag to specify your chain")
		}
		return chainId, nil
	}
	return chainId, nil
}

// Parse QOS and QSCs from string
// str example : 100qos,100qstar
func ParseCoins(str string) (btypes.BigInt, types.QSCs, error) {
	reDnm := `[[:alpha:]][[:alnum:]]{2,15}`
	reAmt := `[[:digit:]]+`
	reSpc := `[[:space:]]*`
	reCoin := regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))

	arr := strings.Split(str, ",")
	qos := btypes.ZeroInt()
	qscs := types.QSCs{}
	for _, q := range arr {
		coin := reCoin.FindStringSubmatch(q)
		coin[2] = strings.TrimSpace(coin[2])
		amount, err := strconv.ParseInt(strings.TrimSpace(coin[1]), 10, 64)
		if err != nil {
			return btypes.ZeroInt(), nil, err
		}
		if strings.ToLower(coin[2]) == "qos" {
			qos = btypes.NewInt(amount)
		} else {
			qscs = append(qscs, &types.QSC{
				coin[2],
				btypes.NewInt(amount),
			})
		}

	}

	return qos, qscs, nil
}
