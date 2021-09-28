package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

type LinearRegression struct {
	quota.UnimplementedIndicator
	Source       quota.Source `mapstructure:"source"`
	InTimePeriod int          `mapstructure:"period"`
}

func (lr *LinearRegression) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - lr.InTimePeriod
		if startIndex < 0 {
			return false
		}

		quote := (*q)[startIndex : i+1]

		values := talib.LinearReg(quote.Get(lr.Source), lr.InTimePeriod)
		c.AddIndicator(lr.Tag(), values[len(values)-1])

		return true
	}

	if len(*q) < lr.InTimePeriod {
		return false
	}

	values := talib.LinearReg(q.Get(lr.Source), lr.InTimePeriod)
	err := q.AddIndicator(lr.Tag(), values)

	return err == nil
}
