package client

import (
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	gtxs "github.com/QOSGroup/qos/module/guardian/txs"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type CreateGuardiansReq struct {
	rpc.BaseRequest `json:"base"`
	Description     string `json:"description"`
}

type DeleteGuardiansReq struct {
	rpc.BaseRequest `json:"base"`
}

type HaltNetworkReq struct {
	rpc.BaseRequest `json:"base"`
	Reason          string `json:"reason"`
}

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
	registerTxRoutes(ctx, r)
}

func registerQueryRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/gov/guardians", QueryGuardiansHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/gov/guardians/{address}", QueryGuardianHandleFn(ctx)).Methods("GET")
}

func registerTxRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/gov/guardians/{address}/guardians", CreateGuardiansHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/guardians/{address}/deletes", DeleteGuardiansHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/gov/guardians/halt_network", HaltNetWorkHandleFn(ctx)).Methods("POST")
}

func HaltNetWorkHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, e error) {
			hnrReq := req.(HaltNetworkReq)
			itx := gtxs.TxHaltNetwork{
				Guardian: from,
				Reason:   hnrReq.Reason,
			}

			e = itx.ValidateInputs()
			tx = itx

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(HaltNetworkReq{}), fn)
	}
}

func DeleteGuardiansHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, e error) {
			addr, e := types.AccAddressFromBech32(vars["address"])
			if e != nil {
				return nil, e
			}

			itx := gtxs.TxDeleteGuardian{
				Address:   addr,
				DeletedBy: from,
			}

			e = itx.ValidateInputs()
			tx = itx

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(DeleteGuardiansReq{}), fn)
	}

}

func CreateGuardiansHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, e error) {
			cgrReq := req.(CreateGuardiansReq)

			addr, e := types.AccAddressFromBech32(vars["address"])
			if e != nil {
				return nil, e
			}

			itx := gtxs.TxAddGuardian{
				Description: cgrReq.Description,
				Address:     addr,
				Creator:     from,
			}

			e = itx.ValidateInputs()
			tx = itx

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateGuardiansReq{}), fn)
	}
}

func QueryGuardiansHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := queryAllGuardians(ctx)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		if result == nil || len(result) == 0 {
			rpc.Write40XErrorResponse(writer, context.RecordsNotFoundError)
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func QueryGuardianHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		addr, err := types.AccAddressFromBech32(vars["address"])
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		result, err := getGuardian(ctx, addr)
		if err != nil {
			rpc.Write40XErrorResponse(writer, err)
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}
