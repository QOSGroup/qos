package client

import (
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	qtxs "github.com/QOSGroup/qos/module/qsc/txs"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
	registerTxRoutes(ctx, r)
}

type IssueQscReq struct {
	rpc.BaseRequest `json:"base"`
	Amount          types.BigInt `json:"amount"`
}

func registerQueryRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/qsc/qscs/{qsc}", QueryQscInfoHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/qsc/qscs", QueryAllQscsHandleFn(ctx)).Methods("GET")
}

func registerTxRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/qsc/qscs", CreateQscHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/qsc/qscs/{qsc}/issues", IssueQscHandleFn(ctx)).Methods("POST")
}

func QueryAllQscsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := queryAllQscs(ctx)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		if result == nil || len(result) == 0 {
			rpc.WriteErrorResponse(writer, http.StatusNotFound, "result not found")
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func QueryQscInfoHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		info, err := queryQscInfo(vars["qsc"], ctx)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, info)
	}
}

func CreateQscHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, err := rpc.ParseRequestForm(request)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		ctx := br.Setup(cliContext)

		accounts := request.FormValue("accounts")
		rate := request.FormValue("exchange_rate")
		description := request.FormValue("description")

		caFile, _, err := request.FormFile("ca_file")
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		caBytes, err := ioutil.ReadAll(caFile)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		var crt cert.Certificate
		err = ctx.Codec.UnmarshalJSON(caBytes, &crt)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		subj, ok := crt.CSR.Subj.(cert.QSCSubject)
		if !ok {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid crt file")
			return
		}

		accs, err := parseAccountStr(accounts, subj.Name, types.AccAddressFromBech32)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		from, _ := types.AccAddressFromBech32(br.From)
		tx := qtxs.TxCreateQSC{
			Creator:      from,
			ExchangeRate: strings.TrimSpace(rate),
			QSCCA:        &crt,
			Description:  strings.TrimSpace(description),
			Accounts:     accs,
		}

		err = tx.ValidateInputs()
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		rpc.WriteGenStdTxResponse(writer, ctx, br, tx)
	}
}

func IssueQscHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, e error) {
			qsc := vars["qsc"]
			iqrReq := req.(IssueQscReq)

			return qtxs.TxIssueQSC{
				QSCName: strings.TrimSpace(qsc),
				Amount:  iqrReq.Amount,
				Banker:  from,
			}, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(IssueQscReq{}), fn)
	}
}
