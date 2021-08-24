package money

import "github.com/shopspring/decimal"

const (
	maxPrecision   = 4
	precision      = 2
	hundredPercent = 100
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

func Procent(number, procent float64) float64 {
	n := decimal.NewFromFloat(number)
	p := decimal.NewFromFloat(procent)
	hp := decimal.NewFromInt(hundredPercent)

	res, _ := n.Div(hp).Mul(p).Round(precision).Float64()

	return res
}
