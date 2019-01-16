package staking

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/QOSGroup/qbase/client/context"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	"strings"
	"time"

	qcliacc "github.com/QOSGroup/qbase/client/account"
	stakemapper "github.com/QOSGroup/qos/module/eco/mapper"
	staketypes "github.com/QOSGroup/qos/module/eco/types"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
)

const (
	flagActive   = "active"
	activeDesc   = "active"
	inactiveDesc = "inactive"

	inactiveRevokeDesc        = "Revoked"
	inactiveMissVoteBlockDesc = "Kicked"
	inactiveMaxValidatorDesc  = "Replaced"
)

type validatorDisplayInfo struct {
	Name            string         `json:"name"`
	Owner           btypes.Address `json:"owner"`
	ValidatorAddr   string         `json:"validatorAddress"`
	ValidatorPubKey crypto.PubKey  `json:"validatorPubkey"`
	BondTokens      uint64         `json:"bondTokens"` //不能超过int64最大值
	Description     string         `json:"description"`

	Status         string    `json:"status"`
	InactiveDesc   string    `json:"InactiveDesc"`
	InactiveTime   time.Time `json:"inactiveTime"`
	InactiveHeight uint64    `json:"inactiveHeight"`

	BondHeight uint64 `json:"bondHeight"`
}

func toValidatorDisplayInfo(validator staketypes.Validator) validatorDisplayInfo {
	info := validatorDisplayInfo{
		Name:            validator.Name,
		Owner:           validator.Owner,
		ValidatorPubKey: validator.ValidatorPubKey,
		BondTokens:      validator.BondTokens,
		Description:     validator.Description,
		InactiveTime:    validator.InactiveTime,
		InactiveHeight:  validator.InactiveHeight,
		BondHeight:      validator.BondTokens,
	}

	if validator.Status == staketypes.Active {
		info.Status = activeDesc
	} else {
		info.Status = inactiveDesc
	}

	if validator.InactiveCode == staketypes.Revoke {
		info.InactiveDesc = inactiveRevokeDesc
	} else if validator.InactiveCode == staketypes.MissVoteBlock {
		info.InactiveDesc = inactiveMissVoteBlockDesc
	} else if validator.InactiveCode == staketypes.MaxValidator {
		info.InactiveDesc = inactiveMaxValidatorDesc
	}

	info.ValidatorAddr = strings.ToUpper(hex.EncodeToString(validator.ValidatorPubKey.Address()))

	return info
}

func queryValidatorInfoCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Query validator's info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			ownerAddress, err := qcliacc.GetAddrFromFlag(cliCtx, flagOwner)
			if err != nil {
				return err
			}

			validator, err := getValidator(cliCtx, ownerAddress)
			if err != nil {
				return err
			}
			return cliCtx.PrintResult(toValidatorDisplayInfo(validator))
		},
	}

	cmd.Flags().String(flagOwner, "", "validator's owner address")
	cmd.MarkFlagRequired(flagOwner)

	return cmd
}

func queryAllValidatorsCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query all validators info",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			opts := buildQueryOptions()

			subspace := "/store/validator/subspace"
			result, err := node.ABCIQueryWithOptions(subspace, staketypes.BulidValidatorPrefixKey(), opts)

			if err != nil {
				return err
			}

			valueBz := result.Response.GetValue()
			if len(valueBz) == 0 {
				return errors.New("response empty value")
			}

			var validators []validatorDisplayInfo

			var vKVPair []store.KVPair
			cdc.UnmarshalBinaryLengthPrefixed(valueBz, &vKVPair)
			for _, kv := range vKVPair {
				var validator staketypes.Validator
				cdc.UnmarshalBinaryBare(kv.Value, &validator)
				validators = append(validators, toValidatorDisplayInfo(validator))
			}

			cliCtx.PrintResult(validators)

			return nil
		},
	}

	// cmd.Flags().Bool(flagActive, false, "only query active status validator")

	return cmd
}

func queryValidatorMissedVoteInfoCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-miss-vote",
		Short: "Query validator miss vote info in the nearest voting window",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ownerAddress, err := qcliacc.GetAddrFromFlag(cliCtx, flagOwner)
			if err != nil {
				return err
			}

			voteSummaryDisplay, err := queryVotesInfoByOwner(cliCtx, ownerAddress)
			if err != nil {
				return err
			}

			cliCtx.PrintResult(voteSummaryDisplay)
			return nil
		},
	}

	cmd.Flags().String(flagOwner, "", "validator's owner address")
	cmd.MarkFlagRequired(flagOwner)

	return cmd
}

type voteSummary struct {
	StartHeight int64            `json:"startHeight"`
	EndHeight   int64            `json:"endHeight"`
	MissCount   int8             `json:"missCount"`
	Votes       []voteInfoDetail `json:"voteDetail"`
}

type voteInfoDetail struct {
	Height int64
	Vote   bool
}

