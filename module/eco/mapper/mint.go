package mapper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
	"github.com/QOSGroup/qbase/mapper"
	"github.com/QOSGroup/qbase/store"
	minttypes "github.com/QOSGroup/qos/module/eco/types"
)

type MintMapper struct {
	*mapper.BaseMapper
}

func NewMintMapper() *MintMapper {
	var qscMapper = MintMapper{}
	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, minttypes.MintMapperName)
	return &qscMapper
}

func (mapper *MintMapper) Copy() mapper.IMapper {
	qscMapper := &MintMapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

// 获取当前Inflation Phrase的键值
func (mapper *MintMapper) GetCurrentInflationPhraseKey(newPhrase bool) ([]byte, error) {
	// 使用KVStorePrefixIterator，当前应该是key最小的也就是第一个
	iter := store.KVStorePrefixIterator(mapper.BaseMapper.GetStore(), []byte(minttypes.MintParamsKey))
	if !iter.Valid() {
		return nil, errors.New("No more coins to come, sad!")
	}
	inflationPhraseKey := iter.Key()
	endtimesecBytes := inflationPhraseKey[len([]byte(minttypes.MintParamsKey)):]
	var endtimesec uint64
	binary.Read(bytes.NewBuffer(endtimesecBytes), binary.BigEndian, &endtimesec)

	nowsec := uint64(time.Now().UTC().Unix())

	// 当前时间已经超过endtime，需要进入下一phrase
	if (nowsec >= endtimesec) {
		if (newPhrase) {
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
func (mapper *MintMapper) GetCurrentInflationPhrase() (inflationPhrase minttypes.InflationPhrase, exist bool) {
	inflationPhrase = minttypes.InflationPhrase{}
	currentInflationPhraseKey, err := mapper.GetCurrentInflationPhraseKey(false)
	if err == nil {
		exist = mapper.Get(currentInflationPhraseKey, &inflationPhrase)
	}
	// TODO: dealing with errors
	return
}

func (mapper *MintMapper) AddInflationPhrase(phrase minttypes.InflationPhrase) {
	endsec := uint64(phrase.EndTime.UTC().Unix())

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, endsec)

	keylen := len([]byte(minttypes.MintParamsKey))
	bz := make([]byte, keylen + 8)

	copy(bz[0:keylen], []byte(minttypes.MintParamsKey))
	copy(bz[keylen:keylen+8], secBytes)

	mapper.Set(bz, phrase)
}

// 设置SPOConfig
// key:		SPOConfigKey+endtime
// value: 	InflationPhrase
func (mapper *MintMapper) SetMintParams(config minttypes.MintParams) {
	for _, inflation_phrase := range config.Phrases {
		endsec := uint64(inflation_phrase.EndTime.UTC().Unix())

		secBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(secBytes, endsec)

		keylen := len([]byte(minttypes.MintParamsKey))
		bz := make([]byte, keylen + 8)

		copy(bz[0:keylen], []byte(minttypes.MintParamsKey))
		copy(bz[keylen:keylen+8], secBytes)

		mapper.Set(bz, inflation_phrase)
	}
}

func (mapper *MintMapper) GetMintParams() minttypes.MintParams {
	var phrases []minttypes.InflationPhrase
	iter := store.KVStorePrefixIterator(mapper.BaseMapper.GetStore(), []byte(minttypes.MintParamsKey))
	for {
		var inflationPhrase minttypes.InflationPhrase
		mapper.DecodeObject(iter.Value(), &inflationPhrase)
		phrases = append(phrases, inflationPhrase)
		iter.Next()
	}
	return minttypes.MintParams{phrases}
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
	inflationPhrase := minttypes.InflationPhrase{}
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
