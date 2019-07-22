package gov

import (
	"errors"
	"fmt"
	qcliacc "github.com/QOSGroup/qbase/client/account"
	"github.com/QOSGroup/qbase/client/context"
	qcltx "github.com/QOSGroup/qbase/client/tx"
	"github.com/QOSGroup/qbase/txs"
	gtxs "github.com/QOSGroup/qos/module/gov/txs"
	gtypes "github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"strings"
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

				deposit := viper.GetInt64(flagDeposit)
				if deposit <= 0 {
					return nil, errors.New("deposit must be positive")
				}

				switch proposalType {
				case gtypes.ProposalTypeText:
					return gtxs.NewTxProposal(title, description, proposer, uint64(deposit)), nil
				case gtypes.ProposalTypeTaxUsage:
					destAddress, err := qcliacc.GetAddrFromFlag(ctx, flagDestAddress)
					if err != nil {
						return nil, err
					}
					percent := viper.GetFloat64(flagPercent)
					if percent <= 0 {
						return nil, errors.New("deposit must be positive")
					}
					ps, _ := types.NewDecFromStr(fmt.Sprintf("%f", percent))
					return gtxs.NewTxTaxUsage(title, description, proposer, uint64(deposit), destAddress, ps), nil
				case gtypes.ProposalTypeParameterChange:
					params, err := parseParams(viper.GetString(flagParams))
					if err != nil {
						return nil, err
					}
					return gtxs.NewTxParameterChange(title, description, proposer, uint64(deposit), params), nil
				}

				return nil, errors.New("unknown proposal-type")
			})
		},
	}

	cmd.Flags().String(flagTitle, "", "Proposal title")
	cmd.Flags().String(flagDescription, "", "Proposal description")
	cmd.Flags().String(flagProposalType, gtypes.ProposalTypeText.String(), "")
	cmd.Flags().String(flagProposer, "", "Proposer who submit the proposal")
	cmd.Flags().Uint64(flagDeposit, 0, "Initial deposit paid by proposer. Must be strictly positive")
	cmd.Flags().String(flagDestAddress, "", "Address to receive QOS")
	cmd.Flags().Float64(flagPercent, 0, "Percent of QOS in fee pool send to dest-address")
	cmd.Flags().String(flagParams, "", "params, format:<module>/<key>:<value>,<module>/<key>:<value>")
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