func queryVotesInfoByOwner(ctx context.CLIContext, ownerAddress btypes.Address) (voteSummary, error) {

	voteSummaryDisplay := voteSummary{}
	stakeConfig, err := getStakeConfig(ctx)
	if err != nil {
		return voteSummaryDisplay, err
	}

	windownLength := uint64(stakeConfig.ValidatorVotingStatusLen)
	votesInfo := make([]voteInfoDetail, 0, windownLength)

	val, err := getValidator(ctx, ownerAddress)
	if err != nil {
		return voteSummaryDisplay, err
	}

	validatorAddress := btypes.Address(val.ValidatorPubKey.Address())

	voteInfo, err := getValidatorVoteInfo(ctx, validatorAddress)
	if err != nil {
		return voteSummaryDisplay, err
	}

	voteInfoMap, _, err := queryValidatorVotesInWindow(ctx, validatorAddress)
	if err != nil {
		return voteSummaryDisplay, err
	}

	votedBlockLength := voteInfo.IndexOffset - 1

	endWindowHeight := voteInfo.StartHeight + votedBlockLength
	startWindowHeight := uint64(1)
	if votedBlockLength <= windownLength {
		startWindowHeight = voteInfo.StartHeight
	} else {
		startWindowHeight = endWindowHeight - windownLength + 1
	}

	voteSummaryDisplay.StartHeight = int64(startWindowHeight)
	voteSummaryDisplay.EndHeight = int64(endWindowHeight)

	i := int8(0)
	for h := endWindowHeight; h >= startWindowHeight; h-- {
		index := h % windownLength
		voted := true

		if v, ok := voteInfoMap[index]; ok {
			voted = v
		}

		if !voted {
			i++
			votesInfo = append(votesInfo, voteInfoDetail{int64(h), voted})
		}
	}

	voteSummaryDisplay.Votes = votesInfo
	voteSummaryDisplay.MissCount = i
	return voteSummaryDisplay, nil
}

func getStakeConfig(ctx context.CLIContext) (staketypes.StakeParams, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return staketypes.StakeParams{}, err
	}

	path := "/store/stake/key"
	key := []byte(staketypes.BuildStakeParamsKey())

	result, err := node.ABCIQueryWithOptions(path, key, buildQueryOptions())
	if err != nil {
		return staketypes.StakeParams{}, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return staketypes.StakeParams{}, errors.New("response empty value. getStakeConfig is empty")
	}

	var stakeConfig staketypes.StakeParams
	{
	}
	ctx.Codec.UnmarshalBinaryBare(valueBz, &stakeConfig)

	return stakeConfig, nil
}

func getValidatorVoteInfo(ctx context.CLIContext, validatorAddr btypes.Address) (staketypes.ValidatorVoteInfo, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return staketypes.ValidatorVoteInfo{}, err
	}

	path := string(stakemapper.BuildVoteInfoStoreQueryPath())
	key := stakemapper.BuildValidatorVoteInfoKey(validatorAddr)

	result, err := node.ABCIQueryWithOptions(path, key, buildQueryOptions())
	if err != nil {
		return staketypes.ValidatorVoteInfo{}, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return staketypes.ValidatorVoteInfo{}, errors.New("response empty value. validatorVoteInfo is empty")
	}

	var voteInfo staketypes.ValidatorVoteInfo
	ctx.Codec.UnmarshalBinaryBare(valueBz, &voteInfo)

	return voteInfo, nil
}

func queryValidatorVotesInWindow(ctx context.CLIContext, validatorAddr btypes.Address) (map[uint64]bool, int64, error) {

	voteInWindowInfo := make(map[uint64]bool)

	node, err := ctx.GetNode()
	if err != nil {
		return voteInWindowInfo, 0, err
	}

	storePath := "/" + strings.Join([]string{"store", stakemapper.VoteInfoMapperName, "subspace"}, "/")
	key := stakemapper.BuildValidatorVoteInfoInWindowPrefixKey(validatorAddr)

	result, err := node.ABCIQueryWithOptions(storePath, key, buildQueryOptions())
	if err != nil {
		return nil, 0, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return voteInWindowInfo, result.Response.Height, nil
	}

	var vKVPair []store.KVPair
	ctx.Codec.UnmarshalBinaryLengthPrefixed(valueBz, &vKVPair)

	for _, kv := range vKVPair {
		k := kv.Key
		var vote bool
		index := binary.LittleEndian.Uint64(k[(len(k) - 8):])
		ctx.Codec.UnmarshalBinaryBare(kv.Value, &vote)
		voteInWindowInfo[index] = vote
	}

	return voteInWindowInfo, result.Response.Height, nil
}

func getValidator(ctx context.CLIContext, ownerAddress btypes.Address) (staketypes.Validator, error) {

	node, err := ctx.GetNode()
	if err != nil {
		return staketypes.Validator{}, err
	}

	result, err := node.ABCIQueryWithOptions(string(staketypes.BuildValidatorStoreQueryPath()), staketypes.BuildOwnerWithValidatorKey(ownerAddress), buildQueryOptions())
	if err != nil {
		return staketypes.Validator{}, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return staketypes.Validator{}, errors.New("owner does't have validator")
	}

	var address btypes.Address
	ctx.Codec.UnmarshalBinaryBare(valueBz, &address)

	key := staketypes.BuildValidatorKey(address)
	result, err = node.ABCIQueryWithOptions(string(staketypes.BuildValidatorStoreQueryPath()), key, buildQueryOptions())
	if err != nil {
		return staketypes.Validator{}, err
	}

	valueBz = result.Response.GetValue()
	if len(valueBz) == 0 {
		return staketypes.Validator{}, errors.New("response empty value")
	}

	var validator staketypes.Validator
	ctx.Codec.UnmarshalBinaryBare(valueBz, &validator)
	return validator, nil
}

func buildQueryOptions() client.ABCIQueryOptions {
	height := viper.GetInt64(bctypes.FlagHeight)
	if height <= 0 {
		height = 0
	}

	trust := viper.GetBool(bctypes.FlagTrustNode)

	return client.ABCIQueryOptions{
		Height: height,
		Prove:  trust,
	}
}
