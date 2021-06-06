package broker

import "github.com/shopspring/decimal"

func calcSub(a, b float64) float64 {
	r, _ := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).Float64()

	return r
}

// формула расчета целевой цены для продажи
//
// ценаПокупки+(ценаПокупки/100*комиссия) = затраты
// затраты + (затраты / 100 * маржа%) = ЦенаПродажиБезКомиссии
// ЦенаПродажиБезКомиссии+(ЦенаПродажиБезКомиссии/100*комиссия) = ЦенаПродажи
func calcTargetPrice(commission, buyingPrice, margin float64) float64 {
	c := decimal.NewFromFloat(commission)
	bp := decimal.NewFromFloat(buyingPrice)
	m := decimal.NewFromFloat(margin)

	spent := bp.Add(bp.Div(decimal.NewFromInt(100)).Mul(c).Round(2))
	gm := spent.Add(spent.Div(decimal.NewFromInt(100)).Mul(m).Round(2))

	target, _ := gm.Add(gm.Div(decimal.NewFromInt(100)).Mul(c).Round(2)).Float64()

	return target
}
