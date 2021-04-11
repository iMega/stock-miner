package helpers

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

func AddCredentials(graphQLUrl, name, apiUrl, token string) (bool, error) {
	var (
		req struct {
			MarketCredentials bool `graphql:"marketCredentials(creds: $in)"`
		}
		ctx    = context.Background()
		client = graphql.NewClient(graphQLUrl, nil)
	)

	type MarketCredentialsInput map[string]interface{}

	variables := map[string]interface{}{
		"in": MarketCredentialsInput{
			"name":   name,
			"apiUrl": apiUrl,
			"token":  token,
		},
	}

	err := client.Mutate(ctx, &req, variables)

	return req.MarketCredentials, err
}
