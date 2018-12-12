package validator

import (
	"fmt"

	"github.com/QOSGroup/qbase/context"

	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx context.Context) (res abci.ResponseEndBlock) {

	log := ctx.Logger()

	if ctx.BlockHeight() > 1 {
		validatorMapper := ctx.Mapper(ValidatorMapperName).(*ValidatorMapper)

		if validatorMapper.IsValidatorChanged() {
			log.Debug(fmt.Sprintf("validators changed on height: %d ", ctx.BlockHeight()))
			allValidators := validatorMapper.GetValidators()
			abciValidators := make([]abci.Validator, len(allValidators))

			for i, val := range allValidators {
				abciValidators[i] = val.ToABCIValidator()
			}

			res.ValidatorUpdates = abciValidators
			validatorMapper.SetValidatorUnChanged()
		}
	}

	return
}
