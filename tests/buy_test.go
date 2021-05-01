package acceptance

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/imega/stock-miner/tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Automatically buy", func() {
	var (
		ctx    = context.Background()
		client = graphql.NewClient(GraphQLUrl, nil)
		figi   = "BBG000B9XRY4"
		ticker = "AAPL"
	)

	It("create creds", func() {
		defer GinkgoRecover()

		res, err := helpers.AddCredentials(
			GraphQLUrl,
			"test",
			"http://acceptance",
			"token",
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeTrue())
	})

	It("set settings slot", func() {
		defer GinkgoRecover()

		var req struct {
			Slot bool `graphql:"slot(global: $in)"`
		}
		type SlotSettingsInput map[string]interface{}

		variables := map[string]interface{}{
			"in": SlotSettingsInput{
				"volume":              10,
				"modificatorMinPrice": 2,
			},
		}
		err := client.Mutate(ctx, &req, variables)
		Expect(err).NotTo(HaveOccurred())

		Expect(req.Slot).To(BeTrue())
	})

	It("set settings commission and margin", func() {
		defer GinkgoRecover()

		var req struct {
			RulePrice bool `graphql:"rulePrice(global: $in)"`
		}
		type RulePriceInput map[string]interface{}

		variables := map[string]interface{}{
			"in": RulePriceInput{
				"marketCommission": 0.3,
				"grossMargin":      0.2,
			},
		}
		err := client.Mutate(ctx, &req, variables)
		Expect(err).NotTo(HaveOccurred())

		Expect(req.RulePrice).To(BeTrue())
	})

	It("add approved stock items", func() {
		defer GinkgoRecover()

		helpers.MockHTTPServer.AddHandler(func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{
				"status": "OK",
				"payload": map[string]interface{}{
					"total": 1,
					"instruments": []interface{}{
						map[string]interface{}{
							"ticker":            ticker,
							"figi":              figi,
							"isin":              "",
							"name":              "",
							"minPriceIncrement": 0.01,
							"lot":               1,
							"currency":          "USD",
							"type":              "Stock",
						},
					},
				},
			}
			b, _ := json.Marshal(data)

			w.Header().Add("Content-Type", "application/json")
			w.Write(b)
		})

		var reqMarketStockItems struct {
			MarketStockItems []struct {
				Ticker graphql.String
				Figi   graphql.String
			}
		}
		variables := map[string]interface{}{}
		err := client.Query(ctx, &reqMarketStockItems, variables)
		Expect(err).NotTo(HaveOccurred())

		type StockItemInput map[string]interface{}
		var reqAddStockItemApproved struct {
			AddStockItemApproved bool `graphql:"addStockItemApproved(items: $in)"`
		}
		variables = map[string]interface{}{
			"in": []StockItemInput{
				{
					"ticker":           reqMarketStockItems.MarketStockItems[0].Ticker,
					"figi":             reqMarketStockItems.MarketStockItems[0].Figi,
					"amountLimit":      0,
					"transactionLimit": 0,
				},
			},
		}
		err = client.Mutate(ctx, &reqAddStockItemApproved, variables)
		Expect(err).NotTo(HaveOccurred())

		Expect(reqAddStockItemApproved.AddStockItemApproved).To(BeTrue())
	})

	It("start mining", func() {
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
				err := r.Body.Close()
				Expect(err).NotTo(HaveOccurred())

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
								"quantity":         4,
								"quantityExecuted": 4,
								"price":            30.09,
								"payment":          -120.36,
								"currency":         sdk.USD,
								"status":           sdk.OperationStatusDone,
								"id":               strconv.Itoa(orderID),
								"commission": map[string]interface{}{
									"currency": "USD",
									"value":    -0.36,
								},
								"trades": []interface{}{
									map[string]interface{}{
										"tradeId":  "3535068930",
										"date":     "2021-03-01T23:39:29.507+03:00",
										"quantity": 1,
										"price":    30.09,
									},
									// map[string]interface{}{
									// 	"tradeId":  "3535068940",
									// 	"date":     "2021-03-01T23:39:29.507+03:00",
									// 	"quantity": 1,
									// 	"price":    30.09,
									// },
									// map[string]interface{}{
									// 	"tradeId":  "3535068950",
									// 	"date":     "2021-03-01T23:39:29.507+03:00",
									// 	"quantity": 2,
									// 	"price":    30.09,
									// },
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

		var reqGlobalMiningStart struct {
			GlobalMiningStart bool `graphql:"globalMiningStart"`
		}
		variables := map[string]interface{}{}
		err := client.Mutate(ctx, &reqGlobalMiningStart, variables)
		Expect(err).NotTo(HaveOccurred())

		<-time.After(13 * time.Second)

		var reqGlobalMiningStop struct {
			GlobalMiningStop bool `graphql:"globalMiningStop"`
		}
		variables = map[string]interface{}{}
		err = client.Mutate(ctx, &reqGlobalMiningStop, variables)
		Expect(err).NotTo(HaveOccurred())
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

		expected := Slots{
			Slots: []Slot{
				{
					Id:           "",
					Ticker:       "AAPL",
					Figi:         "BBG000B9XRY4",
					StartPrice:   95,
					ChangePrice:  94,
					BuyingPrice:  30.09,
					TargetPrice:  30.33,
					Profit:       0.24,
					Qty:          1,
					AmountSpent:  120.72,
					TargetAmount: 30.33,
					TotalProfit:  0,
				},
			},
		}

		reqSlots := Slots{}
		variables := map[string]interface{}{}
		err := client.Query(ctx, &reqSlots, variables)
		Expect(err).NotTo(HaveOccurred())

		for idx := range reqSlots.Slots {
			reqSlots.Slots[idx].Id = ""
		}

		Expect(reqSlots).To(Equal(expected))
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

		expected := Dealings{
			Dealings: []Deal{
				{
					Id:          "",
					Ticker:      "AAPL",
					Figi:        "BBG000B9XRY4",
					StartPrice:  95,
					ChangePrice: 94,
					BuyingPrice: 30.09,
					TargetPrice: 30.33,
					Profit:      0.24,
					Qty:         1,
					AmountSpent: 120.72,
				},
			},
		}

		reqDealings := Dealings{}
		variables := map[string]interface{}{}
		err := client.Query(ctx, &reqDealings, variables)
		Expect(err).NotTo(HaveOccurred())

		for idx := range reqDealings.Dealings {
			reqDealings.Dealings[idx].Id = ""
			reqDealings.Dealings[idx].BuyAt = ""
			reqDealings.Dealings[idx].SellAt = ""
		}

		Expect(reqDealings).To(Equal(expected))
	})
})
