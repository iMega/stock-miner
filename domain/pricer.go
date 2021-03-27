package domain

import "context"

type Pricer interface {
	GetPrice(context.Context, PriceReceiptMessage) (PriceReceiptMessage, error)
}
