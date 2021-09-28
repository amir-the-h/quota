package indicators

import (
	"fmt"
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

// BollingerBandsB is an advanced bollinger bands indicator.
type BollingerBandsB struct {
	quota.UnimplementedIndicator
	Std StandardDeviation `mapstructure:"standardDeviation"`
}

// Add will calculate and add BollingerBandsB into the candle or whole quota.
func (bbb *BollingerBandsB) Add(q *quota.Quota, c *quota.Candle) bool {
	qu, valid := InTimePeriodValidator(bbb.Std.InTimePeriod, q, c)
	if !valid {
		return false
	}
	if c != nil {
		deviation, ok := c.Get(quota.Source(bbb.Std.Tag()))
		if !ok {
			if !bbb.Std.Add(qu, c) {
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
			if !sma.Add(qu, c) {
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
		c.AddIndicator(bbb.Tag(), bbr)

		return true
	}

	for _, candle := range *q {
		if !bbb.Add(q, candle) {
			return false
		}
	}

	return true
}
