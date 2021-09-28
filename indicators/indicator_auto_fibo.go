package indicators

import (
	"fmt"
	"github.com/amir-the-h/quota"
	"github.com/amir-the-h/quota/utils"
)

// AutoFibo is basically fibonacci series applied on the best spot to find some resistance/support lines.
type AutoFibo struct {
	quota.UnimplementedIndicator
	Ratios    []float64 `mapstructure:"ratios"`
	Deviation float64   `mapstructure:"deviation"`
	Depth     int       `mapstructure:"depth"`
}

// Add will calculate and add AutoFibo into the candle or whole quota.
func (af *AutoFibo) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		quote := (*q)[:i+1]
		fibos := utils.AutoFiboRectracement(quote.Get(quota.SourceHigh), quote.Get(quota.SourceLow), quote.Get(quota.SourceClose), af.Ratios, af.Depth, af.Deviation)
		for ratio, fibo := range fibos[len(fibos)-1] {
			c.AddIndicator(quota.IndicatorTag(fmt.Sprintf("%s:%.2f", af.Tag(), ratio)), fibo)
		}
		(*q)[i] = c

		return true
	}

	for _, candle := range *q {
		if !af.Add(q, candle) {
			return false
		}
	}

	return true
}
