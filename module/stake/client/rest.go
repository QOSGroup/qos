package client

import (
	"errors"
	"fmt"
	"github.com/QOSGroup/qbase/client/context"
	"github.com/QOSGroup/qbase/client/rpc"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	stxs "github.com/QOSGroup/qos/module/stake/txs"
	"github.com/QOSGroup/qos/module/stake/types"
	qtypes "github.com/QOSGroup/qos/types"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

func RegisterRoutes(ctx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(ctx, r)
	registerTxRoutes(ctx, r)
}

type CreateValidatorReq struct {
	rpc.BaseRequest `json:"base"`
	ConsPubKey      string        `json:"cons_pub_key"`
	BondTokens      btypes.BigInt `json:"bond_tokens"`
	IsCompound      bool          `json:"is_compound"`
	Moniker         string        `json:"moniker"`
	Logo            string        `json:"logo"`
	Website         string        `json:"website"`
	Details         string        `json:"details"`
	Rate            qtypes.Dec    `json:"rate"`
	MaxRate         qtypes.Dec    `json:"max_rate"`
	MaxChangeRate   qtypes.Dec    `json:"max_change_rate"`
}

type ModifyValidatorReq struct {
	rpc.BaseRequest `json:"base"`
	Moniker         string     `json:"moniker"`
	Logo            string     `json:"logo"`
	Website         string     `json:"website"`
	Details         string     `json:"details"`
	CommissionRate  qtypes.Dec `json:"rate"`
}

type RevokeValidatorReq struct {
	rpc.BaseRequest `json:"base"`
}

type ActiveValidatorReq struct {
	rpc.BaseRequest `json:"base"`
	BondTokens      btypes.BigInt `json:"bond_tokens"`
}

type CreateDelegationReq struct {
	rpc.BaseRequest `json:"base"`
	Amount          btypes.BigInt `json:"amount"`
	IsCompound      bool          `json:"is_compound"`
}

type ModifyDelegationReq struct {
	rpc.BaseRequest `json:"base"`
	IsCompound      bool `json:"is_compound"`
}

type UnbondDelegationReq struct {
	rpc.BaseRequest `json:"base"`
	UnbondAmount    btypes.BigInt `json:"unbond_amount"`
	UnbondAll       bool          `json:"unbond_all"`
}

type ReDelegationReq struct {
	rpc.BaseRequest `json:"base"`
	Amount          btypes.BigInt `json:"amount"`
	RedelegateAll   bool          `json:"redelegate_all"`
	Compound        bool          `json:"compound"`
}

func registerQueryRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/stake/validators", QueryAllValidatorsHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/stake/validators/{validatorAddr}", QueryValidatorHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/stake/delegators/{validatorAddr}", QueryDelegationsToHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/stake/delegators/{delegatorAddr}/delegations", QueryDelegationsHandleFn(ctx)).Methods("GET")
	r.HandleFunc("/stake/delegators/{delegatorAddr}/validators/{validatorAddr}", QueryDelegationHandleFn(ctx)).Methods("GET")
}

func registerTxRoutes(ctx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/stake/validators", CreateValidatorHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/stake/validators/{validatorAddr}/modify_validators", ModifyValidatorHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/stake/validators/{validatorAddr}/revoke_validators", RevokeValidatorHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/stake/validators/{validatorAddr}/active_validators", ActiveValidatorHandleFn(ctx)).Methods("POST")

	r.HandleFunc("/stake/delegators/{validatorAddr}/delegations", CreateDelegationHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/stake/delegators/{validatorAddr}/modify_delegations", ModifyDelegationHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/stake/delegators/{validatorAddr}/unbond_delegations", UnbondDelegationHandleFn(ctx)).Methods("POST")
	r.HandleFunc("/stake/delegators/{fromValidatorAddr}/redelegations/{toValidatorAddr}", RedelegationHandleFn(ctx)).Methods("POST")
}

func CreateValidatorHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			cvr := req.(CreateValidatorReq)

			cr := types.NewCommissionRates(cvr.Rate, cvr.MaxRate, cvr.MaxChangeRate)
			erro := cr.Validate()
			if erro != nil {
				return nil, errors.New(erro.Error())
			}

			pk, err := btypes.GetConsensusPubKeyBech32(cvr.ConsPubKey)
			if err != nil {
				return nil, err
			}

			tx := stxs.TxCreateValidator{
				Owner:      from,
				ConsPubKey: pk,
				BondTokens: cvr.BondTokens,
				IsCompound: cvr.IsCompound,
				Description: types.Description{
					Moniker: cvr.Moniker,
					Logo:    cvr.Logo,
					Website: cvr.Website,
					Details: cvr.Details,
				},
				Commission:  cr,
				Delegations: nil,
			}

			return &tx, tx.ValidateInputs()
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateValidatorReq{}), fn)
	}
}

func ModifyValidatorHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			mvr := req.(ModifyValidatorReq)

			tx := stxs.TxModifyValidator{
				Owner:         from,
				ValidatorAddr: val,
				Description: types.Description{
					Moniker: mvr.Moniker,
					Logo:    mvr.Logo,
					Website: mvr.Website,
					Details: mvr.Details,
				},
				CommissionRate: &mvr.CommissionRate,
			}

			return &tx, tx.ValidateInputs()
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(ModifyValidatorReq{}), fn)
	}
}

func MustParseValidatorAddress(w http.ResponseWriter, addrStr string) (btypes.ValAddress, bool) {
	addr, err := btypes.ValAddressFromBech32(addrStr)
	if err != nil {
		rpc.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("addr not valid. addr: %s", addrStr))
		return nil, false
	}

	return addr, true
}

func MustParseAccountAddress(w http.ResponseWriter, addrStr string) (btypes.AccAddress, bool) {
	addr, err := btypes.AccAddressFromBech32(addrStr)
	if err != nil {
		rpc.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("addr not valid. addr: %s", addrStr))
		return nil, false
	}

	return addr, true
}

func RevokeValidatorHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			tx := stxs.NewRevokeValidatorTx(from, val)
			return tx, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(RevokeValidatorReq{}), fn)
	}
}

func ActiveValidatorHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			avrReq := req.(ActiveValidatorReq)
			tx := stxs.NewActiveValidatorTx(from, val, avrReq.BondTokens)
			return tx, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(ActiveValidatorReq{}), fn)
	}
}

func UnbondDelegationHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			udrReq := req.(UnbondDelegationReq)
			tx := stxs.TxUnbondDelegation{
				Delegator:     from,
				ValidatorAddr: val,
				UnbondAmount:  udrReq.UnbondAmount,
				UnbondAll:     udrReq.UnbondAll,
			}

			return &tx, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(UnbondDelegationReq{}), fn)
	}
}

func RedelegationHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		fromVal, ok := MustParseValidatorAddress(writer, vars["fromValidatorAddr"])
		if !ok {
			return
		}

		toVal, ok := MustParseValidatorAddress(writer, vars["toValidatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			rrReq := req.(ReDelegationReq)

			tx := stxs.TxCreateReDelegation{
				Delegator:         from,
				FromValidatorAddr: fromVal,
				ToValidatorAddr:   toVal,
				Amount:            rrReq.Amount,
				RedelegateAll:     rrReq.RedelegateAll,
				Compound:          rrReq.Compound,
			}

			return &tx, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(ReDelegationReq{}), fn)
	}
}

func ModifyDelegationHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			mvrReq := req.(ModifyDelegationReq)

			tx := stxs.TxModifyCompound{
				Delegator:     from,
				ValidatorAddr: val,
				IsCompound:    mvrReq.IsCompound,
			}

			return &tx, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(ModifyDelegationReq{}), fn)

	}
}

func CreateDelegationHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		fn := func(req interface{}, from btypes.AccAddress, _ map[string]string) (txs.ITx, error) {
			avrReq := req.(CreateDelegationReq)
			tx := stxs.TxCreateDelegation{
				Delegator:     from,
				ValidatorAddr: val,
				Amount:        avrReq.Amount,
				IsCompound:    avrReq.IsCompound,
			}

			return &tx, nil
		}

		rpc.BuildStdTxAndResponse(writer, request, cliContext, reflect.TypeOf(CreateDelegationReq{}), fn)
	}
}

func QueryAllValidatorsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		result, err := queryAllValidators(ctx)
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

func QueryValidatorHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		result, err := getValidator(ctx, val)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)
	}
}

func QueryDelegationsToHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		vars := mux.Vars(request)
		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		result, err := queryDelegationsTo(ctx, val)
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

func QueryDelegationsHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		vars := mux.Vars(request)
		deleAddr, ok := MustParseAccountAddress(writer, vars["delegatorAddr"])
		if !ok {
			return
		}

		result, err := queryDelegations(ctx, deleAddr)
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

func QueryDelegationHandleFn(cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		br, _ := rpc.ParseRequestForm(request)
		ctx := br.Setup(cliContext)

		vars := mux.Vars(request)
		deleAddr, ok := MustParseAccountAddress(writer, vars["delegatorAddr"])
		if !ok {
			return
		}

		val, ok := MustParseValidatorAddress(writer, vars["validatorAddr"])
		if !ok {
			return
		}

		result, err := getDelegationInfo(ctx, deleAddr, val)
		if err != nil {
			rpc.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		rpc.PostProcessResponseBare(writer, ctx, result)

	}
}
