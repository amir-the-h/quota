// Package quota will provide fundamental concept and utils to
// operate with financial time-series data, specially supports crypto-currencies.
package quota

import (
	"errors"
)

// EnterSignal is a signal for entering into positions
type EnterSignal struct {
	Symbol     string
	Score      float64
	Quote      float64
	TakeProfit float64
	Stoploss   float64
	Cause      string
	Candle     Candle
}

// ExitSignal is a signal for exiting from positions
type ExitSignal struct {
	Trade  *Trade
	Candle *Candle
	Cause  ExitCause
}

// ExitCause indicates why the position has been closed for
type ExitCause string

// PositionType indicates position direction
type PositionType string

// MarketType indicates the market type
type MarketType string

// TradeStatus indicates the trade status
type TradeStatus string

// TradesChannel to pass Trade through it
type TradesChannel chan *Trade

// EnterChannel to pass EnterSignal through it
type EnterChannel chan EnterSignal

// ExitChannel to pass ExitSignal through it
type ExitChannel chan ExitSignal

// CandleChannel to pass Candle through it
type CandleChannel chan *Candle

// CandleError will occur on Candle's operations
type CandleError error

// SourceError will occur on Source's operations
type SourceError error

// TradeError will occur on Trade's operations
type TradeError error

const (
	SourceOpen   = Source("open")
	SourceHigh   = Source("high")
	SourceLow    = Source("low")
	SourceClose  = Source("close")
	SourceVolume = Source("volume")

	SourceOpenHigh  = Source("oh2")
	SourceOpenLow   = Source("ol2")
	SourceOpenClose = Source("oc2")
	SourceHighLow   = Source("hl2")
	SourceHighClose = Source("hc2")
	SourceLowClose  = Source("lc2")

	SourceOpenHighLow   = Source("ohl3")
	SourceOpenHighClose = Source("ohc3")
	SourceOpenLowClose  = Source("olc3")
	SourceHighLowClose  = Source("hlc3")

	SourceOpenHighLowClose = Source("ohlc4")

	PositionBuy  = PositionType("Buy")
	PositionSell = PositionType("Sell")

	TradeStatusOpen  = TradeStatus("Open")
	TradeStatusClose = TradeStatus("Close")

	ExitCauseStopLossTriggered   = ExitCause("Stop loss")
	ExitCauseTakeProfitTriggered = ExitCause("Take profit")
	ExitCauseMarket              = ExitCause("Market")

	MarketSpot     = MarketType("Spot")
	MarketFutures  = MarketType("Futures")
	MarketDelivery = MarketType("Delivery")
)

var (
	ErrInvalidCandleData = errors.New("invalid data provided for candle").(CandleError)
	ErrNotEnoughCandles  = errors.New("not enough candles to operate").(CandleError)
)
