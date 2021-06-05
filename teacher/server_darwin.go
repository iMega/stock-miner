// +build darwin

package teacher

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/daemon/configuring/env"
	"github.com/shopspring/decimal"
)

type teacher struct {
	Data       map[string][][]string
	Cursor     map[string]int
	Operations map[string][]sdk.Operation
}

var (
	dataMutex      = sync.RWMutex{}
	cursorMutex    = sync.RWMutex{}
	operationMutex = sync.RWMutex{}
)

func New(mux *http.ServeMux) {
	rand.Seed(time.Now().UnixNano())
	dir, _ := env.Read("FIXTURE_PATH")
	files, err := filesInDir(dir)
	if err != nil {
		return
	}

	t := &teacher{
		Data:       make(map[string][][]string),
		Cursor:     make(map[string]int),
		Operations: make(map[string][]sdk.Operation),
	}

	for _, filename := range files {
		base := filepath.Base(filename)
		ext := filepath.Ext(filename)
		ticker := base[0 : len(base)-len(ext)]
		file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			return
		}

		recs, err := csv.NewReader(file).ReadAll()
		if err != nil {
			return
		}

		t.Data[instrumentByTicker(ticker).FIGI] = recs
		t.Cursor[instrumentByTicker(ticker).FIGI] = 0
		t.Operations[instrumentByTicker(ticker).FIGI] = []sdk.Operation{}
	}

	mux.HandleFunc("/sandbox/market/stocks", t.marketStocks)
	mux.HandleFunc("/sandbox/market/orderbook", trottlehandler(t.marketOrderBook))
	mux.HandleFunc("/sandbox/orders/market-order", trottlehandler(t.marketOrder))
	mux.HandleFunc("/sandbox/operations", trottlehandler(t.operations))
}

func trottlehandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		min := 100
		max := 1000
		<-time.After(time.Duration(rand.Intn(max-min+1)+min) * time.Millisecond)

		next.ServeHTTP(w, r)
	})
}

func filesInDir(dir string) ([]string, error) {
	var files []string

	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		err := filepath.Walk(
			dir,
			func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			},
		)
		if err != nil {
			return nil, err
		}
	case mode.IsRegular():
		files = append(files, dir)
	}

	return files, nil
}

func instrument(FIGI string) sdk.Instrument {
	var ticker string
	switch FIGI {
	case "BBG000B9XRY4":
		ticker = "AAPL"
	case "BBG000BVPV84":
		ticker = "AMZN"
	case "BBG000MM2P62":
		ticker = "FB"
	case "BBG004730RP0":
		ticker = "GAZP"
	case "BBG009S3NB30":
		ticker = "GOOG"
	case "BBG000LNHHJ9":
		ticker = "KMAZ"
	case "BBG000CL9VN6":
		ticker = "NFLX"
	case "BBG000BQY289":
		ticker = "PDCO"
	case "BBG000N9MNX3":
		ticker = "TSLA"
	}

	return instrumentByTicker(ticker)
}

func instrumentByTicker(ticker string) sdk.Instrument {
	def := sdk.Instrument{
		MinPriceIncrement: 0.01,
		Lot:               1,
		Currency:          sdk.USD,
		Type:              sdk.InstrumentTypeStock,
	}

	switch ticker {
	case "AAPL":
		def.FIGI = "BBG000B9XRY4"
		def.Ticker = "AAPL"
		def.Name = "Apple"
	case "AMZN":
		def.FIGI = "BBG000BVPV84"
		def.Ticker = "AMZN"
		def.Name = "Amazon.com"
	case "FB":
		def.FIGI = "BBG000MM2P62"
		def.Ticker = "FB"
		def.Name = "Facebook"
	case "GAZP":
		def.FIGI = "BBG004730RP0"
		def.Ticker = "GAZP"
		def.Name = "Газпром"
		def.Lot = 10
		def.Currency = sdk.RUB
	case "GOOG":
		def.FIGI = "BBG009S3NB30"
		def.Ticker = "GOOG"
		def.Name = "Alphabet Class C"
	case "KMAZ":
		def.FIGI = "BBG000LNHHJ9"
		def.Ticker = "KMAZ"
		def.Name = "КАМАЗ"
		def.Lot = 10
		def.Currency = sdk.RUB
	case "NFLX":
		def.FIGI = "BBG000CL9VN6"
		def.Ticker = "NFLX"
		def.Name = "Netflix"
	case "PDCO":
		def.FIGI = "BBG000BQY289"
		def.Ticker = "PDCO"
		def.Name = "Patterson"
	case "TSLA":
		def.FIGI = "BBG000N9MNX3"
		def.Ticker = "TSLA"
		def.Name = "Tesla Motors"
	}

	return def
}

