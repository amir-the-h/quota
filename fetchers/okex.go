package fetchers

import (
	"context"
	"errors"
	"fmt"
	"github.com/amir-the-h/okex"
	"github.com/amir-the-h/okex/api"
	publicEvents "github.com/amir-the-h/okex/events/public"
	restMarketRequests "github.com/amir-the-h/okex/requests/rest/market"
	wsPublicRequests "github.com/amir-the-h/okex/requests/ws/public"
	tradeKnife "github.com/amir-the-h/quota"
	"time"
)

// Okex is an Okay-Exchange tradeKnife.Fetcher
type Okex struct {
	api *api.Client
}

// NewOkex returns a pointer to a fresh Okex tradeKnife.Fetcher.
func NewOkex(apiKey, secretKey, passphrase string, dest ...okex.Destination) (*Okex, error) {
	d := okex.NormalServer
	if len(dest) > 0 {
		d = dest[0]
	}
	c, err := api.NewClient(context.Background(), apiKey, secretKey, passphrase, d)
	if err != nil {
		return nil, err
	}
	return &Okex{api: c}, nil
}

// NewQuote fetches quote from okex market.
func (ok *Okex) NewQuote(symbol string, barSize okex.BarSize, timestamps ...int64) (*tradeKnife.Quota, error) {
	q := &tradeKnife.Quota{}

	req := restMarketRequests.GetCandlesticks{
		InstID: symbol,
		Bar:    barSize,
	}
	if len(timestamps) > 0 {
		req.After = timestamps[0]
		if len(timestamps) > 1 {
			req.Before = timestamps[1]
		}
	}

	res, err := ok.api.Rest.Market.GetCandlesticks(req)
	if err != nil {
		return q, err
	}

	if res.Code != 0 {
		return q, fmt.Errorf("okex: %s", res.Msg)
	}

	for _, c := range res.Candles {
		candle, err := createCandleFromOkexKline(c.O, c.H, c.L, c.C, c.Vol, (time.Time)(c.TS).Unix(), symbol, barSize)
		if err != nil {
			return q, err
		}
		*q = append(*q, candle)
	}

	q.Sort()
	return q, nil
}

// Refresh fetches all candles after last candle including itself.
func (ok *Okex) Refresh(q *tradeKnife.Quota) error {
	quote := *q
	if len(*q) == 0 {
		return errors.New("won't be able to refresh an empty quote")
	}

	var (
		lastCandle   = (*q)[len(*q)-1]
		openTime     = lastCandle.OpenTime
		fetchedQuote *tradeKnife.Quota
		err          error
	)
	fetchedQuote, err = ok.NewQuote(q.Symbol(), quote.BarSize(), openTime.Unix())
	if err != nil {
		return err
	}

	q.Merge(fetchedQuote)

	return nil
}

// Sync syncs quote with latest okex kline info.
func (ok *Okex) Sync(q *tradeKnife.Quota, update tradeKnife.CandleChannel) (err error) {
	if len(*q) == 0 {
		return errors.New("won't be able to sync an empty quote")
	}
	req := wsPublicRequests.Candlesticks{
		InstID:  (*q).Symbol(),
		Channel: okex.CandleStick1m,
	}
	cCh := make(chan *publicEvents.Candlesticks)
	go func() {
		for e := range cCh {
			for _, c := range e.Candles {
				ot := time.Time(c.TS).UTC()
				ct := ot.Add(q.BarSize().Duration()).UTC()
				candle, err := q.Sync(c.O, c.H, c.L, c.C, c.Vol, ot, ct)
				if err != nil {
					return
				}
				update <- candle
			}
		}
	}()

	return ok.api.Ws.Public.Candlesticks(req, cCh)
}

func createCandleFromOkexKline(open, high, low, close, volume float64, timestamp int64, symbol string, barSize okex.BarSize) (*tradeKnife.Candle, error) {
	ot := time.Unix(timestamp, 0).UTC()
	ct := ot.Add(barSize.Duration())
	return tradeKnife.NewCandle(open, high, low, close, volume, symbol, barSize, ot, ct, nil, nil)
}
