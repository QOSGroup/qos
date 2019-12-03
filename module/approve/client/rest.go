package client

import (
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	atxs "github.com/QOSGroup/qos/module/approve/txs"
	atypes "github.com/QOSGroup/qos/module/approve/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

type CreateApproveReq struct {
	rpc.BaseRequest `json:"base"`
	QOS             types.BigInt `json:"qos"`
	QSCs            qtypes.QSCs  `json:"qscs"`
}

type IncreaseApproveReq struct {
	rpc.BaseRequest `json:"base"`
	QOS             types.BigInt `json:"qos"`
	QSCs            qtypes.QSCs  `json:"qscs"`
}

type DecreaseApproveReq struct {
	rpc.BaseRequest `json:"base"`
	QOS             types.BigInt `json:"qos"`
	QSCs            qtypes.QSCs  `json:"qscs"`
}

type UseApproveReq struct {
	rpc.BaseRequest `json:"base"`
	QOS             types.BigInt `json:"qos"`
	QSCs            qtypes.QSCs  `json:"qscs"`
}

type CancelApproveReq struct {
	rpc.BaseRequest `json:"base"`
}

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
	registerTxRoutes(ctx, r)
}

func registerQueryRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/approve/approves/{approveAddr}/approve/{beneficiaryAddr}", QueryApproveHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/approve/approves/{approveAddr}/approves", QueryUserApprovesHandleFn(ctx)).Methods("GET")
}

func registerTxRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/approve/approves/{address}/create_approves", CreateApproveHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/approve/approves/{address}/increase_approves", IncreaseApproveHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/approve/approves/{address}/decrease_approves", DecreaseApproveHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/approve/approves/{address}/use_approves", UseApproveHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/approve/approves/{address}/cancel_approves", CancelApproveHandleFn(ctx)).Methods("POST")
}

func QueryUserApprovesHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request) //ignore err
		ctx := br.Setup(cliContext)
		vars := mux.Vars(request)

		approve, err := types.AccAddressFromBech32(vars["approveAddr"])
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("addr: %s is not correct bech32 address", vars["approveAddr"]))
			return
		}

		data, err := queryUserApproves(ctx, approve)

		if len(data) == 0 {
			rpc.WriteErrorResponse(writer, http.StatusNotFound, "records not found")
			return
		}

		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, data)
	}
}

func QueryApproveHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		approve, err := types.AccAddressFromBech32(vars["approveAddr"])
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("addr: %s is not correct bech32 address", vars["approveAddr"]))
			return
		}

		beneficiary, err := types.AccAddressFromBech32(vars["beneficiaryAddr"])
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("addr: %s is not correct bech32 address", vars["beneficiaryAddr"]))
			return
		}

		br, _ := rpc.ParseRequestForm(request) //ignore err
		ctx := br.Setup(cliContext)
		appr, err := getApproveInfo(ctx, approve, beneficiary)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, appr)
	}
}

func CreateApproveHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateApproveReq{}), func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			carReq := req.(CreateApproveReq)
			beneficiary, err := types.AccAddressFromBech32(vars["address"])
			if err != nil {
				return nil, fmt.Errorf("addr: %s is not correct bech32 address", vars["address"])
			}

			itx := atxs.TxCreateApprove{
				Approve: atypes.Approve{
					From: from,
					To:   beneficiary,
					QOS:  carReq.QOS,
					QSCs: carReq.QSCs,
				},
			}

			tx = itx
			err = itx.Valid()

			return
		})
	}
}

func CancelApproveHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			beneficiary, err := types.AccAddressFromBech32(vars["address"])
			if err != nil {
				return nil, fmt.Errorf("addr: %s is not correct bech32 address", vars["address"])
			}

			itx := atxs.TxCancelApprove{
				From: from,
				To:   beneficiary,
			}

			tx = itx
			err = itx.ValidateInputs()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CancelApproveReq{}), fn)
	}
}

func UseApproveHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			uarReq := req.(UseApproveReq)
			approver, err := types.AccAddressFromBech32(vars["address"])
			if err != nil {
				return nil, fmt.Errorf("addr: %s is not correct bech32 address", vars["address"])
			}

			itx := atxs.TxUseApprove{
				Approve: atypes.Approve{
					From: from,
					To:   approver,
					QOS:  uarReq.QOS,
					QSCs: uarReq.QSCs,
				},
			}

			tx = itx
			err = itx.Valid()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(UseApproveReq{}), fn)
	}
}

func DecreaseApproveHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			darReq := req.(DecreaseApproveReq)
			beneficiary, err := types.AccAddressFromBech32(vars["address"])
			if err != nil {
				return nil, fmt.Errorf("addr: %s is not correct bech32 address", vars["address"])
			}

			itx := atxs.TxDecreaseApprove{
				Approve: atypes.Approve{
					From: from,
					To:   beneficiary,
					QOS:  darReq.QOS,
					QSCs: darReq.QSCs,
				},
			}

			tx = itx
			err = itx.Valid()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(DecreaseApproveReq{}), fn)
	}
}

func IncreaseApproveHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fn := func(req interface{}, from types.AccAddress, vars map[string]string) (tx txs.ITx, err error) {
			iarReq := req.(IncreaseApproveReq)
			beneficiary, err := types.AccAddressFromBech32(vars["address"])
			if err != nil {
				return nil, fmt.Errorf("addr: %s is not correct bech32 address", vars["address"])
			}

			itx := atxs.TxIncreaseApprove{
				Approve: atypes.Approve{
					From: from,
					To:   beneficiary,
					QOS:  iarReq.QOS,
					QSCs: iarReq.QSCs,
				},
			}

			tx = itx
			err = itx.Valid()

			return
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(IncreaseApproveReq{}), fn)
	}
}
