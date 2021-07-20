package yahooprovider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/imega/stock-miner/domain"
	"github.com/imega/stock-miner/httpwareclient"
)

var errGettingRange = errors.New("failed getting stock item range")

func (p *pricer) Range(
	ctx context.Context,
	s domain.StockItem,
) (domain.StockItemRange, error) {
	data := &responseRange{}
	period1 := strconv.Itoa(int(time.Now().AddDate(0, 0, -80).Unix()))
	result := domain.StockItemRange{}

	req := &httpwareclient.SendIn{
		Method:   http.MethodGet,
		URL:      p.URL + "/v8/finance/chart/" + s.Ticker + "?period2=9999999999&interval=1d&period1=" + period1,
		BodyRecv: data,
		Coder:    httpwareclient.GetCoder(httpwareclient.JSON),
	}
	if err := httpwareclient.Send(ctx, req); err != nil {
		return result, fmt.Errorf("failed to sent request, %w", err)
	}

	if data.Chart.Error != nil {
		return result, errGettingRange
	}

	if len(data.Chart.Result) == 0 {
		return result, errGettingRange
	}

	response := data.Chart.Result[0]

	if len(response.Indicators.Quote) == 0 {
		return result, errGettingRange
	}

	quote := response.Indicators.Quote[0]

	sort.Float64s(quote.High)
	sort.Float64s(quote.Low)

	result.High = quote.High[len(quote.High)-1]
	result.Low = quote.Low[0]

	return result, nil
}

type responseRange struct {
	Chart yfChart `json:"chart"`
}

type yfChart struct {
	Error  *yfError        `json:"chart,omitempty"`
	Result []yfChartResult `json:"result,omitempty"`
}

type yfError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type yfChartResult struct {
	Indicators yfIndicators `json:"indicators"`
}

type yfIndicators struct {
	Quote []yfQuote `json:"quote"`
}

type yfQuote struct {
	High []float64 `json:"high"`
	Low  []float64 `json:"low"`
}
