package client

import (
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/types"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/distribution/validators/{validatorAddr}/periods", queryValidatorPeriodsHandleFn(cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/delegators/{delegatorAddr}/validators/{validatorAddr}/incomes", queryDelegatorIncomeHandleFn(cliCtx)).Methods("GET")
	r.HandleFunc("/distribution/communityFeePools", getCommunityFeePoolsHandleFn(cliCtx)).Methods("GET")
}

func queryValidatorPeriodsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)
		valAddr, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		result, err := queryValidatorPeriods(ctx, valAddr)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)
	}
}

func queryDelegatorIncomeHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)
		valAddr, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		deleAddr, ok := MustParseAccountAddress(writer, vars["delegatorAddr"])
		if !ok {
			return
		}

		result, err := queryDelegatorIncomes(ctx, deleAddr, valAddr)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)

	}
}

func getCommunityFeePoolsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := getCommunityFeePool(ctx)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)
	}
}

func MustParseValidatorAddress(w http.ResponseWriter, addrStr string) (types.ValAddress, bool) {
	addr, err := types.ValAddressFromBech32(addrStr)
	if err != nil {
		rpc.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("addr not valid. addr: %s", addrStr))
		return nil, false
	}

	return addr, true
}

func MustParseAccountAddress(w http.ResponseWriter, addrStr string) (types.AccAddress, bool) {
	addr, err := types.AccAddressFromBech32(addrStr)
	if err != nil {
		rpc.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("addr not valid. addr: %s", addrStr))
		return nil, false
	}

	return addr, true
}
