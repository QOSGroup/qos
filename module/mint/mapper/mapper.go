package mapper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/module/mint/types"
)

type Mapper struct {
	*mapper.BaseMapper
}

func NewMapper() *Mapper {
	var qscMapper = Mapper{}

	qscMapper.BaseMapper = mapper.NewBaseMapper(nil, types.MapperName)
	return &qscMapper
}

func GetMapper(ctx context.Context) *Mapper {
	return ctx.Mapper(types.MapperName).(*Mapper)
}

func (mapper *Mapper) Copy() mapper.IMapper {
	qscMapper := &Mapper{}
	qscMapper.BaseMapper = mapper.BaseMapper.Copy()
	return qscMapper
}

// 获取当前Inflation Phrase的键值
func (mapper *Mapper) GetCurrentInflationPhraseKey(blockSec uint64, newPhrase bool) ([]byte, error) {
	// 使用KVStorePrefixIterator，当前应该是key最小的也就是第一个
	iter := btypes.KVStorePrefixIterator(mapper.BaseMapper.GetStore(), types.BuildInflationPhrasesKey())
	defer iter.Close()

	if !iter.Valid() {
		return nil, errors.New("No more coins to come, sad!")
	}
	inflationPhraseKey := iter.Key()
	endtimesecBytes := inflationPhraseKey[len(types.BuildInflationPhrasesKey()):]
	var endtimesec uint64
	binary.Read(bytes.NewBuffer(endtimesecBytes), binary.BigEndian, &endtimesec)

	// 当前时间已经超过endtime，需要进入下一phrase
	if blockSec >= endtimesec {
		if newPhrase {
			// 排除设置错误，为啥会刚删过又删？
			return nil, errors.New("Removing Inflation Plans too frequently")
		}
		// 删掉过期的phrase
		mapper.Del(inflationPhraseKey)
		return mapper.GetCurrentInflationPhraseKey(blockSec, true)
	}
	return iter.Key(), nil
}

// 获取当前的Inflation Phrase
func (mapper *Mapper) GetCurrentInflationPhrase(blockSec uint64) (inflationPhrase types.InflationPhrase, exist bool) {
	inflationPhrase = types.InflationPhrase{}
	currentInflationPhraseKey, err := mapper.GetCurrentInflationPhraseKey(blockSec, false)
	if err == nil {
		exist = mapper.Get(currentInflationPhraseKey, &inflationPhrase)
	}
	// TODO: dealing with errors
	return
}

func (mapper *Mapper) AddInflationPhrase(phrase types.InflationPhrase) {

	endsec := uint64(phrase.EndTime.UTC().Unix())

	secBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(secBytes, endsec)

	keylen := len(types.BuildInflationPhrasesKey())
	bz := make([]byte, keylen+8)

	copy(bz[0:keylen], types.BuildInflationPhrasesKey())
	copy(bz[keylen:keylen+8], secBytes)

	mapper.Set(bz, phrase)
}

// 设置Params
// key:		InflationPhrasesKey+endTime
// value: 	InflationPhrase
func (mapper *Mapper) SetMintParams(config types.Params) {
	for _, phrase := range config.Phrases {
		mapper.AddInflationPhrase(phrase)
	}
}

func (mapper *Mapper) GetMintParams() types.Params {
	var phrases []types.InflationPhrase
	iter := btypes.KVStorePrefixIterator(mapper.BaseMapper.GetStore(), types.BuildInflationPhrasesKey())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var inflationPhrase types.InflationPhrase
		mapper.DecodeObject(iter.Value(), &inflationPhrase)
		phrases = append(phrases, inflationPhrase)
	}

	return types.Params{phrases}
}

// 获取当前阶段已分配QOS总数
func (mapper *Mapper) getCurrentPhraseAppliedQOSAmount(blockSec uint64) (v uint64) {
	currentInflationPhrase, exists := mapper.GetCurrentInflationPhrase(blockSec)
	if !exists {
		return 0
	}
	return currentInflationPhrase.AppliedAmount
}

// 设置当前阶段已分配 QOS amount
func (mapper *Mapper) setCurrentPhraseAppliedQOSAmount(blockSec uint64, amount uint64) {
	inflationPhrase := types.InflationPhrase{}
	currentInflationPhraseKey, err := mapper.GetCurrentInflationPhraseKey(blockSec, false)
	if err == nil {
		mapper.Get(currentInflationPhraseKey, &inflationPhrase)
		inflationPhrase.AppliedAmount = amount
		mapper.Set(currentInflationPhraseKey, inflationPhrase)
	}
	// TODO dealing with errors
}

// 增加当前阶段已分配 QOS amount
func (mapper *Mapper) addCurrentPhraseAppliedQOSAmount(blockSec uint64, amount uint64) {
	mined := mapper.getCurrentPhraseAppliedQOSAmount(blockSec)
	mined += amount
	mapper.setCurrentPhraseAppliedQOSAmount(blockSec, mined)
}

//mint处理:
//1. 当前阶段分配数
//2. 总分配
func (mapper *Mapper) MintQOS(blockSec uint64, amount uint64) {
	mapper.addAllTotalMintQOSAmount(amount)
	mapper.addCurrentPhraseAppliedQOSAmount(blockSec, amount)
}

func (mapper *Mapper) SetFirstBlockTime(t int64) {
	mapper.Set(types.BuildFirstBlockTimeKey(), t)
}

func (mapper *Mapper) GetFirstBlockTime() (t int64) {
	mapper.Get(types.BuildFirstBlockTimeKey(), &t)
	return
}

//获取总分配的QOS总数
func (mapper *Mapper) GetAllTotalMintQOSAmount() (amount uint64) {
	mapper.Get(types.BuildAllTotalMintQOSKey(), &amount)
	return
}

func (mapper *Mapper) DelAllTotalMintQOSAmount() {
	mapper.Del(types.BuildAllTotalMintQOSKey())
}

//设置总分配的QOS总数
func (mapper *Mapper) SetAllTotalMintQOSAmount(amount uint64) {
	mapper.Set(types.BuildAllTotalMintQOSKey(), amount)
}

//增加总分配的QOS总数
func (mapper *Mapper) addAllTotalMintQOSAmount(amount uint64) {

	totalAmount := mapper.GetAllTotalMintQOSAmount()
	totalAmount += amount

	mapper.SetAllTotalMintQOSAmount(totalAmount)
}
