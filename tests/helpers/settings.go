package helpers

import (
	"context"

	"github.com/hasura/go-graphql-client"
	. "github.com/onsi/gomega"
)

func SetSlot(graphQLUrl string) {
	var (
		req struct {
			Slot bool `graphql:"slot(global: $in)"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, nil)
	)
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
}

func SetRulePrice(graphQLUrl string) {
	var (
		req struct {
			RulePrice bool `graphql:"rulePrice(global: $in)"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, nil)
	)
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
}
