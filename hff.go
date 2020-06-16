package gohff

import (
	"fmt"
	"math"
)

/*
case class Overshoot(instrument: String, scale: Double, direction: Int, prevExt: Double, ext: Double,
		maxOS: Double, extDist: Double)
*/

// Overshoot gives the internal state of the Overshoot: DirectionChange and current price and Extremum of current move
type Overshoot struct {
	Instrument    string  `json:"instrument"`
	Scale         float64 `json:"scale"`
	Direction     int     `json:"direction"`
	StartExtremum float64 `json:"start"`
	PeakExtremum  float64 `json:"peak"`
	Current       float64 `json:"current"`
}

// MaxOS returns the maximal Overshoot (Total Movement to be precise)
// It has thre sign of the Dirction
func (ovs *Overshoot) MaxOS() float64 {
	return 100 * (ovs.PeakExtremum - ovs.StartExtremum) / ovs.StartExtremum / ovs.Scale
}

// Reset resets the pricing info to provided price and direction to +1
func (ovs *Overshoot) Reset(x float64) {
	ovs.StartExtremum = x
	ovs.PeakExtremum = x
	ovs.Current = x
	ovs.Direction = 1
}

// Update takes a price x and returns an updated Overshoot object
func (ovs *Overshoot) Update(x float64) Overshoot {
	new := *ovs

	cos := 100 * (x - ovs.StartExtremum) / ovs.StartExtremum / ovs.Scale
	eDist := 100 * (x - ovs.PeakExtremum) / ovs.PeakExtremum / ovs.Scale
	maxOS := ovs.MaxOS()
	// if reversal...
	//fmt.Printf("COS: %.5f - EDIST: %.5f - MAXOS: %.5f", cos, eDist, maxOS)
	if cos*eDist < 0 && math.Abs(eDist) > 1.0 {
		new.Direction = -new.Direction
		new.StartExtremum = ovs.PeakExtremum
		new.PeakExtremum = x
		new.Current = x
		return new
	} else if math.Abs(cos) > math.Abs(maxOS) {
		new.PeakExtremum = x
	}
	new.Current = x
	return new
}

// OHLC contains a single simple candle data OHLC
type OHLC struct {
	O float64
	H float64
	L float64
	C float64
}

// UpdateWithCandle takes a candle x and returns an updated Overshoot object
// the firrst one is at the selected intraday extremum, the second one at the close price
// it is approximate as we assume either low or high is the relevant extremum
func (ovs *Overshoot) UpdateWithCandle(x OHLC) (Overshoot, Overshoot) {
	new := *ovs

	cosLow := 100 * (x.L - ovs.StartExtremum) / ovs.StartExtremum / ovs.Scale
	cosHigh := 100 * (x.H - ovs.StartExtremum) / ovs.StartExtremum / ovs.Scale
	eDistLow := 100 * (x.L - ovs.PeakExtremum) / ovs.PeakExtremum / ovs.Scale
	eDistHigh := 100 * (x.H - ovs.PeakExtremum) / ovs.PeakExtremum / ovs.Scale
	//maxOS := ovs.MaxOS()
	fmt.Printf("DEBUG: %.2f     %.2f    %.2f    %.2f\n", cosLow, cosHigh, eDistLow, eDistHigh)
	// First check for direction extension UP
	if ovs.Direction == 1 && eDistHigh > 0 {
		new.PeakExtremum = x.H
		new.Current = x.H
		// Second, direction extension Down
	} else if ovs.Direction == -1 && eDistLow < 0 {
		new.PeakExtremum = x.L
		new.Current = x.L
		// Reversal From UP
	} else if ovs.Direction == 1 && ((cosLow*eDistLow < 0 && math.Abs(eDistLow) > 1.0) || cosLow < 0.0) {
		new.Direction = -new.Direction
		new.StartExtremum = ovs.PeakExtremum
		new.PeakExtremum = x.L
		new.Current = x.L
		// Reversel from DOWN
	} else if ovs.Direction == -1 && ((cosHigh*eDistHigh < 0 && math.Abs(eDistHigh) > 1.0) || cosHigh > 0.0) {
		new.Direction = -new.Direction
		new.StartExtremum = ovs.PeakExtremum
		new.PeakExtremum = x.H
		new.Current = x.H
	}
	close := new

	return new, close.Update(x.C)
}

// InitStateWithCandles tries to find direction and extremum from candle data
func (ovs *Overshoot) InitStateWithCandles(x []OHLC) error {
	i := len(x) - 1
	// we start with the latest available price as max and min
	low, high := x[i].C, x[i].C
	last := x[i].C
	// direction is unknown, extremum is unknown
	dir := 0
	startExtremum := -1.0
	peakExtremum := -1.0
	prevReversal := -1.0
	// first high or min farther than scale distance is defining a direction
	for i >= 0 {
		//fmt.Println(x[i])

		// update high and low
		if x[i].H > high {
			high = x[i].H
			dir = -1
		}
		if x[i].L < low {
			low = x[i].L
			dir++
			if dir > 1 {
				dir = 1
			}
		}
		// test if a direction is found
		// note that dow direction tends has precedence,
		// hence the need to use short time frames to make sure ranges are always lower than
		// the scale
		if dir == -1 && 100*(low-high)/high/ovs.Scale < -1 {
			dir = -1
			startExtremum = high
			peakExtremum = low
			high = low
			low = high
			//fmt.Println("found down direction")
			//fmt.Println(x[i])
			i--
			break
		}

		if dir == 1 && 100*(high-low)/low/ovs.Scale > 1 {
			dir = 1
			startExtremum = low
			peakExtremum = high
			high = low
			low = high
			//fmt.Println("found up direction")
			//fmt.Println(x[i])
			i--
			break
		}
		i--
	}
	// Then next reversal determines amplitude
	for i >= 0 && prevReversal <= 0 {
		// update high and low
		if dir == -1 {
			if x[i].H > high {
				high = x[i].H
				low = x[i].O
			} else if x[i].L < low {
				low = x[i].L
			}
		}
		if dir == 1 {
			if x[i].L < low {
				low = x[i].L
				high = x[i].O
			} else if x[i].H > high {
				high = x[i].H
			}
		}

		// if direction is up, we seek a
		if dir == 1 && 100*(low-high)/high/ovs.Scale < -1.0 {
			prevReversal = high
			startExtremum = low
			//fmt.Println("found reversal")
			//fmt.Println(x[i])
		}
		if dir == -1 && 100*(high-low)/low/ovs.Scale > 1.0 {
			prevReversal = low
			startExtremum = high
			//fmt.Println("found reversal direction")
			//fmt.Println(x[i])
		}
		i--
	}
	if prevReversal <= 0 || dir == 0 {
		return fmt.Errorf("No Overshoot DC found")
	}
	ovs.PeakExtremum = peakExtremum
	//prevReversal
	ovs.StartExtremum = startExtremum
	ovs.Direction = dir
	ovs.Current = last
	return nil
}
