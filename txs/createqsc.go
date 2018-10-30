package txs

import (
	"fmt"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/qcp"
	btxs "github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	qosmapper "github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/types"
	go_amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
)

const BASEGAS_CREATEQSC int64 = 10000 //创建qsc需要的最少qos数

// 功能："创建QSC" 对应的Tx结构
type TxCreateQSC struct {
	QscName     string         `json:"qscname"`     //从CA信息获取
	CreateAddr  btypes.Address `json:"createaddr"`  //QSC创建账户
	QscPubkey   crypto.PubKey  `json:"qscpubkey"`   //从CA信息获取
	Banker      btypes.Address `json:"banker"`      //从CA信息获取
	Extrate     string         `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	CAqsc       []byte         `json:"caqsc"`       //CA信息
	CAbanker    []byte         `json:"cabanker"`    //CA信息
	Description string         `json:"description"` //描述信息
	AccInit     []AddrCoin     `json:"accinit"`     //初始化时接受qsc的账户
}

type AddrCoin struct {
	Address btypes.Address `json:"address"` //用户地址
	Amount  btypes.BigInt  `json:"amount"`  //金额
}

// 功能：检测合法性
// 备注：
//		1,成员字段的合法性
//		2,creator的账户余额是否够gas抵扣
func (tx TxCreateQSC) ValidateData(ctx context.Context) bool {
	if !btypes.CheckQscName(tx.QscName) || !CheckAddr(tx.CreateAddr) || !CheckAddr(tx.Banker) {
		return false
	}

	return true
}

// 功能：tx执行
// 描述：
//		保存链信息并检查是否已经创建(qscname, pubkey, 是否已经注册)
//      查询banker是否存在，若不存在，
//		向账户 AccInit 分发qsc
func (tx TxCreateQSC) Exec(ctx context.Context) (ret btypes.Result, crossTxQcps *btxs.TxQcp) {
	mapper := ctx.Mapper(qosmapper.BaseMapperName).(*qosmapper.MainMapper)
	if mapper == nil {
		ret.Log = "Get qsc mapper error!"
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}

	// qsc已存在,不能创建
	if mapper.GetQsc(tx.QscName) != nil {
		ret.Log = fmt.Sprintf("Error: QSC(%s) exist!", tx.QscName)
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}

	var cabank CA
	err := cdc.UnmarshalBinaryBare(tx.CAbanker, &cabank)
	if err != nil {
		ret.Log = err.Error()
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}

	// 检查banker: 不存在则创建; 存在,验证pubkey
	acc := GetAccount(ctx, tx.Banker)
	if acc == nil {
		acc, _ = CreateAndSaveAccount(ctx, tx.Banker)
		ret.Log += "Account: create banker"
	} else {
		if !acc.GetPubicKey().Equals(cabank.Pubkey) {
			ret.Log = "Error: Exist banker's pubkey is different from cabank's"
			ret = btypes.ErrInternal(ret.Log).Result()
			return
		}
	}

	// 保存qsc 信息
	qscinfo := qosmapper.QscInfo{
		tx.QscName,
		cabank.Pubkey,
	}
	if !mapper.SetQsc(tx.QscName, &qscinfo) {
		ret.Log = "Error: Save qsc info error"
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}

	// 给账户分发qsc
	for _, va := range tx.AccInit {
		vaAcc := GetAccount(ctx, va.Address)
		if &vaAcc == nil {
			// vaAcc, _ = CreateAndSaveAccount(ctx, va.Address)
			vaAcc = account.ProtoQOSAccount().(*account.QOSAccount)
			vaAcc.SetAddress(va.Address)
			ret.Log = "Account: create account :" + va.Address.String()
		}

		vaAcc.SetQSC(&types.QSC{tx.QscName, va.Amount})
		SaveAccount(ctx, vaAcc)
	}

	// 将联盟链的publickey加入(chainid/in/pubkey)
	qcpmapper := ctx.Mapper(qcp.QcpMapperName).(*qcp.QcpMapper)
	if qcpmapper == nil {
		ret.Log = "Error: Get qcpmapper error!"
		ret = btypes.ErrInternal(ret.Log).Result()
		return
	}
	qcpmapper.SetChainInTrustPubKey(tx.QscName, tx.QscPubkey)
	ret.Code = btypes.ABCICodeOK

	return
}

//功能：获取签名者
func (tx TxCreateQSC) GetSigner() (ret []btypes.Address) {
	if tx.CreateAddr == nil {
		return nil
	}

	ret = []btypes.Address{tx.CreateAddr}
	return
}

// 功能：计算gas
// 规则：基准值 + 每个初始化用户收10qos
func (tx TxCreateQSC) CalcGas() btypes.BigInt {
	baseGas := btypes.NewInt(BASEGAS_CREATEQSC)
	var accNum int = len(tx.AccInit)
	return baseGas.Add(btypes.NewInt(int64(accNum * 10)))
}

//gas付费人
func (tx TxCreateQSC) GetGasPayer() (ret btypes.Address) {
	if tx.CreateAddr == nil {
		return nil
	}

	ret = tx.CreateAddr
	return
}

// 获取签名字段
func (tx TxCreateQSC) GetSignData() (ret []byte) {
	ret = append(ret, []byte(tx.QscName)...)
	ret = append(ret, tx.QscPubkey.Bytes()...)
	ret = append(ret, []byte(tx.Banker)...)
	ret = append(ret, []byte(tx.Extrate)...)
	ret = append(ret, tx.CAqsc...)
	ret = append(ret, tx.CAbanker...)
	ret = append(ret, []byte(tx.Description)...)

	for _, acn := range tx.AccInit {
		ret = append(ret, acn.Address...)
		ret = append(ret, btypes.Int2Byte(acn.Amount.Int64())...)
	}

	return
}

// 创建 TxCreateQSC结构体
// 备注：CA提供两个证书，联盟链证书 & Banker证书(banker字段)
//		两种证书通过 qscName 字段关联起来
func NewCreateQsc(cdc *go_amino.Codec, caqsc *[]byte, cabank *[]byte,
	createAddr btypes.Address, accs *[]AddrCoin,
	extrate string, dsp string) (rTx *TxCreateQSC) {

	var dataqsc CA
	cdc.UnmarshalBinaryBare(*caqsc, &dataqsc)
	if dataqsc.Banker {
		//qsc的ca证书中banker == false
		return nil
	}

	var databank CA
	cdc.UnmarshalBinaryBare(*cabank, &databank)
	if !databank.Banker {
		//qsc的ca证书中banker == false
		return nil
	}
	if databank.Qcpname != dataqsc.Qcpname {
		return nil
	}

	if accs == nil {
		accs = &[]AddrCoin{}
	}

	rTx = &TxCreateQSC{
		dataqsc.Qcpname,
		createAddr,
		dataqsc.Pubkey,
		[]byte(databank.Pubkey.Address()),
		extrate,
		*caqsc,
		*cabank,
		dsp,
		*accs,
	}

	return
}
