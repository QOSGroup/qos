package client

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/gov/mapper"
	txs2 "github.com/QOSGroup/qos/module/gov/txs"
	"github.com/QOSGroup/qos/module/gov/types"
	"github.com/QOSGroup/qos/module/mint"
	types2 "github.com/QOSGroup/qos/types"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	moduleData = []string{"stake", "distribution", "gov"}
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
	registerTxRoutes(ctx, r)
}

type CreateTextProposalReq struct {
	rpc.BaseRequest `json:"base"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Deposit         btypes.BigInt `json:"deposit"`
}

type CreateTaxUsageProposalReq struct {
	rpc.BaseRequest `json:"base"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Deposit         btypes.BigInt `json:"deposit"`
	DestAddress     string        `json:"dest_address"`
	Percent         string        `json:"percent"`
}

type CreateParamChangeProposalReq struct {
	rpc.BaseRequest `json:"base"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Deposit         btypes.BigInt `json:"deposit"`
	Params          string        `json:"params"`
}

type CreateModifyInflationProposalReq struct {
	rpc.BaseRequest  `json:"base"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Deposit          btypes.BigInt `json:"deposit"`
	InflationPhrases string        `json:"inflation_phrases"`
	TotalAmount      btypes.BigInt `json:"total_amount"`
}

