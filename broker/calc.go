package broker

import "github.com/shopspring/decimal"

func calcSub(a, b float64) float64 {
	r, _ := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).Float64()

	return r
}

const (
	hundredPercent = 100
	precision      = 2
)

// формула расчета целевой цены для продажи
//
// ценаПокупки+(ценаПокупки/100*комиссия) = затраты
// затраты + (затраты / 100 * маржа%) = ЦенаПродажиБезКомиссии
// ЦенаПродажиБезКомиссии+(ЦенаПродажиБезКомиссии/100*комиссия) = ЦенаПродажи.
func calcTargetPrice(commission, buyingPrice, margin float64) float64 {
	c := decimal.NewFromFloat(commission)
	bp := decimal.NewFromFloat(buyingPrice)
	m := decimal.NewFromFloat(margin)

	spent := bp.Add(bp.Div(decimal.NewFromInt(hundredPercent)).Mul(c).Round(precision))
	gm := spent.Add(spent.Div(decimal.NewFromInt(hundredPercent)).Mul(m).Round(precision))

	target, _ := gm.Add(gm.Div(decimal.NewFromInt(hundredPercent)).Mul(c).Round(precision)).Float64()

	return target
}
