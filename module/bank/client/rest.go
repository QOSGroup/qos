package client

import (
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	qtxs "github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/bank/txs"
	bankTypes "github.com/QOSGroup/qos/module/bank/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type TransferSendReq struct {
	rpc.BaseRequest `json:"base"`
	QOS             types.BigInt `json:"qos"`
	QSCs            qtypes.QSCs  `json:"qscs"`
}

type InvariantCheckReq struct {
	rpc.BaseRequest `json:"base"`
}

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerTxRoutes(ctx, r)
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bank/accounts/{address}/transfers", TransferRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/invariant/checks", InvariantCheckRequestHandleFn(cliCtx)).Methods("POST")
}

func InvariantCheckRequestHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(InvariantCheckReq{}), func(req interface{}, from types.AccAddress, vars map[string]string) (tx qtxs.ITx, e error) {
			tx = txs.TxInvariantCheck{
				Sender: from,
			}
			return
		})
	}
}

func TransferRequestHandlerFn(cliCtx context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rpc.BuildStdTxAndResponse(writer, request, cliCtx, reflect.TypeOf(TransferSendReq{}), func(req interface{}, from types.AccAddress, vars map[string]string) (tx qtxs.ITx, err error) {
			tfsr := req.(TransferSendReq)
			receiver, err := types.AccAddressFromBech32(vars["address"])
			if err != nil {
				return nil, fmt.Errorf("addr: %s is not correct bech32 address", vars["address"])
			}

			itx := txs.TxTransfer{
				Senders: bankTypes.TransItems{
					bankTypes.TransItem{
						Address: from,
						QOS:     tfsr.QOS,
						QSCs:    tfsr.QSCs,
					},
				},
				Receivers: bankTypes.TransItems{
					bankTypes.TransItem{
						Address: receiver,
						QOS:     tfsr.QOS,
						QSCs:    tfsr.QSCs,
					},
				},
			}

			tx = itx
			err = itx.ValidateInputs()
			return
		})
	}
}
