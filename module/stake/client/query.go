package client

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/QOSGroup/qbase/client/context"
	bctypes "github.com/QOSGroup/qbase/client/types"
	"github.com/QOSGroup/qbase/store"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/params"
	"github.com/QOSGroup/qos/module/stake/mapper"
	"github.com/QOSGroup/qos/module/stake/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	"strings"
	"time"

	qcliacc "github.com/QOSGroup/qbase/client/account"
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
	inactiveDoubleDesc        = "DoubleSign"
)

type validatorDisplayInfo struct {
	Owner           btypes.Address    `json:"owner"`
	ValidatorAddr   string            `json:"validatorAddress"`
	ValidatorPubKey crypto.PubKey     `json:"validatorPubkey"`
	BondTokens      uint64            `json:"bondTokens"` //不能超过int64最大值
	Description     types.Description `json:"description"`
	Commission      types.Commission  `json:"commission"`

	Status         string    `json:"status"`
	InactiveDesc   string    `json:"InactiveDesc"`
	InactiveTime   time.Time `json:"inactiveTime"`
	InactiveHeight uint64    `json:"inactiveHeight"`

	BondHeight uint64 `json:"bondHeight"`
}

func toValidatorDisplayInfo(validator types.Validator) validatorDisplayInfo {
	info := validatorDisplayInfo{
		Owner:           validator.Owner,
		ValidatorPubKey: validator.ValidatorPubKey,
		BondTokens:      validator.BondTokens,
		Description:     validator.Description,
		InactiveTime:    validator.InactiveTime,
		InactiveHeight:  validator.InactiveHeight,
		BondHeight:      validator.BondHeight,
		Commission:      validator.Commission,
	}

	if validator.Status == types.Active {
		info.Status = activeDesc
	} else {
		info.Status = inactiveDesc
	}

	if validator.InactiveCode == types.Revoke {
		info.InactiveDesc = inactiveRevokeDesc
	} else if validator.InactiveCode == types.MissVoteBlock {
		info.InactiveDesc = inactiveMissVoteBlockDesc
	} else if validator.InactiveCode == types.MaxValidator {
		info.InactiveDesc = inactiveMaxValidatorDesc
	} else if validator.InactiveCode == types.DoubleSign {
		info.InactiveDesc = inactiveDoubleDesc
	}

	info.ValidatorAddr = strings.ToUpper(hex.EncodeToString(validator.ValidatorPubKey.Address()))

	return info
}

func queryValidatorInfoCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [validator-owner]",
		Short: "Query validator's info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			ownerAddress, err := qcliacc.GetAddrFromValue(cliCtx, args[0])
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

	return cmd
}

func queryDelegationInfoCommand(cdc *go_amino.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delegation",
		Short: "Query delegation info",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var owner btypes.Address
			var delegator btypes.Address

			if o, err := qcliacc.GetAddrFromFlag(cliCtx, flagOwner); err == nil {
				owner = o
			}

			if d, err := qcliacc.GetAddrFromFlag(cliCtx, flagDelegator); err == nil {
				delegator = d
			}

			var path = types.BuildGetDelegationCustomQueryPath(delegator, owner)

			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result mapper.DelegationQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	cmd.Flags().String(flagOwner, "", "validator's owner address")
	cmd.Flags().String(flagDelegator, "", "delegator address")
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagDelegator)

	return cmd
}

func queryDelegationsToCommand(cdc *go_amino.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delegations-to [validator-owner]",
		Short: "Query all delegations made to one validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var owner btypes.Address

			if o, err := qcliacc.GetAddrFromValue(cliCtx, args[0]); err == nil {
				owner = o
			}

			var path = types.BuildQueryDelegationsByOwnerCustomQueryPath(owner)

			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result []mapper.DelegationQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryDelegationsCommand(cdc *go_amino.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delegations [delegator]",
		Short: "Query all delegations made by one delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var delegator btypes.Address

			if d, err := qcliacc.GetAddrFromValue(cliCtx, args[0]); err == nil {
				delegator = d
			}

			var path = types.BuildQueryDelegationsByDelegatorCustomQueryPath(delegator)

			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result []mapper.DelegationQueryResult
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

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
			result, err := node.ABCIQueryWithOptions(subspace, types.BulidValidatorPrefixKey(), opts)

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
				var validator types.Validator
				cdc.UnmarshalBinaryBare(kv.Value, &validator)
				validators = append(validators, toValidatorDisplayInfo(validator))
			}

			cliCtx.PrintResult(validators)

			return nil
		},
	}

	return cmd
}

