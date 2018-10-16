package txs

import (
	"fmt"
	"log"

	"github.com/QOSGroup/qbase/context"
	btypes "github.com/QOSGroup/qbase/types"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
)

const BASEGAS_CREATEQSC int64 = 10000 //创建qsc需要的最少qos数

//功能："创建QSC" 对应的Tx结构
type TxCreateQSC struct {
	QscName     string         `json:"qscname"`     //从CA信息获取
	CreateAddr  btypes.Address `json:"createaddr"`  //QSC创建账户
	QscPubkey   crypto.PubKey  `json:"qscpubkey"`   //从CA信息获取
	Banker      btypes.Address `json:"banker"`      //从CA信息获取
	Extrate     string         `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	CA          []byte         `json:"ca"`          //CA信息
	Description string         `json:"description"` //描述信息
	AccInit     []AddrCoin     `json:"accinit"`     //初始化时接受qsc的账户
}

type AddrCoin struct {
	Address btypes.Address `json:"address"` //用户地址
	Amount  btypes.BigInt  `json:"amount"`  //金额
}

//功能：检测合法性
//备注：
//		1,成员字段的合法性
//		2,creator的账户余额是否够gas抵扣
func (tx *TxCreateQSC) ValidateData() bool {
	if !btypes.CheckQsc(tx.QscName) || !CheckAddr(tx.CreateAddr) || !CheckAddr(tx.Banker) {
		return false
	}

	return true
}

//功能：tx执行
//描述：creater身份校验，是否有足够的qos
//      查询banker是否存在，若不存在，
//		向账户 AccInit 分发qsc
//		扣除creater的qos(以gas形式扣除)
func (tx *TxCreateQSC) Exec(ctx context.Context) (ret btypes.Result) {
	if !tx.ValidateData() {
		ret.Code = btypes.ToABCICode(btypes.CodespaceRoot, btypes.CodeInternal) //todo: which code should be here
		ret.Log = "error: ValidateData error!"
		return
	}

	ctAcc := GetAccount(tx.CreateAddr)
	if ctAcc.GetQOS().LT(tx.CalcGas()) {
		ret.Code = btypes.ToABCICode(btypes.CodespaceRoot, btypes.CodeInternal) //todo: which code should be here
		ret.Log = "error: Create should have more qos!"
		return
	}

	//ctx.KVStore("account")
	acc := GetAccount(tx.Banker)
	if &acc == nil {
		acc = CreateAccount(tx.Banker)
		ret.Log += "banker: create banker"
	}
	for _, va := range tx.AccInit {
		vaAcc := GetAccount(va.Address)
		if &vaAcc == nil {
			vaAcc = CreateAccount(va.Address)
			ret.Log += fmt.Sprintf("Account: create account(%s),amount(%d)", va.Address, va.Amount)
		}
		vaAcc.SetQOS(va.Amount)
	}

	gas := tx.CalcGas()
	accreator := GetAccount(tx.CreateAddr)
	accreator.SetQOS(accreator.GetQOS().Sub(gas))

	// todo: 将联盟链的publickey加入kvstore,(qos/doc/store.md)(chainid/in/pubkey)
	//kvstore := store.KVStoreKey{tx.QscName + "/in/pubkey"}
	//mkvstrore := ctx.KVStore(&kvstore)
	//mkvstrore.Set([]byte(kvstore.String()), tx.QscPubkey.Bytes())

	ret.GasUsed = gas.Int64()
	ret.Code = btypes.ABCICodeOK
	return
}

//功能：获取签名者
func (tx *TxCreateQSC) GetSigner() (ret []btypes.Address) {
	if tx.CreateAddr == nil {
		log.Panic("No signer for create QSC")
		return nil
	}

	ret[0] = tx.CreateAddr
	return
}

//功能：计算gas
//规则：基准值 + 每个初始化用户收10qos
//todo：规则暂定为此，可能调整
func (tx *TxCreateQSC) CalcGas() btypes.BigInt {
	baseGas := btypes.NewInt(BASEGAS_CREATEQSC)
	var accNum int = len(tx.AccInit)
	return baseGas.Add(btypes.NewInt(int64(accNum * 10)))
}

//gas付费人
func (tx *TxCreateQSC) GetGasPayer() (ret btypes.Address) {
	if tx.CreateAddr == nil {
		log.Panic("Can't find creater in tx(createQSC)")
		return nil
	}

	ret = tx.CreateAddr
	return
}

//获取签名字段
func (tx *TxCreateQSC) GetSignData() (ret []byte) {
	ret = append(ret, []byte(tx.QscName)...)
	ret = append(ret, tx.QscPubkey.Bytes()...)
	ret = append(ret, []byte(tx.Banker)...)
	ret = append(ret, []byte(tx.Extrate)...)
	ret = append(ret, tx.CA...)
	ret = append(ret, []byte(tx.Description)...)

	for _, acn := range tx.AccInit {
		ret = append(ret, acn.Address...)
		ret = append(ret, btypes.Int2Byte(acn.Amount.Int64())...)
	}

	return
}

//CA结构体
//todo: CA具体格式确定后会更改
type CA struct {
	Qcpname string
	Banker  bool
	Pubkey  crypto.PubKey
	Info    string
}

//创建 TxCreateQSC结构体
//备注：CA提供两个证书，联盟链证书 & Banker证书(banker字段)
//		两种证书通过 qscName 字段关联起来
func NewCreateQsc(cdc *go_amino.Codec, caqsc *[]byte, cabank *[]byte,
	createAddr btypes.Address, accs *[]AddrCoin,
	extrate string, dsp string) (rTx *TxCreateQSC) {

	var dataqsc CA
	cdc.UnmarshalBinaryBare(*caqsc, &dataqsc)
	if dataqsc.Banker {
		//qsc的ca证书中banker == false
		log.Panic("CA(qcs) error")
		return nil
	}

	var databank CA
	cdc.UnmarshalBinaryBare(*cabank, &databank)
	if !databank.Banker {
		//qsc的ca证书中banker == false
		log.Panic("CA(bank) error")
		return nil
	}
	if databank.Qcpname != dataqsc.Qcpname {
		log.Panic("The two input CA(caqsc, cabank) should have the same qcpname")
		return nil
	}

	rTx = &TxCreateQSC{
		dataqsc.Qcpname,
		createAddr,
		dataqsc.Pubkey,
		[]byte("banker"), //todo: for test, extract from databank.Pubkey
		extrate,
		*caqsc,
		dsp,
		*accs,
	}

	return
}
