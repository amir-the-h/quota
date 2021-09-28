package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

// StandardDeviation represents  standard deviation indicator.
type StandardDeviation struct {
	quota.UnimplementedIndicator
	Source       quota.Source `mapstructure:"source"`
	InTimePeriod int          `mapstructure:"period"`
	Deviation    float64      `mapstructure:"deviation"`
}

// Add will calculate and add StandardDeviation into the candle or whole quota.
func (sd *StandardDeviation) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - sd.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := (*q)[startIndex : i+1]

		values := talib.StdDev(quote.Get(sd.Source), sd.InTimePeriod, sd.Deviation)
		c.AddIndicator(sd.Tag(), values[len(values)-1])

		return true
	}

	if len(*q) < sd.InTimePeriod {
		return false
	}

	values := talib.StdDev(q.Get(sd.Source), sd.InTimePeriod, sd.Deviation)
	err := q.AddIndicator(sd.Tag(), values)

	return err == nil
}
