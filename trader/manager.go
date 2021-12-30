package trader

import (
	"encoding/json"
	"io"
	"log"
	"sync"
)

// Manager allows to safely bulk-update trades
type Manager struct {
	sync.RWMutex
	// a positive value means the trade is active, zero is not yet initialized (waiting for a start price),
	// negative means the trade is done
	OSTrades map[*OvershootTrade]int
	CLTrades map[*CoastlineTrade]int
}

// NewManager is the constructor
func NewManager() *Manager {
	return &Manager{
		OSTrades: make(map[*OvershootTrade]int),
		CLTrades: make(map[*CoastlineTrade]int),
	}
}

// AddOvershootTrade a trade at this entry price, and paramters
func (manager *Manager) AddOvershootTrade(u0 int, dir int, scale float64, lmax float64) {
	manager.Lock()
	defer manager.Unlock()
	trade := NewOvershootTrade(u0, dir, scale, lmax)
	manager.OSTrades[trade] = 0
}

// AddCoastlineTrade a trades at this entry price, and paramters
func (manager *Manager) AddCoastlineTrade(v0 int, dir int, scale float64, ic int) {
	manager.Lock()
	defer manager.Unlock()
	trade := NewCoastlineTrade(v0, dir, scale, ic)
	manager.CLTrades[trade] = 0
}

// Update all trades with a new price, for a spectrum of l values (scale -> L)
func (manager *Manager) Update(price float64, lls map[float64]float64) {
	manager.Lock()
	defer func() {
		manager.Unlock()
		manager.Cleanup()
	}()

	for trade, state := range manager.OSTrades {
		scale := trade.Params.Scale
		l, ok := lls[scale]
		if !ok {
			continue
		}
		if state > 0 {
			changed, _ := trade.Update(price, l, false)
			if changed {
				log.Printf("Changed trade %+v - %+v at price %.5f L = %.2f", trade.Params, trade.Trader, price, l)
			}
		} else if state == 0 {
			trade.Update(price, l, true)
			manager.OSTrades[trade] = 1
			log.Printf("Initiating trade %+v at price %.5f L = %.2f", trade.Params, price, l)
		}
	}

	for trade, state := range manager.CLTrades {
		if state > 0 {
			changed, _ := trade.Update(price, false)
			if changed {
				log.Printf("Changed trade %+v - %+v at price %.5f", trade.Params, trade.Trader, price)
			}
		} else if state == 0 {
			trade.Update(price, true)
			manager.CLTrades[trade] = 1
			log.Printf("Initiating trade %+v at price %.5f", trade.Params, price)
		}
	}

}

// Cleanup removes closed trades (Coastline Trades end if Exposure reaches 0)
func (manager *Manager) Cleanup() {
	manager.Lock()
	defer manager.Unlock()

	for trade := range manager.OSTrades {
		if trade.Trader.Exposure == 0.0 {
			log.Printf("Closed and removed trade: %+v", trade.Params)
			delete(manager.OSTrades, trade)
		}
	}
	for trade := range manager.CLTrades {
		if trade.Trader.Exposure == 0.0 {
			log.Printf("Closed and removed trade: %+v", trade.Params)
			delete(manager.CLTrades, trade)
		}
	}
}

// Exposure aggregated over the set of open Trades
func (manager *Manager) Exposure() int64 {
	manager.RLock()
	defer manager.RUnlock()
	e := 0
	for trade := range manager.OSTrades {
		e += trade.Trader.Exposure
	}
	for trade := range manager.CLTrades {
		e += trade.Trader.Exposure
	}
	return int64(e)
}

// CumProfit aggregated PL over the set of open Trades
func (manager *Manager) CumProfit() float64 {
	manager.RLock()
	defer manager.RUnlock()
	pl := 0.0
	for trade := range manager.OSTrades {
		pl += trade.Trader.CumProfit
	}
	for trade := range manager.CLTrades {
		pl += trade.Trader.CumProfit
	}
	return pl
}

// MultiTrades embeds OS and CL Trades for JSON conversion
type MultiTrades struct {
	OSTrades []OvershootTrade
	CLTrades []CoastlineTrade
}

// ToJSON Marshal Trades as Json Array
func (manager *Manager) ToJSON() (string, error) {
	manager.RLock()
	defer manager.RUnlock()

	oskeys := make([]OvershootTrade, len(manager.OSTrades))
	clkeys := make([]CoastlineTrade, len(manager.CLTrades))
	i := 0
	for k := range manager.OSTrades {
		oskeys[i] = *k
		i++
	}
	i = 0
	for k := range manager.CLTrades {
		clkeys[i] = *k
		i++
	}
	b, err := json.Marshal(MultiTrades{
		OSTrades: oskeys,
		CLTrades: clkeys,
	})

	return string(b), err
}

// ManagerFromJSON is the constructor from a JSON
func ManagerFromJSON(s io.Reader) *Manager {

	multi := MultiTrades{}
	ostrades := make(map[*OvershootTrade]int)
	cltrades := make(map[*CoastlineTrade]int)
	//var arr []trader.OvershootTrade
	/*err := */
	json.NewDecoder(s).Decode(&multi)

	for _, tr := range multi.OSTrades {
		lc := tr
		ostrades[&lc] = 1
	}
	for _, tr := range multi.CLTrades {
		lc := tr
		cltrades[&lc] = 1
	}

	return &Manager{
		OSTrades: ostrades,
		CLTrades: cltrades,
	}
}
