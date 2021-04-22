package acceptance

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
			"in": SlotSettingsInput{"volume": 10},
		}
		err := client.Mutate(ctx, &req, variables)
		Expect(err).NotTo(HaveOccurred())

		Expect(req.Slot).To(BeTrue())
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
		// Expect(err).NotTo(HaveOccurred())

		// Expect(reqAddStockItemApproved.AddStockItemApproved).To(BeTrue())
	})

	It("start mining", func() {
		var (
			requestCount int
			startPrice   = 100
		)

		defer GinkgoRecover()

		helpers.MockHTTPServer.AddHandler(func(w http.ResponseWriter, r *http.Request) {
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
						"status": "Ok",
						"payload": map[string]interface{}{
							"executedLots": 1,
							// orderId
							// operation
							// status
							// rejectReason
							// requestedLots
							// executedLots
							// commission
							// message
						},
					}
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

		<-time.After(14 * time.Second)

		var reqGlobalMiningStop struct {
			GlobalMiningStop bool `graphql:"globalMiningStop"`
		}
		variables = map[string]interface{}{}
		err = client.Mutate(ctx, &reqGlobalMiningStop, variables)
		Expect(err).NotTo(HaveOccurred())
	})
})
