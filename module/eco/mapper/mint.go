package mapper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	ecotypes "github.com/QOSGroup/qos/module/eco/types"
)

type MintMapper struct {
	*mapper.BaseMapper
}

func NewMintMapper() *MintMapper {
	var qscMapper = MintMapper{}

	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, ecotypes.MintMapperName)
	return &qscMapper
}

func GetMintMapper(ctx context.Context) *MintMapper {
	return ctx.Mapper(ecotypes.MintMapperName).(*MintMapper)
}

func (mapper *MintMapper) Copy() mapper.IMapper {
	qscMapper := &MintMapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

// 获取当前Inflation Phrase的键值
func (mapper *MintMapper) GetCurrentInflationPhraseKey(newPhrase bool) ([]byte, error) {
	// 使用KVStorePrefixIterator，当前应该是key最小的也就是第一个
	iter := store.KVStorePrefixIterator(mapper.BaseMapper.GetStore(), ecotypes.BuildMintParamsKey())
	if !iter.Valid() {
		return nil, errors.New("No more coins to come, sad!")
	}
	inflationPhraseKey := iter.Key()
	endtimesecBytes := inflationPhraseKey[len(ecotypes.BuildMintParamsKey()):]
	var endtimesec uint64
	binary.Read(bytes.NewBuffer(endtimesecBytes), binary.BigEndian, &endtimesec)

	nowsec := uint64(time.Now().UTC().Unix())

	// 当前时间已经超过endtime，需要进入下一phrase
	if nowsec >= endtimesec {
		if newPhrase {
			// 排除设置错误，为啥会刚删过又删？
			return nil, errors.New("Removing Inflation Plans too frequently")
		}
		// 删掉过期的phrase
		mapper.Del(inflationPhraseKey)
		return mapper.GetCurrentInflationPhraseKey(true)
	}
	return iter.Key(), nil
}

// 获取当前的Inflation Phrase
func (mapper *MintMapper) GetCurrentInflationPhrase() (inflationPhrase ecotypes.InflationPhrase, exist bool) {
	inflationPhrase = ecotypes.InflationPhrase{}
	currentInflationPhraseKey, err := mapper.GetCurrentInflationPhraseKey(false)
	if err == nil {
		exist = mapper.Get(currentInflationPhraseKey, &inflationPhrase)
	}
	// TODO: dealing with errors
	return
}

// 设置Params
// key:		MintParamsKey+endtime
// value: 	InflationPhrase
func (mapper *MintMapper) SetParams(config ecotypes.MintParams) {
	for _, inflation_phrase := range config.Phrases {
		mapper.AddInflationPhrase(inflation_phrase)
	}
}

func (mapper *MintMapper) AddInflationPhrase(phrase ecotypes.InflationPhrase) {

	endsec := uint64(phrase.EndTime.UTC().Unix())

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, endsec)

	keylen := len(ecotypes.BuildMintParamsKey())
	bz := make([]byte, keylen+8)

	copy(bz[0:keylen], ecotypes.BuildMintParamsKey())
	copy(bz[keylen:keylen+8], secBytes)

	mapper.Set(bz, phrase)
}

// 设置Params
// key:		MintParamsKey+endtime
// value: 	InflationPhrase
func (mapper *MintMapper) SetMintParams(config ecotypes.MintParams) {
	for _, inflation_phrase := range config.Phrases {
		mapper.AddInflationPhrase(inflation_phrase)
	}
}

func (mapper *MintMapper) GetMintParams() ecotypes.MintParams {
	var phrases []ecotypes.InflationPhrase
	iter := store.KVStorePrefixIterator(mapper.BaseMapper.GetStore(), ecotypes.BuildMintParamsKey())

	for ; iter.Valid(); iter.Next() {
		var inflationPhrase ecotypes.InflationPhrase
		mapper.DecodeObject(iter.Value(), &inflationPhrase)
		phrases = append(phrases, inflationPhrase)
	}

	return ecotypes.MintParams{phrases}
}

// 获取已分配QOS总数
func (mapper *MintMapper) GetAppliedQOSAmount() (v uint64) {
	currentInflationPhrase, exists := mapper.GetCurrentInflationPhrase()
	if !exists {
		return 0
	}
	return currentInflationPhrase.AppliedAmount
}

// 设置 已分配 QOS amount
func (mapper *MintMapper) SetAppliedQOSAmount(amount uint64) {
	inflationPhrase := ecotypes.InflationPhrase{}
	currentInflationPhraseKey, err := mapper.GetCurrentInflationPhraseKey(false)
	if err == nil {
		mapper.Get(currentInflationPhraseKey, &inflationPhrase)
		inflationPhrase.AppliedAmount = amount
		mapper.Set(currentInflationPhraseKey, inflationPhrase)
	}
	// TODO dealing with errors
}

// 增加 已分配 QOS amount
func (mapper *MintMapper) AddAppliedQOSAmount(amount uint64) {
	mined := mapper.GetAppliedQOSAmount()
	mined += amount
	mapper.SetAppliedQOSAmount(mined)
}

func (mapper *MintMapper) SetFirstBlockTime(t int64) {
	mapper.Set(ecotypes.BuildFirstBlockTimeKey(), t)
}

func (mapper *MintMapper) GetFirstBlockTime() (t int64) {
	mapper.Get(ecotypes.BuildFirstBlockTimeKey(), &t)
	return
}
