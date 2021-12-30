package tests

import (
	"fmt"
	"testing"

	"github.com/xtordoir/gohff/trader"
)

func TestOvershootAccumulation(t *testing.T) {

	var price = 1.0
	var l = -1.0

	trade := trader.NewOvershootTrade(1000, 1, 1.0, -5.0)
	trade.Update(price, l, true)

	if trade.Trader.Exposure != 1000 {
		t.Errorf("Opened Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 1000)
	}

	price = 0.995
	l = -1.05
	trade.Update(price, l, false)
	if trade.Trader.Exposure != 1000 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 1000)
	}

	price = 0.999
	l = -1.1
	trade.Update(price, l, false)
	if trade.Trader.Exposure != 1105 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 1105)
	}

	price = 0.99
	l = -2.0
	trade.Update(price, l, false)
	if trade.Trader.Exposure != 2718 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 2718)
	}

	price = 0.99
	l = -2.0
	trade.Update(price, l, false)
	if trade.Trader.Exposure != 2718 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 2718)
	}

	price = 1.0
	l = 1.0
	trade.Update(price, l, false)
	if trade.Trader.Exposure != 0 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 0)
	}

	price = 1.0
	l = 1.0

	trade = trader.NewOvershootTrade(1000, -1, 1.0, -5.0)
	trade.Update(price, l, true)

	if trade.Trader.Exposure != -1000 {
		t.Errorf("Opened Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, -1000)
	}

	price = 1.005
	l = 1.05
	trade.Update(price, l, false)
	if trade.Trader.Exposure != -1000 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, -1000)
	}

	price = 1.001
	l = 1.1
	trade.Update(price, l, false)
	if trade.Trader.Exposure != -1105 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, -1105)
	}

	price = 1.01
	l = 2.0
	trade.Update(price, l, false)
	if trade.Trader.Exposure != -2718 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, -2718)
	}

	price = 1.01
	l = 2.0
	trade.Update(price, l, false)
	if trade.Trader.Exposure != -2718 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, -2718)
	}

	price = 1.0
	l = -1.0
	trade.Update(price, l, false)
	if trade.Trader.Exposure != 0 {
		t.Errorf("I = 0 Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 0)
	}

	price = 1.0
	l = 2.5

	trade = trader.NewOvershootTrade(1000, -1, 1.0, 3.6)
	trade.Update(price, l, true)
	fmt.Printf("%+v  %+v\n", trade.Params, trade.Trader)
	if trade.Trader.Exposure != -1000 {
		t.Errorf("Opened Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, -1000)
	}

	price = 1.037
	l = 6.7
	trade.Update(price, l, false)
	fmt.Printf("%+v  %+v\n", trade.Params, trade.Trader)
	if trade.Trader.Exposure != 0 {
		t.Errorf("Trade Exposure is wrong at price %.4f: %d, want: %d.", price, trade.Trader.Exposure, 0)
	}
}
