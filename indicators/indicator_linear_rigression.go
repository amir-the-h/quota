package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

// LinearRegression is the lr indicator.
type LinearRegression struct {
	Tag          quota.IndicatorTag `mapstructure:"tag"`
	Source       quota.Source       `mapstructure:"source"`
	InTimePeriod int                `mapstructure:"period"`
}

// Add will calculate and add LinearRegression into the candle or whole quota.
func (lr *LinearRegression) Add(q *quota.Quota, c *quota.Candle) bool {
	qu, valid := InTimePeriodValidator(lr.InTimePeriod, q, c)
	if !valid {
		return false
	}
	if c != nil {
		values := talib.LinearReg(qu.Get(lr.Source), lr.InTimePeriod)
		c.AddIndicator(lr.Tag, values[len(values)-1])

		return true
	}

	values := talib.LinearReg(q.Get(lr.Source), lr.InTimePeriod)
	err := q.AddIndicator(lr.Tag, values)

	return err == nil
}

// Is determine provided tag belongs to this quota.Indicator or not.
func (lr *LinearRegression) Is(tag quota.IndicatorTag) bool {
	return lr.Tag == tag
}
