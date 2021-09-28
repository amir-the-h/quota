package quota

// ScoreByAbove scores candles by checking if candles are above source.
func (q *Quota) ScoreByAbove(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.IsAbove(source) {
			candle.Score += score
		}
	}
}

// ScoreByTouchUp scores candles by checking if candles are touching up source.
func (q *Quota) ScoreByTouchUp(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.TouchedUp(source) {
			candle.Score += score
		}
	}
}

// ScoreByMiddle scores candles by checking if candles are middle of the source.
func (q *Quota) ScoreByMiddle(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.IsMiddle(source) {
			candle.Score += score
		}
	}
}

// ScoreByTouchDown scores candles by checking if candles are touching down source.
func (q *Quota) ScoreByTouchDown(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.TouchedDown(source) {
			candle.Score += score
		}
	}
}

// ScoreByBelow scores candles by checking if candles are below source.
func (q *Quota) ScoreByBelow(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.IsBelow(source) {
			candle.Score += score
		}
	}
}

// ScoreByCrossOver scores candles by checking sources cross over condition on each of them.
func (q *Quota) ScoreByCrossOver(fastSource, slowSource Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.CrossedOver(fastSource, slowSource) {
			candle.Score += score
		}
	}
}

// ScoreByCrossUnder scores candles by checking sources cross under condition on each of them.
func (q *Quota) ScoreByCrossUnder(fastSource, slowSource Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.CrossedUnder(fastSource, slowSource) {
			candle.Score += score
		}
	}
}

// ScoreByCross scores candles by checking sources both crosses over and under condition on each of them.
func (q *Quota) ScoreByCross(fastSource, slowSource Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.CrossedOver(fastSource, slowSource) {
			candle.Score += score
		}
		if candle.CrossedUnder(fastSource, slowSource) {
			candle.Score -= score
		}
	}
}

// ScoreByBearish scores candles by if candle is bearish
func (q *Quota) ScoreByBearish(score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.IsBearish() {
			candle.Score += score
		}
	}
}

// ScoreByBullish scores candles by checking if candles are bullish.
func (q *Quota) ScoreByBullish(score float64) {
	// loop through quote
	for _, candle := range *q {
		if candle.IsBullish() {
			candle.Score += score
		}
	}
}

// ScoreBySupportResistance scores candles by checking support/resistance line reaction.
func (q *Quota) ScoreBySupportResistance(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		checker := func(c *Candle) float64 {
			if c.IsBearish() {
				score *= -1
			}
			if c.TouchedDown(source) || c.IsMiddle(source) || c.IsMiddle(source) {
				return score
			}
			return 0
		}

		result := checker(candle)
		if result != 0 {
			candle.Score += result
		}
		if candle.Previous == nil {
			continue
		}
		prevResult := checker(candle.Previous)
		if prevResult*-1 == result {
			candle.Score += result
		}
	}
}

// ScoreByBands scores candles by checking touching and turning back into the band.
func (q *Quota) ScoreByBands(source Source, score float64) {
	// loop through quote
	for _, candle := range *q {
		checker := func(c *Candle) float64 {
			if c.IsBearish() {
				score *= -1
			}
			if c.TouchedDown(source) || c.IsMiddle(source) || c.IsMiddle(source) {
				return score
			}
			return 0
		}

		result := checker(candle)
		if result != 0 {
			candle.Score += result
		}
		if candle.Previous == nil {
			continue
		}
		prevResult := checker(candle.Previous)
		if prevResult*-1 == result {
			candle.Score += result
		}
	}
}
