package types

import (
	"errors"
	"fmt"
	"time"
)

type InflationPhrase struct {
	EndTime       time.Time `json:"end_time"`
	TotalAmount   uint64    `json:"total_amount"`
	AppliedAmount uint64    `json:"applied_amount"`
}

func (phrase InflationPhrase) Equals(p InflationPhrase) bool {
	return phrase.EndTime.Equal(p.EndTime) && phrase.TotalAmount == p.TotalAmount
}

type InflationPhrases []InflationPhrase

func DefaultInflationPhrases() InflationPhrases {
	return InflationPhrases{
		{
			time.Date(2023, 10, 20, 0, 0, 0, 0, time.UTC),
			25.5e12,
			0,
		},
		{
			time.Date(2027, 10, 20, 0, 0, 0, 0, time.UTC),
			12.75e12,
			0,
		},
		{
			time.Date(2031, 10, 20, 0, 0, 0, 0, time.UTC),
			6.375e12,
			0,
		},
		{
			time.Date(2035, 10, 20, 0, 0, 0, 0, time.UTC),
			3.1875e12,
			0,
		},
		{
			time.Date(2039, 10, 20, 0, 0, 0, 0, time.UTC),
			1.59375e12,
			0,
		},
		{
			time.Date(2043, 10, 20, 0, 0, 0, 0, time.UTC),
			0.796875e12,
			0,
		},
		{
			time.Date(2047, 10, 20, 0, 0, 0, 0, time.UTC),
			0.796875e12,
			0,
		},
	}
}

func (phrases InflationPhrases) Equals(ips InflationPhrases) bool {
	if len(phrases) != len(ips) {
		return false
	}

	pm := make(map[time.Time]InflationPhrase, len(phrases))
	for _, p := range phrases {
		pm[p.EndTime] = p
	}

	for _, p := range ips {
		phrase, ok := pm[p.EndTime]
		if !ok || !p.Equals(phrase) {
			return false
		}
	}

	return true
}

// 获取时间点对应通胀阶段
func (phrases InflationPhrases) GetPhrase(time time.Time) (phrase *InflationPhrase, exists bool) {
	for _, p := range phrases {
		endTime := p.EndTime.UTC()
		if endTime.After(time) && (phrase == nil || endTime.Before(phrase.EndTime.UTC())) {
			phrase = &InflationPhrase{p.EndTime, p.TotalAmount, p.AppliedAmount}
			exists = true
		}
	}

	return
}

// 获取前一通胀阶段
func (phrases InflationPhrases) GetPrePhrase(time time.Time) (phrase *InflationPhrase, exists bool) {
	for _, p := range phrases {
		endTime := p.EndTime.UTC()
		if !endTime.After(time) && (phrase == nil || endTime.After(phrase.EndTime.UTC())) {
			phrase = &InflationPhrase{p.EndTime, p.TotalAmount, p.AppliedAmount}
			exists = true
		}
	}

	return
}

// 释放通胀
func (phrases InflationPhrases) ApplyQOS(phraseEndTime time.Time, amount uint64) (newPhrase InflationPhrases) {
	for _, p := range phrases {
		if p.EndTime.UTC().Equal(phraseEndTime) {
			p.AppliedAmount += amount
		}
		newPhrase = append(newPhrase, p)
	}

	return
}

// 校验
func (phrases InflationPhrases) Valid() error {
	if len(phrases) == 0 {
		return errors.New("phrases is empty")
	}
	times := map[time.Time]bool{}
	for _, p := range phrases {
		// 通胀时间不能重复
		if _, ok := times[p.EndTime]; !ok {
			times[p.EndTime] = true
		} else {
			return errors.New("duplicate end time in phrases")
		}
		// 通胀量不能为0
		if p.TotalAmount == 0 {
			return fmt.Errorf("total amount not positive in phrase:%v", p.EndTime)
		}
	}

	return nil
}

// 校验新通胀
func (phrases InflationPhrases) ValidNewPhrases(newTotal, totalApplied uint64, newPhrases InflationPhrases) error {
	if phrases.Equals(newPhrases) {
		return errors.New("phrases not change")
	}
	var phrasesApplied uint64 = 0
	currentNewPhrase, _ := newPhrases.GetPhrase(time.Now().UTC())
	for _, p := range phrases {
		phrasesApplied += p.AppliedAmount
	}

	// 新的通胀规则必须包含当前及之前通胀阶段，且对应通胀阶段TotalAmount一致
	var newPhrasesTotal uint64
	for _, np := range newPhrases {
		newPhrasesTotal += np.TotalAmount
		if currentNewPhrase != nil && !np.EndTime.After(currentNewPhrase.EndTime) {
			exists := false
			for _, p := range phrases {
				if p.EndTime == np.EndTime {
					exists = true
					if np.TotalAmount != p.TotalAmount {
						return fmt.Errorf("total amount not equals in phrase:%v", p.EndTime)
					}
				}
			}
			if !exists {
				return fmt.Errorf("new phrases must contain %v", np.EndTime)
			}
		}
	}

	// 新总发行数量 = 总已发行数量-已发行通胀总量+新通胀总量
	if newTotal == totalApplied-phrasesApplied+newPhrasesTotal {
		return nil
	} else {
		return errors.New("total amount not valid")
	}

}

// 适配旧通胀阶段，填充已发行
func (phrases InflationPhrases) Adapt(oldPhrases InflationPhrases) (newPhrase InflationPhrases) {
	newPhrases := InflationPhrases{}
	for _, p := range phrases {
		for _, op := range oldPhrases {
			if p.EndTime.Equal(op.EndTime) {
				p.AppliedAmount = op.AppliedAmount
			}
		}
		newPhrases = append(newPhrases, p)
	}
	return newPhrases
}
