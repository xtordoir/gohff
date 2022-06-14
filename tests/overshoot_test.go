package tests

import (
	"testing"
  "github.com/xtordoir/gohff"
)

func TestOvershootReversal(t *testing.T) {

  var os = gohff.Overshoot{
		Instrument: "XXXYYY",
		Scale: 1.0,
	}

	var price = 1.0


	os.Reset(price)

	if os.Direction != 1 {
		t.Errorf("Inital Direction is not 1 but %d", os.Direction)
	}

	os = os.Update(1.02)
  os = os.Update(1.015)
	if os.Direction != 1 {
		t.Errorf("Inital Direction is not 1 but %d", os.Direction)
	}
	os = os.Update(1.005)
	if os.Direction != -1 {
		t.Errorf("Inital Direction is not -1 but %+v", os)
	}
	// testing a large reversal
	os = os.Update(1.0)
	os = os.Update(1.02)
	os = os.Update(1.019)
  os = os.Update(0.99)
	if os.Direction != -1 {
		t.Errorf("Inital Direction is not -1 but %+v", os)
	}

	os = os.Update(1.0)
	os = os.Update(0.98)
	os = os.Update(0.981)
  os = os.Update(1.02)
	if os.Direction != 1 {
		t.Errorf("Inital Direction is not -1 but %+v", os)
	}
}
