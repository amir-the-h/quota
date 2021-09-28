package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

// Ma is a general wrapper for all talib.MaType supported by talib.
type Ma struct {
	Tag          quota.IndicatorTag `mapstructure:"tag"`
	Source       quota.Source       `mapstructure:"source"`
	Type         talib.MaType       `mapstructure:"type"`
	InTimePeriod int                `mapstructure:"period"`
}

// Add will calculate and add Ma into the candle or whole quota.
func (ma *Ma) Add(q *quota.Quota, c *quota.Candle) bool {
	quote, valid := InTimePeriodValidator(ma.InTimePeriod, q, c)
	if !valid {
		return false
	}
	if c != nil {
		values := talib.Ma(quote.Get(ma.Source), ma.InTimePeriod, ma.Type)
		c.AddIndicator(ma.Tag, values[len(values)-1])

		return true
	}

	values := talib.Ma(q.Get(ma.Source), ma.InTimePeriod, ma.Type)
	err := q.AddIndicator(ma.Tag, values)

	return err == nil
}

// Is determine provided tag belongs to this quota.Indicator or not.
func (ma *Ma) Is(tag quota.IndicatorTag) bool {
	return ma.Tag == tag
}
