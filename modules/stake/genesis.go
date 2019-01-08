package stake

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qos/modules/stake/mapper"
	staketypes "github.com/QOSGroup/qos/modules/stake/types"
)

type GenesisState struct {
	Params     staketypes.Params      `json:"params"` // inflation params
	Validators []staketypes.Validator `json:"validators"`
}

func NewGenesisState(params staketypes.Params, validators []staketypes.Validator) GenesisState {
	return GenesisState{
		Params:     params,
		Validators: validators,
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: staketypes.DefaultParams(),
	}
}

func InitGenesis(ctx context.Context, data GenesisState) {
	initValidators(ctx, data.Validators)
	initParams(ctx, data.Params)
}

func initValidators(ctx context.Context, validators []staketypes.Validator) {
	mapper := ctx.Mapper(mapper.ValidatorMapperName).(*mapper.ValidatorMapper)
	for _, v := range validators {

		if mapper.Exists(v.ValidatorPubKey.Address().Bytes()) {
			panic(fmt.Errorf("validator %s already exists", v.ValidatorPubKey.Address()))
		}
		if mapper.ExistsWithOwner(v.Owner) {
			panic(fmt.Errorf("owner %s already bind a validator", v.Owner))
		}

		mapper.CreateValidator(v)
	}
}

func initParams(ctx context.Context, params staketypes.Params) {
	mapper := ctx.Mapper(mapper.ValidatorMapperName).(*mapper.ValidatorMapper)
	mapper.SetParams(params)
}

func ValidateGenesis(data GenesisState) error {
	err := validateValidators(data.Validators)
	if err != nil {
		return err
	}

	return nil
}

func validateValidators(validators []staketypes.Validator) (err error) {
	addrMap := make(map[string]bool, len(validators))
	for i := 0; i < len(validators); i++ {
		val := validators[i]
		strKey := string(val.ValidatorPubKey.Bytes())
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, Owner %v", val.Description, val.Owner)
		}
		if val.Status != staketypes.Active {
			return fmt.Errorf("validator is bonded and jailed in genesis state: moniker %v, Owner %v", val.Description, val.Owner)
		}
		addrMap[strKey] = true
	}
	return
}
