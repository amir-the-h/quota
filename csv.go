package quota

import (
	"encoding/csv"
	"fmt"
	"github.com/amir-the-h/okex"
	"os"
	"strconv"
	"strings"
	"time"
)

// Csv returns csv row of the candle.
func (c *Candle) Csv(indicators ...IndicatorTag) (csv string) {
	// first basic records
	csv = fmt.Sprintf("%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%s,%s", c.Symbol, c.BarSize, c.Open, c.High, c.Low, c.Close, c.Volume, c.Score, c.OpenTime.Local().Format(time.RFC3339), c.CloseTime.Local().Format(time.RFC3339))

	// get indicators values too
	if len(indicators) > 0 {
		for _, indicator := range indicators {
			csv += fmt.Sprintf(",%.2f", c.Indicators[indicator])
		}
	} else {
		for _, indicator := range c.Indicators {
			csv += fmt.Sprintf(",%.2f", indicator)
		}
	}

	return
}

// Csv returns csv formatted string of whole quote.
func (q *Quota) Csv(indicators ...IndicatorTag) (csv string) {
	if len(*q) == 0 {
		return
	}
	// fix the headers
	headers := []string{"Symbol", "BarSize", "Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time"}
	var indicatorTags []IndicatorTag
	if len(indicators) > 0 {
		indicatorTags = indicators
	} else {
		indicatorTags = q.IndicatorTags()
	}
	for _, indicatorTag := range indicatorTags {
		headers = append(headers, string(indicatorTag))
	}
	csv = strings.Join(headers, ",") + "\n"
	// get each candle csv value
	for _, candle := range *q {
		// and also add the indicators as well
		csv += fmt.Sprintln(candle.Csv(indicatorTags...))
	}

	return
}

// WriteToCsv writes down whole quote into a csv file.
func (q *Quota) WriteToCsv(filename string, indicators ...IndicatorTag) error {
	if len(*q) == 0 {
		return ErrNotEnoughCandles
	}

	// need our file
	if filename == "" {
		filename = fmt.Sprintf("%s:%s.csv", q.Symbol(), q.BarSize())
	}

	// open or create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	// truncate the file
	if err = file.Truncate(0); err != nil {
		return err
	}

	_, err = file.Write([]byte(q.Csv(indicators...)))
	return err
}

// NewQuoteFromCsv reads quote from csv file.
func NewQuoteFromCsv(filename string, symbol string, barSize okex.BarSize) (*Quota, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil, err
	}
	headers := csvLines[0]
	csvLines = csvLines[1:]
	indexMap := make(map[string]int)
	indicatorsMap := make(map[string]int)
	for i, header := range headers {
		switch header {
		case "Symbol", "BarSize", "Open", "High", "Low", "Close", "Volume", "Score", "Open time", "Close time":
			indexMap[header] = i
		default:
			indicatorsMap[header] = i
		}
	}
	quote := &Quota{}
	for _, line := range csvLines {
		openPrice, _ := strconv.ParseFloat(line[indexMap["Open"]], 64)
		highPrice, _ := strconv.ParseFloat(line[indexMap["High"]], 64)
		lowPrice, _ := strconv.ParseFloat(line[indexMap["Low"]], 64)
		closePrice, _ := strconv.ParseFloat(line[indexMap["Close"]], 64)
		volume, _ := strconv.ParseFloat(line[indexMap["Volume"]], 64)
		openTime, _ := time.Parse(time.RFC3339, line[indexMap["Open time"]])
		closeTime, _ := time.Parse(time.RFC3339, line[indexMap["Close time"]])
		candle, err := NewCandle(openPrice, highPrice, lowPrice, closePrice, volume, symbol, barSize, openTime, closeTime, nil, nil)
		if err != nil {
			return nil, err
		}
		candle.Score, _ = strconv.ParseFloat(line[indexMap["Score"]], 64)

		for indicator, index := range indicatorsMap {
			candle.Indicators[IndicatorTag(indicator)], _ = strconv.ParseFloat(line[index], 64)
		}

		*quote = append(*quote, candle)
	}

	quote.Sort()
	return quote, nil
}
