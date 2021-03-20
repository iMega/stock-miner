package broker

import (
	"context"

	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

type Pricer interface {
	Price(context.Context, sdk.Instrument) sdk.RestOrderBook
}
