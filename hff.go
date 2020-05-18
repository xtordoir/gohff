package gohff

import (
  "math"
)
/*
case class Overshoot(instrument: String, scale: Double, direction: Int, prevExt: Double, ext: Double,
		maxOS: Double, extDist: Double)
*/

// Overshoot gives the internal state of the Overshoot: DirectionChange and current price and Extremum of current move
type Overshoot struct {
  Instrument string
  Scale float64
  Direction int
  StartExtremum float64
  PeakExtremum float64
  Current float64
}

// MaxOS returns the maximal Overshoot (Total Movement to be precise)
func (ovs *Overshoot) MaxOS() float64 {
  return 100*(ovs.PeakExtremum - ovs.StartExtremum)/ovs.StartExtremum/ovs.Scale
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

  cos := 100*(x - ovs.StartExtremum)/ovs.StartExtremum/ovs.Scale
  eDist := 100*(x - ovs.PeakExtremum)/ovs.PeakExtremum/ovs.Scale
  maxOS := ovs.MaxOS()
  // if reversal...
  //fmt.Printf("COS: %.5f - EDIST: %.5f - MAXOS: %.5f", cos, eDist, maxOS)
  if (cos*eDist < 0 && math.Abs(eDist) > 1.0) {
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
