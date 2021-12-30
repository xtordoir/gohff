package trader

import (
	"fmt"
	"math"

	"github.com/xtordoir/gohff"
)

// OSParameters of an OvershootTrade
type OSParameters struct {
	// Inital size (Positive number)
	U0 int
	// Direction of Trade: +1 for Long, -1 for Short
	Dir   int
	Scale float64
	// Computed frm Scale and U0
	Target float64
	Lmax   float64

	// the state of the trade
	X0 float64 // initial price
	L  float64 // the current extention of the market
	L0 float64 // initial extension of the market
}

func toUnits(x float64) int64 {
	y, _ := math.Modf(x)
	return int64(y)
}

func (p *OSParameters) nextExposure(l float64, minFrac float64) (int, bool) {
	// check if it is a reversal
	if l*p.L < 0 {
		return 0, true
	}
	// check if it is a stop loss
	if math.Abs(l) > math.Abs(p.Lmax) {
		return 0, true
	}
	// if it is not an extension of l
	if math.Abs(l)-math.Abs(p.L) < minFrac {
		return 0, false
	}

	dm := l - p.L0 // beware that this has a sign!
	posFactor := math.Exp(math.Abs(dm))
	if math.Abs(posFactor) < 1.0 || posFactor < 0.0 {
		posFactor = 0.0
	}
	return p.Dir * int(float64(p.U0)*posFactor), true
}

// OvershootTrade implementation
type OvershootTrade struct {
	Params *OSParameters
	Trader *gohff.TraderProcess
}

// NewOvershootTrade to create a NewOvershootTrade
func NewOvershootTrade(u0 int, dir int, scale float64, lmax float64) *OvershootTrade {
	params := &OSParameters{
		U0:     u0,
		Dir:    dir,
		Scale:  scale,
		Target: float64(u0) * scale / 100.0,
		Lmax:   lmax,
		X0:     -1.0,
		L:      0.0,
		L0:     0.0,
	}
	trader := &gohff.TraderProcess{
		Exposure:     0,
		PriceAverage: 0.0,
		CumProfit:    0.0,
	}
	return &OvershootTrade{
		Params: params,
		Trader: trader,
	}
}

// Update function to update the trader with new l value (computed from frice elsewhere)
// returns true if the trade was changed
func (trade *OvershootTrade) Update(price float64, l float64, init bool) (bool, error) {

	trader := trade.Trader
	p := trade.Params
	// entering the position
	if init {
		trader.IncreaseBy(price, p.Dir*p.U0)
		p.L = l
		p.X0 = price
		p.L0 = l
		// !!! Set the Lmax as distance from l0
		p.Lmax = p.Lmax + p.L0
	}

	// compute the next state:
	//i := p.nextI(price)
	targetExposure, changed := p.nextExposure(l, 0.1)

	// if PL is above target or reversal or Stop  we must liquidate
	pl := trader.TotalProfit(price)
	if pl > p.Target || l*p.L < 0 || math.Abs(l) > math.Abs(p.Lmax) {
		fmt.Printf("Closing: Decrease Exposure by: %d\n", trader.Exposure)
		trader.DecreaseBy(price, trader.Exposure)
		return true, nil
	}
	// if extention has jumped higher
	if changed && math.Abs(float64(targetExposure)) > math.Abs(float64(trader.Exposure)) {
		fmt.Printf("Trading: Increase Exposure by: %d\n", targetExposure-trader.Exposure)
		trader.IncreaseBy(price, targetExposure-trader.Exposure)
		p.L = l
		return true, nil
	}
	// else, leave state alone
	return false, nil
}
