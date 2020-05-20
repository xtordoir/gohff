package gohff

// Spectrum is a range of Overshoots
type Spectrum struct {
	Scales     []float64   `json:"scales"`
	Overshoots []Overshoot `json:"overshoots"`
}

// Update all Overshoots Scales
func (ovs *Spectrum) Update(x float64) Spectrum {
	scales := ovs.Scales
	overshoots := make([]Overshoot, len(scales))
	for i := range overshoots {
		newos := ovs.Overshoots[i].Update(x)
		overshoots[i] = newos
	}
	return Spectrum{
		Scales:     scales,
		Overshoots: overshoots,
	}
}
