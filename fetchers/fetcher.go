package fetchers

import (
	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/quota"
	"time"
)

type Fetcher interface {
	NewQuote(symbol string, barSize okex.BarSize, openTime *time.Time) (*quota.Quota, error)
	Refresh(q *quota.Quota) error
	Sync(q *quota.Quota, update quota.CandleChannel) (err error)
}
