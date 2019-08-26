package types

import (
	"fmt"
	"strings"

	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/client/context"
	clikeys "github.com/QOSGroup/qbase/client/keys"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/pkg/errors"
)

type QOSAccount struct {
	account.BaseAccount `json:"base_account"` // inherits BaseAccount
	QOS                 btypes.BigInt         `json:"qos"`  // coins in public chain
	QSCs                QSCs                  `json:"qscs"` // varied QSCs
}

var _ account.Account = (*QOSAccount)(nil)

func ProtoQOSAccount() account.Account {
	return NewQOSAccountZero()
}

func NewQOSAccountZero() *QOSAccount {
	return &QOSAccount{QOS: btypes.ZeroInt()}
}

func NewQOSAccountWithAddress(addr btypes.Address) *QOSAccount {
	return &QOSAccount{
		BaseAccount: account.BaseAccount{
			AccountAddress: addr,
		},
		QOS: btypes.ZeroInt(),
	}
}

func NewQOSAccount(addr btypes.Address, qos btypes.BigInt, qscs QSCs) *QOSAccount {
	return &QOSAccount{
		BaseAccount: account.BaseAccount{
			AccountAddress: addr,
		},
		QOS:  qos,
		QSCs: qscs,
	}
}

func (account *QOSAccount) GetQOS() btypes.BigInt {
	return account.QOS.NilToZero()
}

// 设置QOS，币值必须大于等于0
func (account *QOSAccount) SetQOS(qos btypes.BigInt) error {
	if qos.LT(btypes.ZeroInt()) {
		return errors.New("qos must gte zero")
	}

	account.QOS = qos

	return nil
}

// 是否有足够QOS
func (account *QOSAccount) EnoughOfQOS(qos btypes.BigInt) bool {
	if account.QOS.LT(qos) {
		return false
	}

	return true
}

// 增加QOS，增加量小于0时返回错误
func (account *QOSAccount) PlusQOS(qos btypes.BigInt) error {
	if qos.LT(btypes.ZeroInt()) {
		return errors.New("qos must gte zero")
	}

	account.QOS = account.QOS.NilToZero().Add(qos)

	return nil
}

// 增加QOS，返回错误时panic
func (account *QOSAccount) MustPlusQOS(qos btypes.BigInt) {
	if err := account.PlusQOS(qos); err != nil {
		panic(err)
	}
}

// 减少QOS，减少量小于0或结果小于0时返回错误
func (account *QOSAccount) MinusQOS(qos btypes.BigInt) error {
	if qos.LT(btypes.ZeroInt()) {
		return errors.New("qos must gte zero")
	}

	qos = account.QOS.NilToZero().Sub(qos)
	if qos.LT(btypes.ZeroInt()) {
		return errors.New("qos must gte zero")
	}

	account.QOS = qos

	return nil
}

// 减少QOS，返回错误panic
func (account *QOSAccount) MustMinusQOS(qos btypes.BigInt) {
	if err := account.MinusQOS(qos); err != nil {
		panic(err)
	}
}

// 返回指定币种币值
func (account *QOSAccount) GetQSC(qscName string) (qsc QSC, exists bool) {
	for _, q := range account.QSCs {
		if q.GetName() == qscName {
			return *q, true
		}
	}
	return QSC{}, false
}

// 设置币种币值
func (account *QOSAccount) SetQSC(qsc QSC) error {
	if qsc.Amount.LT(btypes.ZeroInt()) {
		return errors.New("amount must gte zero")
	}

	for _, q := range account.QSCs {
		if q.GetName() == qsc.GetName() {
			q.SetAmount(qsc.GetAmount())
			return nil
		}
	}

	account.QSCs = append(account.QSCs, &qsc)
	return nil
}

// 是否有足够QSC
func (account *QOSAccount) EnoughOfQSC(qsc QSC) bool {
	for _, q := range account.QSCs {
		if q.Name == qsc.Name && !q.Amount.LT(qsc.Amount) {
			return true
		}
	}

	if qsc.Amount.Equal(btypes.ZeroInt()) {
		return true
	}

	return false
}

// 增加QSC
func (account *QOSAccount) PlusQSC(qsc QSC) error {
	if qsc.Amount.LT(btypes.ZeroInt()) {
		return errors.New("amount must gte zero")
	}

	for _, q := range account.QSCs {
		if q.GetName() == qsc.GetName() {
			amount := q.GetAmount().Add(qsc.GetAmount())
			if amount.LT(btypes.ZeroInt()) {
				return errors.New("result must gte zero")
			}

			q.Amount = amount
			return nil
		}
	}

	account.QSCs = append(account.QSCs, &qsc)

	return nil
}

func (account *QOSAccount) MustPlusQSC(qsc QSC) {
	if err := account.PlusQSC(qsc); err != nil {
		panic(err)
	}
}

// 减少QSC
func (account *QOSAccount) MinusQSC(qsc QSC) error {
	if qsc.Amount.LT(btypes.ZeroInt()) {
		return errors.New("amount must gte zero")
	}

	for _, q := range account.QSCs {
		if q.GetName() == qsc.GetName() {
			amount := q.GetAmount().Sub(qsc.GetAmount())
			if amount.LT(btypes.ZeroInt()) {
				return errors.New("result must gte zero")
			}

			q.Amount = amount
			return nil
		}
	}

	return fmt.Errorf("no %s in account", qsc.Name)
}

func (account *QOSAccount) MustMinusQSC(qsc QSC) {
	if err := account.MinusQSC(qsc); err != nil {
		panic(err)
	}
}

func (account *QOSAccount) GetQSCs() QSCs {
	return account.QSCs
}

