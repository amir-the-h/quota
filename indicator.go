package quota

// Indicator indicates how an UnimplementedIndicator should be implemented.
type Indicator interface {
	Add(q *Quota, c *Candle) (ok bool)
	Is(tag IndicatorTag) bool
	Tag() IndicatorTag
}

// IndicatorTag specifies each UnimplementedIndicator for others.
type IndicatorTag string

// UnimplementedIndicator adds functionality of indicator tags.
type UnimplementedIndicator struct {
	tag IndicatorTag `mapstructure:"tag"`
}

// Is determine provided tag belongs to this UnimplementedIndicator or not.
func (i *UnimplementedIndicator) Is(tag IndicatorTag) bool {
	return i.tag == tag
}

// Tag will return the UnimplementedIndicator's tag.
func (i *UnimplementedIndicator) Tag() IndicatorTag {
	return i.tag
}
