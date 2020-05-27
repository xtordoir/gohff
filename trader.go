package gohff

import "math"

// GearFunc is a function taking a price and computing a position change
type GearFunc func(x float64) (g float64)

// TraderProcess is used to know a managed trade process state
type TraderProcess struct {
	// exposure
	Exposure float64
	// average price of position
	PriceAverage float64
	// cumulated profit (Estimated)
	CumProfit float64
	// cumulated profit (Actual)
	ActualCumProfit float64
	// G is the Exposure increase (decrease) function, it depends basically on transaction prices
	GIncr GearFunc
	GDecr GearFunc
}

// Increase is used to increase the current position of the TraderProcess
func (tp *TraderProcess) Increase(x float64) {
	de := tp.GIncr(x)
	e := tp.Exposure + de
	a := (tp.PriceAverage*math.Abs(tp.Exposure) + x*math.Abs(de)) / math.Abs(e)
	tp.Exposure += e
	tp.PriceAverage = a
}

// IncreaseBy a number of units
func (tp *TraderProcess) IncreaseBy(x float64, units float64) {
	de := units
	e := tp.Exposure + de
	a := (tp.PriceAverage*math.Abs(tp.Exposure) + x*math.Abs(de)) / math.Abs(e)
	tp.Exposure = e
	tp.PriceAverage = a
}

// Decrease is used to Decrease the position on the TraderProcess
func (tp *TraderProcess) Decrease(x float64) {
	de := tp.GDecr(x)
	e := tp.Exposure - de
	pi := de * (x/tp.PriceAverage - 1.0)

	tp.Exposure = e
	tp.CumProfit += pi
}

// DecreaseBy a number of Units
func (tp *TraderProcess) DecreaseBy(x float64, units float64) {
	de := units
	e := tp.Exposure - de
	pi := de * (x/tp.PriceAverage - 1.0)

	tp.Exposure = e
	tp.CumProfit += pi
}

// TotalProfit compute the Process total profit for a given exit price
func (tp *TraderProcess) TotalProfit(x float64) float64 {
	return tp.CumProfit + tp.Exposure*(x/tp.PriceAverage-1.0)
}
