package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/amir-the-h/quota/utils"
)

// Hp filters the given source by Hodrick-Prescott filter.
type Hp struct {
	Tag    quota.IndicatorTag `mapstructure:"tag"`
	Source quota.Source       `mapstructure:"source"`
	Lambda float64            `mapstructure:"lambda"`
	Length int                `mapstructure:"length"`
}

// Add will calculate and add Hp into the candle or whole quota.
func (hp *Hp) Add(q *quota.Quota, c *quota.Candle) bool {
	qu, valid := InTimePeriodValidator(hp.Length, q, c)
	if !valid {
		return false
	}
	if c != nil {
		values := utils.HPFilter(qu.Get(hp.Source), hp.Lambda)
		c.AddIndicator(hp.Tag, values[len(values)-1])

		return true
	}

	for _, candle := range *q {
		if !hp.Add(q, candle) {
			return false
		}
	}

	return true
}

// Is determine provided tag belongs to this quota.Indicator or not.
func (hp *Hp) Is(tag quota.IndicatorTag) bool {
	return hp.Tag == tag
}
