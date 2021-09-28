package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/amir-the-h/quota/utils"
)

type Hp struct {
	quota.UnimplementedIndicator
	Source quota.Source `mapstructure:"source"`
	Lambda float64      `mapstructure:"lambda"`
	Length int          `mapstructure:"length"`
}

func (hp *Hp) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		startIndex := i - hp.Length
		if startIndex < 0 {
			return false
		}

		quote := (*q)[startIndex : i+1]

		values := utils.HPFilter(quote.Get(hp.Source), hp.Lambda)
		c.AddIndicator(hp.Tag(), values[len(values)-1])

		return true
	}

	if len(*q) < hp.Length {
		return false
	}

	for _, candle := range *q {
		if !hp.Add(q, candle) {
			return false
		}
	}

	return true
}