// 是否有足够QSCs
func (account *QOSAccount) EnoughOfQSCs(qscs QSCs) bool {
	for _, qsc := range qscs {
		if !account.EnoughOfQSC(*qsc) {
			return false
		}
	}

	return true
}

// 增加QSCs
func (account *QOSAccount) PlusQSCs(qscs QSCs) error {
	if qscs.IsZero() {
		return nil
	}
	qs := qscs.Plus(account.QSCs)
	for _, qsc := range qs {
		if qsc.Amount.LT(btypes.ZeroInt()) {
			return errors.New("qsc in result must gte zero")
		}
	}

	account.QSCs = qs

	return nil
}

func (account *QOSAccount) MustPlusQSCs(qscs QSCs) {
	if err := account.PlusQSCs(qscs); err != nil {
		panic(err)
	}
}

// 减少QSCs
func (account *QOSAccount) MinusQSCs(qscs QSCs) error {
	if qscs.IsZero() {
		return nil
	}
	qs := account.QSCs.Minus(qscs)
	for _, qsc := range qs {
		if qsc.Amount.LT(btypes.ZeroInt()) {
			return errors.New("qsc in result must gte zero")
		}
	}

	account.QSCs = qs

	return nil
}

func (account *QOSAccount) MustMinusQSCs(qscs QSCs) {
	if err := account.MinusQSCs(qscs); err != nil {
		panic(err)
	}
}

// 增加QOS，QSCs
func (account *QOSAccount) Plus(qos btypes.BigInt, qscs QSCs) error {
	qos = qos.NilToZero()
	if qos.LT(btypes.ZeroInt()) {
		return errors.New("qos must gte zero")
	}
	qosTotal := account.QOS.NilToZero().Add(qos)

	for _, qsc := range qscs {
		if qsc.Amount.LT(btypes.ZeroInt()) {
			return fmt.Errorf("%s must gte zero", qsc.Name)
		}
	}

	qscsTotal := qscs.Plus(account.QSCs)
	for _, qsc := range qscsTotal {
		if qsc.Amount.LT(btypes.ZeroInt()) {
			return fmt.Errorf("%s in result must gte zero", qsc.Name)
		}
	}

	account.QOS = qosTotal
	account.QSCs = qscsTotal

	return nil
}

func (account *QOSAccount) MustPlus(qos btypes.BigInt, qscs QSCs) {
	if err := account.Plus(qos, qscs); err != nil {
		panic(err)
	}
}

// 减少QOS，QSCs
func (account *QOSAccount) Minus(qos btypes.BigInt, qscs QSCs) error {
	qos = qos.NilToZero()
	if qos.LT(btypes.ZeroInt()) {
		return errors.New("qos must gte zero")
	}
	qosTotal := account.QOS.NilToZero().Sub(qos)
	if qosTotal.LT(btypes.ZeroInt()) {
		return errors.New("qos in result must gte zero")
	}

	for _, qsc := range qscs {
		if qsc.Amount.LT(btypes.ZeroInt()) {
			return fmt.Errorf("%s must gte zero", qsc.Name)
		}
	}

	qscsTotal := account.QSCs.Minus(qscs)
	for _, qsc := range qscsTotal {
		if qsc.Amount.LT(btypes.ZeroInt()) {
			return fmt.Errorf("%s in result must gte zero", qsc.Name)
		}
	}

	account.QOS = qosTotal
	account.QSCs = qscsTotal

	return nil
}

func (account *QOSAccount) MustMinus(qos btypes.BigInt, qscs QSCs) {
	if err := account.Minus(qos, qscs); err != nil {
		panic(err)
	}
}

func (account *QOSAccount) EnoughOf(qos btypes.BigInt, qscs QSCs) bool {

	return account.EnoughOfQOS(qos) && account.EnoughOfQSCs(qscs)
}

// 移除QSC
func (account *QOSAccount) RemoveQSC(qscName string) {
	for i, qsc := range account.QSCs {
		if qsc.GetName() == qscName {
			account.QSCs = append(account.QSCs[:i], account.QSCs[i+1:]...)
		}
	}
}

// Parse accounts from string
// address16lwp3kykkjdc2gdknpjy6u9uhfpa9q4vj78ytd,1000000qos,1000000qstars. Multiple accounts separated by ';'
func ParseAccounts(str, clientHome string) ([]*QOSAccount, error) {
	var keybaseAvailable bool
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	keybase, err := clikeys.GetKeyBaseFromDir(cliCtx, clientHome)
	if err == nil {
		keybaseAvailable = true
	}

	accounts := make([]*QOSAccount, 0)
	tis := strings.Split(str, ";")
	for _, ti := range tis {
		if ti == "" {
			continue
		}

		addrAndCoins := strings.Split(ti, ",")
		if len(addrAndCoins) < 2 {
			return nil, fmt.Errorf("`%s` not match rules", ti)
		}

		var addr btypes.Address
		if !strings.HasPrefix(addrAndCoins[0], btypes.PREF_ADD) {
			if keybaseAvailable {
				info, err := keybase.Get(addrAndCoins[0])
				if err != nil {
					return nil, err
				}
				addr = info.GetAddress()
			} else {
				return nil, fmt.Errorf("no %s found in keybase", addrAndCoins[0])
			}
		} else {
			addr, err = btypes.GetAddrFromBech32(addrAndCoins[0])
		}
		if err != nil {
			return nil, err
		}
		qos, qscs, err := ParseCoins(strings.Join(addrAndCoins[1:], ","))
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, NewQOSAccount(addr, qos, qscs))
	}

	return accounts, nil
}
