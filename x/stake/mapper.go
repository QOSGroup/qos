package stake

import (
	"encoding/binary"
	"github.com/QOSGroup/qbase/mapper"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"math"
)

const (
	StakeMapperName = "stake"
	StakeParamsName = "params"
)

var (
	currentTotalPowerKey			= []byte{0x01}	// value is the current Total Power(sum of bonded tokens of active validators)

	validatorVotingStatusLenKey   = []byte{0x02} // param, key: prefix, value: how many latest blocks the will be observed of a validator's voting
	validatorVotingStatusLeastKey = []byte{0x03} // param, key: prefix, value: at least how many blocks must be signed out of validatorVotingStatusLen for a validator to stay active

	maxValidatorCntKey				= []byte{0x04}	// param, no restraint in testnet, hence set to 10000

	validatorVotingStatusKey		= []byte{0x11}	// key: prefix + operator address, value: bits of validator vote, length = validatorVotingStatusLen
	validatorVotingCntKey			= []byte{0x12}	// key: prefix + operator address, value: vote count in validatorVotingStatus
	currentValidatorPowerRankKey	= []byte{0x13}	// key: prefix + ranked by (voting power + validatorVotingCnt + (this time) activate height), value: operator address

	//pendingValidatorTimeSliceKey	= []byte{0x21}	// key: prefix + pending time, value: list of operator address
)

type StakeMapper struct {
	*mapper.BaseMapper
}

var _ mapper.IMapper = (*StakeMapper)(nil)

func NewStakeMapper() *StakeMapper {
	var stakeMapper = StakeMapper{}
	stakeMapper.BaseMapper = mapper.NewBaseMapper(nil, StakeMapperName)
	return &stakeMapper
}

func (mapper *StakeMapper) Copy() mapper.IMapper {
	stakeMapper := &StakeMapper{}
	stakeMapper.BaseMapper = mapper.BaseMapper.Copy()
	return stakeMapper
}

func BuildCurrentTotalPowerKey() []byte{
	return append([]byte(StakeParamsName), currentTotalPowerKey...)
}

func BuildValidatorVotingStatusLenKey() []byte{
	return append([]byte(StakeParamsName), validatorVotingStatusLenKey...)
}

func BuildValidatorVotingStatusLeastKey() []byte{
	return append([]byte(StakeParamsName), validatorVotingStatusLeastKey...)
}

func BuildMaxValidatorCntKey() []byte{
	return append([]byte(StakeParamsName), maxValidatorCntKey...)
}

func BuildValidatorVotingStatusKey(operator btypes.Address) []byte{
	return append(validatorVotingStatusKey, operator...)
}

func BuildValidatorVotingCntKey(operator btypes.Address) []byte{
	return append(validatorVotingCntKey, operator...)
}

func (mapper *StakeMapper) BuildValidatorPowerRankKey(validator types.Validator) []byte{
	votingPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(votingPowerBytes[:], uint64(validator.VotingPower))

	votingCntBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(votingCntBytes[:], uint64(mapper.GetValidatorVotingCnt(validator.Operator)))

	// height should be reverted, earlier activated, higher ranked
	activeHeightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(activeHeightBytes[:], ^uint64(0) - uint64(validator.Height)) // max value of uint64 - validator.height

	powerRankKey := append(currentValidatorPowerRankKey, votingPowerBytes...)
	powerRankKey = append(powerRankKey, votingCntBytes...)
	powerRankKey = append(powerRankKey, activeHeightBytes...) // height should be reverted, earlier activated, higher ranked
	return powerRankKey
}

func (mapper *StakeMapper) GetCurrentTotalPower() uint64 {
	var totalPower uint64
	mapper.Get(BuildCurrentTotalPowerKey(), &totalPower)
	return totalPower
}

func (mapper *StakeMapper) SetCurrentTotalPower(power uint64) {
	mapper.Set(BuildCurrentTotalPowerKey(), power)
}

func (mapper *StakeMapper) GetValidatorVotingStatusLen() uint64 {
	var validatorVotingStatusLen uint64
	mapper.Get(BuildValidatorVotingStatusLenKey(), &validatorVotingStatusLen)
	return validatorVotingStatusLen
}

func (mapper *StakeMapper) GetValidatorVotingStatusLeast() uint64 {
	var validatorVotingStatusLeast uint64
	mapper.Get(BuildValidatorVotingStatusLeastKey(), &validatorVotingStatusLeast)
	return validatorVotingStatusLeast
}

func (mapper *StakeMapper) GetMaxValidatorCnt() uint64 {
	var maxValidatorCnt uint64
	mapper.Get(BuildMaxValidatorCntKey(), &maxValidatorCnt)
	return maxValidatorCnt
}

func (mapper *StakeMapper) GetValidatorVotingStatus(operator btypes.Address) uint64 {
	var validatorVotingStatus uint64
	mapper.Get(BuildValidatorVotingStatusKey(operator), &validatorVotingStatus)
	return validatorVotingStatus
}

func (mapper *StakeMapper) SetValidatorVotingStatus(operator btypes.Address, status uint64) {
	mapper.Set(BuildValidatorVotingStatusKey(operator), status)
}

func (mapper *StakeMapper) GetValidatorVotingCnt(operator btypes.Address) uint64 {
	var validatorVotingCnt uint64
	mapper.Get(BuildValidatorVotingCntKey(operator), &validatorVotingCnt)
	return validatorVotingCnt
}

func (mapper *StakeMapper) SetValidatorVotingCnt(operator btypes.Address, cnt uint64){
	mapper.Set(BuildValidatorVotingCntKey(operator), cnt)
}


// 更新最新的投票信息（无论当前块是否投了票） update latest voting status (no matter having voted in the latest block)
// operator: 当前validator的operator address
// hasVote: 最新块是否投票，投票1，未投0
// 返回：更新后validator在窗口期内的投票数 returns voting cnt of the validator in it's observed voting status duration
func (mapper *StakeMapper) addLatestVoteStatus(operator btypes.Address, hasVote int) uint64{
	status := mapper.GetValidatorVotingStatus(operator)
	// 最高位
	extrusion := ((uint64(status) & uint64(math.Pow(2, float64(mapper.GetValidatorVotingStatusLen() - 1)))) > 0)

	//更新status，左移一位，最低位加最新块投票，截短
	status <<= 1
	status |= uint64(hasVote)
	status &= uint64(math.Pow(2, float64(mapper.GetValidatorVotingStatusLen())) - 1)
	mapper.SetValidatorVotingStatus(operator, status)

	votingCnt := mapper.GetValidatorVotingCnt(operator)
	votingCnt += uint64(hasVote)
	if(extrusion){
		votingCnt -= 1
	}
	mapper.SetValidatorVotingCnt(operator, votingCnt)
	return votingCnt
}