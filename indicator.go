package quota

// Indicator indicates how an UnimplementedIndicator should be implemented.
type Indicator interface {
	Add(q *Quota, c *Candle) (ok bool)
	Is(tag IndicatorTag) bool
}

// IndicatorTag specifies each UnimplementedIndicator for others.
type IndicatorTag string
