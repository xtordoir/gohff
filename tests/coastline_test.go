package tests

import (
	"testing"

	"github.com/xtordoir/gohff/trader"
)

func TestCoastlineAccumulation(t *testing.T) {

	var price = 1.0

	trade := trader.NewCoastlineTrade(1000, 1, 1.0, 5)
	trade.Update(price, true)

	if trade.Trader.Exposure != 1000 {
		t.Errorf("Opened Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 1000)
	}

	price *= 0.995
	trade.Update(price, false)
	if trade.Trader.Exposure != 1000 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 1000)
	}

	price *= 0.99
	trade.Update(price, false)
	if trade.Trader.Exposure != 2000 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 2000)
	}

	price *= 0.99
	trade.Update(price, false)
	if trade.Trader.Exposure != 3000 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 3000)
	}

	price *= 0.99
	trade.Update(price, false)
	if trade.Trader.Exposure != 4000 {
		t.Errorf("I = 0 Trade Exposure is wrong: %d, want: %d.", trade.Trader.Exposure, 4000)
	}
	price *= 1.01
	trade.Update(price, false)
	if trade.Trader.Exposure != 4000 {
		t.Errorf("I = 0 Trade Exposure is wrong: %d, want: %d.", trade.Trader.Exposure, 4000)
	}
	price *= 0.99
	trade.Update(price, false)
	if trade.Trader.Exposure != 4000 {
		t.Errorf("I = 0 Trade Exposure is wrong: %d, want: %d.", trade.Trader.Exposure, 4000)
	}
	price *= 0.99
	trade.Update(price, false)
	if trade.Trader.Exposure != 5000 {
		t.Errorf("I = 0 Trade Exposure is wrong: %d, want: %d.", trade.Trader.Exposure, 5000)
	}
}
