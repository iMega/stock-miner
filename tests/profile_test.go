package acceptance

import (
	"context"

	graphql "github.com/hasura/go-graphql-client"
	"github.com/imega/stock-miner/tests/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Profile", func() {
	var (
		ctx    = context.Background()
		client = graphql.NewClient(GraphQLUrl, helpers.GetHTTPClient())

		expected = struct {
			Name   graphql.String
			Token  graphql.String
			ApiUrl graphql.String
		}{
			Name:   "creds",
			Token:  "token",
			ApiUrl: "apiUrl",
		}
	)

	It("set creds", func() {
		defer GinkgoRecover()

		res, err := helpers.AddCredentials(
			GraphQLUrl,
			string(expected.Name),
			string(expected.ApiUrl),
			string(expected.Token),
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeTrue())
	})

	It("getting creds", func() {
		defer GinkgoRecover()
		var req struct {
			Settings struct {
				MarketCredentials []struct {
					Name   graphql.String
					Token  graphql.String
					ApiUrl graphql.String
				} `graphql:"marketCredentials"`
			}
		}
		variables := map[string]interface{}{}
		err := client.Query(ctx, &req, variables)
		Expect(err).NotTo(HaveOccurred())

		Expect(req.Settings.MarketCredentials).Should(ContainElement(expected))
	})
})
