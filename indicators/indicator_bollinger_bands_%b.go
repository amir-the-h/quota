package indicators

import (
	"fmt"
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

type BollingerBandsB struct {
	quota.UnimplementedIndicator
	Std StandardDeviation `mapstructure:"standardDeviation"`
}

func (bbb *BollingerBandsB) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - bbb.Std.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := (*q)[startIndex : i+1]

		deviation, ok := c.Get(quota.Source(bbb.Std.Tag()))
		if !ok {
			if !bbb.Std.Add(&quote, c) {
				return false
			}

			deviation, ok = c.Get(quota.Source(bbb.Std.Tag()))
			if !ok {
				return false
			}
		}

		sma := &Ma{
			UnimplementedIndicator: quota.UnimplementedIndicator{UTag: quota.IndicatorTag(fmt.Sprintf("bbb:sma:%s:%d", bbb.Std.Source, bbb.Std.InTimePeriod))},
			Source:                 bbb.Std.Source,
			Type:                   talib.SMA,
			InTimePeriod:           bbb.Std.InTimePeriod,
		}
		basis, ok := c.Get(quota.Source(sma.Tag()))
		if !ok {
			if !sma.Add(&quote, c) {
				return false
			}
			basis, ok = c.Get(quota.Source(sma.Tag()))
			if !ok {
				return false
			}
		}

		upper := basis + deviation
		lower := basis - deviation
		bbr, ok := c.Get(bbb.Std.Source)
		if !ok {
			return false
		}
		bbr = (bbr - lower) / (upper - lower)
		candle.AddIndicator(bbb.Tag(), bbr)
		(*q)[i] = candle

		return true
	}

	if len(*q) < bbb.Std.InTimePeriod {
		return false
	}

	for _, candle := range *q {
		if !bbb.Add(q, candle) {
			return false
		}
	}

	return true
}
