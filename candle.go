package quota

import (
	"github.com/amir-the-h/okex"
	"time"
)

// Candle is the main structure which contains a group of useful candlestick data.
type Candle struct {
	Open       float64                             `json:"open"`
	High       float64                             `json:"high"`
	Low        float64                             `json:"low"`
	Close      float64                             `json:"close"`
	Volume     float64                             `json:"volume"`
	Score      float64                             `json:"score"`
	Symbol     string                              `json:"symbol"`
	BarSize    okex.BarSize             `json:"barSize"`
	Indicators map[IndicatorTag]float64 `json:"indicators"`
	OpenTime   time.Time                `json:"open_time"`
	CloseTime  time.Time                           `json:"close_time"`
	Next       *Candle                             `json:"-"`
	Previous   *Candle                             `json:"-"`
}

// NewCandle returns a pointer to a fresh candle with provided data.
func NewCandle(open, high, low, close, volume float64, symbol string, barSize okex.BarSize, openTime, closeTime time.Time, previous, next *Candle) (candle *Candle, err CandleError) {
	if high < low || high < open || high < close || low > open || low > close {
		err = ErrInvalidCandleData
		return
	}

	candle = &Candle{
		Open:       open,
		High:       high,
		Low:        low,
		Close:      close,
		Volume:     volume,
		Symbol:     symbol,
		BarSize:    barSize,
		OpenTime:   openTime,
		CloseTime:  closeTime,
		Indicators: make(map[IndicatorTag]float64),
		Previous:   previous,
		Next:       next,
	}

	return
}

// AddIndicator adds unimplementedIndicator value by the given name into the candle.
func (c *Candle) AddIndicator(tag IndicatorTag, value float64) {
	c.Indicators[tag] = value
}
