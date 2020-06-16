package trader

import (
	"fmt"
	"math"

	"github.com/xtordoir/gohff"
)

// Parameters of a CoastlineTrade
type Parameters struct {
	// Inital size (Positive number)
	V0 int
	// Direction of Trade: +1 for Long, -1 for Short
	Dir   int
	Scale float64
	// Computed frm Scale and V0
	Target float64
	Ic     int

	// the state of the trade
	X0         float64 // initial price
	I          int     // the current extention of the market
	LastCrossI int
}

func (p *Parameters) nextI(x float64) int {
	// compute direction move size with sign
	z := -100.0 * float64(p.Dir) * (x - p.X0) / p.X0 / p.Scale
	// extract integer part
	y, _ := math.Modf(z)
	return int(y)
}

// CoastlineTrade implementation
type CoastlineTrade struct {
	Params *Parameters
	Trader *gohff.TraderProcess
}

// NewCoastlineTrade to create a CoastlineTrade
func NewCoastlineTrade(v0 int, dir int, scale float64, ic int) *CoastlineTrade {
	params := &Parameters{
		V0:     v0,
		Dir:    dir,
		Scale:  scale,
		Target: float64(v0) * scale / 100.0,
		Ic:     ic,
		X0:     -1.0,
		I:      -1,
	}
	trader := &gohff.TraderProcess{
		Exposure:     0,
		PriceAverage: 0.0,
		CumProfit:    0.0,
	}
	return &CoastlineTrade{
		Params: params,
		Trader: trader,
	}
}

// Update function to update the trader with new price
func (trade *CoastlineTrade) Update(price float64, init bool) (float64, error) {

	trader := trade.Trader
	p := trade.Params
	// entering the position
	if init {
		trader.IncreaseBy(price, p.Dir*p.V0)
		p.I = 0
		p.LastCrossI = 0
		p.X0 = price
	}

	// compute the next state:
	i := p.nextI(price)

	// if PL is above target we must liquidate
	pl := trader.TotalProfit(price)
	if pl > p.Target {
		fmt.Printf("Closing: Decrease Exposure by: %d\n", trader.Exposure)
		trader.DecreaseBy(price, trader.Exposure)
		return trader.CumProfit, nil
	}
	// if extention has jumped higher
	// and this cross is at a different level than last one
	if i > p.I && p.LastCrossI != i {
		fmt.Printf("Trading: Increase Exposure by: %d\n", p.Dir*p.V0*(i-p.I))
		trader.IncreaseBy(price, p.Dir*p.V0*(i-p.I))
		p.I = i
		p.LastCrossI = i
		return trader.CumProfit, nil
	}
	// if extension has jumped lower by at least one step
	if i <= p.I-1 {
		fmt.Printf("Trading: Decrease Exposure by: %d\n", p.Dir*p.V0*(p.I-i-1))
		trader.DecreaseBy(price, p.Dir*p.V0*(p.I-i-1))
		p.I = i
		p.LastCrossI = i + 1 // we set the last cross level right above
		return trader.CumProfit, nil
	}
	// else, leave state alone
	return trader.TotalProfit(price), nil
}
