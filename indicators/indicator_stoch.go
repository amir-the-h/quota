package indicators

import (
	"github.com/amir-the-h/quota"
	"github.com/markcheno/go-talib"
)

// Stoch is the stochastic indicator.
type Stoch struct {
	Tag           quota.IndicatorTag `mapstructure:"tag"`
	KTag          quota.IndicatorTag `mapstructure:"kTag"`
	DTag          quota.IndicatorTag `mapstructure:"dTag"`
	InFastKPeriod int                `mapstructure:"kLength"`
	InSlowKPeriod int                `mapstructure:"kSmoothing"`
	InKMaType     talib.MaType       `mapstructure:"kMaType"`
	InSlowDPeriod int                `mapstructure:"dSmoothing"`
	InDMaType     talib.MaType       `mapstructure:"dMaType"`
}

// Add will calculate and add Stoch into the candle or whole quota.
func (s *Stoch) Add(q *quota.Quota, c *quota.Candle) bool {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return false
		}

		quote := (*q)[:i+1]

		k, d := talib.Stoch(quote.Get(quota.SourceHigh), quote.Get(quota.SourceLow), quote.Get(quota.SourceClose), s.InFastKPeriod, s.InSlowKPeriod, s.InKMaType, s.InSlowDPeriod, s.InDMaType)
		c.AddIndicator(s.KTag, k[len(k)-1])
		c.AddIndicator(s.DTag, d[len(d)-1])

		return true
	}

	k, d := talib.Stoch(q.Get(quota.SourceHigh), q.Get(quota.SourceLow), q.Get(quota.SourceClose), s.InFastKPeriod, s.InSlowKPeriod, s.InKMaType, s.InSlowDPeriod, s.InDMaType)
	err := q.AddIndicator(s.KTag, k)
	if err != nil {
		return false
	}
	err = q.AddIndicator(s.DTag, d)

	return err == nil
}

// Is determine provided tag belongs to this quota.Indicator or not.
func (s *Stoch) Is(tag quota.IndicatorTag) bool {
	return s.Tag == tag
}
