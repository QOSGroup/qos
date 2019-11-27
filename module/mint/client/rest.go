package client

import (
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/mint/inflationPhrases", queryInflationPhrasesHandleFn(cliCtx)).Methods("GET")
	r.HandleFunc("/mint/inflations", queryTotalInflationsHandleFn(cliCtx)).Methods("GET")
	r.HandleFunc("/mint/applies", queryTotalAppliesHandleFn(cliCtx)).Methods("GET")
}

func queryInflationPhrasesHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := queryInflationPhrases(ctx)

		if len(result) == 0 {
			rpc.Write40XErrorResponse(writer, context.RecordsNotFoundError)
			return
		}

		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)
	}
}

func queryTotalInflationsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := queryTotal(ctx)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)

	}
}

func queryTotalAppliesHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := queryApplied(ctx)

		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, cliContext, result)
	}
}
