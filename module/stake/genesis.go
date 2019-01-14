package stake

import (
	"fmt"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/stake/mapper"
	staketypes "github.com/QOSGroup/qos/module/stake/types"
	"github.com/QOSGroup/qos/types"
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
	validatorMapper := ctx.Mapper(mapper.ValidatorMapperName).(*mapper.ValidatorMapper)
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	for _, v := range validators {

		if validatorMapper.Exists(v.ValidatorPubKey.Address().Bytes()) {
			panic(fmt.Errorf("validator %s already exists", v.ValidatorPubKey.Address()))
		}
		if validatorMapper.ExistsWithOwner(v.Owner) {
			panic(fmt.Errorf("owner %s already bind a validator", v.Owner))
		}

		owner := accountMapper.GetAccount(v.Owner).(*types.QOSAccount)
		owner.MustMinusQOS(btypes.NewInt(int64(v.BondTokens)))
		accountMapper.SetAccount(owner)
		validatorMapper.CreateValidator(v)
	}
}

func initParams(ctx context.Context, params staketypes.Params) {
	mapper := ctx.Mapper(mapper.ValidatorMapperName).(*mapper.ValidatorMapper)
	mapper.SetParams(params)
}

func ValidateGenesis(genesisAccounts []*types.QOSAccount, data GenesisState) error {
	err := validateValidators(genesisAccounts, data.Validators)
	if err != nil {
		return err
	}

	return nil
}

func validateValidators(genesisAccounts []*types.QOSAccount, validators []staketypes.Validator) (err error) {
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

		var ownerExists bool
		for _, acc := range genesisAccounts {
			if acc.AccountAddress.EqualsTo(val.Owner) {
				ownerExists = true

				if uint64(acc.QOS.Int64()) < val.BondTokens {
					return fmt.Errorf("owner of %s no enough QOS", val.Name)
				}
			}
		}

		if !ownerExists {
			return fmt.Errorf("owner of %s not exists", val.Name)
		}
	}
	return nil
}
