package txs

import (
	"fmt"
	baccount "github.com/QOSGroup/qbase/account"
	bcontext "github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/common"
	"time"
)

// 通过地址获取QOSAccount
func GetAccount(ctx bcontext.Context, addr btypes.Address) (acc *account.QOSAccount) {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return nil
	}

	accbase := mapper.GetAccount(addr)
	if accbase == nil {
		return nil
	}
	acc = accbase.(*account.QOSAccount)

	return
}

// 通过地址创建QOSAccount
// 若账户存在，返回账户 & false
func CreateAndSaveAccount(ctx bcontext.Context, addr btypes.Address) (acc *account.QOSAccount, success bool) {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return nil, false
	}

	accfind := mapper.GetAccount(addr).(*account.QOSAccount)
	if accfind != nil {
		return accfind, false
	}

	acc = mapper.NewAccountWithAddress(addr).(*account.QOSAccount)
	mapper.SetAccount(acc)

	return acc, true
}

func SaveAccount(ctx bcontext.Context, acc *account.QOSAccount) bool {
	mapper := ctx.Mapper(baccount.AccountMapperName).(*baccount.AccountMapper)
	if mapper == nil {
		return false
	}
	mapper.SetAccount(acc)

	return true
}

func GetRootPubkey(ctx bcontext.Context) crypto.PubKey {
	mainmapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.MainMapper)
	if mainmapper == nil {
		return nil
	}

	return mainmapper.GetRoot()
}

func GetCrtcodec() *amino.Codec {
	cdccrt := amino.NewCodec()
	CrtRegisterCodecSingle(cdccrt)

	return cdccrt
}

func FetchCA(crtfile string) (caQsc *[]byte) {
	var ca Certificate
	cdccrt := GetCrtcodec()

	byfile := common.MustReadFile(crtfile)
	err := cdccrt.UnmarshalBinaryBare(byfile, &ca)
	if err != nil {
		panic(fmt.Sprintf("error: Decode %s", crtfile))
	}
	caQsc = &byfile

	return
}

type Subject struct {
	// TODO: Compatible with the openssl
	CN string `json:"cn"`
}

type CertificateSigningRequest struct {
	Subj      Subject               `json:"subj"`
	IsCa      bool                  `json:"is_ca"`
	IsBanker  bool                  `json:"is_banker"`
	NotBefore time.Time             `json:"not_before"`
	NotAfter  time.Time             `json:"not_after"`
	PublicKey ed25519.PubKeyEd25519 `json:"public_key"`
}

type Issuer struct {
	Subj      Subject               `json:"subj"`
	PublicKey ed25519.PubKeyEd25519 `json:"public_key"`
}

type Certificate struct {
	CSR       CertificateSigningRequest `json:"csr"`
	CA        Issuer                    `json:"ca"`
	Signature []byte                    `json:"signature"`
}

func CrtRegisterCodec(cdc *amino.Codec) {
	const (
		CsrAminoRoute = "certificate/csr"
		CrtAminoRoute = "certificate/crt"
	)

	cdc.RegisterConcrete(CertificateSigningRequest{}, CsrAminoRoute, nil)
	cdc.RegisterConcrete(Certificate{}, CrtAminoRoute, nil)
}

func CrtRegisterCodecSingle(cdc *amino.Codec) {
	CrtRegisterCodec(cdc)

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.Ed25519PubKeyAminoRoute, nil)

	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.Ed25519PrivKeyAminoRoute, nil)
}

func VerityCrt(caPublicKeys []ed25519.PubKeyEd25519, crt Certificate) bool {
	ok := false

	// Check issuer
	cdc := GetCrtcodec()
	signbyte, err := cdc.MarshalBinaryBare(crt.CSR)
	if err != nil {
		fmt.Print("Decode crt byte error!")
		return false
	}

	for _, value := range caPublicKeys {
		if value.Equals(crt.CA.PublicKey) {
			ok = crt.CA.PublicKey.VerifyBytes(signbyte, crt.Signature)
			break
		}
	}

	// Check timestamp
	now := time.Now().Unix()
	if now <= crt.CSR.NotBefore.Unix() || now >= crt.CSR.NotAfter.Unix() {
		ok = false
	}

	return ok
}

func NewCertificate(cdc *amino.Codec, qscname string, isbanker bool, pubkey crypto.PubKey, rootkey crypto.PubKey) []byte {
	crt := Certificate{
		CertificateSigningRequest{
			Subject{qscname},
			false,
			isbanker,
			time.Now(),
			time.Now().AddDate(1, 0, 0),
			pubkey.(ed25519.PubKeyEd25519),
		},
		Issuer{
			Subject{"root"},
			rootkey.(ed25519.PubKeyEd25519),
		},
		[]byte{},
	}
	crtbyte, err := cdc.MarshalBinaryBare(crt)
	if err != nil {
		return nil
	}

	return crtbyte
}
