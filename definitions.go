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
	// SourceOpen determines open Source
	SourceOpen = Source("open")
	// SourceHigh determines open Source
	SourceHigh = Source("high")
	// SourceLow determines open Source
	SourceLow = Source("low")
	// SourceClose determines open Source
	SourceClose = Source("close")
	// SourceVolume determines open Source
	SourceVolume = Source("volume")

	// SourceOpenHigh determines oh2 Source
	SourceOpenHigh = Source("oh2")
	// SourceOpenLow determines ol2 Source
	SourceOpenLow = Source("ol2")
	// SourceOpenClose determines oc2 Source
	SourceOpenClose = Source("oc2")
	// SourceHighLow determines hl2 Source
	SourceHighLow = Source("hl2")
	// SourceHighClose determines hc2 Source
	SourceHighClose = Source("hc2")
	// SourceLowClose determines lc2 Source
	SourceLowClose = Source("lc2")

	// SourceOpenHighLow determines ohl3 Source
	SourceOpenHighLow = Source("ohl3")
	// SourceOpenHighClose determines ohc3 Source
	SourceOpenHighClose = Source("ohc3")
	// SourceOpenLowClose determines olc3 Source
	SourceOpenLowClose = Source("olc3")
	// SourceHighLowClose determines hlc3 Source
	SourceHighLowClose = Source("hlc3")

	// SourceOpenHighLowClose determines ohlc4 Source
	SourceOpenHighLowClose = Source("ohlc4")

	// PositionBuy determines buy PositionType
	PositionBuy = PositionType("Buy")
	// PositionSell determines sell PositionType
	PositionSell = PositionType("Sell")

	// TradeStatusOpen determines open TradeStatus
	TradeStatusOpen = TradeStatus("Open")
	// TradeStatusClose determines close TradeStatus
	TradeStatusClose = TradeStatus("Close")

	// ExitCauseStopLossTriggered determines stop loss triggered on trade.
	ExitCauseStopLossTriggered = ExitCause("Stop loss")
	// ExitCauseTakeProfitTriggered determines take profit triggered on trade.
	ExitCauseTakeProfitTriggered = ExitCause("Take profit")
	// ExitCauseMarket determines trade closed by market.
	ExitCauseMarket = ExitCause("Market")
)

var (
	// ErrInvalidCandleData occurs on Candle operations.
	ErrInvalidCandleData = errors.New("invalid data provided for candle").(CandleError)
	// ErrNotEnoughCandles occurs when there is not enough Candle in Quota to operate.
	ErrNotEnoughCandles = errors.New("not enough candles to operate").(CandleError)
)