func (t *teacher) marketStocks(w http.ResponseWriter, r *http.Request) {
	var instruments []sdk.Instrument

	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", http.StatusInternalServerError)
		return
	}

	for k := range t.Data {
		instruments = append(instruments, instrument(k))
	}

	data := map[string]interface{}{
		"status": "Ok",
		"payload": map[string]interface{}{
			"total":       len(instruments),
			"instruments": instruments,
		},
	}
	b, _ := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (t *teacher) marketOrderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", http.StatusInternalServerError)
		return
	}

	v, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	figi := v.Get("figi")

	cursorMutex.Lock()
	t.Cursor[figi]++
	cursorMutex.Unlock()

	cursorMutex.RLock()
	cursor, ok := t.Cursor[figi]
	cursorMutex.RUnlock()
	if !ok {
		http.Error(w, "end list", http.StatusInternalServerError)
		return
	}

	dataMutex.RLock()
	if cursor >= len(t.Data[figi]) {
		http.Error(w, "end list", http.StatusInternalServerError)
		return
	}
	dataMutex.RUnlock()

	dataMutex.RLock()
	rawPrice := t.Data[figi][cursor][1]
	dataMutex.RUnlock()

	d, err := decimal.NewFromString(rawPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	price, _ := d.Float64()

	data := map[string]interface{}{
		"status": "Ok",
		"payload": sdk.RestOrderBook{
			FIGI:        instrument(figi).FIGI,
			Depth:       1,
			TradeStatus: sdk.NormalTrading,
			LastPrice:   price,
		},
	}
	b, _ := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (t *teacher) marketOrder(w http.ResponseWriter, r *http.Request) {
	var requestOrder struct {
		Lots      int    `json:"lots"`
		Operation string `json:"operation"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "unsupported method", http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	if err := json.Unmarshal(b, &requestOrder); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	figi := v.Get("figi")
	cursorMutex.RLock()
	cursor, ok := t.Cursor[figi]
	cursorMutex.RUnlock()
	if !ok {
		http.Error(w, "not exists cursor: "+figi, http.StatusInternalServerError)
		return
	}

	dataMutex.RLock()
	rawPrice := t.Data[figi][cursor][1]
	dataMutex.RUnlock()

	d, err := decimal.NewFromString(rawPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	price, _ := d.Float64()

	operationMutex.Lock()
	t.Operations[figi] = append(t.Operations[figi], sdk.Operation{
		ID:               figi + "-" + strconv.Itoa(cursor),
		Status:           sdk.OperationStatusDone,
		Currency:         instrument(figi).Currency,
		Payment:          price,
		Price:            price,
		Quantity:         requestOrder.Lots,
		QuantityExecuted: requestOrder.Lots,
		Commission: sdk.MoneyAmount{
			Currency: instrument(figi).Currency,
			Value:    0,
		},
		FIGI:           figi,
		InstrumentType: sdk.InstrumentTypeStock,
		IsMarginCall:   false,
		DateTime:       time.Now(),
		OperationType:  sdk.OperationType(requestOrder.Operation),
		Trades: []sdk.Trade{
			{
				ID:       "123",
				Price:    price,
				Quantity: requestOrder.Lots,
			},
		},
	})
	operationMutex.Unlock()

	data := map[string]interface{}{
		"status": "Ok",
		"payload": sdk.PlacedOrder{
			ID:            figi + "-" + strconv.Itoa(cursor),
			Operation:     sdk.OperationType(requestOrder.Operation),
			Status:        sdk.OrderStatusFill,
			RequestedLots: requestOrder.Lots,
			ExecutedLots:  requestOrder.Lots,
		},
	}
	b, _ = json.Marshal(data)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (t *teacher) operations(w http.ResponseWriter, r *http.Request) {
	var result []sdk.Operation

	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", http.StatusInternalServerError)
		return
	}

	// if rand.Intn(2) > 0 {
	// 	data := map[string]interface{}{
	// 		"status": "Ok",
	// 		"payload": map[string]interface{}{
	// 			"operations": result,
	// 		},
	// 	}
	// 	b, _ := json.Marshal(data)

	// 	w.Header().Add("Content-Type", "application/json")
	// 	w.Write(b)

	// 	return
	// }

	v, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	figi := v.Get("figi")

	operationMutex.RLock()
	for _, op := range t.Operations[figi] {
		result = append(result, op)
	}
	operationMutex.RUnlock()

	data := map[string]interface{}{
		"status": "Ok",
		"payload": map[string]interface{}{
			"operations": result,
		},
	}
	b, _ := json.Marshal(data)

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
