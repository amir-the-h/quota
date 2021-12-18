package quota

import (
	"errors"
	"sort"
	"time"

	"github.com/amir-the-h/okex"
)

// Quota is the group of candles and make time-series.
type Quota []*Candle

// Symbol will return the symbol of candles.
func (q *Quota) Symbol() string {
	for _, candle := range *q {
		return candle.Symbol
	}

	return ""
}

// BarSize will return the okex.BarSize of candles.
func (q *Quota) BarSize() okex.BarSize {
	for _, candle := range *q {
		return candle.BarSize
	}

	return ""
}

// IndicatorTags returns a list of indicators used in quota.
func (q *Quota) IndicatorTags() []IndicatorTag {
	var tags []IndicatorTag
	for _, candle := range *q {
		for indicator := range candle.Indicators {
			hasIndicator := false
			for _, tag := range tags {
				if tag == indicator {
					hasIndicator = true
					break
				}
			}
			if !hasIndicator {
				tags = append(tags, indicator)
			}
		}
	}

	return tags
}

// Find searches for a candle, and its index among Quota by its symbol and provided timestamp.
func (q *Quota) Find(timestamp int64) (*Candle, int) {
	for i, candle := range *q {
		if candle.OpenTime.Unix() == timestamp || (candle.OpenTime.Unix() < timestamp && candle.CloseTime.Unix() > timestamp) {
			return candle, i
		}
	}

	return nil, -1
}

// Sort runs through the quota and reorder candles by the open time.
func (q *Quota) Sort() {
	sort.Slice(*q, func(i, j int) bool { return (*q)[i].OpenTime.Before((*q)[j].OpenTime) })
	for i, candle := range *q {
		if i > 0 {
			candle.Previous = (*q)[i-1]
		}
		if i < len(*q)-1 {
			candle.Next = (*q)[i+1]
		}
	}
}

// Sync searches the quota for provided candle and update it if it exists, otherwise, it will append to end of the quota.
//
// If you want to update a candle directly then pass sCandle.
func (q *Quota) Sync(open, high, low, close, volume float64, openTime, closeTime time.Time, sCandle ...*Candle) (candle *Candle, err CandleError) {
	var lc *Candle
	checker := func(candle *Candle, openTime, closeTime time.Time) bool {
		return candle.OpenTime.Equal(openTime) && candle.CloseTime.Equal(closeTime)
	}

	// try last candle first
	if len(*q) > 0 {
		lc = (*q)[len(*q)-1]
		candle = lc
	}

	// if any suspicious candle provided try it then.
	if len(sCandle) > 0 {
		candle = sCandle[0]
	}

	if candle == nil || !checker(candle, openTime, closeTime) {
		candle, _ = q.Find(openTime.Unix())
		if candle == nil {
			candle, err = NewCandle(open, high, low, close, volume, q.Symbol(), q.BarSize(), openTime, closeTime, lc, nil)
			if err != nil {
				return
			}
			*q = append(*q, candle)
			if lc != nil {
				lc.Next = candle
			}
		}
	}
	candle.Open = open
	candle.High = high
	candle.Low = low
	candle.Close = close
	candle.Volume = volume

	return
}

// Merge target quota into the current quota, rewrite duplicates and sort it.
func (q *Quota) Merge(target *Quota) {
	for _, candle := range *target {
		c, i := q.Find(candle.OpenTime.Unix())
		if c != nil {
			(*q)[i] = candle
		} else {
			*q = append(*q, candle)
		}
	}
}

// AddIndicator adds unimplementedIndicator values by the given tag into the quota.
func (q *Quota) AddIndicator(tag IndicatorTag, values []float64) error {
	quota := *q
	if len(values) != len(quota) {
		return errors.New("count mismatched")
	}

	for i := range values {
		(*q)[i].Indicators[tag] = values[i]
	}

	return nil
}
