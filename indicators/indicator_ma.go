package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

type Ma struct {
	quota.UnimplementedIndicator
	Source       quota.Source `mapstructure:"source"`
	Type         talib.MaType `mapstructure:"type"`
	InTimePeriod int          `mapstructure:"period"`
}

func (ma *Ma) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - ma.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := (*q)[startIndex : i+1]

		values := talib.Ma(quote.Get(ma.Source), ma.InTimePeriod, ma.Type)
		c.AddIndicator(ma.Tag(), values[len(values)-1])

		return true
	}

	if len(*q) < ma.InTimePeriod {
		return false
	}

	values := talib.Ma(q.Get(ma.Source), ma.InTimePeriod, ma.Type)
	err := q.AddIndicator(ma.Tag(), values)
	if err != nil {
		return false
	}

	return true
}
