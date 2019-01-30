package init

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/app"
	"github.com/QOSGroup/qos/module/distribution"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"path/filepath"
	"strings"
)

const (
	flagName        = "name"
	flagOwner       = "owner"
	flagPubKey      = "pubkey"
	flagBondTokens  = "tokens"
	flagDescription = "description"
	flagCompound    = "compound"
)

func AddGenesisValidator(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-validator",
		Short: "Add genesis validator to genesis.json",
		Long: `
pubkey is a tendermint validator pubkey. the public key of the validator used in
Tendermint consensus.

home node's home directory.

owner is account address.

ex: pubkey: {"type":"tendermint/PubKeyEd25519","value":"VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA="}

example:

	 qosd add-genesis-validator --home "/.qosd/" --name validatorName --owner address1vdp54s5za8tl4dmf9dcldfzn62y66m40ursfsa --pubkey "VOn2rPx+t7Njdgi+eLb+jBuF175T1b7LAcHElsmIuXA=" --tokens 100

		`,
		RunE: func(_ *cobra.Command, args []string) error {

			home := viper.GetString(cli.HomeFlag)
			genFile := filepath.Join(home, "config", "genesis.json")
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `qosd init` first", genFile)
			}

			name := viper.GetString(flagName)
			if len(name) == 0 {
				return errors.New("name is empty")
			}

			ownerStr := viper.GetString(flagOwner)
			if !strings.HasPrefix(ownerStr, "address") {
				return errors.New("owner is invalid")
			}

			owner, err := btypes.GetAddrFromBech32(ownerStr)
			if err != nil {
				return err
			}
			valPubkey := viper.GetString(flagPubKey)
			if len(valPubkey) == 0 {
				return errors.New("pubkey is empty")
			}
			tokens := uint64(viper.GetInt64(flagBondTokens))
			if tokens <= 0 {
				return errors.New("tokens lte zero")
			}
			desc := viper.GetString(flagDescription)

			bz, err := base64.StdEncoding.DecodeString(valPubkey)
			if err != nil {
				return err
			}
			var cKey ed25519.PubKeyEd25519
			copy(cKey[:], bz)

			val := ecotypes.Validator{
				Name:            name,
				ValidatorPubKey: cKey,
				Owner:           owner,
				BondTokens:      uint64(tokens),
				Status:          ecotypes.Active,
				BondHeight:      1,
				Description:     desc,
			}

			genDoc, err := loadGenesisDoc(cdc, genFile)
			if err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, v := range appState.StakeData.Validators {
				if v.ValidatorPubKey.Equals(val.ValidatorPubKey) {
					return errors.New("validator already exists")
				}
				if bytes.Equal(v.Owner, val.Owner) {
					return fmt.Errorf("owner %s already bind a validator", val.Owner)
				}
			}

			addValidator(&appState, val, viper.GetBool(flagCompound))

			rawMessage, _ := cdc.MarshalJSON(appState)
			genDoc.AppState = rawMessage

			err = genDoc.ValidateAndComplete()
			if err != nil {
				return err
			}

			err = genDoc.SaveAs(genFile)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagName, "", "name for validator")
	cmd.Flags().String(flagOwner, "", "account address")
	cmd.Flags().String(flagPubKey, "", "tendermint consensus validator public key")
	cmd.Flags().Int64(flagBondTokens, 0, "bond tokens amount")
	cmd.Flags().String(flagDescription, "", "description")
	cmd.Flags().String(cli.HomeFlag, types.DefaultNodeHome, "node's home directory")
	cmd.Flags().Bool(flagCompound, false, "whether the income is calculated as compound interest")

	cmd.MarkFlagRequired(flagName)
	cmd.MarkFlagRequired(flagOwner)
	cmd.MarkFlagRequired(flagPubKey)
	cmd.MarkFlagRequired(flagBondTokens)

	return cmd
}

func addValidator(appState *app.GenesisState, validator ecotypes.Validator, isCompound bool) {
	accIndex := -1
	var acc *types.QOSAccount

	for i, qosAcc := range appState.Accounts {
		if qosAcc.GetAddress().EqualsTo(validator.Owner) {
			accIndex = i
			acc = qosAcc
			break
		}
	}

	if accIndex == -1 {
		panic(fmt.Sprintf("owner: %s not exsits", validator.Owner.String()))
	}

	//owner账户扣减
	minusQOS := btypes.NewInt(int64(validator.BondTokens))
	acc.MustMinusQOS(minusQOS)

	//stake:
	appState.StakeData.Validators = append(appState.StakeData.Validators, validator)
	appState.StakeData.DelegatorsInfo = append(appState.StakeData.DelegatorsInfo, ecotypes.DelegationInfo{
		DelegatorAddr: validator.Owner,
		ValidatorAddr: validator.GetValidatorAddress(),
		Amount:        validator.BondTokens,
		IsCompound:    isCompound,
	})

	//distribution
	appState.DistributionData.ValidatorHistoryPeriods = append(appState.DistributionData.ValidatorHistoryPeriods, distribution.ValidatorHistoryPeriodState{
		ValAddress: validator.GetValidatorAddress(),
		Period:     uint64(0),
		Summary:    types.ZeroFraction(),
	})

	appState.DistributionData.ValidatorCurrentPeriods = append(appState.DistributionData.ValidatorCurrentPeriods, distribution.ValidatorCurrentPeriodState{
		ValAddress: validator.GetValidatorAddress(),
		CurrentPeriodSummary: ecotypes.ValidatorCurrentPeriodSummary{
			Fees:   btypes.ZeroInt(),
			Period: uint64(1),
		},
	})

	appState.DistributionData.DelegatorEarningInfos = append(appState.DistributionData.DelegatorEarningInfos, distribution.DelegatorEarningStartState{
		ValAddress:  validator.GetValidatorAddress(),
		DeleAddress: validator.Owner,
		DelegatorEarningsStartInfo: ecotypes.DelegatorEarningsStartInfo{
			PreviousPeriod:        uint64(0),
			BondToken:             validator.BondTokens,
			CurrentStartingHeight: uint64(1),
			FirstDelegateHeight:   uint64(1),
			HistoricalRewardFees:  btypes.ZeroInt(),
		},
	})

	incomeHeight := appState.DistributionData.Params.DelegatorsIncomePeriodHeight + uint64(1)
	appState.DistributionData.DelegatorIncomeHeights = append(appState.DistributionData.DelegatorIncomeHeights, distribution.DelegatorIncomeHeightState{
		ValAddress:  validator.GetValidatorAddress(),
		DeleAddress: validator.Owner,
		Height:      incomeHeight,
	})

}
