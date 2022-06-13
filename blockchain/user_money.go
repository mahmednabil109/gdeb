package blockchain

import (
	"sync"

	"github.com/mahmednabil109/gdeb/data"
)

type UserMoney struct {
	Money map[string]float64
	Mux   sync.Mutex
}

func NewUserMoney(stakeDist *StakeDistribution) *UserMoney {
	money := make(map[string]float64)
	total := stakeDist.TotalCoins
	for user, stake := range stakeDist.Distribution {
		m := stake * float64(total)
		money[user] = m
	}
	return &UserMoney{
		Money: money,
	}
}

func (u *UserMoney) HasEnoughMoney(user string, amount float64) bool {
	u.Mux.Lock()
	defer u.Mux.Unlock()

	val, ok := u.Money[user]
	if !ok {
		return false
	}
	if val >= amount {
		return true
	}
	return false

}
func (u *UserMoney) Update(transactions []data.Transaction) {
	u.Mux.Lock()
	defer u.Mux.Unlock()
	m := totalAmount(transactions)

	for user, amount := range m {
		u.Money[user] += amount
		u.Money[user] += amount
	}
}

func totalAmount(transactions []data.Transaction) map[string]float64 {
	m := make(map[string]float64)
	for _, t := range transactions {
		fee := 0.0
		if t.ConsumedGas != 0 {
			fee = float64(t.ConsumedGas) * 0.4
		}
		amount := float64(t.Amount)
		_, ok := m[t.From]
		if ok {
			m[t.From] -= (amount + fee)
		} else {
			m[t.From] = -1 * (amount + fee)
		}

		_, ok = m[t.To]
		if ok {
			m[t.To] += amount
		} else {
			m[t.To] = amount
		}
	}
	return m
}
