package blockchain

import (
	"fmt"
	"sync"

	"github.com/mahmednabil109/gdeb/data"
)

type StakeDistribution struct {
	TotalCoins   uint64
	Distribution map[string]float64
	Mux          sync.Mutex
}

func NewStakeDist(total uint64, dist map[string]float64) StakeDistribution {
	return StakeDistribution{
		TotalCoins:   total,
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
func (sd *StakeDistribution) Valid(trans data.Transaction) (bool, error) {
	sd.Mux.Lock()
	defer sd.Mux.Unlock()

	amount, from := float64(trans.Amount), trans.From
	val, ok := sd.Distribution[from]
	if !ok {
		return false, fmt.Errorf("Not enough money!")
	} else if val*float64(sd.TotalCoins) < amount {
		return false, fmt.Errorf("Sender %s does not have enough money of amount %f", from, amount)
	}
	return true, nil
}

func (sd *StakeDistribution) Update(trans data.Transaction) error {
	sd.Mux.Lock()
	defer sd.Mux.Unlock()
	amount, from, to := float64(trans.Amount), trans.From, trans.To

	sd.Distribution[from] -= amount
	sd.Distribution[to] += amount

	return nil
}