func queryValidatorMissedVoteInfoCommand(cdc *go_amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-miss-vote",
		Short: "Query validator miss vote info in the nearest voting window",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ownerAddress, err := qcliacc.GetAddrFromValue(cliCtx, args[0])
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

	windownLength, err := getStakeConfig(ctx)
	if err != nil {
		return voteSummaryDisplay, err
	}

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

func getStakeConfig(ctx context.CLIContext) (uint64, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return 0, err
	}

	path := "/store/params/key"
	key := params.BuildParamKey(types.ParamSpace, types.KeyValidatorVotingStatusLen)

	result, err := node.ABCIQueryWithOptions(path, key, buildQueryOptions())
	if err != nil {
		return 0, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return 0, errors.New("response empty value. getStakeConfig is empty")
	}

	var length uint64
	ctx.Codec.UnmarshalBinaryBare(valueBz, &length)

	return length, nil

	return 0, nil
}

func getValidatorVoteInfo(ctx context.CLIContext, validatorAddr btypes.Address) (types.ValidatorVoteInfo, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return types.ValidatorVoteInfo{}, err
	}

	path := string(types.BuildStakeStoreQueryPath())
	key := types.BuildValidatorVoteInfoKey(validatorAddr)

	result, err := node.ABCIQueryWithOptions(path, key, buildQueryOptions())
	if err != nil {
		return types.ValidatorVoteInfo{}, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return types.ValidatorVoteInfo{}, errors.New("response empty value. validatorVoteInfo is empty")
	}

	var voteInfo types.ValidatorVoteInfo
	ctx.Codec.UnmarshalBinaryBare(valueBz, &voteInfo)

	return voteInfo, nil
}

func queryValidatorVotesInWindow(ctx context.CLIContext, validatorAddr btypes.Address) (map[uint64]bool, int64, error) {

	voteInWindowInfo := make(map[uint64]bool)

	node, err := ctx.GetNode()
	if err != nil {
		return voteInWindowInfo, 0, err
	}

	storePath := "/" + strings.Join([]string{"store", types.MapperName, "subspace"}, "/")
	key := types.BuildValidatorVoteInfoInWindowPrefixKey(validatorAddr)

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

func getValidator(ctx context.CLIContext, ownerAddress btypes.Address) (types.Validator, error) {

	node, err := ctx.GetNode()
	if err != nil {
		return types.Validator{}, err
	}

	result, err := node.ABCIQueryWithOptions(string(types.BuildStakeStoreQueryPath()), types.BuildOwnerWithValidatorKey(ownerAddress), buildQueryOptions())
	if err != nil {
		return types.Validator{}, err
	}

	valueBz := result.Response.GetValue()
	if len(valueBz) == 0 {
		return types.Validator{}, errors.New("owner does't have validator")
	}

	var address btypes.Address
	ctx.Codec.UnmarshalBinaryBare(valueBz, &address)

	key := types.BuildValidatorKey(address)
	result, err = node.ABCIQueryWithOptions(string(types.BuildStakeStoreQueryPath()), key, buildQueryOptions())
	if err != nil {
		return types.Validator{}, err
	}

	valueBz = result.Response.GetValue()
	if len(valueBz) == 0 {
		return types.Validator{}, errors.New("response empty value")
	}

	var validator types.Validator
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

func queryUnbondingsCommand(cdc *go_amino.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "unbondings [delegator]",
		Short: "Query all unbonding delegations by delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var delegator btypes.Address

			if o, err := qcliacc.GetAddrFromValue(cliCtx, args[0]); err == nil {
				delegator = o
			}

			var path = types.BuildQueryUnbondingsByDelegatorCustomQueryPath(delegator)

			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result []types.UnbondingDelegationInfo
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}

func queryRedelegationsCommand(cdc *go_amino.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "redelegations [delegator]",
		Short: "Query all redelegations by delegator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var delegator btypes.Address

			if o, err := qcliacc.GetAddrFromValue(cliCtx, args[0]); err == nil {
				delegator = o
			}

			var path = types.BuildQueryRedelegationsByDelegatorCustomQueryPath(delegator)

			res, err := cliCtx.Query(path, []byte(""))
			if err != nil {
				return err
			}

			var result []types.RedelegationInfo
			cliCtx.Codec.UnmarshalJSON(res, &result)
			return cliCtx.PrintResult(result)
		},
	}

	return cmd
}
