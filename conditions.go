package quota

// IsAbove checks if candle is above of the source.
//
// O,H,L,C > source
func (c *Candle) IsAbove(source Source) bool {
	value, _ := c.Get(source)

	return c.Open > value && c.High > value && c.Low > value && c.Close > value
}

// IsBelow checks if candle is below of the source.
//
// O,H,L,C < source
func (c *Candle) IsBelow(source Source) bool {
	value, _ := c.Get(source)

	return c.Open < value && c.High < value && c.Low < value && c.Close < value
}

// IsMiddle checks if source is passed through middle of the candle.
//
// O >= source
// H > source
// L < source
// C <= source
func (c *Candle) IsMiddle(source Source) bool {
	value, _ := c.Get(source)

	return c.Open >= value && c.High > value && c.Low < value && c.Close <= value
}

// TouchedDown checks if candle closed above the source but the Low shadow touched the source.
//
// O,H,C > source
// L <= source
func (c *Candle) TouchedDown(source Source) bool {
	value, _ := c.Get(source)

	return c.Open > value && c.High > value && c.Low <= value && c.Close > value
}

// TouchedUp checks if candle close above the source but the High shadow touched the source.
//
// H >= source
// O,L,C < source
func (c *Candle) TouchedUp(source Source) bool {
	value, _ := c.Get(source)

	return c.Open < value && c.High >= value && c.Low < value && c.Close < value
}

// CrossedOver checks if fast source crossed over the slow source.
//
// fastSource > slowSource
// prevFastSource <= prevSlowSource
func (c *Candle) CrossedOver(fastSource, slowSource Source) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastValue, ok := c.Get(fastSource)
	if !ok {
		return false
	}
	slowValue, ok := c.Get(slowSource)
	if !ok {
		return false
	}
	previousFastValue, ok := previousCandle.Get(fastSource)
	if !ok {
		return false
	}
	previousSlowValue, ok := previousCandle.Get(slowSource)
	if !ok {
		return false
	}

	return fastValue > slowValue && previousFastValue <= previousSlowValue
}

// CrossedUnder checks if fast source crossed under the slow source.
//
// fastSource < slowSource
// prevFastSource >= prevSlowSource
func (c *Candle) CrossedUnder(fastSource, slowSource Source) bool {
	previousCandle := c.Previous
	if previousCandle == nil {
		return false
	}

	fastValue, ok := c.Get(fastSource)
	if !ok {
		return false
	}
	slowValue, ok := c.Get(slowSource)
	if !ok {
		return false
	}
	previousFastValue, ok := previousCandle.Get(fastSource)
	if !ok {
		return false
	}
	previousSlowValue, ok := previousCandle.Get(slowSource)
	if !ok {
		return false
	}

	return fastValue < slowValue && previousFastValue >= previousSlowValue
}

// IsBearish checks if candle is bearish.
//
// O < C,
func (c *Candle) IsBearish() bool {
	return c.Open < c.Close
}

// IsBullish checks if candle is bullish.
//
// O > C,
func (c *Candle) IsBullish() bool {
	return c.Open > c.Close
}
