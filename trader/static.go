package trader

//
//import (
//	"fmt"
//	"math"
//
//	"github.com/xtordoir/gohff"
//)
//
//// OSParameters of an OvershootTrade
//type StaticParameters struct {
//	// Inital size (Positive number)
//	U0 int
//	// Direction of Trade: +1 for Long, -1 for Short
//	Dir   int
//	Scale float64
//	Lmax  float64 // positive number for largest scale extension before stop
//	// Computed from Scale and U0
//	Target  float64
//	LossMax float64
//
//	// the state of the trade
//	X0 float64 // initial price
//}
//
//func toUnits(x float64) int64 {
//	y, _ := math.Modf(x)
//	return int64(y)
//}
//
//// StaticTrade implementation
//type StaticTrade struct {
//	Params *StaticParameters
//	Trader *gohff.TraderProcess
//}
//
//// NewStaticTrade to create a NewOvershootTrade
//func NewStaticTrade(u0 int, dir int, scale float64, lmax float64) *StaticTrade {
//	params := &OSParameters{
//		U0:      u0,
//		Dir:     dir,
//		Scale:   scale,
//		Target:  float64(u0) * scale / 100.0,
//		LossMax: -float64(u0) * scale * lmax / 100.0,
//		Lmax:    lmax,
//		X0:      -1.0,
//	}
//	trader := &gohff.TraderProcess{
//		Exposure:     0,
//		PriceAverage: 0.0,
//		CumProfit:    0.0,
//	}
//	return &OvershootTrade{
//		Params: params,
//		Trader: trader,
//	}
//}
//
//// Update function to update the trader with new l value (computed from frice elsewhere)
//// returns true if the trade was changed
//func (trade *StaticTrade) Update(price float64, l float64, init bool) (bool, error) {
//
//	trader := trade.Trader
//	p := trade.Params
//	// entering the position
//	if init {
//		trader.IncreaseBy(price, p.Dir*p.U0)
//		p.X0 = price
//	}
//
//	// if PL is above target or reversal or Stop  we must liquidate
//	pl := trader.TotalProfit(price)
//	if pl > p.Target || pl < p.LossMax {
//		fmt.Printf("Closing: Decrease Exposure by: %d\n", trader.Exposure)
//		trader.DecreaseBy(price, trader.Exposure)
//		return true, nil
//	}
//	// else, leave state alone
//	return false, nil
//}
