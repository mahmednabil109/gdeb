package blockchain

import (
	"sync"
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
