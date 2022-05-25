package blockchain

import (
	"fmt"
	"sync"

	"github.com/mahmednabil109/gdeb/data"
)

type StakeDistribution struct {
	// totalCoins uint64
	Distribution map[string]float64
	Mux          sync.Mutex
}

func New(dist map[string]float64) StakeDistribution {
	return StakeDistribution{
		Distribution: dist,
	}
}

func (sd *StakeDistribution) Get(key string) float64 {
	sd.Mux.Lock()
	defer sd.Mux.Unlock()

	val, ok := sd.Distribution[key]
	if ok {
		return val
	}
	return -1

}
func (sd *StakeDistribution) Update(trans data.Transaction) error {
	sd.Mux.Lock()
	defer sd.Mux.Unlock()

	amount, from, to := float64(trans.Amount), trans.From, trans.To

	val, ok := sd.Distribution[from]
	if !ok {
		return fmt.Errorf("Not enough money!")
	} else if val >= amount {
		return fmt.Errorf("Sender %s does not have enough money of amount %f", from, amount)
	}

	sd.Distribution[from] -= amount
	sd.Distribution[to] += amount

	return nil
}
