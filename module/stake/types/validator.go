package types

import (
	"time"

	"github.com/QOSGroup/qos/codec"
	"github.com/tendermint/tendermint/crypto"

	btypes "github.com/QOSGroup/qbase/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type InactiveCode int8

const (
	//Active 可获得挖矿奖励状态
	Active int8 = iota

	//Inactive
	Inactive

	//Inactive Code
	Revoke        InactiveCode = iota // 2
	MissVoteBlock                     // 3
	MaxValidator                      // 4
	DoubleSign                        // 5
)

// Description - description fields for a validator
type Description struct {
	Moniker string `json:"moniker"` // name
	Logo    string `json:"logo"`    // optional logo link
	Website string `json:"website"` // optional website link
	Details string `json:"details"` // optional details
}

type Validator struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	Creator         btypes.AccAddress `json:"creator"`
	ConsPubKey      crypto.PubKey     `json:"consensus_pubkey"`
	BondTokens      uint64            `json:"bond_tokens"` //不能超过int64最大值
	Description     Description       `json:"description"`
	Commission      Commission        `json:"commission"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactive_code"`
	InactiveTime   time.Time    `json:"inactive_time"`
	InactiveHeight uint64       `json:"inactive_height"`

	MinPeriod  uint64 `json:"min_period"`
	BondHeight uint64 `json:"bond_height"`
}

type jsonifyValidator struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	Creator         btypes.AccAddress `json:"creator"`
	ConsPubKey      string            `json:"consensus_pubkey"`
	BondTokens      uint64            `json:"bond_tokens"`
	Description     Description       `json:"description"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactive_code"`
	InactiveTime   time.Time    `json:"inactive_time"`
	InactiveHeight uint64       `json:"inactive_height"`

	MinPeriod  uint64 `json:"min_period"`
	BondHeight uint64 `json:"bond_height"`
}

func (val Validator) ConsAddress() btypes.ConsAddress {
	return btypes.ConsAddress(val.ConsPubKey.Address())
}

func (val Validator) ConsensusPower() int64 {
	return int64(val.BondTokens)
}

func (val Validator) ToABCIValidator() (abciVal abci.Validator) {
	abciVal.Power = val.ConsensusPower()
	abciVal.Address = val.ConsAddress()
	return
}

func (val Validator) ToABCIValidatorUpdate(isRemoved bool) (abciVal abci.ValidatorUpdate) {
	abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ConsPubKey)
	if isRemoved {
		abciVal.Power = int64(0)
	} else {
		abciVal.Power = val.ConsensusPower()
	}
	return
}

func (val Validator) IsActive() bool {
	return val.Status == Active
}

func (val Validator) GetBondTokens() uint64 {
	return val.BondTokens
}

func (val Validator) GetConsensusPubKey() crypto.PubKey {
	return val.ConsPubKey
}

func (val Validator) GetCreator() btypes.AccAddress {
	return val.Creator
}

func (val Validator) MarshalJSON() ([]byte, error) {
	bechPubKey, err := btypes.ConsensusPubKeyString(val.ConsPubKey)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSON(jsonifyValidator{
		OperatorAddress: val.OperatorAddress,
		Creator:         val.Creator,
		ConsPubKey:      bechPubKey,
		BondTokens:      val.BondTokens,
		Description:     val.Description,

		Status:         val.Status,
		InactiveCode:   val.InactiveCode,
		InactiveTime:   val.InactiveTime,
		InactiveHeight: val.InactiveHeight,

		MinPeriod:  val.MinPeriod,
		BondHeight: val.BondHeight,
	})
}

func (val *Validator) UnmarshalJSON(data []byte) error {

	jv := &jsonifyValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, jv); err != nil {
		return err
	}

	consPubKey, err := btypes.GetConsensusPubKeyBech32(jv.ConsPubKey)
	if err != nil {
		return err
	}

	*val = Validator{
		OperatorAddress: jv.OperatorAddress,
		Creator:         jv.Creator,
		ConsPubKey:      consPubKey,
		BondTokens:      jv.BondTokens,
		Description:     jv.Description,

		Status:         jv.Status,
		InactiveCode:   jv.InactiveCode,
		InactiveTime:   jv.InactiveTime,
		InactiveHeight: jv.InactiveHeight,

		MinPeriod:  jv.MinPeriod,
		BondHeight: jv.BondHeight,
	}

	return nil
}

func (val Validator) GetValidatorAddress() btypes.ValAddress {
	return val.OperatorAddress
}
