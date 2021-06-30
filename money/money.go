package money

import "github.com/shopspring/decimal"

const maxPrecision = 4

func Sum(a, b float64) float64 {
	dA := decimal.NewFromFloat(a).Truncate(maxPrecision)
	dB := decimal.NewFromFloat(b).Truncate(maxPrecision)

	f, _ := dA.Add(dB).Float64()

	return f
}
