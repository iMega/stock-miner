package domain

import "context"

type Stack interface {
	Slot(context.Context) ([]Slot, error)
}
