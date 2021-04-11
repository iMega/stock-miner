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
		client = graphql.NewClient(GraphQLUrl, nil)
	)

	Context("cred is empty", func() {
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

		Expect(req.Settings.MarketCredentials).To(BeEmpty())
	})

	Context("set creds", func() {
		defer GinkgoRecover()

		res, err := helpers.AddCredentials(GraphQLUrl, "name", "apiUrl", "token")
		Expect(err).NotTo(HaveOccurred())
		Expect(res).To(BeTrue())
	})

	Context("getting creds", func() {
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

		cred := req.Settings.MarketCredentials[0]
		Expect(cred.Token).To(Equal(graphql.String("token")))
		Expect(cred.ApiUrl).To(Equal(graphql.String("apiUrl")))
		Expect(cred.Name).To(Equal(graphql.String("name")))
	})
})
