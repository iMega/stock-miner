package helpers

import (
	"context"

	"github.com/hasura/go-graphql-client"
	. "github.com/onsi/gomega"
)

func AddStockItemApproved(graphQLUrl, ticker, figi string) {
	type StockItemInput map[string]interface{}
	var (
		reqAddStockItemApproved struct {
			AddStockItemApproved bool `graphql:"addStockItemApproved(items: $in)"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, nil)
	)
	variables := map[string]interface{}{
		"in": []StockItemInput{
			{
				"ticker":           ticker,
				"figi":             figi,
				"amountLimit":      0,
				"transactionLimit": 0,
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
		client = graphql.NewClient(graphQLUrl, nil)
	)
	variables := map[string]interface{}{
		"in": []StockItemInput{
			{
				"ticker":           ticker,
				"figi":             "",
				"amountLimit":      0,
				"transactionLimit": 0,
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
}

func StockItemApproved(graphQLUrl string) []StockItem {
	var (
		ctx                  = context.Background()
		client               = graphql.NewClient(graphQLUrl, nil)
		reqStockItemApproved struct {
			StockItemApproved []StockItem
		}
	)

	variables := map[string]interface{}{}
	err := client.Query(ctx, &reqStockItemApproved, variables)
	Expect(err).NotTo(HaveOccurred())

	return reqStockItemApproved.StockItemApproved
}
