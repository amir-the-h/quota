package indicators

import (
	"github.com/amir-the-h/quota"
)

// InTimePeriodValidator will validate the quota length for the operator indicator.
func InTimePeriodValidator(period int, q *quota.Quota, c *quota.Candle) (*quota.Quota, bool) {
	if c != nil {
		candle, i := q.Find(c.OpenTime.Unix())
		if candle == nil {
			return q, false
		}

		startIndex := i - period
		if startIndex < 0 {
			return q, false
		}

		quota := (*q)[startIndex : i+1]
		return &quota, true
	}

	return q, len(*q) >= period
}
