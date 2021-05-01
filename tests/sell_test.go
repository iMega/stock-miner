package acceptance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/imega/stock-miner/tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Automatically buy and sell", func() {
	var (
		figi   = "BBG000BQY289"
		ticker = "PDCO"
	)
	It("set environment", func() {
		defer GinkgoRecover()

		res, err := helpers.AddCredentials(
			GraphQLUrl,
			"test",
			"http://acceptance",
			"token",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeTrue())

		helpers.SetSlot(GraphQLUrl)
		helpers.SetRulePrice(GraphQLUrl)
		items := helpers.StockItemApproved(GraphQLUrl)
		for _, item := range items {
			helpers.RemoveStockItemApproved(GraphQLUrl, string(item.Ticker))
		}
		helpers.AddStockItemApproved(GraphQLUrl, ticker, figi)
	})

	It("buy stock item", func() {
		var (
			requestCount int
			startPrice   = 100
			orderID      = 235774468340
		)
		defer GinkgoRecover()

		helpers.MockHTTPServer.AddHandler(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()

			requestCount++
			data := map[string]interface{}{
				"status": "OK",
				"payload": map[string]interface{}{
					"figi":  figi,
					"depth": 20,
					"bids": []map[string]interface{}{
						{
							"price": startPrice - requestCount,
							"qty":   1,
						},
					},
					"asks": []map[string]interface{}{
						{
							"price": startPrice - requestCount,
							"qty":   1,
						},
					},
					"tradeStatus":       sdk.NormalTrading,
					"minPriceIncrement": 0,
					"lastPrice":         startPrice - requestCount,
					"closePrice":        0,
					"limitUp":           0,
					"limitDown":         0,
					"faceValue":         0,
				},
			}

			if "/orders/market-order" == r.URL.Path {
				orderID = orderID + requestCount
				var requestOrderAdd struct {
					Lots      int    `json:"lots"`
					Operation string `json:"operation"`
				}

				b, _ := ioutil.ReadAll(r.Body)
				json.Unmarshal(b, &requestOrderAdd)

				data = map[string]interface{}{
					"status": "Fail",
					"payload": map[string]string{
						"message": "string",
						"code":    "string",
					},
				}
				if requestOrderAdd.Operation == string(sdk.BUY) && requestOrderAdd.Lots == 1 {
					data = map[string]interface{}{
						"trackingId": "dbb781ba4e984bd9",
						"status":     "Ok",
						"payload": map[string]interface{}{
							"orderId":       strconv.Itoa(orderID),
							"operation":     "Buy",
							"status":        "Fill",
							"executedLots":  1,
							"requestedLots": 1,
							"commission": map[string]interface{}{
								"currency": "USD",
								"value":    0,
							},
						},
					}
				}
			}

			if "/operations" == r.URL.Path {
				v, err := url.ParseQuery(r.URL.RawQuery)
				Expect(err).NotTo(HaveOccurred())

				actual, err := time.Parse("2006-01-02T15:04:05-07:00", v.Get("from"))
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).Should(BeTemporally("~", time.Now(), 2*time.Second))

				actual, err = time.Parse("2006-01-02T15:04:05-07:00", v.Get("to"))
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).Should(BeTemporally("~", time.Now(), 2*time.Minute))

				data = map[string]interface{}{
					"trackingId": "4f48d98e8040c23a",
					"status":     "Ok",
					"payload": map[string]interface{}{
						"operations": []interface{}{
							map[string]interface{}{
								"operationType":    "Buy",
								"date":             "2021-03-01T23:39:29.507+03:00",
								"isMarginCall":     false,
								"instrumentType":   sdk.InstrumentTypeStock,
								"figi":             figi,
								"quantity":         1,
								"quantityExecuted": 1,
								"price":            95,
								"payment":          -95,
								"currency":         sdk.USD,
								"status":           sdk.OperationStatusDone,
								"id":               strconv.Itoa(orderID),
								"commission": map[string]interface{}{
									"currency": "USD",
									"value":    -0.5,
								},
								"trades": []interface{}{
									map[string]interface{}{
										"tradeId":  "3535068930",
										"date":     "2021-03-01T23:39:29.507+03:00",
										"quantity": 1,
										"price":    95,
									},
								},
							},
						},
					},
				}
			}

			b, _ := json.Marshal(data)

			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		})

		helpers.GlobalMiningStart(GraphQLUrl)
		<-time.After(13 * time.Second)
		helpers.GlobalMiningStop(GraphQLUrl)
	})

	It("sell stock item", func() {
		var (
			requestCount int
			startPrice   = 94
			// orderID      = 235774468340
		)
		helpers.MockHTTPServer.AddHandler(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()

			fmt.Printf("========== %s + %d\n", r.URL.Path, startPrice+requestCount)
			b, _ := ioutil.ReadAll(r.Body)
			err := r.Body.Close()
			Expect(err).NotTo(HaveOccurred())
			fmt.Printf("========== %q\n", b)

			requestCount++
			data := map[string]interface{}{
				"status": "OK",
				"payload": map[string]interface{}{
					"figi":  figi,
					"depth": 20,
					"bids": []map[string]interface{}{
						{
							"price": startPrice + requestCount,
							"qty":   1,
						},
					},
					"asks": []map[string]interface{}{
						{
							"price": startPrice + requestCount,
							"qty":   1,
						},
					},
					"tradeStatus":       sdk.NormalTrading,
					"minPriceIncrement": 0,
					"lastPrice":         startPrice + requestCount,
					"closePrice":        0,
					"limitUp":           0,
					"limitDown":         0,
					"faceValue":         0,
				},
			}

			b, _ = json.Marshal(data)

			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		})

		helpers.GlobalMiningStart(GraphQLUrl)
		<-time.After(13 * time.Second)
		helpers.GlobalMiningStop(GraphQLUrl)
	})
})