type CreateSoftwareUpgradeProposalReq struct {
	rpc.BaseRequest `json:"base"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Deposit         btypes.BigInt `json:"deposit"`
	Version         string        `json:"version"`
	DataHeight      int64         `json:"data_height"`
	GenesisFile     string        `json:"genesis_file"`
	GenesisMD5      string        `json:"genesis_md5"`
	ForZeroHeight   bool          `json:"for_zero_height"`
}

type DepositProposalReq struct {
	rpc.BaseRequest `json:"base"`
	Amount          btypes.BigInt `json:"amount"`
}

type VoteProposalReq struct {
	rpc.BaseRequest `json:"base"`
	Option          string `json:"option"`
}

func registerQueryRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/gov/proposals", QueryAllProposalsHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposalId}", GetProposalHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposalId}/deposits", QueryAllDepositsHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposalId}/deposits/{address}", GetDepositHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposalId}/votes", QueryAllVotesHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposalId}/votes/{address}", GetVoteHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/proposal/{proposalId}/tallies", GetProposalTallyHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/params/modules/{module}/keys/{key}", GetParamsKeyHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/params/modules/{module}/keys", GetModuleKeyHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/params/modules", GetModulesHandleFn(ctx)).Methods("GET")
}

func registerTxRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/gov/proposals/texts", CreateTextProposalHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/proposals/taxes", CreateTaxUsageHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/proposals/params", CreateParamsChangeHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/proposals/inflations", CreateModifyInflationHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/proposals/upgrades", CreateSoftwareUpgradeHandleFn(ctx)).Methods("POST")

	r.HandleFunc("/gov/proposals/{proposalId}/deposits", DepositProposalHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/proposals/{proposalId}/votes", VoteProposalHandleFn(ctx)).Methods("POST")
}

func CreateTextProposalHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (tx txs.ITx, err error) {
			ctprReq := req.(CreateTextProposalReq)
			itx := txs2.NewTxProposal(ctprReq.Title, ctprReq.Description, from, ctprReq.Deposit)

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateTextProposalReq{}), fn)
	}
}

func CreateTaxUsageHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (tx txs.ITx, err error) {
			ctuprReq := req.(CreateTaxUsageProposalReq)

			destAddr, err := btypes.AccAddressFromBech32(ctuprReq.DestAddress)
			if err != nil {
				return nil, errors.New("invalid dest address")
			}

			percent, err := types2.NewDecFromStr(ctuprReq.Percent)
			if err != nil {
				return nil, err
			}

			itx := txs2.NewTxTaxUsage(ctuprReq.Title, ctuprReq.Description, from, ctuprReq.Deposit, destAddr, percent)

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateTaxUsageProposalReq{}), fn)
	}
}

func CreateParamsChangeHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (tx txs.ITx, err error) {
			cpcprReq := req.(CreateParamChangeProposalReq)

			params, err := parseParams(cpcprReq.Params)
			if err != nil {
				return nil, err
			}

			itx := txs2.NewTxParameterChange(cpcprReq.Title, cpcprReq.Description, from, cpcprReq.Deposit, params)

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateParamChangeProposalReq{}), fn)
	}
}

func CreateModifyInflationHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (tx txs.ITx, err error) {
			cmiptReq := req.(CreateModifyInflationProposalReq)

			inflationStr := cmiptReq.InflationPhrases
			if strings.TrimSpace(inflationStr) == "" {
				return nil, errors.New("empty inflation phrases")
			}

			var inflationPhrases mint.InflationPhrases
			err = cliContext.Codec.UnmarshalJSON([]byte(inflationStr), &inflationPhrases)
			if err != nil {
				return nil, fmt.Errorf("invalid inflation phrases. err: %s", err.Error())
			}

			itx := txs2.NewTxModifyInflation(cmiptReq.Title, cmiptReq.Description, from, cmiptReq.Deposit, cmiptReq.TotalAmount, inflationPhrases)

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateModifyInflationProposalReq{}), fn)
	}
}

func CreateSoftwareUpgradeHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (tx txs.ITx, err error) {
			r := req.(CreateSoftwareUpgradeProposalReq)
			itx := txs2.NewTxSoftwareUpgrade(r.Title, r.Description, from, r.Deposit, r.Version, r.DataHeight, r.GenesisFile, r.GenesisMD5, r.ForZeroHeight)

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateSoftwareUpgradeProposalReq{}), fn)
	}
}

func VoteProposalHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			vprReq := req.(VoteProposalReq)
			pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
			if err != nil {
				return nil, errors.New("invalid proposalId")
			}

			vp, err := types.VoteOptionFromString(vprReq.Option)
			if err != nil {
				return nil, err
			}
			itx := txs2.NewTxVote(pid, from, vp)

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(VoteProposalReq{}), fn)
	}
}

func DepositProposalHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from btypes.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			dprReq := req.(DepositProposalReq)
			pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
			if err != nil {
				return nil, errors.New("invalid proposalId")
			}

			tx = txs2.NewTxDeposit(pid, from, dprReq.Amount)

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(DepositProposalReq{}), fn)
	}
}

func GetParamsKeyHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		result, err := queryModuleParams(ctx, vars["module"], vars["key"])
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func GetModuleKeyHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		result, err := queryModuleParams(ctx, vars["module"], "")
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func GetModulesHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		rpc.PostProcessResponseBare(writer, cliContext, moduleData)
	}
}

func GetVoteHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid proposal id")
			return
		}

		addr, err := btypes.AccAddressFromBech32(vars["address"])
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid address")
			return
		}

		result, err := getProposalVote(ctx, pid, addr)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}
		rpc.PostProcessResponseBare(writer, cliContext, result)

	}
}

func QueryAllVotesHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid proposal id")
			return
		}

		result, err := queryProposalVotes(ctx, pid)

		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		if len(result) == 0 {
			rpc.Write40XErrorResponse(writer, context.RecordsNotFoundError)
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)

	}
}

func GetDepositHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid proposal id")
			return
		}

		addr, err := btypes.AccAddressFromBech32(vars["address"])
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid address")
			return
		}

		result, err := getProposalDeposit(ctx, pid, addr)

		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}
		rpc.PostProcessResponseBare(writer, cliContext, result)
	}
}

func QueryAllDepositsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid proposal id")
			return
		}

		result, err := queryProposalVotes(ctx, pid)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		if len(result) == 0 {
			rpc.Write40XErrorResponse(writer, context.RecordsNotFoundError)
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)

	}
}

func GetProposalHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid proposal id")
			return
		}

		result, err := getProposal(ctx, pid)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func GetProposalTallyHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		pid, err := strconv.ParseInt(vars["proposalId"], 10, 64)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid proposal id")
			return
		}

		result, err := getProposalTally(ctx, pid)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func QueryAllProposalsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		limit := int64(10)
		if l, err := strconv.ParseInt(request.FormValue("limit"), 10, 64); err == nil {
			limit = l
		}
		var depositorAddr btypes.AccAddress
		var voterAddr btypes.AccAddress
		var status types.ProposalStatus
		if d, err := btypes.AccAddressFromBech32(request.FormValue("depositor")); err != nil {
			depositorAddr = d
		}
		if d, err := btypes.AccAddressFromBech32(request.FormValue("voter")); err != nil {
			voterAddr = d
		}

		status = toProposalStatus(request.FormValue("status"))

		queryParam := mapper.QueryProposalsParam{
			Depositor: depositorAddr,
			Voter:     voterAddr,
			Status:    status,
			Limit:     limit,
		}

		result, err := queryProposalsByParams(ctx, queryParam)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		if len(result) == 0 {
			rpc.Write40XErrorResponse(writer, context.RecordsNotFoundError)
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)
	}
}
