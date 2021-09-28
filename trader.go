package quota

// Trader determines how a trader should be implemented.
type Trader interface {
	Open(coin, base string, position PositionType, quote, entry, sl, tp float64, openCandle *Candle) *Trade
	Close(id string, exit float64, closeCandle *Candle)
	Start() TradeError
	EntryWatcher()
	ExitWatcher()
	CloseWatcher()
}
