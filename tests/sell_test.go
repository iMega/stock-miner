package acceptance

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/imega/stock-miner/tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Automatically buy and sell", func() {
	var (
		figi   = "BBG000BQY289"
		ticker = "PDCO"
		ctx    = context.Background()
		client = graphql.NewClient(GraphQLUrl, helpers.GetHTTPClient())
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
		helpers.EnableStockItemApproved(GraphQLUrl)
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

			if strings.Contains(r.URL.Path, "chart") {
				data = map[string]interface{}{
					"chart": map[string]interface{}{
						"result": []map[string]interface{}{
							{
								"indicators": map[string]interface{}{
									"quote": []map[string][]float64{
										{
											"high": []float64{149.75},
											"low":  []float64{47.08},
										},
									},
								},
							},
						},
					},
				}
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

				v, err := url.ParseQuery(r.URL.RawQuery)
				Expect(err).NotTo(HaveOccurred())

				if v.Get("figi") == figi && requestOrderAdd.Operation == string(sdk.BUY) && requestOrderAdd.Lots == 1 {
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
				Expect(actual).Should(BeTemporally("~", time.Now(), 5*time.Minute))

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
			startPrice   = 96
			orderID      = 774468340
		)
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

			if "/orders/market-order" == r.URL.Path {
				var requestOrder struct {
					Lots      int    `json:"lots"`
					Operation string `json:"operation"`
				}

				b, _ := ioutil.ReadAll(r.Body)
				json.Unmarshal(b, &requestOrder)
				r.Body.Close()

				data = map[string]interface{}{
					"status": "Fail",
					"payload": map[string]string{
						"message": "string",
						"code":    "string",
					},
				}

				if requestOrder.Operation == string(sdk.SELL) && requestOrder.Lots == 1 {
					data = map[string]interface{}{
						"trackingId": "dbb781ba4e984bd9",
						"status":     "Ok",
						"payload": map[string]interface{}{
							"orderId":       strconv.Itoa(orderID),
							"operation":     sdk.SELL,
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
				Expect(actual).Should(BeTemporally("~", time.Now(), time.Minute))

				actual, err = time.Parse("2006-01-02T15:04:05-07:00", v.Get("to"))
				Expect(err).NotTo(HaveOccurred())
				Expect(actual).Should(BeTemporally("~", time.Now(), 6*time.Minute))

				data = map[string]interface{}{
					"trackingId": "4f48d98e8040c23a",
					"status":     "Ok",
					"payload": map[string]interface{}{
						"operations": []interface{}{
							map[string]interface{}{
								"operationType":    sdk.SELL,
								"date":             "2021-03-01T23:39:29.507+03:00",
								"isMarginCall":     false,
								"instrumentType":   sdk.InstrumentTypeStock,
								"figi":             figi,
								"quantity":         1,
								"quantityExecuted": 1,
								"price":            96.50,
								"payment":          96.50,
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
										"price":    96.50,
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

	It("check slots", func() {
		type Slot struct {
			Id     graphql.String
			Ticker graphql.String
			Figi   graphql.String

			StartPrice  graphql.Float
			ChangePrice graphql.Float
			BuyingPrice graphql.Float
			TargetPrice graphql.Float
			Profit      graphql.Float

			Qty          graphql.Int
			AmountSpent  graphql.Float
			TargetAmount graphql.Float
			TotalProfit  graphql.Float
		}
		type Slots struct {
			Slots []Slot
		}

		reqSlots := Slots{}
		variables := map[string]interface{}{}
		err := client.Query(ctx, &reqSlots, variables)
		Expect(err).NotTo(HaveOccurred())

		for _, slot := range reqSlots.Slots {
			Expect(slot.Figi).NotTo(Equal(figi))
		}
	})

	It("check dealings", func() {
		type Deal struct {
			Id     graphql.String
			Ticker graphql.String
			Figi   graphql.String

			StartPrice  graphql.Float
			ChangePrice graphql.Float
			BuyingPrice graphql.Float
			TargetPrice graphql.Float
			Profit      graphql.Float

			SalePrice   graphql.Float
			Qty         graphql.Int
			AmountSpent graphql.Float

			AmountIncome graphql.Float
			TotalProfit  graphql.Float

			BuyAt    graphql.String
			Duration graphql.Int
			SellAt   graphql.String
		}
		type Dealings struct {
			Dealings []Deal
		}

		expected := Deal{
			Id:           "",
			Ticker:       graphql.String(ticker),
			Figi:         graphql.String(figi),
			StartPrice:   94,
			ChangePrice:  93,
			BuyingPrice:  95,
			TargetPrice:  95.77,
			Profit:       0.77,
			Qty:          1,
			AmountSpent:  95.5,
			AmountIncome: 96,
			SalePrice:    96.5,
			TotalProfit:  0.5,
		}

		reqDealings := Dealings{}
		variables := map[string]interface{}{}
		err := client.Query(ctx, &reqDealings, variables)
		Expect(err).NotTo(HaveOccurred())

		for idx := range reqDealings.Dealings {
			reqDealings.Dealings[idx].Id = ""
			reqDealings.Dealings[idx].BuyAt = ""
			reqDealings.Dealings[idx].SellAt = ""
			reqDealings.Dealings[idx].Duration = 0
		}

		Expect(reqDealings.Dealings).Should(ContainElement(expected))
	})
})
