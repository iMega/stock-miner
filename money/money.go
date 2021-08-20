package money

import "github.com/shopspring/decimal"

const (
	maxPrecision = 4
	precision    = 2
)

func Sum(a, b float64) float64 {
	dA := decimal.NewFromFloat(a).Truncate(maxPrecision)
	dB := decimal.NewFromFloat(b).Truncate(maxPrecision)

	f, _ := dA.Add(dB).Float64()

	return f
}

func Sub(a, b float64) float64 {
	dA := decimal.NewFromFloat(a).Truncate(maxPrecision)
	dB := decimal.NewFromFloat(b).Truncate(maxPrecision)

	f, _ := dA.Sub(dB).Round(precision).Float64()

	return f
}
