package helpers

import (
	"context"

	"github.com/hasura/go-graphql-client"
	. "github.com/onsi/gomega"
)

func GlobalMiningStart(graphQLUrl string) {
	var (
		reqGlobalMiningStart struct {
			GlobalMiningStart bool `graphql:"globalMiningStart"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, nil)
	)
	variables := map[string]interface{}{}
	err := client.Mutate(ctx, &reqGlobalMiningStart, variables)
	Expect(err).NotTo(HaveOccurred())
}

func GlobalMiningStop(graphQLUrl string) {
	var (
		reqGlobalMiningStop struct {
			GlobalMiningStop bool `graphql:"globalMiningStop"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, nil)
	)
	variables := map[string]interface{}{}
	err := client.Mutate(ctx, &reqGlobalMiningStop, variables)
	Expect(err).NotTo(HaveOccurred())
}
