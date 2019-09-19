package client

import (
	"errors"
	"strings"

	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	gtxs "github.com/QOSGroup/qos/module/gov/txs"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/mint"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
)

const (
	layoutISO = "2006-01-02"
)

func ProposalCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return qcltx.BroadcastTxAndPrintResult(cdc, func(ctx context.CLIContext) (txs.ITx, error) {
				title := viper.GetString(flagTitle)
				description := viper.GetString(flagDescription)

				proposalType, err := gtypes.ProposalTypeFromString(viper.GetString(flagProposalType))
				if err != nil {
					return nil, err
				}

				proposer, err := qcliacc.GetAddrFromFlag(ctx, flagProposer)
				if err != nil {
					return nil, err
				}

				deposit, err := types.GetIntFromFlag(flagDeposit, false)
				if err != nil {
					return nil, err
				}

				switch proposalType {
				case gtypes.ProposalTypeText:
					return gtxs.NewTxProposal(title, description, proposer, deposit), nil
				case gtypes.ProposalTypeTaxUsage:
					destAddress, err := qcliacc.GetAddrFromFlag(ctx, flagDestAddress)
					if err != nil {
						return nil, err
					}
					percent := strings.TrimSpace(viper.GetString(flagPercent))
					if len(percent) == 0 {
						return nil, errors.New("empty percent")
					}
					ps, err := types.NewDecFromStr(percent)
					if err != nil {
						return nil, err
					}
					if ps.GT(types.OneDec()) || ps.LTE(types.ZeroDec()) {
						return nil, errors.New("percent ranges (0, 1]")
					}
					return gtxs.NewTxTaxUsage(title, description, proposer, deposit, destAddress, ps), nil
				case gtypes.ProposalTypeParameterChange:
					params, err := parseParams(viper.GetString(flagParams))
					if err != nil {
						return nil, err
					}
					return gtxs.NewTxParameterChange(title, description, proposer, deposit, params), nil
				case gtypes.ProposalTypeModifyInflation:
					inflationPhrasesStr := viper.GetString(flagInflationPhrases)
					if len(inflationPhrasesStr) == 0 {
						return nil, errors.New("inflation-phrases incorrect")
					}

					var inflationPhrases mint.InflationPhrases
					err := cdc.UnmarshalJSON([]byte(inflationPhrasesStr), &inflationPhrases)
					if err != nil {
						return nil, err
					}
					totalAmount, err := types.GetIntFromFlag(flagTotalAmount, false)
					if err != nil {
						return nil, err
					}
					return gtxs.NewTxModifyInflation(title, description, proposer, deposit, totalAmount, inflationPhrases), nil
				case gtypes.ProposalTypeSoftwareUpgrade:
					version := viper.GetString(flagVersion)
					if len(version) == 0 {
						return nil, errors.New("version is empty")
					}
					forZeroHeight := viper.GetBool(flagForZeroHeight)
					dataHeight := viper.GetInt64(flagDataHeight)
					genesisFile := viper.GetString(flagGenesisFile)
					genesisMD5 := viper.GetString(flagGenesisMD5)
					if forZeroHeight {
						if dataHeight <= 0 {
							return nil, errors.New("data-height must be positive")
						}
						if len(genesisFile) == 0 {
							return nil, errors.New("genesis-file is empty")
						}
						if len(genesisMD5) == 0 {
							return nil, errors.New("genesis-md5 is empty")
						}
					}
					return gtxs.NewTxSoftwareUpgrade(title, description, proposer, deposit,
						version, dataHeight, genesisFile, genesisMD5, forZeroHeight), nil
				}

				return nil, errors.New("unknown proposal-type")
			})
		},
	}

	cmd.Flags().String(flagTitle, "", "Proposal title")
	cmd.Flags().String(flagDescription, "", "Proposal description")
	cmd.Flags().String(flagProposalType, gtypes.ProposalTypeText.String(), "")
	cmd.Flags().String(flagProposer, "", "Proposer who submit the proposal")
	cmd.Flags().String(flagDeposit, "0", "Initial deposit paid by proposer. Must be strictly positive")
	cmd.Flags().String(flagDestAddress, "", "Address to receive QOS, for TaxUsage proposal")
	cmd.Flags().String(flagPercent, "0", "Percent of QOS in fee pool send to dest-address, for TaxUsage proposal")
	cmd.Flags().String(flagParams, "", "params, format:<module>/<key>:<value>,<module>/<key>:<value>, for ParameterChange proposal")
	cmd.Flags().String(flagInflationPhrases, "", "Inflation phrases, json marshaled")
	cmd.Flags().String(flagTotalAmount, "0", "Total QOS amount")
	cmd.Flags().String(flagVersion, "", "qosd version, for software upgrade")
	cmd.Flags().Uint64(flagDataHeight, 0, "data version, for software upgrade")
	cmd.Flags().String(flagGenesisFile, "", "url of genesis file, for software upgrade")
	cmd.Flags().String(flagGenesisMD5, "", "signature of genesis.json, for software upgrade")
	cmd.Flags().Bool(flagForZeroHeight, false, "restart from zero height")

	cmd.MarkFlagRequired(flagTitle)
	cmd.MarkFlagRequired(flagDescription)
	cmd.MarkFlagRequired(flagProposalType)
	cmd.MarkFlagRequired(flagProposer)
	cmd.MarkFlagRequired(flagDeposit)

	return cmd
}

func parseParams(paramsStr string) (params []gtypes.Param, err error) {
	if len(paramsStr) == 0 {
		return nil, errors.New("params is empty")
	}
	items := strings.Split(paramsStr, ",")
	for _, item := range items {
		param := strings.Split(item, ":")
		params = append(params, gtypes.NewParam(strings.TrimSpace(param[0]), strings.TrimSpace(param[1]), strings.TrimSpace(param[2])))
	}

	return params, err
}
