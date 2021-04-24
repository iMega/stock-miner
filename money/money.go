package money

import "github.com/shopspring/decimal"

func Sum(a, b float64) float64 {
	dA := decimal.NewFromFloat(a).Truncate(4)
	dB := decimal.NewFromFloat(b).Truncate(4)

	f, _ := dA.Add(dB).Float64()

	return f
}
