package client

import (
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/qcp/txs"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
	registerTxRoutes(ctx, r)
}

func registerQueryRoutes(ctx context.CLIContext, r *mux.Router) {

}

func registerTxRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/qcp/qcps", InitQCPHandlerFn(ctx)).Methods("POST")
}

func InitQCPHandlerFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		bq, err := rpc.ParseRequestForm(request)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

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

		var crt = cert.Certificate{}
		err = cliContext.Codec.UnmarshalJSON(caBytes, &crt)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("crt file unmarshal err. err: %s", err.Error()))
			return
		}

		_, ok := crt.CSR.Subj.(cert.QCPSubject)
		if !ok {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, "invalid crt file")
			return
		}

		fromAddr, _ := types.AccAddressFromBech32(bq.From)
		tx := txs.TxInitQCP{
			Creator: fromAddr,
			QCPCA:   &crt,
		}

		err = tx.ValidateInputs()
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		rpc.WriteGenStdTxResponse(writer, cliContext, bq, tx)
	}
}
