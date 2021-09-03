package helpers

import (
	"context"

	"github.com/hasura/go-graphql-client"
	. "github.com/onsi/gomega"
)

func AddStockItemApproved(graphQLUrl, ticker, figi string, maxPrice float64) {
	type StockItemInput map[string]interface{}
	var (
		reqAddStockItemApproved struct {
			AddStockItemApproved bool `graphql:"addStockItemApproved(items: $in)"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, GetHTTPClient())
	)
	variables := map[string]interface{}{
		"in": []StockItemInput{
			{
				"ticker":           ticker,
				"figi":             figi,
				"amountLimit":      0,
				"transactionLimit": 0,
				"currency":         "USD",
				"startTime":        0,
				"endTime":          24,
				"maxPrice":         maxPrice,
				"active":           true,
			},
		},
	}
	err := client.Mutate(ctx, &reqAddStockItemApproved, variables)
	Expect(err).NotTo(HaveOccurred())

	Expect(reqAddStockItemApproved.AddStockItemApproved).To(BeTrue())
}

func RemoveStockItemApproved(graphQLUrl, ticker string) {
	type StockItemInput map[string]interface{}
	var (
		reqRemoveStockItemApproved struct {
			RemoveStockItemApproved bool `graphql:"removeStockItemApproved(items: $in)"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, GetHTTPClient())
	)
	variables := map[string]interface{}{
		"in": []StockItemInput{
			{
				"ticker":           ticker,
				"figi":             "",
				"amountLimit":      0,
				"transactionLimit": 0,
				"currency":         "USD",
				"startTime":        0,
				"endTime":          0,
				"maxPrice":         0,
				"active":           true,
			},
		},
	}
	err := client.Mutate(ctx, &reqRemoveStockItemApproved, variables)
	Expect(err).NotTo(HaveOccurred())

	Expect(reqRemoveStockItemApproved.RemoveStockItemApproved).To(BeTrue())
}

type StockItem struct {
	Ticker           graphql.String
	Figi             graphql.String
	AmountLimit      graphql.Float
	TransactionLimit graphql.Float
	MaxPrice         graphql.Float
	Active           graphql.Boolean
}

func StockItemApproved(graphQLUrl string) []StockItem {
	var (
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, GetHTTPClient())

		reqStockItemApproved struct {
			StockItemApproved []StockItem
		}
	)

	variables := map[string]interface{}{}
	err := client.Query(ctx, &reqStockItemApproved, variables)
	Expect(err).NotTo(HaveOccurred())

	return reqStockItemApproved.StockItemApproved
}

func EnableStockItemApproved(graphQLUrl string) {
	var (
		req struct {
			EnableStockItemsApproved bool `graphql:"enableStockItemsApproved"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, GetHTTPClient())
	)
	variables := map[string]interface{}{}
	err := client.Mutate(ctx, &req, variables)
	Expect(err).NotTo(HaveOccurred())
}

func DisableStockItemApproved(graphQLUrl string) {
	var (
		req struct {
			DisableStockItemsApproved bool `graphql:"disableStockItemsApproved"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, GetHTTPClient())
	)
	variables := map[string]interface{}{}
	err := client.Mutate(ctx, &req, variables)
	Expect(err).NotTo(HaveOccurred())
}
